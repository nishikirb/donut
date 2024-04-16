package config

import (
	"reflect"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"

	"github.com/nishikirb/donut/test/helper"
)

func Test_ExpandEnvFunc(t *testing.T) {
	dir := t.TempDir()
	helper.SetDirEnv(t, dir)
	f := ExpandEnvFunc()

	type args struct {
		from reflect.Value
		to   reflect.Value
	}
	tests := []struct {
		name      string
		args      args
		want      interface{}
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "OK/ContainsEnv",
			args: args{
				from: reflect.ValueOf("${HOME}"),
				to:   reflect.ValueOf(dir),
			},
			want:      dir,
			assertion: assert.NoError,
		},
		{
			name: "OK/NotContainsEnv",
			args: args{
				from: reflect.ValueOf(dir),
				to:   reflect.ValueOf(dir),
			},
			want:      dir,
			assertion: assert.NoError,
		},
		{
			name: "OK/Integer",
			args: args{
				from: reflect.ValueOf(1),
				to:   reflect.ValueOf(1),
			},
			want:      1,
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mapstructure.DecodeHookExec(f, tt.args.from, tt.args.to)
			assert.Equal(t, tt.want, got)
			tt.assertion(t, err)
		})
	}
}
