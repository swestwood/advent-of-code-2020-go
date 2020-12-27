package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func loadGrid() [][]string {
	file, err := os.Open("./eleven-input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	grid := make([][]string, 0)
	i := 0
	for scanner.Scan() {
		// Example: #.LL.
		line := scanner.Text()
		if err != nil {
			log.Fatal(err)
		}
		grid = append(grid, make([]string, 0))
		for _, r := range line {
			grid[i] = append(grid[i], string(r))
		}
		i++
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return grid
}

func numOccupiedNearby(rowI int, colI int, grid [][]string) (numNearby int) {
	for i := rowI - 1; i <= rowI+1; i++ {
		for j := colI - 1; j <= colI+1; j++ {
			if i < 0 || i >= len(grid) || j < 0 || j >= len(grid[i]) ||
				(i == rowI && j == colI) {
				continue
			}
			if grid[i][j] == "#" {
				numNearby++
			}
		}
	}
	return
}

func simulateRoundOne(prev [][]string) [][]string {
	future := make([][]string, 0)
	for rowI, row := range prev {
		future = append(future, make([]string, 0))
		for colI, val := range row {
			next := val
			occupiedNearby := numOccupiedNearby(rowI, colI, prev)
			if val == "L" && occupiedNearby == 0 {
				next = "#"
			} else if next == "#" && occupiedNearby >= 4 {
				next = "L"
			}
			future[rowI] = append(future[rowI], next)
		}
	}
	return future
}

func occupiedInDirection(rowI int, colI int, rowDir int, colDir int, grid [][]string) bool {
	for {
		rowI += rowDir
		colI += colDir
		if rowI < 0 || colI < 0 || rowI >= len(grid) || colI >= len(grid[0]) {
			return false // Ran off the edge of the map
		}
		val := grid[rowI][colI]
		if val == "L" {
			return false // Can only see an empty seat
		}
		if val == "#" {
			return true // Hit an occupied seat
		}
		// Keep going, it's just the floor
	}
}

func numOccupiedVisible(rowI int, colI int, grid [][]string) (numVisible int) {
	for rowDir := -1; rowDir <= 1; rowDir++ {
		for colDir := -1; colDir <= 1; colDir++ {
			if rowDir == 0 && colDir == 0 {
				continue
			}
			if occupiedInDirection(rowI, colI, rowDir, colDir, grid) {
				numVisible++
			}
		}
	}
	return
}

func simulateRoundTwo(prev [][]string) [][]string {
	future := make([][]string, 0)
	for rowI, row := range prev {
		future = append(future, make([]string, 0))
		for colI, val := range row {
			next := val
			occupiedNearby := numOccupiedVisible(rowI, colI, prev)
			if val == "L" && occupiedNearby == 0 {
				next = "#"
			} else if next == "#" && occupiedNearby >= 5 {
				next = "L"
			}
			future[rowI] = append(future[rowI], next)
		}
	}
	return future
}

func same(one [][]string, two [][]string) bool {
	// assume they're not empty
	if len(one) != len(two) || len(one[0]) != len(two[0]) {
		return false
	}
	for row := 0; row < len(one); row++ {
		for col := 0; col < len(one[0]); col++ {
			if one[row][col] != two[row][col] {
				return false
			}
		}
	}
	return true
}

func print(grid [][]string) {
	for _, row := range grid {
		fmt.Println(row)
	}
	fmt.Println()
}

func totalOccupied(grid [][]string) (total int) {
	for _, row := range grid {
		for _, val := range row {
			if val == "#" {
				total++
			}
		}
	}
	return
}

func main() {
	grid := loadGrid()
	loops := 0
	for {
		nextGrid := simulateRoundTwo(grid) // simulateRoundOne
		loops++
		if same(grid, nextGrid) {
			break
		}
		grid = nextGrid
	}

	fmt.Println("Same after", loops, "iterations with", totalOccupied(grid), "occupied")
}
