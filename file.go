package donut

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

type File struct {
	Path     string
	NotExist bool
	FileInfo fs.FileInfo
}

func NewFile(path string) (*File, error) {
	var notExist bool
	f, err := os.Lstat(path)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		notExist = true
	}
	return &File{
		Path:     path,
		NotExist: notExist,
		FileInfo: f,
	}, nil
}

func (f *File) IsSymLink() bool {
	return f.FileInfo.Mode()&os.ModeSymlink != 0
}

func (f *File) IsSame(path string) (bool, error) {
	if !f.IsSymLink() {
		return f.Path == path, nil
	}
	l, err := os.Readlink(f.Path)
	if err != nil {
		return false, fmt.Errorf("%s: %w", f.Path, err)
	}
	return l == path, nil
}
