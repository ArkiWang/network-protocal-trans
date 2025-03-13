package module

// Response 接口定义了网络协议转换模块的响应行为
type Response interface{
    // GetStatus 获取响应状态
    GetStatus() int
    
    // GetBody 获取响应主体数据
    GetBody() []byte
    
    // GetHeaders 获取响应头信息
    GetHeaders() map[string][]string
}

// BaseResponse 实现了 Response 接口的基础响应结构体
type BaseResponse struct {
    // 响应状态码
    Status int
    // 响应主体数据
    Body []byte
    // 响应头信息
    Headers map[string][]string
}

// NewBaseResponse 创建一个新的 BaseResponse 实例
func NewBaseResponse(status int, body []byte, headers map[string][]string) *BaseResponse {
    return &BaseResponse{
        Status: status,
        Body: body,
        Headers: headers,
    }
}

// GetStatus 实现 Response 接口的获取响应状态方法
func (r *BaseResponse) GetStatus() int {
    return r.Status
}

// GetBody 实现 Response 接口的获取响应主体数据方法
func (r *BaseResponse) GetBody() []byte {
    return r.Body
}

// GetHeaders 实现 Response 接口的获取响应头信息方法
func (r *BaseResponse) GetHeaders() map[string][]string {
    return r.Headers
}
