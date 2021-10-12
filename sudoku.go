package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

var possibleNumbers [10][10]string  // Keeps track of all the numbers that could possibly fit in any particular sqyuare
var cellsToCheck [10][10][25][2]int // For each cell, gives (x,y) co-ordnates for the cells that need to be checked for eliminating numbers

var grid [10][10]int // This is the working grid

type stack struct {
	lock sync.Mutex // you don't have to do this if you don't want thread safety
	s    []int
}

var s *stack
var stackOperations, stackDepth, maxStackDepth int

const puzzleNumber = 5 // Which of the preset puzzles will we do?

func NewStack() *stack {
	stackOperations = 0
	stackDepth = 0
	return &stack{sync.Mutex{}, make([]int, 0)}
}

func (s *stack) Push(v int) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.s = append(s.s, v)
}

func (s *stack) Pop() (int, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	l := len(s.s)
	if l == 0 {
		return 0, errors.New("emptystack")
	}

	res := s.s[l-1]
	s.s = s.s[:l-1]
	return res, nil
}

func pushGrid() {

	for row := 1; row < 10; row++ {
		for col := 1; col < 10; col++ {
			s.Push(grid[row][col])
		}
	}
	stackOperations++
	stackDepth++
	if stackDepth > maxStackDepth {
		maxStackDepth = stackDepth
	}
}

func popGrid() {

	for row := 9; row > 0; row-- {
		for col := 9; col > 0; col-- {
			grid[row][col], _ = s.Pop()
		}
	}
	stackOperations++
	stackDepth--
}

func calculateWhichCellsToCheck() {

	var index int

	// calculate for each cell which are the 24 other cells that need to be checked for eliminating numbers
	for bigRow := 1; bigRow < 4; bigRow++ {
		for bigCol := 1; bigCol < 4; bigCol++ {
			for littleRow := 1; littleRow < 4; littleRow++ {
				for littleCol := 1; littleCol < 4; littleCol++ {
					rowNum := (bigRow-1)*3 + littleRow
					colNum := (bigCol-1)*3 + littleCol
					index = 0

					//					for i := 1; i < 10; i++ {
					//						for j := 1; j < 10; j++ {
					//							grid[i][j] = 0
					//						}
					//					}

					//					grid[rowNum][colNum] = 9

					// Horizontals first - need to check this row in every column except the one we're in
					for col := 1; col < 10; col++ {
						if col != colNum {
							cellsToCheck[rowNum][colNum][index][0] = rowNum
							cellsToCheck[rowNum][colNum][index][1] = col
							//							grid[rowNum][col] = 1
							index++
						}
					}
					// Then verticals - check this column in every row except the one we're in
					for row := 1; row < 10; row++ {
						if row != rowNum {
							cellsToCheck[rowNum][colNum][index][0] = row
							cellsToCheck[rowNum][colNum][index][1] = colNum
							//							grid[row][colNum] = 2
							index++
						}
					}

					// Then all the squares in our block of 9 except the one we are looking at
					for localRow := 1; localRow < 4; localRow++ {
						for localCol := 1; localCol < 4; localCol++ {
							comparatorRow := (bigRow-1)*3 + localRow
							comparatorCol := (bigCol-1)*3 + localCol
							if (comparatorCol == colNum) && (comparatorRow == rowNum) {
								// don't check the cell you are
							} else {
								cellsToCheck[rowNum][colNum][index][0] = comparatorRow
								cellsToCheck[rowNum][colNum][index][1] = comparatorCol
								//								grid[comparatorRow][comparatorCol] = 3
								index++
							}

						}
					}
					//					printGrid("Test")
					if index != 24 {
						fmt.Printf("Index is %d, big (%d, %d) little (%d, %d)\n", index, bigRow, bigCol, littleRow, littleCol)
					}
				}
			}
		}
	}

}

func populateInitialGrid(gridNum int) {
	exampleGrids := [][]string{
		{
			"068309070",
			"042000001",
			"107050600",
			"005070120",
			"700001580",
			"000030740",
			"000190205",
			"801620390",
			"900543010",
		},
		{
			"063500000",
			"027000094",
			"094000510",
			"000000000",
			"000009080",
			"000048730",
			"000450020",
			"200036108",
			"051002307",
		},
		{
			"080400002",
			"009000004",
			"730000008",
			"000913000",
			"042600000",
			"007050600",
			"050006309",
			"020300000",
			"400001000",
		},
		{
			"105900760",
			"020003000",
			"006000080",
			"004000510",
			"800040000",
			"050000800",
			"009000100",
			"000560090",
			"000020003",
		},
		{
			"090040030",
			"600750009",
			"007000600",
			"000000010",
			"230090057",
			"050000000",
			"005000800",
			"300021006",
			"070030090",
		},
		{
			"800000000",
			"003600000",
			"070090200",
			"050007000",
			"000045700",
			"000100030",
			"001000068",
			"008500010",
			"090000400",
		},
		{
			"001000000",
			"800000020",
			"070010500",
			"400005300",
			"010070006",
			"003200080",
			"060500009",
			"004000030",
			"000009700",
		},
		{
			"005300000",
			"800000020",
			"070010500",
			"400005300",
			"010070006",
			"003200080",
			"060500009",
			"004000030",
			"000009700",
		},
	}

	gridline := exampleGrids[gridNum]

	for i := 1; i < 10; i++ {
		for j := 1; j < 10; j++ {
			if digit, err := strconv.Atoi(gridline[i-1][j-1 : j]); err == nil {
				grid[i][j] = digit
			} else {
				fmt.Printf("Error: %s fromn string %s\n", err, gridline[i][j-1:j])
			}
		}
	}

}

func printGrid(title string) {
	fmt.Println(title)
	for i := 1; i < 10; i++ {
		fmt.Print(" ")
		for j := 1; j < 10; j++ {

			if j == 4 || j == 7 {
				fmt.Print("| ")
			}

			if grid[i][j] == 0 {
				fmt.Print(" ")
			} else {
				fmt.Print(grid[i][j])
			}

			fmt.Print(" ")
		}
		fmt.Println("")

		if i == 3 || i == 6 {
			fmt.Println("-------+-------+-------")
		}
		//		fmt.Println()
	}
}

func doTheEasyOnes() (int, int, int, string, bool) {

	// returns are:
	// stillToDo int, shortestOneRow int, shortestOneCol int,shortestOnePossibilities string, impossible bool

	var stillToDo int
	var shortestSoFar, shortestOneRow, shortestOneCol int
	var shortestOnePossibilities string
	var impossible bool
	var gotOne bool

	shortestSoFar = 100
	impossible = false
	gotOne = true
	for !impossible && gotOne {
		stillToDo = 81
		gotOne = false
		for rownum := 1; rownum < 10; rownum++ {
			for colnum := 1; colnum < 10; colnum++ {
				if grid[rownum][colnum] == 0 { // If this grid slot is not alreasdy populated
					possibleNumbers[rownum][colnum] = "123456789" // Could be anything right now

					for k := 0; k < 24; k++ { // Look at alll the cells we need to check and eliminate any number that is already represented
						rowToCheck := cellsToCheck[rownum][colnum][k][0]
						colToCheck := cellsToCheck[rownum][colnum][k][1]
						knockout := strconv.Itoa(grid[rowToCheck][colToCheck])
						possibleNumbers[rownum][colnum] = strings.Replace(possibleNumbers[rownum][colnum], knockout, "", 1)
					}
					if len(possibleNumbers[rownum][colnum]) == 1 { // There is ony one number it could be - put that in
						grid[rownum][colnum], _ = strconv.Atoi(possibleNumbers[rownum][colnum])
						stillToDo--
						gotOne = true

					} else if len(possibleNumbers[rownum][colnum]) == 0 { // There are no number it can be. There is no solution
						impossible = true
						return stillToDo, shortestOneRow, shortestOneCol, shortestOnePossibilities, impossible
					} else { // There is more than one to do. Let's see if this is the shortest one so far
						if len(possibleNumbers[rownum][colnum]) < shortestSoFar {
							shortestSoFar = len(possibleNumbers[rownum][colnum])
							shortestOneRow = rownum
							shortestOneCol = colnum
							shortestOnePossibilities = possibleNumbers[rownum][colnum]
						}

					}
				} else { // The grid slot already had a number in it
					stillToDo--

				}
			}
		}
	}

	return stillToDo, shortestOneRow, shortestOneCol, shortestOnePossibilities, impossible
}

func solver() bool {

	var stillToDo int
	var shortestOneRow, shortestOneCol int
	var shortestOnePossibilities string
	var impossible bool

	stillToDo, shortestOneRow, shortestOneCol, shortestOnePossibilities, impossible = doTheEasyOnes()

	//	fmt.Printf("Still to do: %n, (%n, %n) Possibilities: %n, impossible: %n\n", stillToDo, shortestOneRow, shortestOneCol, shortestOnePossibilities, impossible)
	//	printGrid("")

	if stillToDo > 0 && !impossible {
		for i := 1; i <= len(shortestOnePossibilities); i++ { // Loop round all the possibilities for the easiest grid spot to guess
			impossible = false // Starting again - it's not impossible yet
			pushGrid()         // Save the current state of the grid. We're going to try some stuff
			grid[shortestOneRow][shortestOneCol], _ = strconv.Atoi(shortestOnePossibilities[i-1 : i])
			impossible = solver() // recurse round until this one becomes impossible
			if impossible {       // This hasn't worked - pop the grid back to where it was and go round again with the next guess
				popGrid()
			} else {
				return false // we have solved it so its not impossible
			}

			// if we get this far it must be impossible because none of the guesses we tried worked

		}
		impossible = true // We have tried every possibility and it hasn't worked so we'll say its impossible
	}
	return impossible
}

func initialise() {

	calculateWhichCellsToCheck()
	populateInitialGrid(puzzleNumber)
	s = NewStack()

}

func main() {
	fmt.Println("Sudoku")

	initialise()

	printGrid("Here is our starting grid\n")

	start := time.Now()

	if !solver() {

		printGrid("\n\nHere is the answer:\n")
	} else {
		fmt.Println("That was impossible")
	}

	elapsed := time.Since(start)
	fmt.Printf("That took %s. There were %d stack operations with a max depth of %d.", elapsed, stackOperations, maxStackDepth)

}
