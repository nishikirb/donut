package donut

import (
	"sync"
)

type FileSystem struct {
	FileEntryMap sync.Map
}

var fileSystem = &FileSystem{}

func (m *FileSystem) Get(path string) (*FileEntry, error) {
	if v, ok := m.FileEntryMap.Load(path); ok {
		return v.(*FileEntry), nil
	}
	if f, err := NewFileEntry(path); err != nil {
		return nil, err
	} else {
		m.Set(path, f)
		return f, nil
	}
}

func (m *FileSystem) GetSum(path string) ([]byte, error) {
	e, err := m.Get(path)
	if err != nil {
		return nil, err
	}
	return e.GetSum()
}

func (m *FileSystem) Set(path string, f *FileEntry) {
	m.FileEntryMap.Store(path, f)
}

func (m *FileSystem) Reload(path string) (*FileEntry, error) {
	if f, err := NewFileEntry(path); err != nil {
		return nil, err
	} else {
		m.Set(path, f)
		return f, nil
	}
}
