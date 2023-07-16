package config

import (
	"errors"
	"os"
	"runtime"

	"github.com/spf13/viper"
)

type Config struct {
	Source      string   `mapstructure:"source"`
	Destination string   `mapstructure:"destination"`
	Excludes    []string `mapstructure:"excludes"`
	Editor      []string `mapstructure:"editor"`
	Pager       []string `mapstructure:"pager"`
	Diff        []string `mapstructure:"diff"`
	Merge       []string `mapstructure:"merge"`
	Concurrency int
	File        string
}

func New(opts ...ConfigOption) (*Config, error) {
	v, err := NewViper(opts...)
	if err != nil {
		return nil, err
	}

	c := Config{
		Concurrency: runtime.NumCPU(),
		File:        v.ConfigFileUsed(),
	}
	if err := decode(v, &c); err != nil {
		return nil, err
	}

	return &c, nil
}

func NewViper(opts ...ConfigOption) (*viper.Viper, error) {
	v := viper.NewWithOptions(viper.KeyDelimiter("::"))
	for _, opt := range opts {
		if err := opt(v); err != nil {
			return nil, err
		}
	}
	return v, nil
}

func decode(v *viper.Viper, c *Config) error {
	if v == nil {
		return errors.New("config cannot be nil")
	}
	if err := v.Unmarshal(&c, viper.DecodeHook(defaultDecodeHookFunc)); err != nil {
		return err
	}
	if err := validate(c); err != nil {
		return err
	}
	return nil
}

func validate(c *Config) error {
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
