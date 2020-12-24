package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func initializeMap() [][]bool {
	mapChunk := make([][]bool, 0)
	file, err := os.Open("./three-input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	currentRow := 0
	for scanner.Scan() {
		// Example: ....#...##.#.........#....#....
		line := scanner.Text()
		if err != nil {
			log.Fatal(err)
		}
		mapChunk = append(mapChunk, make([]bool, len(line))) // Values default to false
		for i, char := range line {
			if string(char) == "#" {
				mapChunk[currentRow][i] = true
			}
		}
		currentRow++
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return mapChunk
}

func getTreesHit(colChange int, rowChange int, mapChunk [][]bool) int {
	numTreesHit := 0
	row := 0
	col := 0
	for {
		if row >= len(mapChunk) {
			break
		}
		if mapChunk[row][col%len(mapChunk[row])] {
			numTreesHit++
		}
		col += colChange
		row += rowChange
	}
	return numTreesHit
}

func roundOne() {
	mapChunk := initializeMap()
	fmt.Println(getTreesHit(3, 1, mapChunk))
}
func roundTwo() {
	mapChunk := initializeMap()
	oneOne := getTreesHit(1, 1, mapChunk)
	threeOne := getTreesHit(3, 1, mapChunk)
	fiveOne := getTreesHit(5, 1, mapChunk)
	sevenOne := getTreesHit(7, 1, mapChunk)
	oneTwo := getTreesHit(1, 2, mapChunk)
	fmt.Println(oneOne * threeOne * fiveOne * sevenOne * oneTwo)
}

func main() {
	roundTwo()
}
