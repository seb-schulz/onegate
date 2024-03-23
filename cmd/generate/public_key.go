package generate

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"os"

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

		var block *pem.Block
		if block, _ = pem.Decode(rawPrivKey); block == nil {
			return fmt.Errorf("cannot parse private key")
		}

		privateKey, err := x509.ParseECPrivateKey(block.Bytes)
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
