package cli

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/shinzonetwork/shinzo-view-creator/core/models"
	schemastore "github.com/shinzonetwork/shinzo-view-creator/core/schema/store"
	"github.com/shinzonetwork/shinzo-view-creator/core/schema/store/fileschema"
	viewstore "github.com/shinzonetwork/shinzo-view-creator/core/view/store"
	"github.com/shinzonetwork/shinzo-view-creator/core/view/store/local"
	"github.com/spf13/cobra"
)

type contextKey string

var (
	viewStoreContextKey   = contextKey("view")
	schemaStoreContextKey = contextKey("schema")
)

func mustGetContextViewStore(cmd *cobra.Command) viewstore.ViewStore {
	return cmd.Context().Value(viewStoreContextKey).(viewstore.ViewStore)
}

func mustGetContextSchemaStore(cmd *cobra.Command) schemastore.SchemaStore {
	return cmd.Context().Value(schemaStoreContextKey).(schemastore.SchemaStore)
}

func setContextViewStore(cmd *cobra.Command) error {
	store, err := local.NewLocalStore()
	if err != nil {
		return err
	}
	ctx := context.WithValue(cmd.Context(), viewStoreContextKey, store)
	cmd.SetContext(ctx)
	return nil
}

func setContextSchemaStore(cmd *cobra.Command) error {
	store, err := fileschema.NewFileSchemaStore()
	if err != nil {
		return err
	}
	ctx := context.WithValue(cmd.Context(), schemaStoreContextKey, store)
	cmd.SetContext(ctx)
	return nil
}

func WithViewStore(ctx context.Context, s viewstore.ViewStore) context.Context {
	return context.WithValue(ctx, viewStoreContextKey, s)
}

func WithSchemaStore(ctx context.Context, s schemastore.SchemaStore) context.Context {
	return context.WithValue(ctx, schemaStoreContextKey, s)
}

func printViewPretty(cmd *cobra.Command, view models.View, verbose bool, jsonOutput bool) {
	if jsonOutput {
		var output any
		if verbose {
			output = view
		} else {
			output = struct {
				Name      string           `json:"name"`
				Query     *string          `json:"query"`
				Sdl       *string          `json:"sdl"`
				Transform models.Transform `json:"transform"`
				Metadata  struct {
					Version   int    `json:"_v"`
					Total     int    `json:"_t"`
					CreatedAt string `json:"createdAt"`
					UpdatedAt string `json:"updatedAt"`
				} `json:"metadata"`
			}{
				Name:      view.Name,
				Query:     view.Query,
				Sdl:       view.Sdl,
				Transform: view.Transform,
				Metadata: struct {
					Version   int    `json:"_v"`
					Total     int    `json:"_t"`
					CreatedAt string `json:"createdAt"`
					UpdatedAt string `json:"updatedAt"`
				}{
					Version:   view.Metadata.Version,
					Total:     view.Metadata.Total,
					CreatedAt: view.Metadata.CreatedAt,
					UpdatedAt: view.Metadata.UpdatedAt,
				},
			}
		}

		encoded, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			cmd.Printf("❌ Failed to encode view: %v\n", err)
			return
		}
		cmd.Println(string(encoded))
		return
	}

	// === Pretty Output ===

	cmd.Printf("📄 View: %s\n", view.Name)

	if view.Query != nil && *view.Query != "" {
		cmd.Printf("🔍 Query:\n%s\n\n", *view.Query)
	} else {
		cmd.Println("🔍 Query: <none>")
	}

	if view.Sdl != nil && *view.Sdl != "" {
		cmd.Printf("📐 SDL:\n%s\n\n", *view.Sdl)
	} else {
		cmd.Println("📐 SDL: <none>")
	}

	cmd.Println("🔧 Lenses:")
	if len(view.Transform.Lenses) == 0 {
		cmd.Println(" - (empty)")
	} else {
		for _, lens := range view.Transform.Lenses {
			cmd.Printf(" - %s (%s)\n", lens.Label, lens.Path)
			if len(lens.Arguments) > 0 {
				cmd.Println("   Arguments:")
				for k, v := range lens.Arguments {
					cmd.Printf("     %s: %v\n", k, v)
				}
			}
		}
	}
	cmd.Println()

	createdAt, _ := strconv.ParseInt(view.Metadata.CreatedAt, 10, 64)
	updatedAt, _ := strconv.ParseInt(view.Metadata.UpdatedAt, 10, 64)

	cmd.Printf("🗂  Metadata:\n - Version: %d\n - Total: %d\n - Created At: %s\n - Updated At: %s\n",
		view.Metadata.Version,
		view.Metadata.Total,
		time.Unix(createdAt, 0).UTC(),
		time.Unix(updatedAt, 0).UTC(),
	)

	if verbose && len(view.Metadata.Revisions) > 0 {
		cmd.Printf("📝 Revisions (%d):\n", len(view.Metadata.Revisions))
		for i, rev := range view.Metadata.Revisions {
			cmd.Printf(" - Revision %d:\n   %s\n", i+1, rev.Diff)
		}
	}
}
