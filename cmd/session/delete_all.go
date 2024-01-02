package session

import (
	"fmt"
	"time"

	"github.com/seb-schulz/onegate/internal/config"
	"github.com/seb-schulz/onegate/internal/model"
	"github.com/seb-schulz/onegate/internal/utils"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var (
	softDelete  bool
	dryRun      bool
	inactive    bool
	withoutUser bool
)

func init() {
	sessionCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolVar(&softDelete, "soft-delete", false, "Only soft-delete entries")
	deleteCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print SQL statement instead of executing it")
	deleteCmd.Flags().BoolVar(&inactive, "inactive", false, fmt.Sprintf("Delete all sessions older than %v", config.Config.Session.ActiveFor))
	deleteCmd.Flags().BoolVar(&withoutUser, "without-user", false, "exclude session with logged in users")
}

var deleteCmd = &cobra.Command{
	Use:     "delete-all",
	Aliases: []string{"rma"},
	Short:   "(Soft-)delete user",
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := utils.OpenDatabase(utils.WithDebugOption(debug))
		if err != nil {
			return err
		}

		tx := db.Session(&gorm.Session{DryRun: dryRun, AllowGlobalUpdate: true})

		if !softDelete {
			tx = tx.Unscoped()
		}

		if withoutUser {
			tx.Where("user_id IS NULL")
		}

		if inactive {
			tx = tx.Where("updated_at <= ?", time.Now().Add(-config.Config.Session.ActiveFor))
		}

		tx.Delete(&model.Session{})

		if dryRun {
			sql := tx.ToSQL(func(tx *gorm.DB) *gorm.DB {
				return tx
			})
			fmt.Println(sql)
		}

		return nil
	},
}
