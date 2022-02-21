package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

const (
	DEAD  = "▅"
	ALIVE = "○"
)

type coordinate struct {
	x int
	y int
}

type Map struct {
	size  int
	board [][]string
}

func (board Map) PrintMap() {
	for row := 0; row < board.size; row++ {
		for col := 0; col < board.size; col++ {
			fmt.Print(board.board[row][col])
		}
		fmt.Print("\n")
	}
}

func newMap(size int) Map {
	board := make([][]string, size)
	for row := range board {
		board[row] = make([]string, size)
	}
	for row := 0; row < size; row++ {
		for col := 0; col < size; col++ {
			board[row][col] = "▅"
		}
	}
	return Map{
		size:  size,
		board: board,
	}
}

func readCoordinates(nCoordinates, size int) []coordinate {
	var coordinates []coordinate
	for pair := 0; pair < nCoordinates; pair++ {
		var x, y int
		fmt.Printf("pair %d\nx:", pair+1)
		fmt.Scanln(&x)
		fmt.Printf("pair %d\ny:", pair+1)
		fmt.Scanln(&y)
		if x >= size && y >= size {
			log.Fatal("supplied coordinate is out of bounds for size ", size)
		}
		coordinate := coordinate{
			x: x,
			y: y,
		}
		coordinates = append(coordinates, coordinate)
	}
	return coordinates
}

func countAlive(board Map, coordinate coordinate) int {
	aliveNeighbours := 0
	// Looking up
	if coordinate.x > 0 && board.board[coordinate.x-1][coordinate.y] == ALIVE {
		aliveNeighbours++
	}
	// Looking down
	if coordinate.x < board.size-1 && board.board[coordinate.x+1][coordinate.y] == ALIVE {
		aliveNeighbours++
	}
	// Looking left
	if coordinate.y > 0 && board.board[coordinate.x][coordinate.y-1] == ALIVE {
		aliveNeighbours++
	}
	// Looking right
	if coordinate.y < board.size-1 && board.board[coordinate.x][coordinate.y+1] == ALIVE {
		aliveNeighbours++
	}
	// Looking top left
	if coordinate.x > 0 && coordinate.y > 0 && board.board[coordinate.x-1][coordinate.y-1] == ALIVE {
		aliveNeighbours++
	}
	// Looking top right
	if coordinate.x > 0 && coordinate.y < board.size-1 && board.board[coordinate.x-1][coordinate.y+1] == ALIVE {
		aliveNeighbours++
	}
	// Looking bottom left
	if coordinate.x < board.size-1 && coordinate.y > 0 && board.board[coordinate.x+1][coordinate.y-1] == ALIVE {
		aliveNeighbours++
	}
	// Looking bottom right
	if coordinate.x < board.size-1 && coordinate.y < board.size-1 && board.board[coordinate.x+1][coordinate.y+1] == ALIVE {
		aliveNeighbours++
	}
	return aliveNeighbours
}

func evaluateCell(board Map, coordinate coordinate) string {
	aliveNeighbours := countAlive(board, coordinate)
	if board.board[coordinate.x][coordinate.y] == ALIVE {
		// The following rules can be put into one line, but it's clearer to have them separate for theoretical purposes
		// 1. Underpopulation rule
		if aliveNeighbours < 2 {
			return DEAD
		}
		// 2. Lives to next generation rule
		if aliveNeighbours == 2 || aliveNeighbours == 3 {
			return ALIVE
		}
		// 3. Overpopulation rule
		if aliveNeighbours > 3 {
			return DEAD
		}
	} else {
		// 4. Reproduction rule
		if aliveNeighbours == 3 {
			return ALIVE
		}
	}
	return DEAD
}

func hasAlive(board Map) bool {
	for row := 0; row < board.size; row++ {
		for col := 0; col < board.size; col++ {
			if board.board[row][col] == ALIVE {
				return true
			}
		}
	}
	return false
}

func (board Map) deepCopy() Map {
	newBoard := make([][]string, board.size)
	for row := range newBoard {
		newBoard[row] = make([]string, board.size)
	}
	for row := 0; row < board.size; row++ {
		for col := 0; col < board.size; col++ {
			newBoard[row][col] = board.board[row][col]
		}
	}
	return Map{
		size:  board.size,
		board: newBoard,
	}
}

func run(board Map, coordinates []coordinate) {
	// Initialize board with ALIVE cells
	for _, coordinate := range coordinates {
		board.board[coordinate.x][coordinate.y] = ALIVE
	}
	board.PrintMap()

	// Run main loop
	epoch := 0
	for {
		fmt.Println("Epoch: ", epoch+1)
		// 1. Get Snapshot of current board
		snapshot := board.deepCopy()

		// 2. Create new board based on Conway rules on current board
		for row := 0; row < snapshot.size; row++ {
			for col := 0; col < snapshot.size; col++ {
				coordinate := coordinate{
					x: row,
					y: col,
				}
				board.board[row][col] = evaluateCell(snapshot, coordinate)
			}
		}

		// 3. Check if new board has ALIVE cells
		if hasAlive(board) {
			// If yes, print board and repeat
			board.PrintMap()
			epoch++
		} else {
			// 4. Otherwise, Terminate program
			board.PrintMap()
			fmt.Println("Simulation over.")
			os.Exit(0)
		}
		time.Sleep(1 * time.Second)
	}

}

func main() {
	// Initialize conway map
	fmt.Println("Square size of map?")
	var size int
	if _, err := fmt.Scanln(&size); err != nil {
		log.Fatal("unable to read size\n", err)
	}
	board := newMap(size)
	board.PrintMap()

	// Get coordinate list as input
	var nCoordinates int
	fmt.Println("How many coordinate pairs?")
	if _, err := fmt.Scanln(&nCoordinates); err != nil {
		log.Fatal("unable to read nCoordinates\n", err)
	}
	coordinates := readCoordinates(nCoordinates, size)
	fmt.Println(coordinates)

	// Run the simulation
	run(board, coordinates)
}
