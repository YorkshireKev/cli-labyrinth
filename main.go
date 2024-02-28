package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"
)

var maze *Maze
var screen tcell.Screen

func main() {

	screen = initScreen()

	showTitle()

	//rows, cols := 21, 37
	rows, cols := 9, 11
	maze = NewMaze(rows, cols)
	maze.generateMaze()

	game := &Game{}
	game.init(screen)
	game.gameLoop()

	game.screen.Fini()
}

func initScreen() tcell.Screen {
	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Printf("Error creating screen: %s\n", err)
		os.Exit(99)
	}
	//defer screen.Fini()

	if err := screen.Init(); err != nil {
		fmt.Printf("Error initializing screen: %s\n", err)
		os.Exit(99)
	}
	screen.SetStyle(tcell.StyleDefault.
		Background(tcell.ColorWhite).
		Foreground(tcell.ColorBlack))

	screen.Clear()

	return screen
}

func showTitle() {
	PrintString(29, 0, "Command Line Interface")
	PrintString(6, 2, "#          #    ######  #     # ######  ### #     # ####### #     #")
	PrintString(6, 3, "#         # #   #     #  #   #  #     #  #  ##    #    #    #     #")
	PrintString(6, 4, "#        #   #  #     #   # #   #     #  #  # #   #    #    #     #")
	PrintString(6, 5, "#       #     # ######     #    ######   #  #  #  #    #    #######")
	PrintString(6, 6, "#       ####### #     #    #    #   #    #  #   # #    #    #     #")
	PrintString(6, 7, "#       #     # #     #    #    #    #   #  #    ##    #    #     #")
	PrintString(6, 8, "####### #     # ######     #    #     # ### #     #    #    #     #")
	PrintString(31, 10, "©2024 By Kev Ellis")
	PrintString(2, 16, "Please select your maze size: Press S for Small, M for Medium or L for Large")
	PrintString(15, 23, "Press Esc or Ctrl-C to Quit to exit this program")
	screen.Show()

	for {
		//Handle Keyboard Input
		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyCtrlC || ev.Key() == tcell.KeyEsc {
				//Quit the game
				return
			}
		}
	}
}

func PrintString(x, y int, str string) {
	ix := 0 //Set the index for printing on the x axis. This is because runes can be more than 8 bits.
	for _, ch := range str {
		screen.SetContent(x+ix, y, ch, nil, tcell.StyleDefault)
		ix++
	}
}
