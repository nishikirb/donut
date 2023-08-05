package donut

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"time"

	"github.com/gleamsoda/donut/system"
)

type Entry struct {
	Path      string      `json:"-"`
	Empty     bool        `json:"empty"`
	Mode      fs.FileMode `json:"mode"`
	ModTime   time.Time   `json:"mod_time"`
	sum       []byte
	content   []byte
	isFetched bool `json:"-"`
}

var _ json.Marshaler = (*Entry)(nil)
var _ json.Unmarshaler = (*Entry)(nil)

func NewEntry(path string) (*Entry, error) {
	f, err := system.Lstat(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return &Entry{
				Path:  path,
				Empty: true,
			}, nil
		} else {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
	}

	return &Entry{
		Path:    path,
		Empty:   false,
		Mode:    f.Mode(),
		ModTime: f.ModTime(),
	}, nil
}

func (e *Entry) GetSum() ([]byte, error) {
	if e == nil || e.Empty {
		return nil, nil
	}
	if e.sum == nil && !e.isFetched {
		if err := e.loadSum(); err != nil {
			return nil, err
		}
	}
	return e.sum, nil
}

func (e *Entry) GetContent() ([]byte, error) {
	if e == nil || e.Empty {
		return nil, nil
	}
	if e.content == nil && !e.isFetched {
		if err := e.loadContent(); err != nil {
			return nil, err
		}
	}
	return e.content, nil
}

// func (e *Entry) isDir() bool {
// 	return e.Mode.IsDir()
// }

// func (e *Entry) isSymLink() bool {
// 	return e.Mode&os.ModeSymlink != 0
// }

// func (f *Entry) isSame(path string) (bool, error) {
// 	if !f.isSymLink() {
// 		return f.Path == path, nil
// 	}
// 	l, err := os.Readlink(f.Path)
// 	if err != nil {
// 		return false, fmt.Errorf("%s: %w", f.Path, err)
// 	}
// 	return l == path, nil
// }

// MarshalJSON implements json.Marshaler interface.
func (e *Entry) MarshalJSON() ([]byte, error) {
	type Alias Entry // エイリアスを作成して、再帰的な呼び出しを避ける

	sum, err := e.GetSum()
	if err != nil {
		return nil, err
	}
	return json.Marshal(&struct {
		*Alias
		Sum []byte `json:"sum"`
	}{
		Alias: (*Alias)(e),
		Sum:   sum,
	})
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (e *Entry) UnmarshalJSON(value []byte) error {
	type Alias Entry // エイリアスを作成して、再帰的な呼び出しを避ける

	aux := &struct {
		*Alias
		Sum []byte `json:"sum"`
	}{
		Alias: (*Alias)(e),
	}
	if err := json.Unmarshal(value, &aux); err != nil {
		return err
	}

	e.sum = aux.Sum
	e.isFetched = true
	return nil
}

func (e *Entry) loadSum() error {
	file, err := system.Open(e.Path)
	if err != nil {
		return err
	}
	defer file.Close()
	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		return err
	}
	e.sum = h.Sum(nil)
	return nil
}

func (e *Entry) loadContent() error {
	r, err := system.ReadFile(e.Path)
	if err != nil {
		return err
	}
	h := sha256.New()
	if _, err := h.Write(r); err != nil {
		return err
	}
	e.content = r
	e.sum = h.Sum(nil)
	return nil
}
