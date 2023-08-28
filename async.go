package main

import (
	"fmt"
	"log"
	"runtime/debug"
)

func Execute(fn func()) {
    go func() {
        defer recoverPanic()
        fn()
    }()
}

func recoverPanic() {
    if r := recover(); r != nil {
        err, ok := r.(error)
        if !ok {
            err = fmt.Errorf("%v", r)
        }
        log.Printf("panic catched: %v", err)
        debug.PrintStack()
    }
}

func WithoutError[T interface{}](val T, err error) T {
	if err != nil {
		panic(err)
	}
	return val
}
