package donut

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/viper"
)

const AppName string = "donut"

type App struct {
	config *Config
	out    io.Writer
	err    io.Writer
}

type AppOption func(*App) error

func New(opts ...AppOption) (*App, error) {
	app := &App{
		out: os.Stdout,
		err: os.Stderr,
	}
	for _, opt := range opts {
		if err := opt(app); err != nil {
			return nil, err
		}
	}

	return app, nil
}

func WithConfig(v *viper.Viper) AppOption {
	return func(app *App) error {
		var cfg Config
		if err := v.Unmarshal(&cfg, viper.DecodeHook(defaultDecodeHookFunc)); err != nil {
			return err
		}
		if err := validateConfig(&cfg); err != nil {
			return err
		}
		app.config = &cfg
		return nil
	}
}

func WithOut(r io.Writer) AppOption {
	return func(app *App) error {
		app.out = r
		return nil
	}
}

func WithErr(r io.Writer) AppOption {
	return func(app *App) error {
		app.err = r
		return nil
	}
}

func validateConfig(cfg *Config) error {
	for _, path := range []string{cfg.Source, cfg.Destination} {
		if err := isDir(path); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) Echo() error {
	fmt.Println("Source:", a.config.Source, "Destination:", a.config.Destination)
	return nil
}
