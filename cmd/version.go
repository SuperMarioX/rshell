package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

const VERSION = "v0.1"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of rshell",
	Long:  `All software has versions. This is rshell's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(VERSION)
	},
}
