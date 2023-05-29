package donut

import (
	"errors"
	"io"
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

// AbsPath returns if path is relative, joins baseDir. if path is absolute, clean the path.
func AbsPath(path, baseDir string) string {
	if filepath.IsAbs(path) {
		return filepath.Clean(path)
	}
	return filepath.Join(baseDir, path)
}

func copyFile(sourcePath, destinationPath string) error {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()
	destinationFile, err := os.Create(destinationPath)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	return err
}
