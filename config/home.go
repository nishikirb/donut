package config

import (
	"fmt"
	"os"
	"path/filepath"
)

const AppName string = "donut"

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

func DefaultStateDir() string {
	return filepath.Join(UserHomeDir, ".local", "state", AppName)
}

func DefaultConfigFile() string {
	return filepath.Join(UserHomeDir, ".config", AppName, fmt.Sprintf("%s.%s", AppName, "toml"))
}

func defaultConfigDirs() []string {
	return []string{
		"$XDG_CONFIG_HOME",
		filepath.Join("$XDG_CONFIG_HOME", AppName),
		filepath.Join(UserHomeDir, ".config"),
		filepath.Join(UserHomeDir, ".config", AppName),
	}
}

func defaultSourceDir() string {
	return filepath.Join(UserHomeDir, ".local", "share", AppName)
}
