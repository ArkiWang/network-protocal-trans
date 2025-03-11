package start

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LogConfig 日志配置结构体
type LogConfig struct {
    LogLevel  string     `toml:"log_level"`
    Log       LogFile    `toml:"log"`
    AccessLog AccessFile `toml:"access_log"`
}

// LogFile 普通日志配置
type LogFile struct {
    File       string `toml:"file"`
    Console    bool   `toml:"console"`
    MaxSize    int    `toml:"max_size"`
    MaxBackups int    `toml:"max_backups"`
    MaxAge     int    `toml:"max_age"`
    Compress   bool   `toml:"compress"`
}

// AccessFile 访问日志配置
type AccessFile struct {
    File       string `toml:"file"`
    Enable     bool   `toml:"enable"`
    MaxSize    int    `toml:"max_size"`
    MaxBackups int    `toml:"max_backups"`
    MaxAge     int    `toml:"max_age"`
    Compress   bool   `toml:"compress"`
}

var (
    // Logger 全局日志实例
    Logger *log.Logger
    // AccessLogger 访问日志实例
    AccessLogger *log.Logger
)

// InitLogger 初始化日志配置
func InitLogger(ctx context.Context, configPath string) error {
    var config LogConfig
    if _, err := toml.DecodeFile(configPath, &config); err != nil {
        return err
    }

    // 初始化应用日志
    appLogWriter := &lumberjack.Logger{
        Filename:   config.Log.File,
        MaxSize:    config.Log.MaxSize,
        MaxBackups: config.Log.MaxBackups,
        MaxAge:     config.Log.MaxAge,
        Compress:   config.Log.Compress,
    }

    // 如果需要同时输出到控制台
    var writers []io.Writer
    writers = append(writers, appLogWriter)
    if config.Log.Console {
        writers = append(writers, os.Stdout)
    }

    // 创建多输出writer
    multiWriter := io.MultiWriter(writers...)
    Logger = log.New(multiWriter, "", log.Ldate|log.Ltime|log.Lshortfile)

    // 初始化访问日志（如果启用）
    if config.AccessLog.Enable {
        accessLogWriter := &lumberjack.Logger{
            Filename:   config.AccessLog.File,
            MaxSize:    config.AccessLog.MaxSize,
            MaxBackups: config.AccessLog.MaxBackups,
            MaxAge:     config.AccessLog.MaxAge,
            Compress:   config.AccessLog.Compress,
        }
        AccessLogger = log.New(accessLogWriter, "", log.Ldate|log.Ltime)
    }

    return nil
}