package logger

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/cektrendstudio/cektrend-engine-go/pkg/serror"
	"github.com/cektrendstudio/cektrend-engine-go/pkg/utils/utstring"
	"github.com/cektrendstudio/cektrend-engine-go/pkg/utils/uttime"
)

type (
	LogInterceptorTranslateArguments struct {
		Level   ErrorLevel
		Tags    map[string]string
		Payload interface{}
	}

	// LogInterceptor interface
	LogInterceptor interface {
		Translate(args LogInterceptorTranslateArguments) string
		Process(lvl ErrorLevel, msg string)
	}

	defaultInterceptorObj struct{}
)

func (defaultInterceptorObj) Translate(args LogInterceptorTranslateArguments) string {
	return DefaultTranslate(args, 2)
}

func (defaultInterceptorObj) Process(lvl ErrorLevel, msg string) {
	DefaultProcess(lvl, msg)
}

// DefaultInterceptor for default interceptor
func DefaultInterceptor() LogInterceptor {
	return defaultInterceptorObj{}
}

// DefaultTranslate for default translation
func DefaultTranslate(args LogInterceptorTranslateArguments, traceSkip int) string {
	plainMsg, coloredMsg := DefaultTransform(args, traceSkip+1)

	if !serror.IsLocal() {
		coloredMsg = plainMsg
	}

	tagMsg := "{"
	if layerName, ok := args.Tags[LogSquadTagLayerName]; ok {
		tagMsg += fmt.Sprintf("[%s]", layerName)
	}

	for tagName, tagValue := range args.Tags {
		if tagName == LogSquadTagLayerName {
			continue
		}

		if tagMsg != "{" {
			tagMsg += " "
		}

		tagMsg += fmt.Sprintf("%s:%s;", tagName, tagValue)
	}

	coloredMsg = tagMsg + "} " + coloredMsg

	// formating
	var (
		now                     = time.Now()
		logLevelLabel           = "?"
		errorLevelLogConfigMapx = map[ErrorLevel]struct {
			Label string
			Color utstring.Color
		}{
			ErrorLevelInfo:     {"INFO", utstring.LIGHT_BLUE},
			ErrorLevelLog:      {"LOG", utstring.LIGHT_GRAY},
			ErrorLevelWarning:  {"WARN", utstring.LIGHT_YELLOW},
			ErrorLevelCritical: {"ERR", utstring.RED},
		}
	)
	if cur, ok := errorLevelLogConfigMapx[args.Level]; ok {
		logLevelLabel = cur.Label

		if serror.IsLocal() {
			logLevelLabel = utstring.ApplyForeColor(logLevelLabel, cur.Color)
		}
	}

	return fmt.Sprintf("[%s] %s: %s", uttime.Format(uttime.DefaultDateTimeFormat, now), logLevelLabel, coloredMsg)
}

// DefaultTransform for default transforming
func DefaultTransform(args LogInterceptorTranslateArguments, traceSkip int) (plainMsg string, colorMsg string) {
	plainMsg = fmt.Sprintf("%v", args.Payload)
	colorMsg = plainMsg

	switch args.Level {
	case ErrorLevelCritical, ErrorLevelWarning:
		switch vx := args.Payload.(type) {
		case serror.SError:
			plainMsg = vx.String()
			colorMsg = vx.ColoredString()

		case error:
			pc, fn, line, _ := runtime.Caller(1 + traceSkip)
			plainMsg = fmt.Sprintf(serror.StandardFormat(), runtime.FuncForPC(pc).Name(), fn, line, plainMsg)
			colorMsg = fmt.Sprintf(serror.StandardColorFormat(), runtime.FuncForPC(pc).Name(), fn, line, colorMsg)
		}
	}

	if Environment() != "local" {
		plainMsg = strings.ReplaceAll(plainMsg, "\n", "↩")
		colorMsg = strings.ReplaceAll(colorMsg, "\n", "↩")
	}

	return
}

// DefaultProcess for default processing
func DefaultProcess(lvl ErrorLevel, msg string) {
	if msg == "" {
		return
	}

	switch lvl {
	case ErrorLevelCritical, ErrorLevelWarning:
		DefaultStderr(msg)

	default:
		DefaultStdout(msg)
	}
}

// DefaultStdout for default stdout print
func DefaultStdout(msg string) {
	fmt.Fprintln(os.Stdout, msg)
}

// DefaultStdout for default stderr print
func DefaultStderr(msg string) {
	fmt.Fprintln(os.Stderr, msg)
}

func Environment() string {
	return utstring.Env("APP_ENV", "production")
}
