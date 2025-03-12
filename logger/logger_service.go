package logger

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// AccessFile 访问日志配置
type AccessFile struct {
	File       string `toml:"file"`
	Enable     bool   `toml:"enable"`
	MaxSize    int    `toml:"max_size"`
	MaxBackups int    `toml:"max_backups"`
	MaxAge     int    `toml:"max_age"`
	Compress   bool   `toml:"compress"`
}

// LogEntry 定义日志条目结构
type LogEntry struct {
	Time        string         `json:"time"`
	Level       string         `json:"level"`
	Message     string         `json:"message"`
	ContextInfo map[string]any `json:"context_info,omitempty"`
	File        string         `json:"file,omitempty"`
	Line        int            `json:"line,omitempty"`
	Data        any            `json:"data,omitempty"`
}

var (
	once          = sync.Once{}
	DefaultLogger *AppLogger
)

// 修改配置结构体
type LogConfig struct {
	LogLevel  string     `toml:"log_level"`
	Log       LogFile    `toml:"log"`
	AccessLog AccessFile `toml:"access_log"`
}

type LogFile struct {
	BaseDir    string            `toml:"base_dir"`
	Levels     map[string]string `toml:"levels"`
	Console    bool              `toml:"console"`
	MaxSize    int               `toml:"max_size"`
	MaxBackups int               `toml:"max_backups"`
	MaxAge     int               `toml:"max_age"`
	Compress   bool              `toml:"compress"`
}

type AppLogger struct {
	loggers map[string]*logrus.Logger
}

// InitLogger 初始化日志配置
func InitLogger(ctx context.Context, configPath string) *AppLogger {
	if DefaultLogger != nil {
		return DefaultLogger
	}

	once.Do(func() {
		var config LogConfig
		if _, err := toml.DecodeFile(configPath, &config); err != nil {
			return
		}

		// 创建基础目录
		if err := os.MkdirAll(config.Log.BaseDir, 0755); err != nil {
			return
		}

		// 初始化不同级别的日志记录器
		loggers := make(map[string]*logrus.Logger)

		// 为每个级别创建logger
		for level, filename := range config.Log.Levels {
			logger := logrus.New()

			// 配置JSON格式输出
			logger.SetFormatter(&logrus.JSONFormatter{
				TimestampFormat: "2006-01-02 15:04:05",
				FieldMap: logrus.FieldMap{
					logrus.FieldKeyTime:  "time",
					logrus.FieldKeyLevel: "level",
					logrus.FieldKeyMsg:   "message",
				},
			})

			// 配置日志输出
			writer := &lumberjack.Logger{
				Filename:   filepath.Join(config.Log.BaseDir, filename),
				MaxSize:    config.Log.MaxSize,
				MaxBackups: config.Log.MaxBackups,
				MaxAge:     config.Log.MaxAge,
				Compress:   config.Log.Compress,
			}

			if config.Log.Console {
				logger.SetOutput(io.MultiWriter(writer, os.Stdout))
			} else {
				logger.SetOutput(writer)
			}

			// 设置日志级别
			switch level {
			case "debug":
				logger.SetLevel(logrus.DebugLevel)
			case "info":
				logger.SetLevel(logrus.InfoLevel)
			case "warn":
				logger.SetLevel(logrus.WarnLevel)
			case "error":
				logger.SetLevel(logrus.ErrorLevel)
			case "fatal":
				logger.SetLevel(logrus.FatalLevel)
			}

			loggers[level] = logger
		}

		// 初始化访问日志
		if config.AccessLog.Enable {
			accessLogger := logrus.New()
			writer := &lumberjack.Logger{
				Filename:   config.AccessLog.File,
				MaxSize:    config.AccessLog.MaxSize,
				MaxBackups: config.AccessLog.MaxBackups,
				MaxAge:     config.AccessLog.MaxAge,
				Compress:   config.AccessLog.Compress,
			}
			accessLogger.SetOutput(writer)
			loggers["access"] = accessLogger
		}

		DefaultLogger = &AppLogger{
			loggers: loggers,
		}
	})
	if DefaultLogger== nil {
		panic("Failed to initialize logger")
	}
	return DefaultLogger
}

func (a *AppLogger) entry(ctx context.Context, level string, kvs ...any) *logrus.Entry {
	if logger, ok := a.loggers[level]; ok {
		entry := logger.WithContext(ctx)
		// 添加上下文信息
		if requestID := ctx.Value("request_id"); requestID != nil {
			entry = entry.WithField("request_id", requestID)
		}
		if traceID := ctx.Value("trace_id"); traceID != nil {
			entry = entry.WithField("trace_id", traceID)
		}
		// 添加键值对到日志条目
		if len(kvs) > 1 {
			for i := 0; i < len(kvs)-1; i += 2 {
				if key, ok := kvs[i].(string); ok {
					entry = entry.WithField(key, kvs[i+1])
				}
			}
		}
		return entry
	}
	return nil
}

func (a *AppLogger) LogInfof(ctx context.Context, message string, args ...any) {
	if entry := a.entry(ctx, "info"); entry != nil {
		entry.Infof(message, args...)
	}
}

// 其他日志级别方法类似
func (a *AppLogger) LogWarnf(ctx context.Context, message string, args ...any) {
	if entry := a.entry(ctx, "warn"); entry != nil {
		entry.Warnf(message, args...)
	}
}

func (a *AppLogger) LogErrorf(ctx context.Context, message string, args ...any) {
	if entry := a.entry(ctx, "error"); entry != nil {
		entry.Errorf(message, args...)
	}
}

func (a *AppLogger) LogDebugf(ctx context.Context, message string, args ...any) {
	if entry := a.entry(ctx, "debug"); entry != nil {
		entry.Debugf(message, args...)
	}
}

func (a *AppLogger) LogFatalf(ctx context.Context, message string, args ...any) {
	if entry := a.entry(ctx, "fatal"); entry != nil {
		entry.Fatalf(message, args...)
	}
}

func (a *AppLogger) LogInfo(ctx context.Context, message string, kvs ...any) {
	if entry := a.entry(ctx, "info", kvs...); entry != nil {
		entry.Info(message)
	}
}

// 其他日志级别方法类似
func (a *AppLogger) LogWarn(ctx context.Context, message string, kvs ...any) {
	if entry := a.entry(ctx, "warn", kvs...); entry != nil {
		entry.Warn(message)
	}
}

func (a *AppLogger) LogError(ctx context.Context, message string, kvs ...any) {
	if entry := a.entry(ctx, "error", kvs...); entry != nil {
		entry.Error(message)
	}
}

func (a *AppLogger) LogDebug(ctx context.Context, message string, kvs ...any) {
	if entry := a.entry(ctx, "debug", kvs...); entry != nil {
		entry.Debug(message)
	}
}

func (a *AppLogger) LogFatal(ctx context.Context, message string, kvs ...any) {
	if entry := a.entry(ctx, "fatal", kvs...); entry != nil {
		entry.Fatal(message)
	}
}
