package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/shinzonetwork/shinzo-view-creator/core/service"
	"github.com/shinzonetwork/viewbundle-go"
	"github.com/spf13/cobra"
)

func MakeViewSizeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "size [name]",
		Short: "Show the bundled wire size of a view",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			storeImpl := mustGetContextViewStore(cmd)
			viewName := args[0]

			view, err := storeImpl.Load(viewName)
			if err != nil {
				return err
			}

			if view.Query == nil || view.Sdl == nil {
				return fmt.Errorf("view missing query or sdl")
			}

			home, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("unable to get home directory: %w", err)
			}

			assetsDir := filepath.Join(home, ".shinzo", "views", viewName, "assets")

			vb := viewbundle.View{
				Query: *view.Query,
				Sdl:   *view.Sdl,
				Transform: viewbundle.Transform{
					Lenses: make([]viewbundle.Lens, 0, len(view.Transform.Lenses)),
				},
			}

			fmt.Println()
			for i, l := range view.Transform.Lenses {
				wasmBase64, err := storeImpl.GetAssetBlob(viewName, l.Label)
				if err != nil {
					return fmt.Errorf("lens[%d] %q: %w", i, l.Label, err)
				}

				wasmPath := filepath.Join(assetsDir, l.Label+".wasm")
				info, _ := os.Stat(wasmPath)
				if info != nil {
					fmt.Printf("  %-25s %s (raw) -> %s (base64)\n",
						l.Label+".wasm",
						formatSize(info.Size()),
						formatSize(int64(len(wasmBase64))),
					)
				}

				args := ""
				switch v := any(l.Arguments).(type) {
				case string:
					args = v
				default:
					bz, _ := json.Marshal(l.Arguments)
					args = string(bz)
				}

				vb.Transform.Lenses = append(vb.Transform.Lenses, viewbundle.Lens{
					Path:      wasmBase64,
					Arguments: args,
				})
			}

			bd := viewbundle.NewBundler()
			wireBytes, err := bd.BundleView(vb)
			if err != nil {
				return fmt.Errorf("failed to bundle: %w", err)
			}

			txData := service.EncodeRegisterBytesCalldata(wireBytes)

			fmt.Println()
			fmt.Printf("  %-25s %s\n", "Wire bundle", formatSize(int64(len(wireBytes))))
			fmt.Printf("  %-25s %s\n", "Tx calldata", formatSize(int64(len(txData))))
			fmt.Println()

			return nil
		},
	}

	return cmd
}

func formatSize(bytes int64) string {
	const (
		kb = 1024
		mb = 1024 * kb
	)

	switch {
	case bytes >= mb:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(mb))
	case bytes >= kb:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(kb))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}
