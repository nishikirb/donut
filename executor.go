package donut

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
