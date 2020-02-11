package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version of muikku",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("muikku version 0.0.1")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
