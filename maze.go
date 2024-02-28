package main

import (
	"math/rand"
)

type Maze struct {
	rows, cols int
	grid       [][]bool
	start, end Cell
}

func NewMaze(rows, cols int) *Maze {
	maze := &Maze{rows: rows, cols: cols}
	maze.initGrid()
	return maze
}

func (m *Maze) initGrid() {
	m.grid = make([][]bool, m.rows)
	for i := range m.grid {
		m.grid[i] = make([]bool, m.cols)
		for j := range m.grid[i] {
			m.grid[i][j] = true
		}
	}
}

func (m *Maze) generateMaze() {
	startRow, startCol := 1, 1
	m.start = Cell{startRow, startCol}
	m.grid[startRow][startCol] = false

	stack := []Cell{m.start}
	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		neighbors := m.getNeighbors(current.row, current.col)
		rand.Shuffle(len(neighbors), func(i, j int) {
			neighbors[i], neighbors[j] = neighbors[j], neighbors[i]
		})

		for _, neighbor := range neighbors {
			if m.grid[neighbor.row][neighbor.col] {
				m.grid[current.row+(neighbor.row-current.row)/2][current.col+(neighbor.col-current.col)/2] = false
				m.grid[neighbor.row][neighbor.col] = false
				stack = append(stack, neighbor)
			}
		}
	}

	// Set end position
	endRow, endCol := m.rows-2, m.cols-2
	m.end = Cell{endRow, endCol}
	m.grid[endRow][endCol] = false

	//m.grid[m.rows-1][m.cols-2] = false
	m.grid[0][1] = false //Exit always top left
}

type Cell struct {
	row, col int
}

func (m *Maze) getNeighbors(row, col int) []Cell {
	neighbors := []Cell{
		{row - 2, col},
		{row + 2, col},
		{row, col - 2},
		{row, col + 2},
	}

	validNeighbors := make([]Cell, 0)
	for _, neighbor := range neighbors {
		if m.isValidMove(neighbor.row, neighbor.col) {
			validNeighbors = append(validNeighbors, neighbor)
		}
	}

	return validNeighbors
}

func (m *Maze) isValidMove(row, col int) bool {
	return row > 0 && row < m.rows && col > 0 && col < m.cols && m.grid[row][col]
}

func (m *Maze) getMazeBlock(row, col int) bool {
	if row >= 0 && row < len(m.grid) && col >= 0 && col < len(m.grid[0]) {
		return m.grid[row][col]
	} else {
		//if the array indexes are our of bounds of the array, just return false.
		return false
	}
}
