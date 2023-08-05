package system

import (
	"io/fs"
	"os"
	"os/exec"

	"github.com/google/renameio/v2"
	"github.com/spf13/viper"

	"github.com/gleamsoda/donut/logger"
)

func MkdirAll(path string, perm fs.FileMode) error {
	err := os.MkdirAll(path, os.ModePerm)
	logger.Info().Str("directory", path).Err(err).Msg("Create")
	return err
}

func Stat(name string) (fs.FileInfo, error) {
	info, err := os.Stat(name)
	logger.Info().Str("entry", name).Err(err).Msg("Stat")
	return info, err
}

func Lstat(name string) (fs.FileInfo, error) {
	info, err := os.Lstat(name)
	logger.Info().Str("entry", name).Err(err).Msg("Lstat")
	return info, err
}

func Open(name string) (*os.File, error) {
	file, err := os.Open(name)
	logger.Info().Str("entry", name).Err(err).Msg("Open")
	return file, err
}

func ReadFile(name string) ([]byte, error) {
	bytes, err := os.ReadFile(name)
	logger.Info().Str("entry", name).Err(err).Msg("Read")
	return bytes, err
}

func Remove(name string) error {
	err := os.Remove(name)
	logger.Info().Str("entry", name).Err(err).Msg("Remove")
	return err
}

func Run(cmd *exec.Cmd) error {
	err := cmd.Run()
	logger.Info().Str("command", cmd.Path).Strs("args", cmd.Args[1:]).Err(err).Msg("Execute")
	return err
}

func Output(cmd *exec.Cmd) ([]byte, error) {
	bytes, err := cmd.Output()
	logger.Info().Str("command", cmd.Path).Strs("args", cmd.Args[1:]).Err(err).Msg("Execute")
	return bytes, err
}

func Overwrite(filename string, data []byte, perm fs.FileMode, opts ...renameio.Option) error {
	err := renameio.WriteFile(filename, data, perm)
	logger.Info().Str("entry", filename).Err(err).Msg("Write")
	return err
}

func WriteConfig(v *viper.Viper, path string) error {
	err := v.SafeWriteConfigAs(path)
	logger.Info().Str("entry", path).Err(err).Msg("Write")
	return err
}
