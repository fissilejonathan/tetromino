package cmd

import (
	"fmt"
	"os"

	"github.com/fissilejonathan/tetromino/internals/game"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tetro",
	Short: "cli tetromino",
	Run: func(cmd *cobra.Command, args []string) {
		game := game.New()

		game.Start()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
