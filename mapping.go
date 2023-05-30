package donut

import (
	"io/fs"
	"path/filepath"
	"strings"

	"golang.org/x/exp/slices"
)

type Mapping struct {
	Source      File
	Destination File
}

type Mappings []Mapping

type MappingsBuilder struct {
	Source      string
	Destination string
	Excludes    []string
}

type MappingsBuilderOption func(b *MappingsBuilder)

var mustExcludes = []string{".git"}

func NewMapping(src, dst string) (*Mapping, error) {
	var sf, df *File
	var err error
	if sf, err = NewFile(src); err != nil {
		return nil, err
	}
	if df, err = NewFile(dst); err != nil {
		return nil, err
	}

	return &Mapping{
		Source:      *sf,
		Destination: *df,
	}, nil
}

func NewMappingsBuilder(s, d string, funcs ...MappingsBuilderOption) *MappingsBuilder {
	b := &MappingsBuilder{
		Source:      s,
		Destination: d,
		Excludes:    mustExcludes,
	}

	for _, fn := range funcs {
		fn(b)
	}

	return b
}

func WithExcludes(s ...string) MappingsBuilderOption {
	return func(b *MappingsBuilder) {
		b.Excludes = append(b.Excludes, s...)
	}
}

func (b *MappingsBuilder) Build() ([]Mapping, error) {
	var mps []Mapping
	err := filepath.WalkDir(b.Source, func(path string, d fs.DirEntry, _ error) error {
		prefixTrimmed := strings.TrimPrefix(strings.TrimPrefix(path, b.Source), string(filepath.Separator))

		eq := func(s string) bool {
			ok, _ := filepath.Match(s, prefixTrimmed)
			return ok
		}
		if d.IsDir() {
			// skip directoires if it is in b.excludes
			if slices.IndexFunc(b.Excludes, eq) != -1 {
				return fs.SkipDir
			}
			return nil
		}
		// ignore b.excludes
		if slices.IndexFunc(b.Excludes, eq) != -1 {
			return nil
		}
		// Specify the destination path
		dPath := strings.Replace(path, b.Source, b.Destination, 1)
		m, err := NewMapping(path, dPath)
		if err != nil {
			return err
		}
		mps = append(mps, *m)
		return nil
	})

	return mps, err
}
