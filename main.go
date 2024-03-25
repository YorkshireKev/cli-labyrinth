package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/gdamore/tcell/v2"
)

var maze *Maze
var screen tcell.Screen
var operatingSystem string

func main() {
	//Read the first command line parameter to check for the parameter: --cheat
	//which will allow the player to move while viewing trhe map.
	cheat := false
	if len(os.Args) > 1 && os.Args[1] == "--cheat" {
		cheat = true
	}

	operatingSystem = runtime.GOOS

	screen = initScreen()
	for {
		//rows, cols := 9, 11
		rows, cols := showTitle(cheat)

		maze = NewMaze(rows, cols)
		maze.generateMaze()

		game := &Game{}
		game.init(screen, cheat)
		if game.gameLoop() {
			//Player has escaped the maze!
			showEscapedScreen(game, cheat)
		}
	}

	//game.screen.Fini()
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

func showTitle(cheat bool) (rows int, cols int) {
	screen.Clear()
	PrintString(29, 0, "Command Line Interface")
	PrintString(6, 2, "#          #    ######  #     # ######  ### #     # ####### #     #")
	PrintString(6, 3, "#         # #   #     #  #   #  #     #  #  ##    #    #    #     #")
	PrintString(6, 4, "#        #   #  #     #   # #   #     #  #  # #   #    #    #     #")
	PrintString(6, 5, "#       #     # ######     #    ######   #  #  #  #    #    #######")
	PrintString(6, 6, "#       ####### #     #    #    #   #    #  #   # #    #    #     #")
	PrintString(6, 7, "#       #     # #     #    #    #    #   #  #    ##    #    #     #")
	PrintString(6, 8, "####### #     # ######     #    #     # ### #     #    #    #     #")
	PrintString(22, 10, "©2024 By Kev Ellis - www.kevssite.com")
	PrintString(37, 11, "(v0.2)")
	PrintString(2, 16, "Please select your maze size: Press S for Small, M for Medium or L for Large")
	if operatingSystem != "js" {
		PrintString(15, 23, "Press Esc or Ctrl-C to Quit to exit this program")
	}
	if cheat {
		PrintString(73, 23, "(cheat)")
	}
	screen.Show()

	for {
		//Handle Keyboard Input
		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			if operatingSystem != "js" {
				if ev.Key() == tcell.KeyCtrlC || ev.Key() == tcell.KeyEsc {
					//Quit the game
					screen.Fini()
					os.Exit(0)
				}
			}
			if ev.Key() == tcell.KeyRune {
				if ev.Rune() == 's' {
					//Small Map
					return 9, 11
				}
				if ev.Rune() == 'm' {
					//Medium Map
					return 15, 23
				}
				if ev.Rune() == 'l' {
					//Large Map
					return 23, 37
				}
			}
		case *tcell.EventResize:
			screen.Sync()
		}
	}
}

func showEscapedScreen(g *Game, cheat bool) {
	screen.Clear()
	PrintString(29, 0, "** Congratulations **")
	PrintString(9, 2, "#     # ####### #     #     #     #    #    #     # #######")
	PrintString(9, 3, " #   #  #     # #     #     #     #   # #   #     # #")
	PrintString(9, 4, "  # #   #     # #     #     #     #  #   #  #     # # ")
	PrintString(9, 5, "   #    #     # #     #     ####### #     # #     # #####")
	PrintString(9, 6, "   #    #     # #     #     #     # #######  #   #  #")
	PrintString(9, 7, "   #    #     # #     #     #     # #     #   # #   #")
	PrintString(9, 8, "   #    #######  #####      #     # #     #    #    #######")
	PrintString(9, 10, "#######  #####   #####     #    ######  ####### ######  ###")
	PrintString(9, 11, "#       #     # #     #   # #   #     # #       #     # ###")
	PrintString(9, 12, "#       #       #        #   #  #     # #       #     # ###")
	PrintString(9, 13, "#####    #####  #       #     # ######  #####   #     #  #")
	PrintString(9, 14, "#             # #       ####### #       #       #     #")
	PrintString(9, 15, "#       #     # #     # #     # #       #       #     # ###")
	PrintString(9, 16, "#######  #####   #####  #     # #       ####### ######  ###")
	PrintString(12, 18, "You took "+fmt.Sprint(g.steps)+" steps and you looked at the map "+fmt.Sprint(g.mapShown)+" times")
	PrintString(20, 21, "Poke a key to return to the title screen")
	if operatingSystem != "js" {
		PrintString(15, 23, "Press Esc or Ctrl-C to Quit to exit this program")
	}
	if cheat {
		PrintString(73, 23, "(cheat)")
	}
	screen.Show()

	for {
		//Handle Keyboard Input
		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyCtrlC || ev.Key() == tcell.KeyEsc {
				//Quit the game
				screen.Fini()
				os.Exit(0)
			}
			if ev.Key() == tcell.KeyRune {
				return
			}
		case *tcell.EventResize:
			screen.Sync()
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
