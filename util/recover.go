package util

import (
	"context"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/networkProtocalTrans/logger"
)

// HandlePanic recovers from panics with detailed stack trace
func HandlePanic(fn func() error) {
    if r := recover(); r != nil {
        // Get stack trace
        stack := debug.Stack()
        
        // Log the panic and stack trace
        logger.DefaultLogger.LogFatalf(context.Background(),"PANIC: %v\n%s", r, string(stack))
        // Optionally re-panic if you want the program to terminate
        // panic(r)
		fn()
    }
}

func GinPanicHandler(c *gin.Context){
    if r:= recover();r !=nil{

    }
}
