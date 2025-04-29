package storage

import (
	"context"

	"github.com/cektrendstudio/cektrend-engine-go/pkg/serror"
)

type ConfigMap map[interface{}]interface{}

type Option func(config *ConfigMap) (errx serror.SError)

func WithBucketName(bucketName string) Option {
	return func(config *ConfigMap) (errx serror.SError) {
		(*config)[symbolBucketName] = bucketName
		return
	}
}

var (
	symbolBucketName = new(struct{})
	symbolTeam       = new(struct{})
)

func buildConfigFromOption(ctx context.Context, output *ConfigMap, options []Option) (errx serror.SError) {
	if output == nil {
		*output = make(ConfigMap)
	}

	if len(options) == 0 {
		return
	}

	for _, option := range options {
		if errx = option(output); errx != nil {
			errx.AddComments("while resolve option")
			return
		}
	}

	return
}

func (ox *ConfigMap) GetValue(symbol interface{}, defaultValues ...interface{}) (value interface{}) {
	var isExists bool
	if value, isExists = (*ox)[symbol]; !isExists {
		for _, defaultValue := range defaultValues {
			if defaultValue == nil {
				continue
			}

			value = defaultValue
		}
	}

	return
}
