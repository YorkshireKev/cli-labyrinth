package main

import (
	"fmt"

	"github.com/gdamore/tcell"
)

type Game struct {
	screen   tcell.Screen
	colPos   int
	rowPos   int
	dir      int //1 North, 2 East, 3 South, 4 West
	mapView  bool
	mapShown int //How many times the map was shown
	steps    uint
	cheat    bool
}

func (g *Game) init(screen tcell.Screen, cheat bool) {
	//Set initial character postitions etc
	g.colPos = maze.cols - 2
	g.rowPos = maze.rows - 2
	g.dir = 1
	g.mapView = false
	g.mapShown = 0
	g.steps = 0
	g.cheat = cheat

	g.screen = screen
	g.screen.Clear()

	//Default text on main game screen.
	PrintString(2, 22, "Left Arrow: Rotate Left.  Right Arrow: Rotate Right.  Up Arrow: Move Forward")
	PrintString(0, 23, "Press M to cheat and display the maze map")
	PrintString(54, 23, "Press Esc to exit the game")
	g.DrawScreen(*maze)

}

func (g *Game) gameLoop() bool {
	for {
		//Handle Keyboard Input
		ev := g.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyCtrlC || ev.Key() == tcell.KeyEsc {
				//Quit the game
				return false
			}
			switch ev.Key() {
			case tcell.KeyRune:
				if ev.Rune() == 'm' {
					//Show the maze in full
					if !g.mapView {
						PrintString(2, 22, "       ** Study Carefully, then press a key to return to the game. **       ")
						g.mapShown++
					}
					g.mapView = true
				} else {
					if g.mapView {
						PrintString(2, 22, "Left Arrow: Rotate Left.  Right Arrow: Rotate Right.  Up Arrow: Move Forward")
					}
					g.mapView = false
				}
			case tcell.KeyLeft:
				if !g.mapView || g.cheat {
					g.dir-- //Rotate left
					if g.dir == 0 {
						g.dir = 4
					}
				}
			case tcell.KeyRight:
				if !g.mapView || g.cheat {
					g.dir++ //Rotate right
					if g.dir == 5 {
						g.dir = 1
					}
				}
			case tcell.KeyUp:
				if !g.mapView || g.cheat {
					switch g.dir {
					//Move forward one square if not blocked.
					case 1:
						if !maze.getMazeBlock(g.rowPos-1, g.colPos) {
							g.rowPos--
							g.steps++
						}
					case 2:
						if !maze.getMazeBlock(g.rowPos, g.colPos+1) {
							g.colPos++
							g.steps++
						}
					case 3:
						if !maze.getMazeBlock(g.rowPos+1, g.colPos) {
							g.rowPos++
							g.steps++
						}
					case 4:
						if !maze.getMazeBlock(g.rowPos, g.colPos-1) {
							g.colPos--
							g.steps++
						}
					}
					//Ensure the number of steps isn't too big to display
					if g.steps > 999999 {
						g.steps = 999999
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

		//If the player gets to column one, but row -1 (outside the maze)
		//then they have escaped!
		if g.rowPos == -1 && g.colPos == 1 {
			return true
		}

		//No need to add a pause here as the event polling will wait for a keypress.
	}
}

func (g Game) DrawScreen(m Maze) {
	//Draw the 3d view of the maze based on where the player is standing
	//and the direction they are facing.

	//style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)
	//g.screen.Clear()

	var viewPort [6][3]bool

	//Populate the viewPort grid with the view from the player.
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
	g.ClearScreen()

	//FOR DEBUG ONLY - Draw the rotated viewport
	// Iterate over rows
	/*var block rune
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
	}*/
	//END - Draw the rotated viewport

	//Can only see so far...
	PrintString(33, 10, "...")
	PrintString(33, 11, "...")

	//====================================
	if viewPort[1][0] {
		for ix := 10; ix < 12; ix++ {
			PrintString(30, ix, "███")
		}
	} else {
		PrintString(30, 10, "░░░")
		PrintString(30, 11, "░░░")
	}

	if viewPort[1][2] {
		for ix := 10; ix < 12; ix++ {
			PrintString(36, ix, "███")
		}
	} else {
		PrintString(36, 10, "░░░")
		PrintString(36, 11, "░░░")
	}

	if viewPort[1][1] {
		PrintString(33, 10, "░░░")
		PrintString(33, 11, "░░░")
	}
	//====================================

	//====================================
	if viewPort[2][0] {
		for ix := 8; ix < 14; ix++ {
			PrintString(24, ix, "███")
			if ix > 8 && ix < 13 {
				PrintString(27, ix, "███")
			}
		}
	} else {
		for ix := 10; ix < 12; ix++ {
			PrintString(24, ix, "░░░░░░")
		}
	}

	if viewPort[2][2] {
		for ix := 8; ix < 14; ix++ {
			PrintString(42, ix, "███")
			if ix > 8 && ix < 13 {
				PrintString(39, ix, "███")
			}
		}
	} else {
		for ix := 10; ix < 12; ix++ {
			PrintString(39, ix, "░░░░░░")
		}
	}

	if viewPort[2][1] {
		for ix := 9; ix < 13; ix++ {
			PrintString(30, ix, "░░░░░░░░░")
		}
	}
	//====================================

	//====================================
	if viewPort[3][0] {
		for ix := 6; ix < 16; ix++ {
			PrintString(18, ix, "███")
			if ix > 6 && ix < 15 {
				PrintString(21, ix, "███")
			}
		}
	} else {
		for ix := 8; ix < 14; ix++ {
			PrintString(18, ix, "░░░░░░")
		}
	}

	if viewPort[3][2] {
		for ix := 6; ix < 16; ix++ {
			PrintString(48, ix, "███")
			if ix > 6 && ix < 15 {
				PrintString(45, ix, "███")
			}
		}
	} else {
		for ix := 8; ix < 14; ix++ {
			PrintString(45, ix, "░░░░░░")
		}
	}

	if viewPort[3][1] {
		for ix := 7; ix < 15; ix++ {
			PrintString(24, ix, "░░░░░░░░░░░░░░░░░░░░░")
		}
	}
	//====================================

	//====================================
	if viewPort[4][0] {
		for ix := 3; ix < 19; ix++ {
			PrintString(9, ix, "███")
			if ix > 3 && ix < 18 {
				PrintString(12, ix, "███")
			}
			if ix > 4 && ix < 17 {
				PrintString(15, ix, "███")
			}
		}
	} else {
		for ix := 6; ix < 16; ix++ {
			PrintString(9, ix, "░░░░░░░░░")
		}
	}

	if viewPort[4][2] {
		for ix := 3; ix < 19; ix++ {
			PrintString(57, ix, "███")
			if ix > 3 && ix < 18 {
				PrintString(54, ix, "███")
			}
			if ix > 4 && ix < 17 {
				PrintString(51, ix, "███")
			}
		}
	} else {
		for ix := 6; ix < 16; ix++ {
			PrintString(51, ix, "░░░░░░░░░")
		}
	}

	if viewPort[4][1] {
		for ix := 3; ix < 19; ix++ {
			PrintString(9, ix, "░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░")
		}
	}
	//====================================

	//====================================
	//Closest wall left
	if viewPort[5][0] {
		for ix := 0; ix < 22; ix++ {
			PrintString(0, ix, "███")
			if ix > 0 && ix < 21 {
				PrintString(3, ix, "███")
			}
			if ix > 1 && ix < 20 {
				PrintString(6, ix, "███")
			}
		}
	} else {
		for ix := 3; ix < 19; ix++ {
			PrintString(0, ix, "░░░░░░░░░")
		}
	}

	//Closest wall middle is where the player is standing, so nothing to draw.

	//Closest wall right
	if viewPort[5][2] {
		for ix := 0; ix < 22; ix++ {
			PrintString(66, ix, "███")
			if ix > 0 && ix < 21 {
				PrintString(63, ix, "███")
			}
			if ix > 1 && ix < 20 {
				PrintString(60, ix, "███")
			}
		}
	} else {
		for ix := 3; ix < 19; ix++ {
			PrintString(60, ix, "░░░░░░░░░")
		}
	}
	//====================================

	//If near the exit then cull the side walls that were drawn beyond the exit.
	g.CullExitSideWalls(*maze)

	g.DrawSidebar()
	g.screen.Show()
}

func (g *Game) PrintMaze(m Maze) {
	//Draw the full map on the screen as a 2d grid.
	style := tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorWhite)

	g.ClearScreen()

	for i := 0; i < m.rows; i++ {
		for j := 0; j < m.cols; j++ {
			if m.grid[i][j] {
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

/*func (g *Game) PrintString(x, y int, str string) {
	ix := 0 //Set the index for printing on the x axis. This is because runes can be more than 8 bits.
	for _, ch := range str {
		g.screen.SetContent(x+ix, y, ch, nil, tcell.StyleDefault)
		ix++
	}
}*/

func (g *Game) ClearScreen() {
	for ix := 0; ix < 22; ix++ {
		// Iterate over columns
		for iy := 0; iy < 80; iy++ {
			g.screen.SetContent(iy, ix, ' ', nil, tcell.StyleDefault)
		}
	}
}

func (g *Game) DrawSidebar() {
	PrintString(74, 4, "N")
	PrintString(74, 5, "|")
	PrintString(71, 6, "<--+-->")
	PrintString(74, 7, "|")
	PrintString(74, 8, "S")
	switch g.dir {
	case 1:
		PrintString(74, 3, "*")
	case 2:
		PrintString(78, 6, "*")
	case 3:
		PrintString(74, 9, "*")
	case 4:
		PrintString(70, 6, "*")
	}

	if g.mapShown > 0 {
		PrintString(70, 17, "Map Show")
		if g.mapShown > 25 {
			PrintString(70, 18, "25+ Times")
		} else {
			PrintString(70, 18, fmt.Sprint(g.mapShown)+" Times")
		}
	}
	if g.mapShown > 50 {
		g.mapShown = 50 //Prevent int out of range errors
	}

	PrintString(72, 12, "Steps")
	steps := fmt.Sprint(g.steps)
	PrintString(72, 13, steps)

}

func (g Game) CullExitSideWalls(m Maze) {
	if g.dir == 1 && g.colPos == 1 {
		switch g.rowPos {
		case 4:
			PrintString(33, 10, "   ")
			PrintString(33, 11, "   ")
		case 3:
			PrintString(30, 10, "         ")
			PrintString(30, 11, "         ")
		case 2:
			PrintString(24, 10, "                     ")
			PrintString(24, 11, "                     ")
		case 1:
			for ix := 8; ix < 14; ix++ {
				PrintString(18, ix, "                                 ")
			}
		case 0:
			for ix := 2; ix < 20; ix++ {
				PrintString(9, ix, "                                                   ")
			}
		case -1:
			g.ClearScreen()
		}
	}
}
