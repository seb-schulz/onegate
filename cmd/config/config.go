package config

import (
	"github.com/seb-schulz/onegate/cmd"
	"github.com/spf13/cobra"
)

func init() {
	cmd.RootCmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Verify configuration",
	RunE:  runShowCmd,
}
