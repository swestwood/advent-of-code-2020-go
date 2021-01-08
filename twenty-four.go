package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

var directions = []string{"e", "se", "sw", "w", "nw", "ne"}

type coord struct {
	x int
	y int
}

func load() []string {
	file, err := os.Open("./twenty-four-input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	tilesToFlip := make([]string, 0)
	for scanner.Scan() {
		// Example: neeswseenwwswnwswswnw
		line := scanner.Text()
		if err != nil {
			log.Fatal(err)
		}
		tilesToFlip = append(tilesToFlip, line)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return tilesToFlip
}

func move(curr coord, dir string) coord {
	var next coord
	// 0, 0 is in the upper left, north and west is lower
	switch dir {
	case "e":
		next = coord{curr.x + 2, curr.y}
	case "w":
		next = coord{curr.x - 2, curr.y}
	case "se":
		next = coord{curr.x + 1, curr.y + 1}
	case "nw":
		next = coord{curr.x - 1, curr.y - 1}
	case "sw":
		next = coord{curr.x - 1, curr.y + 1}
	case "ne":
		next = coord{curr.x + 1, curr.y - 1}
	default:
		panic("unknown direction" + dir)
	}
	return next
}

const debug = false

func blackAdjacent(c coord, blackGrid map[coord]bool) int {
	count := 0

	_, isBlack := blackGrid[c]
	if debug {
		fmt.Println("COORD:", c, "is black? ", isBlack)
	}
	for _, dir := range directions {
		check := move(c, dir)
		_, neighborIsBlack := blackGrid[check]
		if debug {
			fmt.Println("\t neighbor", dir, "is", check, "\tis black?", neighborIsBlack)
		}
		if blackGrid[check] {
			count++
		}
	}
	return count
}

func update(c coord, blackGrid map[coord]bool, nextGrid map[coord]bool, processedWhites map[coord]bool) {
	adjacent := blackAdjacent(c, blackGrid)
	wasBlack, _ := blackGrid[c]
	nextIsBlack := wasBlack
	if wasBlack && adjacent == 0 || adjacent > 2 {
		nextIsBlack = false
		if debug {
			fmt.Println("->Flipped to white, c", c)
		}
	} else if wasBlack {
		if debug {
			fmt.Println("->Stayed black, c", c)
		}

	}
	if !wasBlack && adjacent == 2 {
		nextIsBlack = true
		if debug {
			fmt.Println("->Flipped to black", c)
		}
	} else if !wasBlack {
		if debug {
			fmt.Println("->Stayed white, c", c)
		}
	}

	_, alreadyExists := nextGrid[c]
	if alreadyExists {
		panic("already processed this one")
	}
	_, alreadyExists = processedWhites[c]
	if alreadyExists {
		panic("already processed this one")
	}
	if nextIsBlack {
		nextGrid[c] = nextIsBlack
		if debug {
			fmt.Println("\t ->->entered", c, "as black in nextGrid")
		}
	}
	if !wasBlack {
		processedWhites[c] = true
	}
}

func main() {
	tilesToFlip := load()
	// Odd rows only have odd cols filled in
	// Even cols only have even cols filled in
	blackGrid := make(map[coord]bool, 0) // false if white, true if black
	center := coord{10, 10}
	numBlack := 0
	for _, flip := range tilesToFlip {
		curr := center
		if debug {
			fmt.Println("--")
		}
		for len(flip) > 0 {
			for _, dir := range directions {
				if strings.HasPrefix(flip, dir) {
					curr = move(curr, dir)
					flip = flip[len(dir):]
				}
			}
		}
		oldIsBlack, _ := blackGrid[curr]
		if oldIsBlack {
			// flipping from white=true to white=false aka black
			numBlack--
			delete(blackGrid, curr)
		} else {
			numBlack++
			blackGrid[curr] = true
		}
	}
	fmt.Println("num black at the end", numBlack, len(blackGrid))

	// Round 2
	// go look at everything currently black and the adjacent tiles to see if they should flip
	totalDays := 100
	fmt.Println("Start:", blackGrid)
	for days := 0; days < totalDays; days++ {
		nextGrid := make(map[coord]bool, 0)
		processedWhites := make(map[coord]bool, 0)
		for c := range blackGrid {
			// fmt.Println("--")
			update(c, blackGrid, nextGrid, processedWhites)
			// Update all white neighbors of black tiles if needed, since they aren't in the blackGrid
			for _, dir := range directions {
				neighbor := move(c, dir)
				_, done := processedWhites[neighbor]
				_, inGrid := blackGrid[neighbor]
				if !inGrid && !done {
					update(neighbor, blackGrid, nextGrid, processedWhites)
				}
			}
		}
		blackGrid = nextGrid

	}
	fmt.Println("After", totalDays, "days,", len(blackGrid), "tiles are black")
}
