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
	"text/template"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type App struct {
	Config   *Config
	Store    *Store
	Logger   zerolog.Logger
	template *template.Template
	in       io.Reader
	out      io.Writer
	err      io.Writer
}

func New(opts ...AppOption) (*App, error) {
	app := &App{
		in:  os.Stdin,
		out: os.Stdout,
		err: os.Stderr,
	}
	for _, opt := range opts {
		if err := opt(app); err != nil {
			return nil, err
		}
	}

	if app.Config != nil {
		if err := app.createTemplates(); err != nil {
			return nil, err
		}
	}

	return app, nil
}

func (a *App) createTemplates() error {
	var tmpl = template.New("")
	var err error

	diffArgs := strings.Join(a.Config.Diff[1:], " ")
	tmpl, err = tmpl.New("diff").Parse(diffArgs)
	if err != nil {
		return err
	}
	mergeArgs := strings.Join(a.Config.Merge[1:], " ")
	tmpl, err = tmpl.New("merge").Parse(mergeArgs)
	if err != nil {
		return err
	}

	a.template = tmpl

	return nil
}

type templateData struct {
	Source      string
	Destination string
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
	a.Logger.Info().Str("name", filepath.Dir(path)).Msg("Created")

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
	diffCmdName := a.Config.Diff[0]
	for _, pm := range mapper.Mapping {
		ss, err := fileSystem.GetSum(pm.Source)
		if err != nil {
			return err
		}
		ds, err := fileSystem.GetSum(pm.Destination)
		if err != nil {
			return err
		}
		if bytes.Equal(ss, ds) {
			continue
		}

		argsBuilder := strings.Builder{}
		if err := a.template.ExecuteTemplate(&argsBuilder, "diff", templateData(pm)); err != nil {
			return err
		}
		args := strings.Split(argsBuilder.String(), " ")
		cmd := exec.Command(diffCmdName, args...)
		out, _ := cmd.Output()
		a.Logger.Info().Str("command", diffCmdName).Strs("args", args).Msg("Executed")

		diff = append(diff, out...)
	}

	pagerCmdName := a.Config.Pager[0]
	pagerCmdArgs := a.Config.Pager[1:]
	cmd := exec.Command(pagerCmdName, pagerCmdArgs...)
	cmd.Stdin = bytes.NewBuffer(diff)
	cmd.Stdout = a.out
	if err := cmd.Run(); err != nil {
		return err
	}
	a.Logger.Info().Str("command", pagerCmdName).Strs("args", pagerCmdArgs).Msg("Executed")
	return nil
}

func (a *App) Merge() error {
	mapper, err := NewPathMapperBuilder(
		a.Config.Source, a.Config.Destination, WithExcludes(a.Config.Excludes...),
	).Build()
	if err != nil {
		return err
	}

	mergeCmdName := a.Config.Merge[0]
	for _, pm := range mapper.Mapping {
		ss, err := fileSystem.GetSum(pm.Source)
		if err != nil {
			return err
		}
		ds, err := fileSystem.GetSum(pm.Destination)
		if err != nil {
			return err
		}
		if bytes.Equal(ss, ds) {
			continue
		}

		argsBuilder := strings.Builder{}
		if err := a.template.ExecuteTemplate(&argsBuilder, "merge", templateData(pm)); err != nil {
			return err
		}
		args := strings.Split(argsBuilder.String(), " ")
		cmd := exec.Command(mergeCmdName, args...)
		cmd.Stdin = a.in
		cmd.Stdout = a.out
		if err := cmd.Run(); err != nil {
			return err
		}
		a.Logger.Info().Str("command", mergeCmdName).Strs("args", args).Msg("Executed")
	}

	return nil
}

func (a *App) Where(dir string) error {
	switch dir {
	case "source":
		fmt.Fprint(a.out, a.Config.Source)
	case "destination":
		fmt.Fprint(a.out, a.Config.Destination)
	case "config":
		fmt.Fprint(a.out, filepath.Dir(GetConfig().ConfigFileUsed()))
	default:
	}
	return nil
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
	editorCmdName := a.Config.Editor[0]
	editorCmdArgs := a.Config.Editor[1:]
	args := append(editorCmdArgs, v.ConfigFileUsed())
	cmd := exec.Command(editorCmdName, args...)
	cmd.Stdin = a.in
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
		ss, err := fileSystem.GetSum(pm.Source)
		if err != nil {
			return err
		}
		ds, err := fileSystem.GetSum(pm.Destination)
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
		a.Logger.Info().Str("name", dir).Msg("Created")
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

	e, err := fileSystem.Get(src)
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
	a.Logger.Info().Str("name", dst).Msg("Updated")

	return nil
}
