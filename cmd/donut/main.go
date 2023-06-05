package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/gleamsoda/donut"
)

var file string
var debug bool

func main() {
	app := donut.NewApp()
	root := NewCmdRoot(app)

	root.AddCommand(
		NewCmdInit(app),
		NewCmdList(app),
		NewCmdDiff(app),
		NewCmdMerge(app),
		NewCmdWhere(app),
		NewCmdConfig(app),
		NewCmdApply(app),
	)

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

func NewCmdRoot(app *donut.App, _ ...donut.Option) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "donut",
		Version:      donut.GetVersion(),
		Short:        "Tiny dotfiles management tool",
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if store, err := donut.OpenStore(); err != nil {
				return err
			} else {
				cobra.OnFinalize(func() { store.Close() })
				app.AddOptions(donut.WithLogger(donut.NewLogger(os.Stdout, debug)), donut.WithStore(store))
			}
			return nil
		},
	}

	cmd.PersistentFlags().StringVarP(&file, "file", "f", "", "Specify the configuration file")
	cmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Enable debug mode")

	return cmd
}

func NewCmdInit(app *donut.App) *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Create a default configuration file",
		Args:  cobra.NoArgs,
		RunE:  run(app),
	}
}

func NewCmdList(app *donut.App) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Display a list of source files",
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			app.AddOptions(donut.WithConfigLoader(donut.WithPath(file)...))
			return nil
		},
		RunE: run(app),
	}
}

func NewCmdDiff(app *donut.App) *cobra.Command {
	return &cobra.Command{
		Use:   "diff",
		Short: "Display a list of differences between source and destination files",
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			app.AddOptions(donut.WithConfigLoader(donut.WithPath(file)...))
			return nil
		},
		RunE: run(app),
	}
}

func NewCmdMerge(app *donut.App) *cobra.Command {
	return &cobra.Command{
		Use:   "merge",
		Short: "Merge the source file into the destination file",
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			app.AddOptions(donut.WithConfigLoader(donut.WithPath(file)...))
			return nil
		},
		RunE: run(app),
	}
}

func NewCmdWhere(app *donut.App) *cobra.Command {
	return &cobra.Command{
		Use:       "where",
		Short:     "Display the location of the source or destination directory",
		Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		ValidArgs: []string{"source", "destination", "config"},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			app.AddOptions(donut.WithConfigLoader(donut.WithPath(file)...))
			return nil
		},
		RunE: run(app),
	}
}

func NewCmdConfig(app *donut.App) *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Edit the configuration file",
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			app.AddOptions(donut.WithConfigLoader(donut.WithPath(file)...))
			return nil
		},
		RunE: run(app),
	}
}

func NewCmdApply(app *donut.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apply",
		Short: "Apply the content of the source file to the destination file",
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			app.AddOptions(donut.WithConfigLoader(donut.WithPath(file)...))
			return nil
		},
		RunE: run(app),
	}

	cmd.Flags().BoolP("overwrite", "o", false, "Overwrite the destination file with the source file")

	return cmd
}

func run(app *donut.App) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		return app.Run(cmd.Context(), cmd.Name(), args, cmd.Flags())
	}
}

// func chain(funcs ...func(cmd *cobra.Command, args []string) error) func(cmd *cobra.Command, args []string) error {
// 	return func(cmd *cobra.Command, args []string) error {
// 		for _, f := range funcs {
// 			if err := f(cmd, args); err != nil {
// 				return err
// 			}
// 		}
// 		return nil
// 	}
// }
