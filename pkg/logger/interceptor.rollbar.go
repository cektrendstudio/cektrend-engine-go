package logger

import (
	"github.com/cektrendstudio/cektrend-engine-go/pkg/serror"
	"github.com/cektrendstudio/cektrend-engine-go/pkg/utils/utarray"
	"github.com/cektrendstudio/cektrend-engine-go/pkg/utils/utstring"

	"github.com/go-playground/validator/v10"
	"github.com/lunixbochs/vtclean"
	rollbar "github.com/rollbar/rollbar-go"
)

type (
	// RollbarOptions type
	RollbarOptions struct {
		Key     string     `json:"key" valid:"required"`
		Name    string     `json:"name" valid:"required"`
		Token   string     `json:"token" valid:"required"`
		Version string     `json:"version" valid:"required"`
		Level   ErrorLevel `json:"level" valid:"required"`
	}

	// Rollbar interceptor
	Rollbar interface {
		LogInterceptor
		IsEnabled() bool
		Enable()
		Disable()
	}

	rollbarInterceptorObj struct {
		Level   ErrorLevel
		Enabled bool
	}
)

// RollbarInterceptor to create rollbar interceptor
func RollbarInterceptor(opt RollbarOptions) (obj Rollbar, errx serror.SError) {
	validate := validator.New()
	if err := validate.Struct(opt); err != nil {
		errx = serror.NewFromErrorc(err, "Invalid rollbar options")
		return obj, errx
	}

	rollbar.SetToken(opt.Token)
	rollbar.SetEnvironment(utstring.Env("APP_ENV", "local"))
	rollbar.SetCodeVersion(opt.Version)
	rollbar.SetServerHost(opt.Key)

	obj = &rollbarInterceptorObj{
		Level:   opt.Level,
		Enabled: true,
	}
	return obj, errx
}

func (rollbarInterceptorObj) Translate(args LogInterceptorTranslateArguments) string {
	return DefaultTranslate(args, 2)
}

func (ox rollbarInterceptorObj) Process(lvl ErrorLevel, msg string) {
	DefaultProcess(lvl, msg)

	rlvl := map[ErrorLevel]string{
		ErrorLevelLog:      rollbar.DEBUG,
		ErrorLevelInfo:     rollbar.INFO,
		ErrorLevelWarning:  rollbar.WARN,
		ErrorLevelCritical: rollbar.CRIT,
	}

	if ox.Enabled {
		allow := []ErrorLevel{
			ErrorLevelDebug,
		}

		switch ox.Level {
		case ErrorLevelLog:
			allow = append(allow, ErrorLevelLog)
			fallthrough

		case ErrorLevelInfo:
			allow = append(allow, ErrorLevelInfo)
			fallthrough

		case ErrorLevelWarning:
			allow = append(allow, ErrorLevelWarning)
			fallthrough

		case ErrorLevelCritical:
			allow = append(allow, ErrorLevelCritical)

		default:
			return
		}

		if utarray.IsExist(lvl, allow) {
			rollbar.Log(rlvl[lvl], vtclean.Clean(msg, false))
		}
	}
}

func (ox rollbarInterceptorObj) IsEnabled() bool {
	return ox.Enabled
}

func (ox *rollbarInterceptorObj) Enable() {
	ox.Enabled = true
}

func (ox *rollbarInterceptorObj) Disable() {
	ox.Enabled = false
}
