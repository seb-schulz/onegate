package config

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/spf13/cobra"
)

var byteLength uint16

func init() {
	configCmd.AddCommand(randomKeyCmd)
	randomKeyCmd.Flags().Uint16VarP(&byteLength, "bytes", "b", 32, "byte size of generated key ")
}

var randomKeyCmd = &cobra.Command{
	Use:     "randomKey",
	Aliases: []string{"randKey", "rk"},
	Short:   "Show current configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		r := make([]byte, byteLength)
		_, err := rand.Read(r)
		if err != nil {
			return err
		}

		fmt.Println(base64.StdEncoding.EncodeToString(r))
		return nil
	},
}
