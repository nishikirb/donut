package donut

import (
	"io/fs"
	"path/filepath"
	"strings"

	"golang.org/x/exp/slices"
)

type PathMapperBuilder struct {
	Source      string
	Destination string
	Excludes    []string
}

type PathMapperBuilderOption func(b *PathMapperBuilder)

var mustExcludes = []string{".git"}

func NewPathMapperBuilder(s, d string, funcs ...PathMapperBuilderOption) *PathMapperBuilder {
	b := &PathMapperBuilder{
		Source:      s,
		Destination: d,
		Excludes:    mustExcludes,
	}

	for _, fn := range funcs {
		fn(b)
	}

	return b
}

func WithExcludes(s ...string) PathMapperBuilderOption {
	return func(b *PathMapperBuilder) {
		b.Excludes = append(b.Excludes, s...)
	}
}

func (b *PathMapperBuilder) Build() (*PathMapper, error) {
	mapper := NewPathMapper(b.Source, b.Destination)

	err := filepath.WalkDir(b.Source, func(path string, d fs.DirEntry, _ error) error {
		rel, _ := filepath.Rel(b.Source, path)
		eq := func(s string) bool {
			ok, _ := filepath.Match(s, rel)
			return ok
		}
		if d.IsDir() {
			if slices.ContainsFunc(b.Excludes, eq) {
				return fs.SkipDir
			}
			return nil
		}
		if slices.ContainsFunc(b.Excludes, eq) {
			return nil
		}

		// Specify the destination path
		dPath := strings.Replace(path, b.Source, b.Destination, 1)
		mapper.AddMapping(path, dPath)
		return nil
	})

	return mapper, err
}
