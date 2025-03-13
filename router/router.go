package router

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/networkProtocalTrans/logger"
	"github.com/networkProtocalTrans/module"
	"github.com/networkProtocalTrans/services"
)

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
	ws := r.Group("/ws")
	{
		// ws.GET("/test", RequestPanicHandler(controller.WSTestHandler))
		/*
		  /ws接收http协议后，前端js转为websocket协议，然后发送给后端/ws/echo
		*/
		ws.GET("/", services.WsServer.Test)
		ws.GET("/echo", services.WsServer.HandleConnections)
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
		logger.DefaultLogger.LogFatalf(ctx, "gin router handler panic +%v", r)
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
			logger.DefaultLogger.LogErrorf(ctx, "Failed to parse request parameters: %v", err)
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

type RequestHandler func(c *gin.Context) (res module.Response, err error)

func RequestPanicHandler(fn RequestHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		if res, err := fn(c); err != nil {
			logger.DefaultLogger.LogErrorf(ctx, "RequestHandler failed with error +%v", err)
		} else {
			// 返回结果
			c.JSON(res.GetStatus(), res.GetBody())
		}
	}
}
