package donut

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/viper"
)

type App struct {
	Config *Config
	Store  *Store
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
	mapper, err := NewPathMapperBuilder(
		a.Config.Source, a.Config.Destination, WithExcludes(a.Config.Excludes...),
	).Build()
	if err != nil {
		return err
	}

	for _, relSourcePath := range mapper.RelSourcePaths() {
		fmt.Fprintln(a.out, relSourcePath)
	}
	return nil
}

func (a *App) Diff() error {
	mapper, err := NewPathMapperBuilder(
		a.Config.Source, a.Config.Destination, WithExcludes(a.Config.Excludes...),
	).Build()
	if err != nil {
		return err
	}

	var diff []byte
	diffConfig := a.Config.Diff
	for _, pm := range mapper.Mapping {
		ss, err := fileEntriesMap.GetSum(pm.Source)
		if err != nil {
			return err
		}
		ds, err := fileEntriesMap.GetSum(pm.Destination)
		if err != nil {
			return err
		}
		if bytes.Equal(ss, ds) {
			continue
		}

		args := append(diffConfig.Args, pm.Destination, pm.Source)
		cmd := exec.Command(diffConfig.Name, args...)
		out, _ := cmd.Output()
		diff = append(diff, out...)
	}

	pagerConfig := a.Config.Pager
	cmd := exec.Command(pagerConfig.Name, pagerConfig.Args...)
	cmd.Stdin = bytes.NewBuffer(diff)
	cmd.Stdout = a.out
	return cmd.Run()
}

func (a *App) ConfigShow() error {
	v := GetConfig()
	f, err := os.Open(v.ConfigFileUsed())
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := io.Copy(a.out, f); err != nil {
		return err
	}
	return nil
}

func (a *App) ConfigEdit() error {
	v := GetConfig()
	editorConfig := a.Config.Editor
	args := append(editorConfig.Args, v.ConfigFileUsed())
	cmd := exec.Command(a.Config.Editor.Name, args...)
	cmd.Stdout = a.out
	return cmd.Run()
}

func (a *App) Apply() error {
	mapper, err := NewPathMapperBuilder(
		a.Config.Source, a.Config.Destination, WithExcludes(a.Config.Excludes...),
	).Build()
	if err != nil {
		return err
	}

	for _, pm := range mapper.Mapping {
		ss, err := fileEntriesMap.GetSum(pm.Source)
		if err != nil {
			return err
		}
		ds, err := fileEntriesMap.GetSum(pm.Destination)
		if err != nil {
			return err
		}
		if bytes.Equal(ss, ds) {
			continue
		}

		// If the directory does not exist, create it
		// os.MkdirAll will return nil if directory already exists
		dir := filepath.Dir(pm.Destination)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
		if err := a.Overwrite(pm.Source, pm.Destination); err != nil {
			return err
		}

		fmt.Fprintf(a.out, "Applied a patch to the destination file. %s from %s\n", pm.Destination, pm.Source)
	}

	return nil
}

// Overwrite replaces the contents of dst with the contents of src.
func (a *App) Overwrite(src, dst string) error {
	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer f.Close()

	e, err := fileEntriesMap.Get(src)
	if err != nil {
		return err
	}
	c, err := e.GetContent()
	if err != nil {
		return err
	}
	if _, err := f.Write(c); err != nil {
		return err
	}
	return nil
}
