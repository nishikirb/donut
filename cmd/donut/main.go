package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/gleamsoda/donut"
)

func main() {
	if err := NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(func() {
		if err := donut.InitStore(); err != nil {
			panic(err)
		}
	})
	cobra.OnFinalize(func() {
		if err := donut.GetStore().Close(); err != nil {
			panic(err)
		}
	})
}

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "donut",
		Version:      donut.GetVersion(),
		Short:        "Tiny dotfiles management tool written in Go.",
		SilenceUsage: true,
	}

	cmd.PersistentFlags().StringP("file", "f", "", "Specify the path to the configuration file")

	cmd.AddCommand(
		NewInitCmd(),
		NewListCmd(),
		NewDiffCmd(),
		NewWhereCmd(),
		NewConfigCmd(),
		NewApplyCmd(),
	)

	return cmd
}

func NewInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Create a default configuration file",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			d, _ := donut.New()
			return d.Init()
		},
	}

	return cmd
}

func NewListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Display a list of source files",
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			f, _ := cmd.Flags().GetString("file")
			return donut.InitConfig(f)
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			d, err := donut.New(donut.WithConfig(donut.GetConfig()))
			if err != nil {
				return err
			}
			return d.List()
		},
	}

	return cmd
}

func NewDiffCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "diff",
		Short: "Display a list of differences between source and destination files",
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			f, _ := cmd.Flags().GetString("file")
			return donut.InitConfig(f)
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			d, err := donut.New(donut.WithConfig(donut.GetConfig()), donut.WithStore(donut.GetStore()))
			if err != nil {
				return err
			}
			return d.Diff()
		},
	}

	return cmd
}

func NewWhereCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "where",
		Short: "Display the location of the source or destination directory",
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			f, _ := cmd.Flags().GetString("file")
			return donut.InitConfig(f)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := donut.New(donut.WithConfig(donut.GetConfig()))
			if err != nil {
				return err
			}
			return d.Where(args[0])
		},
		Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		ValidArgs: []string{"source", "destination", "config"},
	}

	return cmd
}

func NewConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Display the contents of the configuration file",
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			f, _ := cmd.Flags().GetString("file")
			return donut.InitConfig(f)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := donut.New(donut.WithConfig(donut.GetConfig()))
			if err != nil {
				return err
			}
			return d.ConfigShow()
		},
	}

	cmd.AddCommand(NewConfigEditCmd())
	return cmd
}

func NewConfigEditCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit",
		Short: "Edit the configuration file",
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			f, _ := cmd.Flags().GetString("file")
			return donut.InitConfig(f)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := donut.New(donut.WithConfig(donut.GetConfig()))
			if err != nil {
				return err
			}
			return d.ConfigEdit()
		},
	}

	return cmd
}

func NewApplyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apply",
		Short: "Apply the content of the source file to the destination file",
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			f, _ := cmd.Flags().GetString("file")
			return donut.InitConfig(f)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := donut.New(donut.WithConfig(donut.GetConfig()), donut.WithStore(donut.GetStore()))
			if err != nil {
				return err
			}
			return d.Apply()
		},
	}

	return cmd
}
