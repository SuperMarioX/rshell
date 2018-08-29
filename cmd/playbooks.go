package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(playbooksCmd)
}

var playbooksCmd = &cobra.Command{
	Use:   "playbooks",
	Short: "run playbooks.",
	Long:  `run playbooks.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("playbooks")
	},
}

