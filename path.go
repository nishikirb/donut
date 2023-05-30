package donut

import (
	"errors"
	"os"
	"path/filepath"
)

var homedir string

func init() {
	var err error
	if homedir, err = os.UserHomeDir(); err != nil {
		panic(err)
	}
}

// AbsPath returns if path is relative, joins baseDir. if path is absolute, clean the path.
func AbsPath(path, baseDir string) string {
	if filepath.IsAbs(path) {
		return filepath.Clean(path)
	}
	return filepath.Join(baseDir, path)
}

// IsDir checks path is exists and is directory
func IsDir(path string) error {
	if path == "" {
		return errors.New("not defined")
	} else if f, err := os.Stat(path); err != nil {
		return err
	} else if !f.IsDir() {
		return errors.New("not directory")
	}
	return nil
}

func SetUserHomeDir(dir string) func() {
	tmp := homedir
	homedir = dir
	return func() { homedir = tmp }
}
