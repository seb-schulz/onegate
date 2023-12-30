package user

import (
	"fmt"

	"github.com/seb-schulz/onegate/internal/config"
	"github.com/seb-schulz/onegate/internal/model"
	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	userCmd.AddCommand(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Soft-delete user",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := gorm.Open(mysql.Open(config.Default.DB.Dsn), &gorm.Config{})
		if err != nil {
			return fmt.Errorf(errDatabaseConnectionFormat, err)
		}

		if debug {
			db = db.Debug()
		}

		user := model.User{}
		if r := db.Where("id = ?", args[0]).First(&user); r.Error != nil {
			return fmt.Errorf(errRetrieveUserFormat, r.Error)
		}
		db.Delete(&user)
		return nil
	},
}
