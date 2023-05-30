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

type App struct {
	Config *Config
	out    io.Writer
	err    io.Writer
}

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

func (a *App) Init() error {
	// check config file exists in default path
	// if exists, no need to init
	_, err := NewConfig(WithDefault(), WithNameAndPath(AppName, DefaultConfigDirs()...))
	if err == nil {
		return errors.New("config file already exists. no need to init")
	} else if !errors.As(err, &viper.ConfigFileNotFoundError{}) {
		return fmt.Errorf("config file already exists, but error: %w", err)
	}

	path := DefaultConfigFile()
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return err
	}
	v, _ := NewConfig(WithDefault())
	if err := v.SafeWriteConfigAs(path); err != nil {
		return err
	}

	fmt.Fprintln(a.out, "Configuration file created in", path)
	return nil
}

func (a *App) List() error {
	mps, err := NewMappingsBuilder(
		a.Config.Source, a.Config.Destination, WithExcludes(a.Config.Excludes...),
	).Build()
	if err != nil {
		return err
	}

	for _, v := range mps {
		prefixTrimmed := strings.TrimPrefix(strings.TrimPrefix(v.Source.Path, a.Config.Source), string(filepath.Separator))
		fmt.Fprintln(a.out, prefixTrimmed)
	}
	return nil
}

func (a *App) Diff() error {
	mps, err := NewMappingsBuilder(
		a.Config.Source, a.Config.Destination, WithExcludes(a.Config.Excludes...),
	).Build()
	if err != nil {
		return err
	}

	var diff []byte
	for _, v := range mps {
		args := append(a.Config.Diff.Args, v.Source.Path, v.Destination.Path)
		cmd := exec.Command(a.Config.Diff.Name, args...)
		b, _ := cmd.Output()
		diff = append(diff, b...)
	}

	cmd := exec.Command(a.Config.Pager.Name, a.Config.Pager.Args...)
	cmd.Stdin = bytes.NewBuffer(diff)
	cmd.Stdout = a.out
	return cmd.Run()
}

func (a *App) Edit(file string, current string) error {
	mps, err := NewMappingsBuilder(
		a.Config.Source, a.Config.Destination, WithExcludes(a.Config.Excludes...),
	).Build()
	if err != nil {
		return err
	}

	abs := AbsPath(file, current)
	var path string
	for _, v := range mps {
		if v.Destination.Path == abs {
			path = v.Source.Path
			break
		}
	}
	if path == "" {
		return fmt.Errorf("file that related to " + file + "  not managed by donut")
	}

	args := append(a.Config.Editor.Args, path)
	cmd := exec.Command(a.Config.Editor.Name, args...)
	cmd.Stdout = a.out
	return cmd.Run()
}

func (a *App) EditConfig() error {
	v := GetConfig()
	args := append(a.Config.Editor.Args, v.ConfigFileUsed())
	cmd := exec.Command(a.Config.Editor.Name, args...)
	cmd.Stdout = a.out
	return cmd.Run()
}

func (a *App) Apply() error {
	mps, err := NewMappingsBuilder(
		a.Config.Source, a.Config.Destination, WithExcludes(a.Config.Excludes...),
	).Build()
	if err != nil {
		return err
	}

	for _, v := range mps {
		// If the directory does not exist, create it
		// os.MkdirAll will return nil if directory already exists
		dir := filepath.Dir(v.Destination.Path)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}

		if !v.Destination.NotExist {
			// os.Remove will unlink the file if the file is a symbolic link
			if err := os.Remove(v.Destination.Path); err != nil {
				return err
			}
		}

		if err := cp(v.Source.Path, v.Destination.Path); err != nil {
			return err
		}

		fmt.Fprintf(a.out, "File copied. %s from %s\n", v.Destination.Path, v.Source.Path)
	}

	return nil
}

func cp(src, dst string) error {
	sf, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sf.Close()
	df, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer df.Close()

	_, err = io.Copy(df, sf)
	return err
}
