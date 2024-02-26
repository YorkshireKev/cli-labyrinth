package main

import (
	"fmt"

	"github.com/gdamore/tcell"
)

type Game struct {
	screen  tcell.Screen
	colPos  int
	rowPos  int
	dir     int //1 North, 2 East, 3 South, 4 West
	mapView bool
}

func (g *Game) init() {
	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Printf("Error creating screen: %s\n", err)
		return
	}
	//defer screen.Fini()
	//Set initial character postitions etc
	g.colPos = 1
	g.rowPos = 1
	g.dir = 1
	g.mapView = false

	if err := screen.Init(); err != nil {
		fmt.Printf("Error initializing screen: %s\n", err)
		return
	}
	screen.SetStyle(tcell.StyleDefault.
		Background(tcell.ColorWhite).
		Foreground(tcell.ColorBlack))
	g.screen = screen
	g.screen.Clear()

	//Default text on main game screen.
	g.PrintString(2, 22, "Left Arrow: Rotate Left.  Right Arrow: Rotate Right.  Up Arrow: Move Forward")
	g.PrintString(0, 23, "Press M to cheat and display the maze map")
	g.PrintString(54, 23, "Press Esc to exit the game")
}

func (g *Game) gameLoop() {
	for {
		//Handle Keyboard Input
		ev := g.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyCtrlC || ev.Key() == tcell.KeyEsc {
				//Quit the game
				return
			}
			switch ev.Key() {
			case tcell.KeyRune:
				if ev.Rune() == 'm' {
					//Show the maze in full
					g.mapView = true
				} else {
					g.mapView = false
				}
			case tcell.KeyLeft:
				g.dir-- //Rotate left
				if g.dir == 0 {
					g.dir = 4
				}
			case tcell.KeyRight:
				g.dir++ //Rotate right
				if g.dir == 5 {
					g.dir = 1
				}
			case tcell.KeyUp:
				switch g.dir {
				//Move forward one square if not blocked.
				case 1:
					if !maze.getMazeBlock(g.rowPos-1, g.colPos) {
						g.rowPos--
					}
				case 2:
					if !maze.getMazeBlock(g.rowPos, g.colPos+1) {
						g.colPos++
					}
				case 3:
					if !maze.getMazeBlock(g.rowPos+1, g.colPos) {
						g.rowPos++
					}
				case 4:
					if !maze.getMazeBlock(g.rowPos, g.colPos-1) {
						g.colPos--
					}
				}
			}
		case *tcell.EventResize:
			g.screen.Sync()
		}

		//Update gamestate
		if g.mapView {
			//Show the maze in full (Cheat!)
			g.PrintMaze(*maze)
		} else {
			// Draw the game screen
			g.DrawScreen(*maze)
		}

		//Need any kind of pause here? to limit FPS?
	}
}

func (g Game) DrawScreen(m Maze) {
	//Draw the 3d view of the maze based on where the player is standing
	//and the direction they are facing.

	//style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)
	//g.screen.Clear()

	var viewPort [6][3]bool
	var block rune

	//TODO - populate the viewPort grid with the view from the player!!
	switch g.dir {
	case 1:
		//Facing North
		for ix := 0; ix < 6; ix++ {
			// Iterate over rows
			for iy := 0; iy < 3; iy++ {
				viewPort[ix][iy] = m.getMazeBlock((ix+g.rowPos)-5, (iy+g.colPos)-1)

			}
		}
	case 2:
		//Facing East
		for ix := 0; ix < 6; ix++ {
			// Iterate over rows
			for iy := 0; iy < 3; iy++ {
				viewPort[5-ix][iy] = m.getMazeBlock((g.rowPos+iy)-1, g.colPos+ix)

			}
		}
	case 3:
		//Facing South
		for ix := 0; ix < 6; ix++ {
			// Iterate over rows
			for iy := 0; iy < 3; iy++ {
				viewPort[5-ix][2-iy] = m.getMazeBlock(ix+g.rowPos, (iy+g.colPos)-1)

			}
		}
	case 4:
		//Facing West
		for ix := 0; ix < 6; ix++ {
			// Iterate over rows
			for iy := 0; iy < 3; iy++ {
				viewPort[ix][2-iy] = m.getMazeBlock((g.rowPos+iy)-1, g.colPos+ix-5)

			}
		}
	}

	//Clear the maze viewport area of the screen
	for ix := 0; ix < 22; ix++ {
		// Iterate over columns
		for iy := 0; iy < 72; iy++ {
			g.screen.SetContent(iy, ix, '.', nil, tcell.StyleDefault)
		}
	}

	//TEMP - Draw 2d representation of the viewport
	// Iterate over rows
	for ix := 0; ix < 6; ix++ {
		// Iterate over columns
		for iy := 0; iy < 3; iy++ {
			if viewPort[ix][iy] {
				block = '█'
			} else {
				block = ' '
			}
			// Set content at position (ix, iy)
			//TODO - This is to render the viewport!!
			g.screen.SetContent((iy*2)+74, ix, block, nil, tcell.StyleDefault)
			g.screen.SetContent((iy*2)+75, ix, block, nil, tcell.StyleDefault)
			g.screen.SetContent(76, 5, 'X', nil, tcell.StyleDefault)
			g.screen.SetContent(77, 5, 'X', nil, tcell.StyleDefault)
		}
	}
	//TEMP END

	//Con only see so far...
	g.PrintString(33, 10, "???")
	g.PrintString(33, 11, "???")

	//====================================
	if viewPort[1][0] {
		for ix := 10; ix < 12; ix++ {
			g.PrintString(30, ix, "███")
		}
	} else {
		g.PrintString(30, 10, "░░░")
		g.PrintString(30, 11, "░░░")
	}

	if viewPort[1][1] {
		g.PrintString(33, 10, "░░░")
		g.PrintString(33, 11, "░░░")
	}
	//====================================

	//====================================
	if viewPort[2][0] {
		for ix := 8; ix < 14; ix++ {
			g.PrintString(24, ix, "███")
			if ix > 8 && ix < 13 {
				g.PrintString(27, ix, "███")
			}
		}
	} else {
		for ix := 9; ix < 13; ix++ {
			g.PrintString(27, ix, "░░░")
		}
	}

	if viewPort[2][1] {
		for ix := 9; ix < 13; ix++ {
			g.PrintString(30, ix, "░░░░░░░░░░░")
		}
	}
	//====================================

	//====================================
	if viewPort[3][0] {
		for ix := 6; ix < 16; ix++ {
			g.PrintString(18, ix, "███")
			if ix > 6 && ix < 15 {
				g.PrintString(21, ix, "███")
			}
		}
	} else {
		for ix := 8; ix < 14; ix++ {
			g.PrintString(18, ix, "░░░░░░")
		}
	}

	if viewPort[3][1] {
		for ix := 7; ix < 15; ix++ {
			g.PrintString(24, ix, "░░░░░░░░░░░░░░░░░░░░░░░")
		}
	}
	//====================================

	//====================================
	if viewPort[4][0] {
		for ix := 3; ix < 19; ix++ {
			g.PrintString(9, ix, "███")
			if ix > 3 && ix < 18 {
				g.PrintString(12, ix, "███")
			}
			if ix > 4 && ix < 17 {
				g.PrintString(15, ix, "███")
			}
		}
	} else {
		for ix := 6; ix < 16; ix++ {
			g.PrintString(9, ix, "░░░░░░░░░")
		}
	}

	if viewPort[4][1] {
		for ix := 3; ix < 19; ix++ {
			g.PrintString(9, ix, "░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░")
		}
	}
	//====================================

	//====================================
	//Closest wall left
	if viewPort[5][0] {
		for ix := 0; ix < 22; ix++ {
			g.PrintString(0, ix, "███")
			if ix > 0 && ix < 21 {
				g.PrintString(3, ix, "███")
			}
			if ix > 1 && ix < 20 {
				g.PrintString(6, ix, "███")
			}
		}
	} else {
		for ix := 3; ix < 19; ix++ {
			g.PrintString(0, ix, "░░░░░░░░░")
		}
	}

	//Closest wall middle is where the player is standing, so nothing to draw.

	//Closest wall right
	if viewPort[5][2] {

	}
	//====================================

	g.screen.Show()
}

func (g *Game) PrintMaze(m Maze) {
	//Draw the full map on the screen as a 2d grid.
	style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)

	for i := 0; i < m.rows; i++ {
		for j := 0; j < m.cols; j++ {
			if m.end.row == i && m.end.col == j {
				g.screen.SetContent(j*2, i, '(', nil, style)
				g.screen.SetContent(j*2+1, i, ')', nil, style)
			} else if m.grid[i][j] {
				g.screen.SetContent(j*2, i, '█', nil, style)
				g.screen.SetContent(j*2+1, i, '█', nil, style)
			} else {
				g.screen.SetContent(j*2, i, ' ', nil, style)
				g.screen.SetContent(j*2+1, i, ' ', nil, style)
			}
		}
	}

	//Print where the player is on the map
	dirRune := '↓'
	switch g.dir {
	case 1:
		dirRune = '↑' // U+2191
	case 2:
		dirRune = '→' //U+2192
	case 3:
		dirRune = '↓' //U+2193
	case 4:
		dirRune = '←' // U+2190
	}
	g.screen.SetContent(g.colPos*2, g.rowPos, 'A', nil, style)
	g.screen.SetContent((g.colPos*2)+1, g.rowPos, dirRune, nil, style)

	g.screen.Show()
}

func (g *Game) PrintString(x, y int, str string) {
	ix := 0 //Set the index for printing on the x axis. This is because runes can be more than 8 bits.
	for _, ch := range str {
		g.screen.SetContent(x+ix, y, ch, nil, tcell.StyleDefault)
		ix++
	}
}
