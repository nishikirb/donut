package donut

import (
	"io"

	"github.com/nishikirb/donut/config"
)

type Option func(*App) error

func WithConfig(c *config.Config) Option {
	return func(a *App) error {
		a.config = c
		return nil
	}
}

func WithConfigLoader(opts ...config.ConfigOption) Option {
	return func(a *App) error {
		if c, err := config.New(opts...); err != nil {
			return err
		} else {
			a.config = c
		}
		return nil
	}
}

func WithIn(r io.Reader) Option {
	return func(a *App) error {
		a.in = r
		return nil
	}
}

func WithOut(w io.Writer) Option {
	return func(a *App) error {
		a.out = w
		return nil
	}
}

func WithErr(w io.Writer) Option {
	return func(a *App) error {
		a.err = w
		return nil
	}
}
