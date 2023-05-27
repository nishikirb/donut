package donut

import (
	"os"
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"
)

var defaultDecodeHookFunc = mapstructure.ComposeDecodeHookFunc(
	expandEnvFunc(),
	mapstructure.StringToTimeDurationHookFunc(),
	mapstructure.StringToSliceHookFunc(","),
)

// mapstructure's DecodeHookFunc that reads environment variables
func expandEnvFunc() mapstructure.DecodeHookFunc {
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
