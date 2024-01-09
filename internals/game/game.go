package game

import (
	"log"
	"os"

	"github.com/gdamore/tcell"
)

type Game struct {
	screen tcell.Screen
}

func New() *Game {
	screen, err := tcell.NewScreen()

	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	defStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	screen.SetStyle(defStyle)

	return &Game{
		screen: screen,
	}
}

func (g *Game) Start() {
	for {
		switch event := g.screen.PollEvent().(type) {
		case *tcell.EventResize:
			g.screen.Sync()
		case *tcell.EventKey:
			key := event.Key()

			if key == tcell.KeyEscape || key == tcell.KeyCtrlC {
				g.screen.Fini()
				os.Exit(0)
			} else if key == tcell.KeyDown {
				g.screen.SetContent(50, 10, ' ', nil, tcell.StyleDefault.Background(tcell.ColorRed))
			} else if key == tcell.KeyLeft {
				g.screen.SetContent(51, 11, ' ', nil, tcell.StyleDefault.Background(tcell.ColorRed))
			} else if key == tcell.KeyRight {
				g.screen.SetContent(52, 12, ' ', nil, tcell.StyleDefault.Background(tcell.ColorRed))
			} else if event.Name() == "Rune[ ]" {
				g.screen.SetContent(53, 13, ' ', nil, tcell.StyleDefault.Background(tcell.ColorRed))
			}

			g.screen.Show()
		}
	}
}
