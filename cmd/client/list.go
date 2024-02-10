package client

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/seb-schulz/onegate/internal/auth"
	"github.com/seb-schulz/onegate/internal/database"
	"github.com/spf13/cobra"
)

func init() {
	clientCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all clients",
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := database.Open(database.WithDebug(debug))
		if err != nil {
			return err
		}

		clients := []auth.Client{}
		if r := db.Unscoped().Find(&clients); r.Error != nil {
			return fmt.Errorf(errRetrieveClientFormat, r.Error)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "ClientID\tDescription\tRedirectURI\tUpdated at\tCreated at\tDeleted at")
		for _, client := range clients {
			deletedAt, _ := client.DeletedAt.Value()
			deletedAtStr := ""
			if deletedAt != nil {
				deletedAtStr = deletedAt.(time.Time).Format(time.DateOnly)
			}

			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n", client.ClientID(), client.Description, client.RedirectURI(), client.CreatedAt.Format(time.DateOnly), client.UpdatedAt.Format(time.DateOnly), deletedAtStr)
		}
		w.Flush()
		return nil
	},
}
