package cmd

import (
	"log"

	"github.com/seb-schulz/onegate/internal/config"
	"github.com/seb-schulz/onegate/internal/model"
	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	RootCmd.AddCommand(migrateCmd)
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate database",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := gorm.Open(mysql.Open(config.Config.DB.Dsn), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}

		if config.Config.DB.Debug {
			db = db.Debug()
		}

		if err := db.AutoMigrate(model.User{}, model.Credential{}, model.Session{}, model.AuthSession{}); err != nil {
			log.Fatalln("Migration failed: ", err)
		}

		// Manual migration was added because tags generated multiple indexes
		if !db.Migrator().HasIndex(&model.User{}, "idx_user_authn_id_uniq") {
			db.Exec("CREATE UNIQUE INDEX idx_user_authn_id_uniq ON users(authn_id(16))")
		}
	},
}
