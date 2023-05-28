package donut

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

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

func (a *App) Echo() error {
	fmt.Println("Source:", a.config.Source, "Destination:", a.config.Destination)
	return nil
}

func (a *App) Init() error {
	// check config file exists in default path
	// if exists, no need to init
	_, err := NewConfig(WithDefault(), WithNameAndPath(AppName, DefaultConfigDirs()...))
	if err == nil {
		return errors.New("config file already exists. no need to init")
	} else if !errors.As(err, &viper.ConfigFileNotFoundError{}) {
		return fmt.Errorf("config file already exists, but error: %w", err)
	}

	configFile := DefaultConfigFile()
	if err := os.MkdirAll(filepath.Dir(configFile), os.ModePerm); err != nil {
		return err
	}
	v, _ := NewConfig(WithDefault())
	if err := v.SafeWriteConfigAs(configFile); err != nil {
		return err
	}

	fmt.Fprintln(a.out, "Configuration file created in", configFile)
	return nil
}

func (a *App) List() error {
	list, err := NewRelationsBuilder(
		a.config.Source,
		a.config.Destination,
		WithExcludes(a.config.Excludes...),
	).Build()
	if err != nil {
		return err
	}

	for _, v := range list {
		prefixTrimmed := strings.TrimPrefix(strings.TrimPrefix(v.Source.Path, a.config.Source), string(filepath.Separator))
		fmt.Fprintln(a.out, prefixTrimmed)
	}
	return nil
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
