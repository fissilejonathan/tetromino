package game

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell"
)

const (
	yLength            = 21
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

	nCurrentPiece := rand.Intn(7)
	nCurrentRotation := 0
	nCurrentX := xLength / 2
	nCurrentY := 0
	nSpeed := 20
	nSpeedCount := 0
	bForceDown := false
	bRotateHold := true
	nPieceCount := 0
	nScore := 0
	vLines := []int{}
	bGameOver := false

	input := make(chan *tcell.EventKey)

	// poll for user input
	go func() {
		for {
			switch event := g.screen.PollEvent().(type) {
			case *tcell.EventResize:
				g.screen.Sync()
			case *tcell.EventKey:
				input <- event
			}
		}
	}()

	for !bGameOver {
		time.Sleep(50 * time.Millisecond)

		nSpeedCount++
		bForceDown = (nSpeedCount == nSpeed)

		select {
		case event := <-input:
			key := event.Key()

			if key == tcell.KeyEscape || key == tcell.KeyCtrlC {
				g.screen.Fini()
				close(input)
				os.Exit(0)
			} else if key == tcell.KeyUp {
				rotate := 0

				if bRotateHold && g.doesPieceFit(nCurrentPiece, nCurrentRotation+1, nCurrentX, nCurrentY) {
					rotate = 1
				}

				nCurrentRotation += rotate
				bRotateHold = false
			} else if key == tcell.KeyDown {
				move := 0

				if g.doesPieceFit(nCurrentPiece, nCurrentRotation, nCurrentX, nCurrentY+1) {
					move = 1
				}

				nCurrentY += move
			} else if key == tcell.KeyLeft {
				move := 0

				if g.doesPieceFit(nCurrentPiece, nCurrentRotation, nCurrentX-1, nCurrentY) {
					move = 1
				}

				nCurrentX -= move
			} else if key == tcell.KeyRight {
				move := 0

				if g.doesPieceFit(nCurrentPiece, nCurrentRotation, nCurrentX+1, nCurrentY) {
					move = 1
				}

				nCurrentX += move
			} else if event.Name() == "Rune[ ]" {
			} else {
				bRotateHold = true
			}
		default:
			// do nothing
		}

		// Force the piece down the playfield if it's time
		if bForceDown {
			// Update difficulty every 50 pieces
			nSpeedCount = 0
			nPieceCount++

			if nPieceCount%50 == 0 {
				if nSpeed >= 10 {
					nSpeed--
				}
			}

			// Test if piece can be moved down
			if g.doesPieceFit(nCurrentPiece, nCurrentRotation, nCurrentX, nCurrentY+1) {
				nCurrentY++
			} else {
				// It can't! Lock the piece in place
				for px := 0; px < 4; px++ {
					for py := 0; py < 4; py++ {
						if g.tetrominos[nCurrentPiece][g.rotate(px, py, nCurrentRotation)] != '.' {
							g.board[(nCurrentY+py)*xLength+(nCurrentX+px)] = rune(nCurrentPiece + 1)
						}
					}
				}

				// Check for lines
				for py := 0; py < 4; py++ {
					if nCurrentY+py < yLength-1 {
						bLine := true
						for px := 1; px < xLength-1; px++ {
							bLine = bLine && (g.board[(nCurrentY+py)*xLength+px]) != 0
						}

						if bLine {
							// Remove Line, set to =
							for px := 1; px < xLength-1; px++ {
								g.board[(nCurrentY+py)*xLength+px] = 8
							}

							vLines = append(vLines, (nCurrentY + py))
						}
					}
				}

				nScore += 25
				if len(vLines) != 0 {
					nScore += (1 << len(vLines)) * 100
				}

				// Pick New Piece
				nCurrentX = xLength / 2
				nCurrentY = 0
				nCurrentRotation = 0
				nCurrentPiece = rand.Intn(7)

				// If piece does not fit straight away, game over!
				bGameOver = !g.doesPieceFit(nCurrentPiece, nCurrentRotation, nCurrentX, nCurrentY)
			}
		}

		// Draw Field
		for x := 0; x < xLength; x++ {
			for y := 0; y < yLength; y++ {
				value := rune(" ABCDEFG=#"[g.board[y*xLength+x]])

				g.screen.SetContent(x, y, value, nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
			}
		}

		// Draw Current Piece
		for px := 0; px < 4; px++ {
			for py := 0; py < 4; py++ {
				if g.tetrominos[nCurrentPiece][g.rotate(px, py, nCurrentRotation)] != '.' {
					x := nCurrentX + px
					y := nCurrentY + py
					value := rune(nCurrentPiece + 65)
					g.screen.SetContent(x, y, value, nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
				}
			}
		}

		g.screen.Sync()

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
}

func (g *Game) doesPieceFit(nTetromino, nRotation, nPosX, nPosY int) bool {
	// All Field cells >0 are occupied
	for px := 0; px < 4; px++ {
		for py := 0; py < 4; py++ {
			// Get index into piece
			pi := g.rotate(px, py, nRotation)

			// Get index into field
			fi := (nPosY+py)*xLength + (nPosX + px)

			// Check that test is in bounds. Note out of bounds does
			// not necessarily mean a fail, as the long vertical piece
			// can have cells that lie outside the boundary, so we'll
			// just ignore them
			if nPosX+px >= 0 && nPosX+px < xLength {
				if nPosY+py >= 0 && nPosY+py < yLength {
					// In Bounds so do collision check
					if g.tetrominos[nTetromino][pi] != '.' && g.board[fi] != rune(0) {
						return false // fail on first hit
					}
				}
			}
		}
	}

	return true
}

func (g *Game) rotate(px, py, r int) int {
	pi := 0

	switch r % 4 {
	case 0: // 0 degrees
		pi = py*4 + px
	case 1: // 90 degrees
		pi = 12 + py - (px * 4)
	case 2: // 180 degrees
		pi = 15 - (py * 4) - px
	case 3: // 270 degrees
		pi = 3 - py + (px * 4)
	}

	return pi
}

func (g *Game) getXYFromIndex(index, numColumns int) (int, int) {
	y := index / numColumns
	x := index % numColumns
	return x, y
}
