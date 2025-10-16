package cmd

import (
	"fmt"

	"github.com/bernd/desk/internal/config"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("◲  desk %s\n", config.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
