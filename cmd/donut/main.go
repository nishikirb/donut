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
	rootCmd := &cobra.Command{
		Use:     "donut",
		Version: donut.GetVersion(),
		Short:   "A brief description of your application",
		Long: `A longer description that spans multiple lines and likely contains
	examples and usage of using your application. For example:
	
	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
		SilenceUsage: true,

		// Uncomment the following line if your bare application
		// has an action associated with it:
		// Run: func(cmd *cobra.Command, args []string) { },
	}

	rootCmd.PersistentFlags().StringP("config", "c", "", "location of config file")

	rootCmd.AddCommand(NewEchoCmd())

	return rootCmd
}

func NewEchoCmd() *cobra.Command {
	echoCmd := &cobra.Command{
		Use:   "echo",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:
	
	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			cfgPath, _ := cmd.Flags().GetString("config")
			return donut.InitConfig(cfgPath)
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			d, err := donut.New(donut.WithConfig(donut.GetConfig()))
			if err != nil {
				return err
			}
			return d.Echo()
		},
	}

	return echoCmd
}
