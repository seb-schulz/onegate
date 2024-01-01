package config

import (
	"errors"
	"fmt"

	"github.com/seb-schulz/onegate/internal/config"
	"github.com/spf13/cobra"
)

func init() {
	configCmd.AddCommand(checkCmd)
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Validate current configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error

		for k, v := range map[string]bool{"relyingParty.name": config.Config.RelyingParty.Name == "", "db.dsn": config.Config.DB.Dsn == "", "session.key": config.Config.Session.Key == "", "urlLogin.key": len(config.Config.UrlLogin.Key) == 0} {
			if v {
				err = errors.Join(err, fmt.Errorf("missing value for %v", k))
			}
		}

		return err
	},
}
