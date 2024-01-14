package session

import (
	"fmt"
	"time"

	"github.com/seb-schulz/onegate/internal/config"
	"github.com/seb-schulz/onegate/internal/database"
	"github.com/seb-schulz/onegate/internal/model"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var (
	softDelete bool
	dryRun     bool
	inactive   bool
	deleted    bool
)

func init() {
	sessionCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolVar(&softDelete, "soft-delete", false, "Only soft-delete entries")
	deleteCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print SQL statement instead of executing it")
	deleteCmd.Flags().BoolVar(&inactive, "inactive", false, fmt.Sprintf("Delete all sessions older than %v", config.Config.Session.ActiveFor))
	deleteCmd.Flags().BoolVar(&deleted, "deleted", false, "Delete all soft-deleted entriess")
}

var deleteCmd = &cobra.Command{
	Use:     "delete-all",
	Aliases: []string{"rma"},
	Short:   "(Soft-)delete user",
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := database.Open(database.WithDebug(debug))
		if err != nil {
			return err
		}

		tx := db.Session(&gorm.Session{DryRun: dryRun, AllowGlobalUpdate: true})

		if softDelete && deleted {
			return fmt.Errorf("soft-delete and delete flags are mutually exclusive")
		} else if softDelete && !deleted {
			// continue with current scope
		} else if !softDelete && deleted {
			tx = tx.Unscoped().Where("deleted_at IS NOT NULL")
		} else if !softDelete && !deleted {
			tx = tx.Unscoped()
		} else {
			panic("all cases should be covered")
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
