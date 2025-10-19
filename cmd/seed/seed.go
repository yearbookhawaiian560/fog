package seed

import (
	"github.com/spf13/cobra"

	"github.com/0div/fog/internal/entry/seed"
)

func Setup() *cobra.Command {
	migrateRootCmd := &cobra.Command{
		Use:   "seed",
		Short: "Seed the database with initial data",
		Run: func(cmd *cobra.Command, args []string) {
			seed.Seed()
		},
	}

	return migrateRootCmd
}
