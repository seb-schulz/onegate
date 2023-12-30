package user

import (
	"fmt"
	"os"
	"time"

	"github.com/seb-schulz/onegate/internal/config"
	"github.com/seb-schulz/onegate/internal/middleware"
	"github.com/seb-schulz/onegate/internal/model"
	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var expiresIn time.Duration

func init() {
	userCmd.AddCommand(loginCmd)
	loginCmd.Flags().DurationVarP(&expiresIn, "expires", "e", config.Default.UrlLogin.ExpiresIn, "Duration when URL will expire")
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Provide login URL for user recovery",
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

		url, err := middleware.GetLoginUrl(user.ID, expiresIn)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot generate URL: %v\n", err)
		}
		fmt.Printf("%v\n", url)
	},
}
