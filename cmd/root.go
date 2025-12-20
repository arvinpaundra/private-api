package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use: "",
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
		HiddenDefaultCmd:  true,
	},
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func Execute() error {
	return rootCmd.Execute()
}
