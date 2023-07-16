package config

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"

	"github.com/gleamsoda/donut/tutil"
)

func TestNew(t *testing.T) {
	home, _, data, _ := tutil.CreateBaseDir(t)
	tutil.SetDirEnv(t, home)
	defer SetUserHomeDir(home)()

	tests := []struct {
		name      string
		opts      []ConfigOption
		want      *Config
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "OK/WithDefault",
			opts: []ConfigOption{WithDefault()},
			want: &Config{
				Source:      data,
				Destination: home,
				Editor:      []string{"vim"},
				Pager:       []string{"less", "-R"},
				Diff:        []string{"diff", "-upN", "{{.Destination}}", "{{.Source}}"},
				Merge:       []string{"vimdiff", "{{.Destination}}", "{{.Source}}"},
			},
			assertion: assert.NoError,
		},
		{
			name: "OK/WithData",
			opts: []ConfigOption{WithData(map[string]interface{}{
				"source":      data,
				"destination": home,
				"editor":      []string{"nvim"},
				"pager":       []string{"delta"},
			})},
			want: &Config{
				Source:      data,
				Destination: home,
				Editor:      []string{"nvim"},
				Pager:       []string{"delta"},
			},
			assertion: assert.NoError,
		},
		{
			name: "OK/WithNameAndPath",
			opts: []ConfigOption{WithNameAndPath("basic", "../testdata/config")},
			want: &Config{
				Source:      data,
				Destination: home,
				Editor:      []string{"nvim"},
			},
			assertion: assert.NoError,
		},
		{
			name: "OK/WithFile",
			opts: []ConfigOption{WithFile("../testdata/config/basic.toml")},
			want: &Config{
				Source:      data,
				Destination: home,
				Editor:      []string{"nvim"},
			},
			assertion: assert.NoError,
		},
		{
			name:      "Error/WithFile/Broken",
			opts:      []ConfigOption{WithFile("../testdata/config/broken.toml")},
			want:      nil,
			assertion: assert.Error,
		},
		{
			name:      "Error/WithFile/Less",
			opts:      []ConfigOption{WithFile("../testdata/config/less.toml")},
			want:      nil,
			assertion: assert.Error,
		},
		{
			name:      "Error/WithFile/Wrong",
			opts:      []ConfigOption{WithFile("../testdata/config/wrong.toml")},
			want:      nil,
			assertion: assert.Error,
		},
		{
			name: "OK/WithPath",
			opts: WithPath("../testdata/config/basic.toml"),
			want: &Config{
				Source:      data,
				Destination: home,
				Editor:      []string{"nvim"},
				Pager:       []string{"less", "-R"},
				Diff:        []string{"diff", "-upN", "{{.Destination}}", "{{.Source}}"},
				Merge:       []string{"vimdiff", "{{.Destination}}", "{{.Source}}"},
			},
			assertion: assert.NoError,
		},
		{
			name:      "Error/WithPath/PathEmpty",
			opts:      WithPath(""),
			want:      nil,
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.opts...)
			tt.assertion(t, err)
			if err == nil {
				opts := cmp.Options{
					cmpopts.IgnoreFields(Config{}, "Concurrency", "File"),
				}
				if diff := cmp.Diff(tt.want, got, opts); diff != "" {
					t.Errorf("mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}
