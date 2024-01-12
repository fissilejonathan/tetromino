package game

import (
	"log"
	"os"

	"github.com/gdamore/tcell"
)

const (
	yLength            = 22
	xLength            = 12
	boardDimensions    = yLength * xLength
	numberOfTetrominos = 7
)

type Game struct {
	screen     tcell.Screen
	board      [boardDimensions]rune
	tetrominos [numberOfTetrominos]string
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
		screen:     screen,
		board:      [boardDimensions]rune{},
		tetrominos: [numberOfTetrominos]string{},
	}
}

func (g *Game) Start() {
	g.setup()

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
			} else if key == tcell.KeyLeft {
			} else if key == tcell.KeyRight {
			} else if event.Name() == "Rune[ ]" {
			}

			g.screen.Show()
		}
	}
}

func (g *Game) setup() {
	///////////////////
	// setup tetrominos
	///////////////////
	g.tetrominos[0] = `..X...X...X...X.`
	g.tetrominos[1] = "..X..XX...X....."
	g.tetrominos[2] = ".....XX..XX....."
	g.tetrominos[3] = "..X..XX..X......"
	g.tetrominos[4] = ".X...XX...X....."
	g.tetrominos[5] = ".X...X...XX....."
	g.tetrominos[6] = "..X...X..XX....."

	//////////////
	// setup board
	//////////////
	for x := 0; x < xLength; x++ {
		for y := 0; y < yLength; y++ {
			value := 0

			if x == 0 || x == xLength-1 || y == yLength-1 {
				value = 9
			}

			g.board[y*xLength+x] = rune(value)
		}
	}

	///////////////
	// setup screen
	///////////////
	defStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	g.screen.SetStyle(defStyle)

	for i := 0; i < boardDimensions; i++ {
		value := g.board[i]

		x, y := getXYFromIndex(i, xLength)

		if value == rune(9) {
			g.screen.SetContent(x, y, value, nil, tcell.StyleDefault.Background(tcell.ColorWhite))

		}
	}

	g.screen.Show()
}

func getXYFromIndex(index, numColumns int) (int, int) {
	y := index / numColumns
	x := index % numColumns
	return x, y
}
