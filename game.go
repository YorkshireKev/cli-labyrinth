package main

import (
	"fmt"

	"github.com/gdamore/tcell"
)

type Game struct {
	screen  tcell.Screen
	xPos    int
	yPos    int
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
	g.xPos = 1
	g.yPos = 1
	g.dir = 1
	g.mapView = false

	if err := screen.Init(); err != nil {
		fmt.Printf("Error initializing screen: %s\n", err)
		return
	}
	g.screen = screen

	g.screen.Clear()
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
					if !maze.getMazeBlock(g.xPos, g.yPos-1) {
						g.yPos--
					}
				case 2:
					if !maze.getMazeBlock(g.xPos+1, g.yPos) {
						g.xPos++
					}
				case 3:
					if !maze.getMazeBlock(g.xPos, g.yPos+1) {
						g.yPos++
					}
				case 4:
					if !maze.getMazeBlock(g.xPos-1, g.yPos) {
						g.xPos--
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
	//TODO - Draw the 3d view of the maze.
	style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)
	g.screen.SetContent(1, 23, 'X', nil, style)
	g.screen.Show()
}

func (g *Game) PrintMaze(m Maze) {
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
	g.screen.SetContent(g.xPos*2, g.yPos, 'A', nil, style)
	g.screen.SetContent((g.xPos*2)+1, g.yPos, dirRune, nil, style)

	g.screen.Show()
}
