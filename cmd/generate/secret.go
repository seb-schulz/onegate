package generate

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/spf13/cobra"
)

var byteLength uint16

func init() {
	generateCmd.AddCommand(secretCmd)
	secretCmd.Flags().Uint16VarP(&byteLength, "bytes", "b", 32, "byte size of generated key ")
}

var secretCmd = &cobra.Command{
	Use:     "secret",
	Aliases: []string{"s"},
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
