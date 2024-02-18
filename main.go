package main

var maze *Maze

func main() {

	rows, cols := 21, 37
	maze = NewMaze(rows, cols)
	maze.generateMaze()

	game := &Game{}
	game.init()
	game.gameLoop()

	game.screen.Fini()
}
