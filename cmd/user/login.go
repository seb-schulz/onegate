package user

import (
	"fmt"
	"os"
	"time"

	"github.com/mdp/qrterminal/v3"
	"github.com/seb-schulz/onegate/internal/config"
	"github.com/seb-schulz/onegate/internal/database"
	"github.com/seb-schulz/onegate/internal/middleware"
	"github.com/seb-schulz/onegate/internal/model"
	"github.com/spf13/cobra"
)

var (
	expiresIn time.Duration
	qrCode    bool
)

func init() {
	userCmd.AddCommand(loginCmd)
	loginCmd.Flags().DurationVarP(&expiresIn, "expires", "e", config.Config.UrlLogin.ExpiresIn, "Duration when URL will expire")
	loginCmd.Flags().BoolVar(&qrCode, "qr", false, "Output link as QR code")
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Provide login URL for user recovery",
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

		url, err := middleware.GetLoginUrl(user.ID, expiresIn)
		if err != nil {
			return fmt.Errorf("cannot generate URL: %v", err)
		}

		if qrCode {
			qrterminal.GenerateHalfBlock(url.String(), qrterminal.L, os.Stdout)
		} else {
			fmt.Printf("%v\n", url)
		}
		return nil
	},
}
