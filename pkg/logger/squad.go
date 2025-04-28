package logger

import (
	"fmt"
	"regexp"

	"github.com/cektrendstudio/cektrend-engine-go/pkg/serror"
	"github.com/cektrendstudio/cektrend-engine-go/pkg/utils/utarray"
	"github.com/cektrendstudio/cektrend-engine-go/pkg/utils/utinterface"
)

type logSquadObj struct {
	logger    *logger
	layerName string
	tags      map[string]string
}

type LogSquad interface {
	logCommon

	// SetTag to set logger tag name and value
	SetTag(name string, value interface{}) (ok bool)
}

const (
	LogSquadTagLayerName      = "layerName"
	LogSquadTagDatadogTraceID = "ddTraceID"
	LogSquadTagDatadogSpanID  = "ddSpanID"
)

var (
	isValidTagName  = regexp.MustCompile(`^[a-zA-Z0-9\._-]{1,50}$`).MatchString
	isValidTagValue = regexp.MustCompile(`^[a-zA-Z0-9 \~\!\@\#\$\%\^\&\*\(\)\-\_\=\+\.\,\?\:\|\/]{1,255}$`).MatchString
)

func (ox *logSquadObj) SetTag(name string, value interface{}) (ok bool) {
	if ox.tags == nil {
		ox.tags = make(map[string]string)
	}

	valueStr := utinterface.ToString(value)
	if ox.isWellKnownTag(name) || !isValidTagName(name) || !isValidTagValue(valueStr) {
		return
	}

	ox.logger.mu.Lock()
	defer ox.logger.mu.Unlock()

	ox.tags[name] = valueStr

	ok = true
	return
}

func (ox *logSquadObj) Info(msg interface{}) {
	ox.process(ErrorLevelInfo, msg)
}

func (ox *logSquadObj) Infof(msg string, args ...interface{}) {
	ox.Info(fmt.Sprintf(msg, args...))
}

func (ox *logSquadObj) Log(msg interface{}) {
	ox.process(ErrorLevelLog, msg)
}

func (ox *logSquadObj) Logf(msg string, args ...interface{}) {
	ox.Log(fmt.Sprintf(msg, args...))
}

func (ox *logSquadObj) Warn(msg interface{}) {
	ox.process(ErrorLevelWarning, msg)
}

func (ox *logSquadObj) Warnf(msg string, args ...interface{}) {
	ox.Warn(fmt.Sprintf(msg, args...))
}

func (ox *logSquadObj) Err(msg interface{}) {
	ox.process(ErrorLevelCritical, msg)
}

func (ox *logSquadObj) Errf(msg string, args ...interface{}) {
	ox.Err(serror.Newsf(1, msg, args...))
}

func (ox *logSquadObj) Panic(msg interface{}) {
	ox.Err(castToSError(msg, 1))
	exit()
}

func (ox logSquadObj) isWellKnownTag(name string) bool {
	return utarray.IsExist(name, []string{
		LogSquadTagLayerName,
		LogSquadTagDatadogTraceID,
		LogSquadTagDatadogSpanID,
	})
}

func (ox logSquadObj) process(lvl ErrorLevel, msg interface{}) {
	finalMsg := ox.logger.activeInterceptor.Translate(LogInterceptorTranslateArguments{
		Level:   lvl,
		Tags:    ox.tags,
		Payload: msg,
	})
	ox.logger.activeInterceptor.Process(lvl, finalMsg)
	_ = ox.logger.write(finalMsg)
}
