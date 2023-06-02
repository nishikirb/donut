package donut

import "path/filepath"

type PathMapper struct {
	Mapping     []PathMapping
	Source      string
	Destination string
}

type PathMapping struct {
	Source      string
	Destination string
}

func NewPathMapper(src, dst string) *PathMapper {
	return &PathMapper{
		Source:      src,
		Destination: dst,
	}
}

func (m *PathMapper) AddMapping(src, dst string) {
	m.Mapping = append(m.Mapping, PathMapping{
		Source:      src,
		Destination: dst,
	})
}

func (m *PathMapper) RelSourcePaths() []string {
	var paths []string
	for _, v := range m.Mapping {
		rel, _ := filepath.Rel(m.Source, v.Source)
		paths = append(paths, rel)
	}
	return paths
}
