package donut

import (
	"io/fs"
	"os"
	"os/exec"

	"github.com/google/renameio/v2"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type Executor struct {
	logger zerolog.Logger
}

func NewExecutor(logger zerolog.Logger) *Executor {
	return &Executor{logger: logger}
}

func (e *Executor) SetLogger(logger zerolog.Logger) {
	e.logger = logger
}

func (e *Executor) MkdirAll(path string, perm fs.FileMode) error {
	err := os.MkdirAll(path, os.ModePerm)
	e.logger.Info().Str("directory", path).Err(err).Msg("Create")
	return err
}

func (e *Executor) Remove(name string) error {
	err := os.Remove(name)
	e.logger.Info().Str("entry", name).Err(err).Msg("Remove")
	return err
}

func (e *Executor) Run(cmd *exec.Cmd) error {
	err := cmd.Run()
	e.logger.Info().Str("command", cmd.Path).Strs("args", cmd.Args[1:]).Err(err).Msg("Execute")
	return err
}

func (e *Executor) Output(cmd *exec.Cmd) ([]byte, error) {
	bytes, err := cmd.Output()
	e.logger.Info().Str("command", cmd.Path).Strs("args", cmd.Args[1:]).Err(err).Msg("Execute")
	return bytes, err
}

func (e *Executor) Overwrite(filename string, data []byte, perm fs.FileMode, opts ...renameio.Option) error {
	err := renameio.WriteFile(filename, data, perm)
	e.logger.Info().Str("entry", filename).Err(err).Msg("Write")
	return err
}

func (e *Executor) WriteConfig(v *viper.Viper, path string) error {
	err := v.SafeWriteConfigAs(path)
	e.logger.Info().Str("entry", path).Err(err).Msg("Write")
	return err
}
