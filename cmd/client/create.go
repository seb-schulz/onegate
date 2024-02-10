package client

import (
	"context"
	"fmt"

	"github.com/seb-schulz/onegate/internal/auth"
	"github.com/seb-schulz/onegate/internal/database"
	"github.com/spf13/cobra"
)

var (
	description string
	redirectURI string
)

func init() {
	clientCmd.AddCommand(createCmd)
	createCmd.Flags().StringVar(&description, "desc", "", "Summery about purpose of this client")
	createCmd.Flags().StringVarP(&redirectURI, "redirect-url", "u", "", "URI callback used by oAuth2/OIDC")
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create client ID and client secret",
	RunE: func(cmd *cobra.Command, args []string) error {
		if description == "" {
			return fmt.Errorf("client must contain a description")
		}
		if redirectURI == "" {
			return fmt.Errorf("client must contain a redirect URI")
		}

		db, err := database.Open(database.WithDebug(debug))
		if err != nil {
			return err
		}

		clientID, clientSecret, err := auth.CreateClient(database.WithContext(context.Background(), db), auth.NewClientSecretHasher(), description, redirectURI)
		if err != nil {
			return fmt.Errorf("cannot create client: %v", err)
		}

		fmt.Printf("Client ID: %s\n", clientID)
		fmt.Printf("Client Secret: %s\n", clientSecret)
		fmt.Println("Take a not of client' secret because it is only visible once.")

		return nil
	},
}
