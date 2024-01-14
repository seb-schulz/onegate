package user

import (
	"fmt"

	"github.com/seb-schulz/onegate/internal/database"
	"github.com/seb-schulz/onegate/internal/model"
	"github.com/spf13/cobra"
)

func init() {
	userCmd.AddCommand(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Soft-delete user",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := database.Open(database.WithDebug(debug))
		if err != nil {
			return err
		}

		user := model.User{}
		if r := db.Where("id = ?", args[0]).First(&user); r.Error != nil {
			return fmt.Errorf(errRetrieveUserFormat, r.Error)
		}
		db.Delete(&user)
		return nil
	},
}
