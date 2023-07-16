package config

import (
	"os"
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"
)

var defaultDecodeHookFunc = mapstructure.ComposeDecodeHookFunc(
	ExpandEnvFunc(),
	mapstructure.StringToTimeDurationHookFunc(),
	mapstructure.StringToSliceHookFunc(","),
)

// mapstructure's DecodeHookFunc that reads environment variables
func ExpandEnvFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Kind, t reflect.Kind, data interface{}) (interface{}, error) {
		if f != reflect.String {
			return data, nil
		}

		raw := data.(string)
		if !strings.Contains(raw, "$") {
			return data, nil
		}

		return os.ExpandEnv(raw), nil
	}
}
