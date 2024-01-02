package session

import (
	"github.com/seb-schulz/onegate/cmd"
	"github.com/spf13/cobra"
)

const (
	errDatabaseConnectionFormat = "failed to connect to database: %v"
	errRetrieveUserFormat       = "cannot retrieve user: %v"
)

var (
	debug bool
)

func init() {
	cmd.RootCmd.AddCommand(sessionCmd)
	sessionCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Verbose output")
}

var sessionCmd = &cobra.Command{
	Use:   "session",
	Short: "Operate with sessions",
}
