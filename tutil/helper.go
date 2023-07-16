package tutil

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/pelletier/go-toml/v2"
)

func ReadToml(t *testing.T, path string) map[string]interface{} {
	t.Helper()

	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	data := map[string]interface{}{}
	if err := toml.Unmarshal(b, &data); err != nil {
		t.Fatal(err)
	}
	return data
}

func WriteToml(t *testing.T, path string, data map[string]interface{}) func() {
	t.Helper()

	b, err := toml.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	WriteFile(t, path, b, os.ModePerm)
	return func() { _ = os.Remove(path) }
}

func WriteFile(t *testing.T, path string, data []byte, perm fs.FileMode) {
	t.Helper()

	if err := os.WriteFile(path, data, perm); err != nil {
		t.Fatal(err)
	}
}

func CopyFile(t *testing.T, src, dest string) func() {
	t.Helper()

	srcFile, err := os.Open(src)
	if err != nil {
		t.Fatal(err)
	}
	defer srcFile.Close()
	destFile, err := os.Create(dest)
	if err != nil {
		t.Fatal(err)
	}
	defer destFile.Close()

	if _, err = io.Copy(destFile, srcFile); err != nil {
		t.Fatal(err)
	}
	return func() { _ = os.Remove(dest) }
}

func SetDirEnv(t *testing.T, home string) {
	t.Helper()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".config"))
}

func CreateBaseDir(t *testing.T) (home, config, data, state string) {
	t.Helper()

	home = t.TempDir()
	config = filepath.Join(home, ".config", "donut")
	if err := os.MkdirAll(config, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	data = filepath.Join(home, ".local", "share", "donut")
	if err := os.MkdirAll(data, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	state = filepath.Join(home, ".local", "state", "donut")
	if err := os.MkdirAll(state, os.ModePerm); err != nil {
		t.Fatal(err)
	}

	return
}

func CreateDirs(t *testing.T, dirs ...string) {
	t.Helper()

	for _, d := range dirs {
		if err := os.MkdirAll(d, os.ModePerm); err != nil {
			t.Fatal(err)
		}
	}
}
