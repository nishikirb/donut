package donut

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const appName string = "donut"
const defaultConfigExt = "toml"

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

type ConfigOption func(v *viper.Viper) error

var defaultConcurrency = 8
var UserHomeDir string

func init() {
	var err error
	if UserHomeDir, err = os.UserHomeDir(); err != nil {
		panic(err)
	}
}

func SetUserHomeDir(dir string) func() {
	tmp := UserHomeDir
	UserHomeDir = dir
	return func() { UserHomeDir = tmp }
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

func WithPath(path string) []ConfigOption {
	var opts []ConfigOption
	if path != "" {
		opts = []ConfigOption{WithDefault(), WithFile(path)}
	} else {
		opts = []ConfigOption{WithDefault(), WithNameAndPath(appName, defaultConfigDirs()...)}
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

func defaultConfigFile() string {
	return filepath.Join(UserHomeDir, ".config", appName, appName+"."+defaultConfigExt)
}

func defaultConfigDirs() []string {
	return []string{
		"$XDG_CONFIG_HOME",
		filepath.Join("$XDG_CONFIG_HOME", appName),
		filepath.Join(UserHomeDir, ".config"),
		filepath.Join(UserHomeDir, ".config", appName),
	}
}

func defaultSourceDir() string {
	return filepath.Join(UserHomeDir, ".local", "share", appName)
}

func defaultStateDir() string {
	return filepath.Join(UserHomeDir, ".local", "state", appName)
}

const FileEntryBucket = "file_entries"

var Buckets = []string{FileEntryBucket}

func DefaultDBFile() string {
	return filepath.Join(defaultStateDir(), "donut.db")
}
