package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/gdamore/tcell"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tetro",
	Short: "cli tetromino",
	Run: func(cmd *cobra.Command, args []string) {
		screen, err := tcell.NewScreen()

		if err != nil {
			log.Fatalf("%+v", err)
		}
		if err := screen.Init(); err != nil {
			log.Fatalf("%+v", err)
		}

		defStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
		screen.SetStyle(defStyle)

		for {
			switch event := screen.PollEvent().(type) {
			case *tcell.EventResize:
				screen.Sync()
			case *tcell.EventKey:
				key := event.Key()

				if key == tcell.KeyEscape || key == tcell.KeyCtrlC {
					screen.Fini()
					os.Exit(0)
				} else if key == tcell.KeyDown {
					fmt.Println("down")
				} else if key == tcell.KeyLeft {
					fmt.Println("left")
				} else if key == tcell.KeyRight {
					fmt.Println("right")
				} else if event.Name() == "Rune[ ]" {
					fmt.Println("space")
				}
			}
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
