package donut

import (
	"io"

	"github.com/spf13/viper"
)

type AppOption func(*App) error

func WithConfig(v *viper.Viper) AppOption {
	return func(a *App) error {
		var c Config
		if err := v.Unmarshal(&c, viper.DecodeHook(defaultDecodeHookFunc)); err != nil {
			return err
		}
		if err := validateConfig(&c); err != nil {
			return err
		}
		a.Config = &c
		return nil
	}
}

func WithOut(r io.Writer) AppOption {
	return func(a *App) error {
		a.out = r
		return nil
	}
}

func WithErr(r io.Writer) AppOption {
	return func(a *App) error {
		a.err = r
		return nil
	}
}

func validateConfig(c *Config) error {
	for _, path := range []string{c.Source, c.Destination} {
		if err := IsDir(path); err != nil {
			return err
		}
	}
	return nil
}
