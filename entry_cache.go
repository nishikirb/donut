package donut

import (
	"sync"
)

type EntryCache struct {
	cache sync.Map
}

var entryCache = &EntryCache{}

func (c *EntryCache) Get(path string) (*Entry, error) {
	if v, ok := c.cache.Load(path); ok {
		return v.(*Entry), nil
	}
	if f, err := NewEntry(path); err != nil {
		return nil, err
	} else {
		c.Set(path, f)
		return f, nil
	}
}

func (c *EntryCache) GetSum(path string) ([]byte, error) {
	e, err := c.Get(path)
	if err != nil {
		return nil, err
	}
	return e.GetSum()
}

func (c *EntryCache) Set(path string, f *Entry) {
	c.cache.Store(path, f)
}

func (c *EntryCache) Reload(path string) (*Entry, error) {
	if f, err := NewEntry(path); err != nil {
		return nil, err
	} else {
		c.Set(path, f)
		return f, nil
	}
}
