package donut

import (
	"path/filepath"
	"sync"

	"github.com/spf13/viper"
)

const AppName string = "donut"
const defaultConfigExt = "toml"

type Config struct {
	Source      string   `mapstructure:"source"`
	Destination string   `mapstructure:"destination"`
	Excludes    []string `mapstructure:"excludes"`
	Editor      Command  `mapstructure:"editor"`
	Diff        Command  `mapstructure:"diff"`
	Pager       Command  `mapstructure:"pager"`
	Merge       Command  `mapstructure:"merge"`
}

type Command struct {
	Name string   `mapstructure:"command"`
	Args []string `mapstructure:"args"`
}

type ConfigOption func(v *viper.Viper) error

var (
	config     = viper.New()
	initConfig sync.Once
)

func InitConfig(path string) error {
	var err error
	initConfig.Do(func() {
		var opts []ConfigOption
		if path != "" {
			opts = append(opts, WithDefault(), WithFile(path))
		} else {
			opts = append(opts, WithDefault(), WithNameAndPath(AppName, DefaultConfigDirs()...))
		}

		config, err = NewConfig(opts...)
	})

	return err
}

func NewConfig(opts ...ConfigOption) (*viper.Viper, error) {
	v := viper.NewWithOptions(viper.KeyDelimiter("::"))

	for _, opt := range opts {
		if err := opt(v); err != nil {
			return nil, err
		}
	}

	return v, nil
}

func GetConfig() *viper.Viper {
	return config
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
		v.SetDefault("destination", homedir)
		v.SetDefault("editor::command", "vim")
		v.SetDefault("editor::args", []string{})
		v.SetDefault("diff::command", "diff")
		v.SetDefault("diff::args", []string{"-upN"})
		v.SetDefault("pager::command", "less")
		v.SetDefault("pager::args", []string{"-R"})
		v.SetDefault("merge::command", "vimdiff")
		v.SetDefault("merge::args", []string{})
		return nil
	}
}

func DefaultConfigFile() string {
	return filepath.Join(homedir, ".config", AppName, AppName+"."+defaultConfigExt)
}

func DefaultConfigDirs() []string {
	return []string{
		"$XDG_CONFIG_HOME",
		filepath.Join("$XDG_CONFIG_HOME", AppName),
		filepath.Join(homedir, ".config"),
		filepath.Join(homedir, ".config", AppName),
	}
}

func defaultSourceDir() string {
	return filepath.Join(homedir, ".local", "share", AppName)
}

func defaultStateDir() string {
	return filepath.Join(homedir, ".local", "state", AppName)
}
