package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/networkProtocalTrans/logger"
	"github.com/networkProtocalTrans/module"
)

func WSTestHandler(c *gin.Context)(res module.Response, err error){
		ctx := c.Request.Context()
		req := &module.BaseRequest{}
		if err := c.ShouldBind(req); err != nil {
			logger.DefaultLogger.LogErrorf(ctx, "gin bind request failed with error +%v", err)
		}
		return &module.BaseResponse{Status: 0},nil
}