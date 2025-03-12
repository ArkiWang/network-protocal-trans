package module

import (
	"github.com/gin-gonic/gin"
)


type Request interface {
	// Validate 验证请求参数
	Validate() error
	GetMethod(c *gin.Context) string
	GetPayload(c *gin.Context) any
	// Decode(any) any
	// Encode(any) any
}

type BaseRequest struct {
	ID         string `json:"id" form:"id" gorm:"column:id;primary_key" binding:"required"`
	MethodName string `json:"method_name" form:"method_name" gorm:"column:method_name" binding:"required"`
}

// 实现 Req 接口
func (b *BaseRequest) Validate() error {
	return nil
}

func (b *BaseRequest) GetMethod(c *gin.Context) string {
	return c.Request.Method
}

func (b *BaseRequest) GetPayload(c *gin.Context) any{
	ctx := c.Request.Context()
	// 创建结构体实例并绑定查询参数
	if err := c.ShouldBindQuery(b); err != nil {
		log.LogErrorf(ctx, "GetPayload faild with error +%v", err)
		return nil
	}
	return b
}
