package cmd

import (
	getspacekeysctrl "atlas-cli/internal/getspacekeys/controller"
	uploadwlsctrl "atlas-cli/internal/uploadwls/controller"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	CompletionOptions: cobra.CompletionOptions{
		HiddenDefaultCmd:  true,
		DisableDefaultCmd: true,
	},
}

func init() {
	rootCmd.AddCommand(uploadwlsctrl.Cmd)
	rootCmd.AddCommand(getspacekeysctrl.Cmd)
}

func Execute() error {
	return rootCmd.Execute()
}
