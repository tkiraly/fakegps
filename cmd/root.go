package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var lorawan, regionversion, region string

func init() {
}

var rootCmd = &cobra.Command{
	Use:           "fakegps",
	Short:         "fakegps simulates a gps receiver",
	SilenceErrors: false,
	SilenceUsage:  true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Error(err)
	}
}
