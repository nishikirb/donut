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
		Use:     "donut",
		Version: donut.GetVersion(),
		Short:   "A brief description of your application",
		Long: `A longer description that spans multiple lines and likely contains
	examples and usage of using your application. For example:
	
	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
		SilenceUsage: true,
	}

	cmd.PersistentFlags().StringP("file", "f", "", "location of config file")

	cmd.AddCommand(
		NewEchoCmd(),
		NewInitCmd(),
		NewListCmd(),
		NewDiffCmd(),
		NewEditCmd(),
	)

	return cmd
}

func NewEchoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "echo",
		Short: "A brief description of your command",
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			configFile, _ := cmd.Flags().GetString("file")
			return donut.InitConfig(configFile)
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			d, err := donut.New(donut.WithConfig(donut.GetConfig()))
			if err != nil {
				return err
			}
			return d.Echo()
		},
	}

	return cmd
}

func NewInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "A brief description of your command",
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
		Short: "A brief description of your command",
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			configFile, _ := cmd.Flags().GetString("file")
			return donut.InitConfig(configFile)
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
		Short: "A brief description of your command",
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			configFile, _ := cmd.Flags().GetString("file")
			return donut.InitConfig(configFile)
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
		Short: "A brief description of your command",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			configFile, _ := cmd.Flags().GetString("file")
			return donut.InitConfig(configFile)
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
