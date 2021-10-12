package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var possibleNumbers [10][10]string  // Keeps track of all the numbers that could possibly fit in any particular sqyuare
var cellsToCheck [10][10][25][2]int // For each cell, gives (x,y) co-ordnates for the cells that need to be checked for eliminating numbers

var grid [10][10]int // This is the working grid

var s *stack
var stackOperations, stackDepth, maxStackDepth int

const puzzleNumber = 6 // Which of the preset puzzles will we do?

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

	var gridline [10]string

	switch gridNum {
	case 1:
		{
			gridline[1] = "068309070"
			gridline[2] = "042000001"
			gridline[3] = "107050600"
			gridline[4] = "005070120"
			gridline[5] = "700001580"
			gridline[6] = "000030740"
			gridline[7] = "000190205"
			gridline[8] = "801620390"
			gridline[9] = "900543010"

		}
	case 2:
		{
			gridline[1] = "063500000"
			gridline[2] = "027000094"
			gridline[3] = "094000510"
			gridline[4] = "000000000"
			gridline[5] = "000009080"
			gridline[6] = "000048730"
			gridline[7] = "000450020"
			gridline[8] = "200036108"
			gridline[9] = "051002307"
		}
	case 3:
		{
			gridline[1] = "080400002"
			gridline[2] = "009000004"
			gridline[3] = "730000008"
			gridline[4] = "000913000"
			gridline[5] = "042600000"
			gridline[6] = "007050600"
			gridline[7] = "050006309"
			gridline[8] = "020300000"
			gridline[9] = "400001000"
		}
	case 4:
		{
			gridline[1] = "105900760"
			gridline[2] = "020003000"
			gridline[3] = "006000080"
			gridline[4] = "004000510"
			gridline[5] = "800040000"
			gridline[6] = "050000800"
			gridline[7] = "009000100"
			gridline[8] = "000560090"
			gridline[9] = "000020003"
		}
	case 5:
		{
			gridline[1] = "090040030"
			gridline[2] = "600750009"
			gridline[3] = "007000600"
			gridline[4] = "000000010"
			gridline[5] = "230090057"
			gridline[6] = "050000000"
			gridline[7] = "005000800"
			gridline[8] = "300021006"
			gridline[9] = "070030090"
		}
	case 6:
		{
			gridline[1] = "800000000"
			gridline[2] = "003600000"
			gridline[3] = "070090200"
			gridline[4] = "050007000"
			gridline[5] = "000045700"
			gridline[6] = "000100030"
			gridline[7] = "001000068"
			gridline[8] = "008500010"
			gridline[9] = "090000400"
		}
	case 7:
		{
			gridline[1] = "001000000"
			gridline[2] = "800000020"
			gridline[3] = "070010500"
			gridline[4] = "400005300"
			gridline[5] = "010070006"
			gridline[6] = "003200080"
			gridline[7] = "060500009"
			gridline[8] = "004000030"
			gridline[9] = "000009700"
		}
	case 8:
		{
			gridline[1] = "005300000"
			gridline[2] = "800000020"
			gridline[3] = "070010500"
			gridline[4] = "400005300"
			gridline[5] = "010070006"
			gridline[6] = "003200080"
			gridline[7] = "060500009"
			gridline[8] = "004000030"
			gridline[9] = "000009700"
		}
	}

	for i := 1; i < 10; i++ {
		for j := 1; j < 10; j++ {
			if digit, err := strconv.Atoi(gridline[i][j-1 : j]); err == nil {
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
