package util

import (
	"context"
)

// SafeGo 函数用于安全地启动一个 goroutine，并在其中处理可能的 panic
func SafeGo(ctx context.Context, fn func()) {
    go func() {
        // 使用 defer 语句确保在函数退出时执行 recover 操作
        defer func() {
            if r := recover(); r != nil {
                // 捕获到 panic 时，记录错误信息
                log.LogFatalf(ctx,"Recovered in SafeGo: %v", r)
            }
        }()
        // 执行传入的函数
        fn()
    }()
}
