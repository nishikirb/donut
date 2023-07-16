package donut

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"

	"github.com/gleamsoda/donut/config"
	"github.com/gleamsoda/donut/store"
	"github.com/gleamsoda/donut/system"
)

type App struct {
	commands map[string]handler
	opts     []Option
	config   *config.Config
	in       io.Reader
	out      io.Writer
	err      io.Writer
}

type handler func(ctx context.Context, args []string, flags *pflag.FlagSet) error

func NewApp(opts ...Option) *App {
	cfg, _ := config.New(config.WithDefault())
	app := &App{
		in:     os.Stdin,
		out:    os.Stdout,
		err:    os.Stderr,
		config: cfg,
	}

	app.handle("init", app.init)
	app.handle("list", app.list)
	app.handle("diff", app.diff)
	app.handle("merge", app.merge)
	app.handle("where", app.where)
	app.handle("config", app.configEdit)
	app.handle("apply", app.apply)
	app.handle("clean", app.clean)

	return app
}

func (a *App) AddOptions(opts ...Option) *App {
	a.opts = append(a.opts, opts...)
	return a
}

func (a *App) Run(ctx context.Context, command string, args []string, flags *pflag.FlagSet) error {
	h, ok := a.commands[command]
	if !ok {
		return fmt.Errorf("unknown command: %s", command)
	}

	if err := a.applyOptions(); err != nil {
		return err
	}

	if err := createTemplateMap(map[string][]string{
		"diff":  a.config.Diff[1:],
		"merge": a.config.Merge[1:],
	}); err != nil {
		return err
	}

	return h(ctx, args, flags)
}

func (a *App) init(_ context.Context, _ []string, _ *pflag.FlagSet) error {
	// check config file exists in default path
	// if exists, no need to init
	_, err := config.New(config.WithPath("")...)
	if err == nil {
		return errors.New("config file already exists. no need to init")
	} else if !errors.As(err, &viper.ConfigFileNotFoundError{}) {
		return fmt.Errorf("config file already exists, but error: %w", err)
	}

	path := config.DefaultConfigFile()
	if err := system.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return err
	}

	v, _ := config.NewViper(config.WithDefault())
	if err := system.WriteConfig(v, path); err != nil {
		return err
	}

	fmt.Fprintf(a.out, "Created: %s\n", path)
	return nil
}

func (a *App) list(_ context.Context, _ []string, _ *pflag.FlagSet) error {
	mapper, err := NewPathMapper(a.config.Source, a.config.Destination, WithExcludes(a.config.Excludes...))
	if err != nil {
		return err
	}

	for _, relSourcePath := range mapper.RelSourcePaths() {
		fmt.Fprintln(a.out, relSourcePath)
	}
	return nil
}

func (a *App) diff(ctx context.Context, _ []string, _ *pflag.FlagSet) error {
	mapper, err := NewPathMapper(a.config.Source, a.config.Destination, WithExcludes(a.config.Excludes...))
	if err != nil {
		return err
	}

	diffCmdName := a.config.Diff[0]
	diffCh := make(chan []byte, len(mapper.Mapping))
	eg, ectx := errgroup.WithContext(context.Background())
	eg.SetLimit(a.config.Concurrency)
	for _, pm := range mapper.Mapping {
		pm := pm
		eg.Go(func() error {
			select {
			case <-ectx.Done():
				return ectx.Err()
			default:
				ss, err := fileEntryCache.GetSum(pm.Source)
				if err != nil {
					return err
				}
				ds, err := fileEntryCache.GetSum(pm.Destination)
				if err != nil {
					return err
				}
				if bytes.Equal(ss, ds) {
					return nil
				}

				argsBuilder := strings.Builder{}
				data := templateParams(pm)
				if err := tmpl.ExecuteTemplate(&argsBuilder, "diff", data); err != nil {
					return err
				}
				args := strings.Split(argsBuilder.String(), " ")
				cmd := exec.CommandContext(ectx, diffCmdName, args...)
				out, _ := system.Output(cmd)
				diffCh <- out
				return nil
			}
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}
	close(diffCh)
	var diff []byte
	for v := range diffCh {
		diff = append(diff, v...)
	}

	pagerCmdName, pagerCmdArgs := a.config.Pager[0], a.config.Pager[1:]
	cmd := exec.CommandContext(ctx, pagerCmdName, pagerCmdArgs...)
	cmd.Stdin = bytes.NewBuffer(diff)
	cmd.Stdout = a.out
	if err := system.Run(cmd); err != nil {
		return err
	}
	return nil
}

func (a *App) merge(ctx context.Context, _ []string, _ *pflag.FlagSet) error {
	mapper, err := NewPathMapper(a.config.Source, a.config.Destination, WithExcludes(a.config.Excludes...))
	if err != nil {
		return err
	}

	mergeCmdName := a.config.Merge[0]
	for _, pm := range mapper.Mapping {
		ss, err := fileEntryCache.GetSum(pm.Source)
		if err != nil {
			return err
		}
		ds, err := fileEntryCache.GetSum(pm.Destination)
		if err != nil {
			return err
		}
		if bytes.Equal(ss, ds) {
			continue
		}

		argsBuilder := strings.Builder{}
		data := templateParams(pm)
		if err := tmpl.ExecuteTemplate(&argsBuilder, "merge", data); err != nil {
			return err
		}
		args := strings.Split(argsBuilder.String(), " ")
		cmd := exec.CommandContext(ctx, mergeCmdName, args...)
		cmd.Stdin = a.in
		cmd.Stdout = a.out
		if err := system.Run(cmd); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) where(_ context.Context, args []string, _ *pflag.FlagSet) error {
	if len(args) != 1 {
		return errors.New("invalid argument")
	}
	switch dir := args[0]; dir {
	case "source":
		fmt.Fprintln(a.out, a.config.Source)
	case "destination":
		fmt.Fprintln(a.out, a.config.Destination)
	case "config":
		fmt.Fprintln(a.out, filepath.Dir(a.config.File))
	default:
	}
	return nil
}

func (a *App) configEdit(_ context.Context, _ []string, _ *pflag.FlagSet) error {
	editorCmdName := a.config.Editor[0]
	editorCmdArgs := a.config.Editor[1:]
	args := append(editorCmdArgs, a.config.File)
	cmd := exec.Command(editorCmdName, args...)
	cmd.Stdin = a.in
	cmd.Stdout = a.out
	return system.Run(cmd)
}

func (a *App) apply(ctx context.Context, _ []string, flags *pflag.FlagSet) error {
	overwrite, _ := flags.GetBool("overwrite")

	mapper, err := NewPathMapper(a.config.Source, a.config.Destination, WithExcludes(a.config.Excludes...))
	if err != nil {
		return err
	}

	eg, ectx := errgroup.WithContext(ctx)
	eg.SetLimit(a.config.Concurrency)
	for _, pm := range mapper.Mapping {
		pm := pm
		eg.Go(func() error {
			select {
			case <-ectx.Done():
				return ectx.Err()
			default:
				ss, err := fileEntryCache.GetSum(pm.Source)
				if err != nil {
					return err
				}
				ds, err := fileEntryCache.GetSum(pm.Destination)
				if err != nil {
					return err
				}
				if bytes.Equal(ss, ds) {
					return nil
				}

				var be *FileEntry
				if err := store.Get(store.FileEntryBucket, pm.Destination, &be); err != nil {
					return err
				}
				bs, err := be.GetSum()
				if err != nil {
					return err
				}

				// skip if the following conditions are met
				// 1. exist in store
				// 2. checksum is equal to destination
				// 3. overwrite is false
				if bs != nil && !bytes.Equal(bs, ds) && !overwrite {
					fmt.Fprintf(a.out, "Skipped: %s has been modified since the last apply. use --overwrite to overwrite\n", pm.Destination)
					return nil
				}

				// If the directory does not exist, create it
				// os.MkdirAll will return nil if directory already exists
				dir := filepath.Dir(pm.Destination)
				if err := system.MkdirAll(dir, os.ModePerm); err != nil {
					return err
				}
				if err := a.overwrite(pm.Source, pm.Destination); err != nil {
					return err
				}

				de, err := fileEntryCache.Reload(pm.Destination)
				if err != nil {
					return err
				}
				if err := store.Set(store.FileEntryBucket, pm.Destination, de); err != nil {
					return err
				}
				fmt.Fprintf(a.out, "Applied: %s from %s\n", pm.Destination, pm.Source)
				return nil
			}
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}

func (a *App) clean(ctx context.Context, _ []string, flags *pflag.FlagSet) error {
	if err := store.Close(); err != nil {
		return err
	}
	if err := system.Remove(store.DefaultDBFile()); err != nil {
		return err
	}
	return nil
}

func (a *App) handle(name string, h handler) {
	if a.commands == nil {
		a.commands = make(map[string]handler)
	}
	a.commands[name] = h
}

func (a *App) applyOptions() error {
	for _, opt := range a.opts {
		if err := opt(a); err != nil {
			return err
		}
	}
	return nil
}

// overwrite replaces the contents of dst with the contents of src.
func (a *App) overwrite(src, dst string) error {
	se, err := fileEntryCache.Get(src)
	if err != nil {
		return err
	}
	sc, err := se.GetContent()
	if err != nil {
		return err
	}
	if err := system.Overwrite(dst, sc, os.ModePerm); err != nil {
		return err
	}
	return nil
}
