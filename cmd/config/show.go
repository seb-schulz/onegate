package config

import (
	"fmt"

	"github.com/seb-schulz/onegate/internal/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func init() {
	configCmd.AddCommand(showCmd)
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	RunE:  runShowCmd,
}

func runShowCmd(cmd *cobra.Command, args []string) error {
	cStr, err := yaml.Marshal(config.Config)
	if err != nil {
		return fmt.Errorf("unable to marshal config to YAML: %v", err)
	}
	fmt.Print(string(cStr))
	return nil
}
