package user

import (
	"github.com/seb-schulz/onegate/cmd"
	"github.com/spf13/cobra"
)

const (
	errRetrieveUserFormat = "cannot retrieve user: %v"
)

var (
	debug bool
)

func init() {
	cmd.RootCmd.AddCommand(userCmd)
	userCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Verbose output")
}

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Operate with users",
}
