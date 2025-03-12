package router

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/networkProtocalTrans/logger"
	"github.com/networkProtocalTrans/module"
)

var log = logger.DefaultLogger

// 初始化路由
func InitRouter(ctx context.Context) *gin.Engine {
	// 创建默认的gin路由引擎
	r := gin.Default()

	// 设置跨域中间件
	r.Use(CORSMiddleware())

	// API版本v1分组
	v1 := r.Group("/api/v1")
	{
		// 协议转换相关路由
		v1.POST("/convert", HandleProtocolConversion)

		// 健康检查
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status": "ok",
			})
		})
	}
	ws := r.Group("/api/ws")
	{
		ws.GET("/test", RequestPanicHandler(func(ctx context.Context, req module.Request) (res module.Response, err error) { return nil, nil }))
	}
	return r
}

// 处理协议转换的函数
func HandleProtocolConversion(c *gin.Context) {
	// TODO: 实现协议转换逻辑
	c.JSON(http.StatusOK, gin.H{
		"message": "Protocol conversion endpoint",
	})
}

func HandlePanic(c *gin.Context) {
	ctx := c.Request.Context()
	if r := recover(); r != nil {
		log.LogFatalf(ctx, "gin router handler panic +%v", r)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		// 解析请求参数
		var requestData struct {
			SourceProtocol      string `json:"source_protocol"`
			DestinationProtocol string `json:"destination_protocol"`
			Data                string `json:"data"`
		}

		if err := c.ShouldBindJSON(&requestData); err != nil {
			log.LogErrorf(ctx, "Failed to parse request parameters: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request parameters",
			})
			return
		}

		// 返回转换结果
		c.JSON(http.StatusOK, gin.H{
			"source_protocol":      requestData.SourceProtocol,
			"destination_protocol": requestData.DestinationProtocol,
			"original_data":        requestData.Data,
			"converted_data":       "转换后的数据", // TODO: 替换为实际转换后的数据
			"status":               "success",
		})
	}
}

type RequestHandler func(ctx context.Context, req module.Request) (res module.Response, err error)

func RequestPanicHandler(fn RequestHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		req := &module.BaseRequest{}
		if err := c.ShouldBind(req); err != nil {
			log.LogErrorf(ctx, "gin bind request failed with error +%v", err)
		}
		if res, err := fn(ctx, req); err != nil {
			log.LogErrorf(ctx, "RequestHandler failed with error +%v", err)
		} else {
			// 返回结果
			c.JSON(res.GetStatus(), res.GetBody())
		}
	}
}
