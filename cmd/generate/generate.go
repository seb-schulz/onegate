package generate

import (
	"github.com/seb-schulz/onegate/cmd"
	"github.com/spf13/cobra"
)

var (
	debug bool
)

func init() {
	cmd.RootCmd.AddCommand(generateCmd)
	generateCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Verbose output")
}

var generateCmd = &cobra.Command{
	Use:     "generate",
	Short:   "Generate random data to setup onegate",
	Aliases: []string{"gen"},
}
