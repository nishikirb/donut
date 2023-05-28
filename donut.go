package donut

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

const AppName string = "donut"

type App struct {
	Config *Config
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
	fmt.Println("Source:", a.Config.Source, "Destination:", a.Config.Destination)
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
		a.Config.Source,
		a.Config.Destination,
		WithExcludes(a.Config.Excludes...),
	).Build()
	if err != nil {
		return err
	}

	for _, v := range list {
		prefixTrimmed := strings.TrimPrefix(strings.TrimPrefix(v.Source.Path, a.Config.Source), string(filepath.Separator))
		fmt.Fprintln(a.out, prefixTrimmed)
	}
	return nil
}

func (a *App) Diff() error {
	list, err := NewRelationsBuilder(
		a.Config.Source,
		a.Config.Destination,
		WithExcludes(a.Config.Excludes...),
	).Build()
	if err != nil {
		return err
	}

	d := a.Config.Diff
	bs := []byte{}
	for _, v := range list {
		args := append(d.Args, v.Source.Path, v.Destination.Path)
		cmd := exec.Command(d.Command, args...)
		b, _ := cmd.Output()
		bs = append(bs, b...)
	}

	p := a.Config.Pager
	cmd := exec.Command(p.Command, p.Args...)
	cmd.Stdin = bytes.NewBuffer(bs)
	cmd.Stdout = a.out
	return cmd.Run()
}

func (a *App) Edit(file string, current string) error {
	list, err := NewRelationsBuilder(
		a.Config.Source,
		a.Config.Destination,
		WithExcludes(a.Config.Excludes...),
	).Build()
	if err != nil {
		return err
	}

	abs := AbsPath(file, current)
	var path string
	for _, v := range list {
		if v.Destination.Path == abs {
			path = v.Source.Path
			break
		}
	}
	if path == "" {
		return fmt.Errorf("file that related to " + file + "  not managed by donut")
	}

	e := a.Config.Editor
	args := append(e.Args, path)
	cmd := exec.Command(e.Command, args...)
	cmd.Stdout = a.out
	return cmd.Run()
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
		app.Config = &cfg
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
