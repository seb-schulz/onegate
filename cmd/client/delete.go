package client

import (
	"fmt"

	"github.com/seb-schulz/onegate/internal/auth"
	"github.com/seb-schulz/onegate/internal/database"
	"github.com/spf13/cobra"
)

var force bool

func init() {
	clientCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolVarP(&force, "force", "f", false, "Remove entry without consideration of soft-deletion flag")
}

var deleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Soft-delete user",
	Aliases: []string{"del", "rm"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := database.Open(database.WithDebug(debug))
		if err != nil {
			return err
		}

		tx := db
		if force {
			tx = tx.Unscoped()
		}

		client := auth.Client{}
		if r := tx.Where("id = ?", args[0]).First(&client); r.Error != nil {
			return fmt.Errorf(errRetrieveClientFormat, r.Error)
		}
		tx.Delete(&client)

		return nil
	},
}
