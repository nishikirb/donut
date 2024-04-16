package donut

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nishikirb/donut/config"
	"github.com/nishikirb/donut/test/helper"
)

func TestNewApp(t *testing.T) {
	home, _, _, _ := helper.CreateBaseDir(t)
	helper.SetDirEnv(t, home)
	defer config.SetUserHomeDir(home)()

	stdin := bytes.NewBuffer([]byte(""))
	stdout, stderr := bytes.NewBuffer([]byte("")), bytes.NewBuffer([]byte(""))

	tests := []struct {
		name           string
		opts           []Option
		want           *App
		applyAssertion assert.ErrorAssertionFunc
		assertion      func(t *testing.T, want, got *App)
	}{
		{
			name:           "OK/WithConfig",
			opts:           []Option{WithConfig(&config.Config{Editor: []string{"nvim"}})},
			want:           &App{config: &config.Config{Editor: []string{"nvim"}}},
			applyAssertion: assert.NoError,
			assertion: func(t *testing.T, want, got *App) {
				assert.Equal(t, want.config.Editor, got.config.Editor)
				assert.Equal(t, want.config.Pager, got.config.Pager)
			},
		},
		{
			name:           "OK/WithConfigLoader",
			opts:           []Option{WithConfigLoader(config.WithFile("./test/testdata/config/basic.toml"))},
			want:           &App{config: &config.Config{Editor: []string{"nvim"}}},
			applyAssertion: assert.NoError,
			assertion: func(t *testing.T, want, got *App) {
				assert.Equal(t, want.config.Editor, got.config.Editor)
			},
		},
		{
			name:           "Error/WithConfigLoader",
			opts:           []Option{WithConfigLoader(config.WithFile("./test/testdata/config/broken.toml"))},
			want:           &App{config: nil},
			applyAssertion: assert.Error,
			assertion: func(t *testing.T, want, got *App) {
				assert.Equal(t, want.config.Editor, got.config.Editor)
			},
		},
		{
			name:           "OK/WithIn",
			opts:           []Option{WithIn(stdin)},
			want:           &App{in: stdin},
			applyAssertion: assert.NoError,
			assertion: func(t *testing.T, want, got *App) {
				assert.Equal(t, want.in, got.in)
			},
		},
		{
			name:           "OK/WithOut",
			opts:           []Option{WithOut(stdout)},
			want:           &App{out: stdout},
			applyAssertion: assert.NoError,
			assertion: func(t *testing.T, want, got *App) {
				assert.Equal(t, want.out, got.out)
			},
		},
		{
			name:           "OK/WithErr",
			opts:           []Option{WithErr(stderr)},
			want:           &App{err: stderr},
			applyAssertion: assert.NoError,
			assertion: func(t *testing.T, want, got *App) {
				assert.Equal(t, want.err, got.err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewApp(tt.opts...)
			err := a.ApplyOptions()
			tt.applyAssertion(t, err)
			if err == nil {
				tt.assertion(t, tt.want, a)
			}
		})
	}
}
