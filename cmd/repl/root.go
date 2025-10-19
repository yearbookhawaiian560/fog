package repl

import (
	"github.com/spf13/cobra"

	"github.com/0div/fog/internal/entry/repl"
)

func Setup() *cobra.Command {
	var replCmd = &cobra.Command{
		Use:   "repl",
		Short: "Starts the Fog REPL",
		Run: func(cmd *cobra.Command, args []string) {
			repl.RunFogREPL()
		},
	}

	return replCmd
}
