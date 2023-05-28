package donut

import (
	"errors"
	"os"
	"path/filepath"
)

const fileType = "toml"

var UserHomeDir string

func init() {
	var err error
	if UserHomeDir, err = os.UserHomeDir(); err != nil {
		panic(err)
	}
}

func DefaultConfigFile() string {
	return filepath.Join(UserHomeDir, ".config", AppName, AppName+"."+fileType)
}

func DefaultConfigDirs() []string {
	return []string{
		"$XDG_CONFIG_HOME",
		filepath.Join("$XDG_CONFIG_HOME", AppName),
		filepath.Join(UserHomeDir, ".config"),
		filepath.Join(UserHomeDir, ".config", AppName),
	}
}

func DefaultSourceDir() string {
	return filepath.Join(UserHomeDir, ".local", "share", AppName)
}

func SetUserHomeDir(dir string) func() {
	tmp := UserHomeDir
	UserHomeDir = dir
	return func() { UserHomeDir = tmp }
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
