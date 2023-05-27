package donut

import (
	"github.com/spf13/viper"
)

type Config struct {
	Source      string   `mapstructure:"source"`
	Destination string   `mapstructure:"destination"`
	Excludes    []string `mapstructure:"excludes"`
}

type ConfigOption func(v *viper.Viper) error

var viperInstance = viper.New()

func InitConfig(cfgPath string) error {
	var opts []ConfigOption
	if cfgPath != "" {
		opts = append(opts, WithDefault(), WithFile(cfgPath))
	} else {
		opts = append(opts, WithDefault(), WithNameAndPath(AppName, DefaultConfigDirs()...))
	}

	viperInstance = viper.NewWithOptions(viper.KeyDelimiter("::"))
	for _, opt := range opts {
		if err := opt(viperInstance); err != nil {
			return err
		}
	}

	return nil
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
