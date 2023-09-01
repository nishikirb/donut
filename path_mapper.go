package donut

import (
	"io/fs"
	"path/filepath"
	"slices"
	"strings"
)

type PathMapper struct {
	Mapping     []PathMapping
	source      string
	destination string
	excludes    []string
}

type PathMapping struct {
	Source      string
	Destination string
}

type PathMapperOption func(m *PathMapper)

var defaultExcludes = []string{".git"}

func NewPathMapper(s, d string, funcs ...PathMapperOption) (*PathMapper, error) {
	m := &PathMapper{
		source:      s,
		destination: d,
		excludes:    defaultExcludes,
	}

	for _, fn := range funcs {
		fn(m)
	}

	err := filepath.WalkDir(m.source, func(path string, d fs.DirEntry, _ error) error {
		rel, _ := filepath.Rel(m.source, path)
		eq := func(s string) bool {
			ok, _ := filepath.Match(s, rel)
			return ok
		}
		if d.IsDir() {
			if slices.ContainsFunc(m.excludes, eq) {
				return fs.SkipDir
			}
			return nil
		}
		if slices.ContainsFunc(m.excludes, eq) {
			return nil
		}

		// Specify the destination path
		dPath := strings.Replace(path, m.source, m.destination, 1)
		m.addMapping(path, dPath)
		return nil
	})

	return m, err
}

func WithExcludes(s ...string) PathMapperOption {
	return func(m *PathMapper) {
		m.excludes = append(m.excludes, s...)
	}
}

func (m *PathMapper) RelSourcePaths() []string {
	var paths []string
	for _, v := range m.Mapping {
		rel, _ := filepath.Rel(m.source, v.Source)
		paths = append(paths, rel)
	}
	return paths
}

func (m *PathMapper) addMapping(src, dst string) {
	m.Mapping = append(m.Mapping, PathMapping{
		Source:      src,
		Destination: dst,
	})
}
