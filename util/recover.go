package util

import (
	"log"
	"runtime/debug"
)

// HandlePanic recovers from panics with detailed stack trace
func HandlePanic(fn func() error, log *log.Logger ) {
    if r := recover(); r != nil {
        // Get stack trace
        stack := debug.Stack()
        
        // Log the panic and stack trace
        log.Printf("PANIC: %v\n%s", r, string(stack))
        // Optionally re-panic if you want the program to terminate
        // panic(r)
		fn()
    }
}