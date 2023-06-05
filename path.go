package donut

import (
	"errors"
	"os"
)

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

// abs returns if path is relative, joins baseDir. if path is absolute, clean the path.
// func abs(path, baseDir string) string {
// 	if filepath.IsAbs(path) {
// 		return filepath.Clean(path)
// 	}
// 	return filepath.Join(baseDir, path)
// }
