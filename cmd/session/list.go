package session

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/seb-schulz/onegate/internal/database"
	"github.com/seb-schulz/onegate/internal/model"
	"github.com/spf13/cobra"
)

func init() {
	sessionCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all users",
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := database.Open(database.WithDebug(debug))
		if err != nil {
			return err
		}

		sessions := []model.Session{}
		if r := db.Unscoped().Find(&sessions); r.Error != nil {
			return fmt.Errorf(errRetrieveUserFormat, r.Error)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "ID\tUser ID\tUpdated at\tCreated at\tDeleted at\tActive?")
		for _, session := range sessions {
			deletedAt, _ := session.DeletedAt.Value()
			deletedAtStr := ""
			if deletedAt != nil {
				deletedAtStr = deletedAt.(time.Time).Format(time.RFC3339)
			}

			isActive := ""
			if session.IsActive() {
				isActive = "\xE2\x9C\x94"
			}

			fmt.Fprintf(w, "%v\t%d\t%s\t%s\t%s\t%v\n", session.ID, session.UserID, session.CreatedAt.Format(time.RFC3339), session.UpdatedAt.Format(time.RFC3339), deletedAtStr, isActive)
		}
		w.Flush()
		return nil
	},
}
