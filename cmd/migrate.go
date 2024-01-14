package cmd

import (
	"fmt"

	"github.com/seb-schulz/onegate/internal/config"
	"github.com/seb-schulz/onegate/internal/database"
	"github.com/seb-schulz/onegate/internal/model"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(migrateCmd)
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate database",
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := database.Open(database.WithDebug(config.Config.DB.Debug))
		if err != nil {
			return err
		}

		if err := db.AutoMigrate(model.User{}, model.Credential{}, model.Session{}, model.AuthSession{}); err != nil {
			return fmt.Errorf("migration failed: %v", err)
		}

		// Manual migration was added because tags generated multiple indexes
		if !db.Migrator().HasIndex(&model.User{}, "idx_user_authn_id_uniq") {
			db.Exec("CREATE UNIQUE INDEX idx_user_authn_id_uniq ON users(authn_id(16))")
		}
		return nil
	},
}
