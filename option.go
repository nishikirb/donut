package donut

import (
	"errors"
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type Option func(*App) error

func WithConfig(v *viper.Viper) Option {
	return func(a *App) error {
		var c Config
		if err := decodeConfig(v, &c); err != nil {
			return err
		}
		a.config = &c
		return nil
	}
}

func WithConfigLoader(opts ...ConfigOption) Option {
	return func(a *App) error {
		var c Config
		if v, err := NewConfig(opts...); err != nil {
			return err
		} else if err := decodeConfig(v, &c); err != nil {
			return err
		}
		a.config = &c
		return nil
	}
}

func WithStore(s *Store) Option {
	return func(a *App) error {
		if s == nil {
			return errors.New("store cannot be nil")
		}
		a.store = s
		return nil
	}
}

func WithLogger(l zerolog.Logger) Option {
	return func(a *App) error {
		a.logger = l
		a.executor.SetLogger(l)
		return nil
	}
}

func WithIn(r io.Reader) Option {
	return func(a *App) error {
		a.in = r
		return nil
	}
}

func WithOut(r io.Writer) Option {
	return func(a *App) error {
		a.out = r
		return nil
	}
}

func WithErr(r io.Writer) Option {
	return func(a *App) error {
		a.err = r
		return nil
	}
}

func decodeConfig(v *viper.Viper, c *Config) error {
	if v == nil {
		return errors.New("config cannot be nil")
	}
	if err := v.Unmarshal(&c, viper.DecodeHook(defaultDecodeHookFunc)); err != nil {
		return err
	}
	if err := validateConfig(c); err != nil {
		return err
	}
	c.Concurrency = defaultConcurrency
	c.File = v.ConfigFileUsed()
	return nil
}

func validateConfig(c *Config) error {
	for _, path := range []string{c.Source, c.Destination} {
		if err := isDir(path); err != nil {
			return err
		}
	}
	return nil
}

// isDir checks path is exists and is directory
func isDir(path string) error {
	if path == "" {
		return errors.New("not defined")
	} else if f, err := os.Stat(path); err != nil {
		return err
	} else if !f.IsDir() {
		return errors.New("not directory")
	}
	return nil
}

// abs returns if path is relative, joins baseDir. if path is absolute, clean the path.
// func abs(path, baseDir string) string {
// 	if filepath.IsAbs(path) {
// 		return filepath.Clean(path)
// 	}
// 	return filepath.Join(baseDir, path)
// }
