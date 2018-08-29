package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(cliCmd)
}

var cliCmd = &cobra.Command{
	Use:   "cli",
	Short: "enter cli interaction mode.",
	Long:  `The cli interaction mode.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("cli")
	},
}

