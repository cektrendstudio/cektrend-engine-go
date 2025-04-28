package logger

import (
	"context"
	"os"
	"syscall"

	"github.com/cektrendstudio/cektrend-engine-go/pkg/serror"
)

var defaultInstance Logger = Construct(Options{})

func init() {
	err := defaultInstance.Startup()
	if err != nil {
		Err(serror.NewFromErrorc(err, "Failed to startup the default logger instance"))
	}
}

// SetInterceptor to set log interceptor
func SetInterceptor(intercept LogInterceptor) {
	if intercept == nil {
		panic("Null interceptor")
	}

	defaultInstance.SetInterceptor(intercept)
}

// CreateSquad returns the logger squad instance
func CreateSquad(ctx context.Context, layerName string) LogSquad {
	return defaultInstance.CreateSquad(ctx, layerName)
}

// Info to logging info level
func Info(msg interface{}) {
	defaultInstance.Info(msg)
}

// Infof to logging info level with function
func Infof(msg string, args ...interface{}) {
	defaultInstance.Infof(msg, args...)
}

// Log to logging log level
func Log(msg interface{}) {
	defaultInstance.Log(msg)
}

// Logf to logging log level with function
func Logf(msg string, args ...interface{}) {
	defaultInstance.Logf(msg, args...)
}

// Warn to logging warning level
func Warn(msg interface{}) {
	defaultInstance.Warn(msg)
}

// Warnf to logging warning level with function
func Warnf(msg string, args ...interface{}) {
	defaultInstance.Warnf(msg, args...)
}

// Err to logging error level
func Err(msg interface{}) {
	defaultInstance.Err(msg)
}

// Errf to logging error level with function
func Errf(msg string, args ...interface{}) {
	defaultInstance.Errf(msg, args...)
}

// Panic to logging error then exit
func Panic(msg interface{}) {
	defaultInstance.Panic(msg)
}

func exit() {
	err := syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	if err != nil {
		os.Exit(1)
	}
}

func castToSError(obj interface{}, skip int) serror.SError {
	var errx serror.SError

	if cur, ok := obj.(serror.SError); ok {
		errx = cur

	} else if cur, ok := obj.(error); ok {
		errx = serror.NewFromErrors(skip+1, cur)

	} else {
		errx = serror.Newsf(skip+1, "%+v", obj)
	}

	return errx
}
