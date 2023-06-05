package donut

import (
	"sync"
)

type FileEntryMap struct {
	Cache sync.Map
}

var fileEntryMap = &FileEntryMap{}

func (m *FileEntryMap) Get(path string) (*FileEntry, error) {
	if v, ok := m.Cache.Load(path); ok {
		return v.(*FileEntry), nil
	}
	if f, err := NewFileEntry(path); err != nil {
		return nil, err
	} else {
		m.Set(path, f)
		return f, nil
	}
}

func (m *FileEntryMap) GetSum(path string) ([]byte, error) {
	e, err := m.Get(path)
	if err != nil {
		return nil, err
	}
	return e.GetSum()
}

func (m *FileEntryMap) Set(path string, f *FileEntry) {
	m.Cache.Store(path, f)
}

func (m *FileEntryMap) Reload(path string) (*FileEntry, error) {
	if f, err := NewFileEntry(path); err != nil {
		return nil, err
	} else {
		m.Set(path, f)
		return f, nil
	}
}
