package logger

import (
	"context"
	"fmt"

	"github.com/cektrendstudio/cektrend-engine-go/pkg/serror"
	"github.com/cektrendstudio/cektrend-engine-go/pkg/utils/utinterface"
	"github.com/cektrendstudio/cektrend-engine-go/pkg/utils/utstring"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// ErrorLevel type
type ErrorLevel string

const (
	// ErrorLevelDebug level
	ErrorLevelDebug ErrorLevel = "debug"

	// ErrorLevelLog level
	ErrorLevelLog ErrorLevel = "log"

	// ErrorLevelInfo level
	ErrorLevelInfo ErrorLevel = "info"

	// ErrorLevelCritical level
	ErrorLevelCritical ErrorLevel = "critical"

	// ErrorLevelWarning level
	ErrorLevelWarning ErrorLevel = "warn"
)

// Mode type
type Mode int

const (
	// ModeDaily mode
	ModeDaily Mode = 1 + iota

	// ModeMonthly mode
	ModeMonthly

	// ModeYearly mode
	ModeYearly

	// ModePermanent mode
	ModePermanent
)

// Options struct
type Options struct {
	Mode        Mode           `valid:"-"`
	Path        string         `valid:"-"`
	Writing     bool           `valid:"-"`
	FileFormat  string         `valid:"-"`
	Interceptor LogInterceptor `valid:"-"`
}

type logCommon interface {
	// Info to logging with info level
	Info(msg interface{})
	// Infof to logging with string binding and info level
	Infof(msg string, args ...interface{})

	// Log to logging
	Log(msg interface{})
	// Logf to logging with string binding
	Logf(msg string, args ...interface{})

	// Warn to logging with warning level
	Warn(msg interface{})
	// Warnf to logging with string binding and warning level
	Warnf(msg string, args ...interface{})

	// Err to logging with error level
	Err(msg interface{})
	// Errf to logging with string binding and error level
	Errf(msg string, args ...interface{})

	// Panic to logging panic error
	Panic(msg interface{})
}

type (
	logger struct {
		logWriterObj
		isReady           bool
		activeInterceptor LogInterceptor
	}

	// Logger service
	Logger interface {
		logCommon
		logWriter

		// Startup to starting up the logger
		Startup() error
		// IsReady returns the logger instance ready status
		IsReady() bool

		// SetInterceptor to set the log interceptor
		SetInterceptor(intercept LogInterceptor)

		// CreateSquad returns the new squad logger instance
		CreateSquad(ctx context.Context, layerName string) LogSquad
	}
)

// public

// Construct function
func Construct(opt Options) Logger {
	if opt.Interceptor == nil {
		opt.Interceptor = DefaultInterceptor()
	}

	log := &logger{
		logWriterObj: logWriterObj{
			Mode:       opt.Mode,
			Path:       opt.Path,
			isWriting:  opt.Writing,
			FileFormat: utstring.Chains(opt.FileFormat, "log-%v.log"),
		},

		activeInterceptor: opt.Interceptor,
	}
	return log
}

func (ox *logger) Startup() (err error) {
	if ox.isReady {
		return
	}

	ox.isReady = true

	err = ox.logWriterObj.boot()
	if err != nil {
		return
	}

	return
}

func (ox *logger) IsReady() bool {
	return ox.isReady
}

func (ox *logger) SetInterceptor(intercept LogInterceptor) {
	ox.activeInterceptor = intercept
}

func (ox *logger) CreateSquad(ctx context.Context, layerName string) LogSquad {
	squad := &logSquadObj{
		logger:    ox,
		layerName: layerName,
		tags:      make(map[string]string),
	}

	if ctx != nil {
		if span, ok := tracer.SpanFromContext(ctx); ok {
			spanCtx := span.Context()
			squad.tags[LogSquadTagLayerName] = layerName
			squad.tags[LogSquadTagDatadogTraceID] = utinterface.ToString(spanCtx.TraceID())
			squad.tags[LogSquadTagDatadogSpanID] = utinterface.ToString(spanCtx.SpanID())
		}
	}

	return squad
}

func (ox *logger) Info(msg interface{}) {
	lvl := ErrorLevelInfo

	m := ox.activeInterceptor.Translate(LogInterceptorTranslateArguments{
		Level:   lvl,
		Payload: msg,
	})
	ox.activeInterceptor.Process(lvl, m)
	_ = ox.write(m)
}

func (ox *logger) Infof(msg string, args ...interface{}) {
	ox.Info(fmt.Sprintf(msg, args...))
}

func (ox *logger) Log(msg interface{}) {
	lvl := ErrorLevelLog

	m := ox.activeInterceptor.Translate(LogInterceptorTranslateArguments{
		Level:   lvl,
		Payload: msg,
	})
	ox.activeInterceptor.Process(lvl, m)
	_ = ox.write(m)
}

func (ox *logger) Logf(msg string, args ...interface{}) {
	ox.Log(fmt.Sprintf(msg, args...))
}

func (ox *logger) Warn(msg interface{}) {
	lvl := ErrorLevelWarning

	m := ox.activeInterceptor.Translate(LogInterceptorTranslateArguments{
		Level:   lvl,
		Payload: msg,
	})
	ox.activeInterceptor.Process(lvl, m)
	_ = ox.write(m)
}

func (ox *logger) Warnf(msg string, args ...interface{}) {
	ox.Warn(fmt.Sprintf(msg, args...))
}

func (ox *logger) Err(msg interface{}) {
	lvl := ErrorLevelCritical

	m := ox.activeInterceptor.Translate(LogInterceptorTranslateArguments{
		Level:   lvl,
		Payload: msg,
	})
	ox.activeInterceptor.Process(lvl, m)
	_ = ox.write(m)
}

func (ox *logger) Errf(msg string, args ...interface{}) {
	ox.Err(serror.Newsf(1, msg, args...))
}

func (ox *logger) Panic(msg interface{}) {
	ox.Err(castToSError(msg, 1))
	exit()
}
