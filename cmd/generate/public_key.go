package generate

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/cobra"
)

func init() {
	generateCmd.AddCommand(publicKeyCmd)
}

var publicKeyCmd = &cobra.Command{
	Use:     "publicKey",
	Aliases: []string{"pubKey"},
	Short:   "Show current configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		rawPrivKey, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("failed to read key: %w", err)
		}

		privateKey, err := jwt.ParseECPrivateKeyFromPEM(rawPrivKey)
		if err != nil {
			return fmt.Errorf("failed to parse private key: %w", err)
		}

		encoded, err := x509.MarshalPKIXPublicKey(privateKey.Public())
		if err != nil {
			return fmt.Errorf("cannot marshal public key: %w", err)
		}
		pem.Encode(os.Stdout, &pem.Block{Type: "PUBLIC KEY", Bytes: encoded})

		return nil
	},
}
