package cmd

import (
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/seb-schulz/onegate/internal/auth"
	"github.com/seb-schulz/onegate/internal/config"
	"github.com/seb-schulz/onegate/internal/server"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(serveCmd)
}

func runServeCmd(cmd *cobra.Command, args []string) error {
	c := server.ServerConfig{
		Router: server.RouterConfig{
			DbDebug: config.Config.DB.Debug,
			Webauthn: webauthn.Config{
				RPDisplayName: config.Config.RelyingParty.Name,
				RPID:          config.Config.RelyingParty.ID,
				RPOrigins:     config.Config.RelyingParty.Origins,
			},
			Limit: server.RouterLimitConfig{
				RequestLimit: config.Config.Server.Limit.RequestLimit, WindowLength: config.Config.Server.Limit.WindowLength,
			},
			SessionKey:              []byte(config.Config.Session.Key),
			UserRegistrationEnabled: config.Config.Features.UserRegistration,
			Login: server.LoginConfig{
				Key:          config.Config.UrlLogin.Key,
				ValidMethods: config.Config.UrlLogin.ValidMethods,
				BaseUrl:      *config.Config.BaseUrl.JoinPath("login"),
			},
			Auth: auth.Config{},
		},
		HttpPort:  config.Config.Server.HttpPort,
		ServeType: server.ServeType(config.Config.Server.Kind),
	}

	return server.Serve(&c)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run server",
	RunE:  runServeCmd,
}
