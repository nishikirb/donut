package donut

import (
	"github.com/spf13/viper"
)

type Config struct {
	Source      string   `mapstructure:"source"`
	Destination string   `mapstructure:"destination"`
	Excludes    []string `mapstructure:"excludes"`
	Editor      string   `mapstructure:"editor"`
	Diff        string   `mapstructure:"diff"`
	Pager       string   `mapstructure:"pager"`
}

type ConfigOption func(v *viper.Viper) error

var viperInstance = viper.New()

func InitConfig(configFile string) error {
	var opts []ConfigOption
	if configFile != "" {
		opts = append(opts, WithDefault(), WithFile(configFile))
	} else {
		opts = append(opts, WithDefault(), WithNameAndPath(AppName, DefaultConfigDirs()...))
	}

	var err error
	viperInstance, err = NewConfig(opts...)

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
	return viperInstance
}

func WriteConfig(filename string) error {
	return viperInstance.WriteConfigAs(filename)
}

func WithFile(file string) ConfigOption {
	return func(v *viper.Viper) error {
		v.SetConfigFile(file)
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
		v.SetDefault("source", DefaultSourceDir())
		v.SetDefault("destination", UserHomeDir)
		return nil
	}
}
