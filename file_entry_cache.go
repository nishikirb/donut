package donut

import (
	"sync"
)

type FileEntryCache struct {
	cache sync.Map
}

var fileEntryCache = &FileEntryCache{}

func (c *FileEntryCache) Get(path string) (*FileEntry, error) {
	if v, ok := c.cache.Load(path); ok {
		return v.(*FileEntry), nil
	}
	if f, err := NewFileEntry(path); err != nil {
		return nil, err
	} else {
		c.Set(path, f)
		return f, nil
	}
}

func (c *FileEntryCache) GetSum(path string) ([]byte, error) {
	e, err := c.Get(path)
	if err != nil {
		return nil, err
	}
	return e.GetSum()
}

func (c *FileEntryCache) Set(path string, f *FileEntry) {
	c.cache.Store(path, f)
}

func (c *FileEntryCache) Reload(path string) (*FileEntry, error) {
	if f, err := NewFileEntry(path); err != nil {
		return nil, err
	} else {
		c.Set(path, f)
		return f, nil
	}
}
