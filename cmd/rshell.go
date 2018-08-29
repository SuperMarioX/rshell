package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "rshell",
	Short: "rshell is a remote shell exec application.",
	Long: `A simple application for exec remote shell command on linux from win/linux, http://github.com/luckywinds/rshell`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("......")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
