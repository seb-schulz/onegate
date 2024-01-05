package cmd

import (
	"github.com/seb-schulz/onegate/internal/server"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(serveCmd)
}

func runServeCmd(cmd *cobra.Command, args []string) error {
	return server.Serve()
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run server",
	RunE:  runServeCmd,
}
