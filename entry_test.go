package donut

import (
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"

	"github.com/gleamsoda/donut/test/helper"
)

func TestNewEntry(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		want      *Entry
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "OK/Exists",
			path: "./test/testdata/dotfiles/.example",
			want: &Entry{
				Path:  "./test/testdata/dotfiles/.example",
				Empty: false,
			},
			assertion: assert.NoError,
		},
		{
			name: "OK/NotExists",
			path: "./test/testdata/dotfiles/not_exists",
			want: &Entry{
				Path:  "./test/testdata/dotfiles/not_exists",
				Empty: true,
			},
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewEntry(tt.path)
			tt.assertion(t, err)
			if err == nil {
				opts := cmp.Options{
					cmpopts.IgnoreUnexported(Entry{}),
					cmpopts.IgnoreFields(Entry{}, "Mode", "ModTime"),
				}
				if diff := cmp.Diff(tt.want, got, opts...); diff != "" {
					t.Errorf("mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestEntry_GetSum(t *testing.T) {
	entry, _ := NewEntry("./test/testdata/dotfiles/.example")
	emptyEntry, _ := NewEntry("./test/testdata/dotfiles/not_exists")
	dir := t.TempDir()
	notReadable := filepath.Join(dir, "not_readable.toml")
	helper.WriteFile(t, notReadable, []byte(`[file]`), 0200)
	notReadableEntry, _ := NewEntry(notReadable)

	tests := []struct {
		name           string
		entry          *Entry
		assertion      assert.ValueAssertionFunc
		errorAssertion assert.ErrorAssertionFunc
	}{
		{
			name:           "OK/Load",
			entry:          entry,
			assertion:      assert.NotEmpty,
			errorAssertion: assert.NoError,
		},
		{
			name:           "OK/NotExists",
			entry:          emptyEntry,
			assertion:      assert.Empty,
			errorAssertion: assert.NoError,
		},
		{
			name:           "Error/NotReadable",
			entry:          notReadableEntry,
			assertion:      assert.Empty,
			errorAssertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.entry.GetSum()
			tt.assertion(t, got)
			tt.errorAssertion(t, err)
		})
	}
}

func TestEntry_GetContent(t *testing.T) {
	entry, _ := NewEntry("./test/testdata/dotfiles/.example")
	emptyEntry, _ := NewEntry("./test/testdata/dotfiles/not_exists")
	dir := t.TempDir()
	notReadable := filepath.Join(dir, "not_readable.toml")
	helper.WriteFile(t, notReadable, []byte(`[file]`), 0200)
	notReadableEntry, _ := NewEntry(notReadable)

	tests := []struct {
		name      string
		entry     *Entry
		want      []byte
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "OK",
			entry:     entry,
			want:      []byte(`# this is example`),
			assertion: assert.NoError,
		},
		{
			name:      "OK/NotExists",
			entry:     emptyEntry,
			want:      nil,
			assertion: assert.NoError,
		},
		{
			name:      "Error/NotReadable",
			entry:     notReadableEntry,
			want:      nil,
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.entry.GetContent()
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
