package user

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/seb-schulz/onegate/internal/config"
	"github.com/seb-schulz/onegate/internal/model"
	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	userCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all users",
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := gorm.Open(mysql.Open(config.Default.DB.Dsn), &gorm.Config{})
		if err != nil {
			return fmt.Errorf(errDatabaseConnectionFormat, err)
		}

		if debug {
			db = db.Debug()
		}

		users := []model.User{}
		if r := db.Unscoped().Find(&users); r.Error != nil {
			return fmt.Errorf(errRetrieveUserFormat, r.Error)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "ID\tName\tDisplay Name\tUpdated at\tCreated at\tDeleted at")
		for _, user := range users {
			deletedAt, _ := user.DeletedAt.Value()
			deletedAtStr := ""
			if deletedAt != nil {
				deletedAtStr = deletedAt.(time.Time).Format(time.DateOnly)
			}

			fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\n", user.ID, user.Name, user.DisplayName, user.CreatedAt.Format(time.DateOnly), user.UpdatedAt.Format(time.DateOnly), deletedAtStr)
		}
		w.Flush()
		return nil
	},
}
