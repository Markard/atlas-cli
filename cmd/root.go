package cmd

import (
	"atlas-cli/internal/uploadwls/controller"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	CompletionOptions: cobra.CompletionOptions{
		HiddenDefaultCmd:  true,
		DisableDefaultCmd: true,
	},
}

func init() {
	rootCmd.AddCommand(controller.Cmd)
}

func Execute() error {
	return rootCmd.Execute()
}
