package user

import (
	"fmt"
	"os"

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
	Run: func(cmd *cobra.Command, args []string) {
		db, err := gorm.Open(mysql.Open(config.Default.DB.Dsn), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}

		if debug {
			db = db.Debug()
		}

		user := model.User{}
		if r := db.Where("id = ?", args[0]).First(&user); r.Error != nil {
			fmt.Fprintf(os.Stderr, "Cannot retrieve user: %v\n", r.Error)
			os.Exit(1)
			return
		}
		db.Delete(&user)
	},
}
