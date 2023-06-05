package donut

import (
	"io/fs"
	"path/filepath"
	"strings"

	"golang.org/x/exp/slices"
)

type PathMapperBuilder struct {
	source      string
	destination string
	excludes    []string
}

type PathMapperBuilderOption func(b *PathMapperBuilder)

var mustExcludes = []string{".git"}

func NewPathMapperBuilder(s, d string, funcs ...PathMapperBuilderOption) *PathMapperBuilder {
	b := &PathMapperBuilder{
		source:      s,
		destination: d,
		excludes:    mustExcludes,
	}

	for _, fn := range funcs {
		fn(b)
	}

	return b
}

func WithExcludes(s ...string) PathMapperBuilderOption {
	return func(b *PathMapperBuilder) {
		b.excludes = append(b.excludes, s...)
	}
}

func (b *PathMapperBuilder) Build() (*PathMapper, error) {
	mapper := NewPathMapper(b.source, b.destination)

	err := filepath.WalkDir(b.source, func(path string, d fs.DirEntry, _ error) error {
		rel, _ := filepath.Rel(b.source, path)
		eq := func(s string) bool {
			ok, _ := filepath.Match(s, rel)
			return ok
		}
		if d.IsDir() {
			if slices.ContainsFunc(b.excludes, eq) {
				return fs.SkipDir
			}
			return nil
		}
		if slices.ContainsFunc(b.excludes, eq) {
			return nil
		}

		// Specify the destination path
		dPath := strings.Replace(path, b.source, b.destination, 1)
		mapper.AddMapping(path, dPath)
		return nil
	})

	return mapper, err
}
