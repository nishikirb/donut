package donut

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/exp/slices"
)

type RelationsBuilder struct {
	Source      string
	Destination string
	Excludes    []string
}

type RelationsBuilderOption func(b *RelationsBuilder)

var MustExcludes = []string{".git"}

func NewRelationsBuilder(source, destination string, funcs ...RelationsBuilderOption) *RelationsBuilder {
	b := &RelationsBuilder{
		Source:      source,
		Destination: destination,
		Excludes:    MustExcludes,
	}

	for _, fn := range funcs {
		fn(b)
	}

	return b
}

func WithExcludes(s ...string) RelationsBuilderOption {
	return func(b *RelationsBuilder) {
		b.Excludes = append(b.Excludes, s...)
	}
}

func (b *RelationsBuilder) Build() ([]Relation, error) {
	var rels []Relation
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

		// define destination path
		dPath := strings.Replace(path, b.Source, b.Destination, 1)

		m, err := NewRelation(path, dPath)
		if err != nil {
			return err
		}
		rels = append(rels, *m)
		return nil
	})

	return rels, err
}

type Relation struct {
	Source      File
	Destination File
}

func NewRelation(src, dest string) (*Relation, error) {
	var sFile, dFile *File
	var err error
	if sFile, err = NewFile(src); err != nil {
		return nil, err
	}
	if dFile, err = NewFile(dest); err != nil {
		return nil, err
	}

	return &Relation{
		Source:      *sFile,
		Destination: *dFile,
	}, nil
}

type File struct {
	Path     string
	NotExist bool
	FileInfo fs.FileInfo
}

func NewFile(path string) (*File, error) {
	var notExist bool
	f, err := os.Lstat(path)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		notExist = true
	}
	return &File{
		Path:     path,
		NotExist: notExist,
		FileInfo: f,
	}, nil
}
