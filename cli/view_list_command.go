package cli

import (
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/spf13/cobra"
)

func MakeViewListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all local views",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			viewstore := mustGetContextViewStore(cmd)
			views, err := viewstore.List()
			if err != nil {
				return err
			}

			if len(views) == 0 {
				cmd.Println("No views found.")
				return nil
			}

			home, _ := os.UserHomeDir()

			type entry struct {
				name    string
				modTime time.Time
			}
			entries := make([]entry, 0, len(views))
			for _, v := range views {
				viewJSON := filepath.Join(home, ".shinzo", "views", v.Name, "view.json")
				var modTime time.Time
				if info, err := os.Stat(viewJSON); err == nil {
					modTime = info.ModTime()
				}
				entries = append(entries, entry{v.Name, modTime})
			}
			sort.Slice(entries, func(i, j int) bool {
				return entries[i].modTime.After(entries[j].modTime)
			})
			for _, e := range entries {
				date := "-"
				if !e.modTime.IsZero() {
					date = e.modTime.Format(time.DateOnly)
				}
				cmd.Printf("📄 %-24s  %s\n", e.name, date)
			}
			return nil
		},
	}
}
