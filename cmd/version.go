package cmd

import (
	"github.com/spf13/cobra"
)

var version string

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version number",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(version) == 0 {
			version = "DEV"
		}
		cmd.Printf("%s\n", version)
		return nil
	},
}
