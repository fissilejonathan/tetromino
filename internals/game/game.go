package game

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell"
)

type Game struct {
	screen tcell.Screen
}

func New() *Game {
	screen, err := tcell.NewScreen()
	screen.DisableMouse()

	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	return &Game{
		screen: screen,
	}
}

func (g *Game) Start() {
	g.setupScreen()

	done := make(chan struct{})
	defer close(done)

	go func() {
		ticker := time.NewTicker(750 * time.Millisecond)

		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				startCodePoint := 0x0041 // 'A'
				endCodePoint := 0x005A   // 'Z'
				randomCodePoint := rand.Intn(endCodePoint-startCodePoint+1) + startCodePoint
				randomRune := rune(randomCodePoint)

				g.screen.SetContent(25, 0, randomRune, nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
				g.screen.Show()
			}
		}
	}()

	for {
		switch event := g.screen.PollEvent().(type) {
		case *tcell.EventResize:
			g.screen.Sync()
		case *tcell.EventKey:
			key := event.Key()

			if key == tcell.KeyEscape || key == tcell.KeyCtrlC {
				g.screen.Fini()
				done <- struct{}{}
				os.Exit(0)
			} else if key == tcell.KeyDown {
			} else if key == tcell.KeyLeft {
			} else if key == tcell.KeyRight {
			} else if event.Name() == "Rune[ ]" {
			}

			g.screen.Show()
		}
	}
}

func (g *Game) updateScreen() {

}

func (g *Game) setupScreen() {
	defStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	g.screen.SetStyle(defStyle)

	for i := 0; i <= 21; i++ {
		// top and bottom
		g.screen.SetContent(i, 0, ' ', nil, tcell.StyleDefault.Background(tcell.ColorWhite))
		g.screen.SetContent(i, 41, ' ', nil, tcell.StyleDefault.Background(tcell.ColorWhite))
	}

	for i := 0; i <= 41; i++ {
		// left and right
		g.screen.SetContent(0, i, ' ', nil, tcell.StyleDefault.Background(tcell.ColorWhite))
		g.screen.SetContent(21, i, ' ', nil, tcell.StyleDefault.Background(tcell.ColorWhite))
	}

	g.screen.Show()
}
