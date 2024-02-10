package client

import (
	"github.com/seb-schulz/onegate/cmd"
	"github.com/spf13/cobra"
)

const (
	errRetrieveClientFormat = "cannot retrieve user: %v"
)

var (
	debug bool
)

func init() {
	cmd.RootCmd.AddCommand(clientCmd)
	clientCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Verbose output")
}

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Operate with clients",
}
