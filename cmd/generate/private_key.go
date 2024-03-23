package generate

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	generateCmd.AddCommand(privateKeyCmd)
}

var privateKeyCmd = &cobra.Command{
	Use:     "privateKey",
	Aliases: []string{"privKey"},
	Short:   "Show current configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return err
		}

		encoded, err := x509.MarshalECPrivateKey(privKey)
		if err != nil {
			return err
		}
		pem.Encode(os.Stdout, &pem.Block{Type: "EC PRIVATE KEY", Bytes: encoded})

		return nil
	},
}
