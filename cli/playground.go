package cli

import (
	"github.com/shinzonetwork/shinzo-view-creator/core/schema/store/fileschema"
	"github.com/shinzonetwork/shinzo-view-creator/core/service"
	"github.com/spf13/cobra"
)

func MakePlaygroundCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "playground",
		Short: "Start a local DefraDB node with schema and mock data",
		Long:  "Starts DefraDB, applies the schema, inserts mock data, and keeps it running for exploration. No views are deployed.",
		RunE: func(cmd *cobra.Command, args []string) error {
			schemastore, err := fileschema.NewFileSchemaStore()
			if err != nil {
				return err
			}
			return service.StartLocalNodePlayground(schemastore, false)
		},
	}

	return cmd
}
