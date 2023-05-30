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
		NewEditCmd(),
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
			d, err := donut.New(donut.WithConfig(donut.GetConfig()))
			if err != nil {
				return err
			}
			return d.Diff()
		},
	}

	return cmd
}

func NewEditCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit",
		Short: "Edit the source file",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			f, _ := cmd.Flags().GetString("file")
			return donut.InitConfig(f)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			current, err := os.Getwd()
			if err != nil {
				return err
			}
			d, err := donut.New(donut.WithConfig(donut.GetConfig()))
			if err != nil {
				return err
			}
			return d.Edit(args[0], current)
		},
	}

	return cmd
}

func NewConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
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
			return d.EditConfig()
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
			d, err := donut.New(donut.WithConfig(donut.GetConfig()))
			if err != nil {
				return err
			}
			return d.Apply()
		},
	}

	return cmd
}
