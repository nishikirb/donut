package config

import (
	"github.com/spf13/viper"
)

type ConfigOption func(v *viper.Viper) error

func WithPath(path string) []ConfigOption {
	var opts []ConfigOption
	if path != "" {
		opts = []ConfigOption{WithDefault(), WithFile(path)}
	} else {
		opts = []ConfigOption{WithDefault(), WithNameAndPath(AppName, defaultConfigDirs()...)}
	}
	return opts
}

func WithFile(path string) ConfigOption {
	return func(v *viper.Viper) error {
		v.SetConfigFile(path)
		if err := v.ReadInConfig(); err != nil {
			return err
		}
		return nil
	}
}

func WithNameAndPath(name string, paths ...string) ConfigOption {
	return func(v *viper.Viper) error {
		v.SetConfigName(name)
		for _, path := range paths {
			v.AddConfigPath(path)
		}
		if err := v.ReadInConfig(); err != nil {
			return err
		}
		return nil
	}
}

func WithData(data map[string]interface{}) ConfigOption {
	return func(v *viper.Viper) error {
		for k, d := range data {
			v.Set(k, d)
		}
		return nil
	}
}

func WithDefault() ConfigOption {
	return func(v *viper.Viper) error {
		v.SetDefault("source", defaultSourceDir())
		v.SetDefault("destination", UserHomeDir)
		v.SetDefault("editor", []string{"vim"})
		v.SetDefault("pager", []string{"less", "-R"})
		v.SetDefault("diff", []string{"diff", "-upN", "{{.Destination}}", "{{.Source}}"})
		v.SetDefault("merge", []string{"vimdiff", "{{.Destination}}", "{{.Source}}"})
		return nil
	}
}
