package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

const tileDim = 10

type side int

const (
	TOP    side = 1
	RIGHT  side = 2
	BOTTOM side = 3
	LEFT   side = 4
)

type solution struct {
	tileGrid  map[coord]tilePermutation
	usedTiles map[int]bool // key is the tile id, value is if it is used
	dim       int          // the dim of tile grid
}

type tile struct {
	grid map[coord]string // indexed by [row][col]
	id   int
}

type tilePermutation struct {
	tileID int
	p      permutation
}

type corner struct {
	top  string
	left string
}

type coord struct {
	// 0, 0 is the top left
	row int
	col int
}

type permutation struct {
	rotations int // 0-3 inclusive
	flipHoriz bool
	flipVert  bool
}

func (t tile) String() string {
	str := "Tile " + strconv.Itoa(t.id) + "\n"
	for r := 0; r < tileDim; r++ {
		for c := 0; c < tileDim; c++ {
			str += (t.grid[coord{r, c}])
		}
		str += "\n"
	}
	return str
}

// To avoid recomputing all the permutations over and over
var cachedPermuted = make(map[tilePermutation]*map[coord]string, 0)

func (t *tile) applyCached(p permutation) map[coord]string {
	val, exists := cachedPermuted[tilePermutation{t.id, p}]
	if !exists {
		panic("doesn't exist" + string(t.id))
	}
	return *val
}

func (t *tile) apply(p permutation) map[coord]string {
	newGrid := make(map[coord]string, 0)
	for pos, val := range t.grid {
		row := pos.row
		col := pos.col
		if p.flipHoriz {
			// with a tile dim of 10, col 0 <-> 9, 1 <-> 8, 2 <-> 7, etc
			col = tileDim - 1 - col
		}
		if p.flipVert {
			row = tileDim - 1 - row
		}
		if p.rotations > 0 {
			// (0, 1) -> (1, 9) -> (9, 8) -> (8, 0) -> (0, 1)
			// (r, c) -> (c, 9 - r) -> (9 - r, 9 - c) -> (9 - c, r)
			switch p.rotations {
			case 1:
				// rotate right
				row, col = col, tileDim-1-row
			case 2:
				// rotate around
				row, col = tileDim-1-row, tileDim-1-col
			case 3:
				// rotate left
				row, col = tileDim-1-col, row
			}
		}
		newGrid[coord{row, col}] = val
	}
	return newGrid
}

var allPermutes []permutation = nil

func initPermutations() []permutation {
	if allPermutes == nil {
		permutes := make([]permutation, 0)
		for r := 0; r < 4; r++ {
			permutes = append(permutes,
				[]permutation{
					permutation{r, true, true},
					permutation{r, true, false},
					permutation{r, false, true},
					permutation{r, false, false},
				}...)
		}
		allPermutes = permutes
	}
	return allPermutes
}

func print(grid map[coord]string, dim int, asImage bool) {
	if asImage {
		for r := 1; r < dim-1; r++ {
			for c := 1; c < dim-1; c++ {
				fmt.Print(grid[coord{r, c}])
			}
			fmt.Println()
		}
	} else {
		for r := 0; r < dim; r++ {
			for c := 0; c < dim; c++ {
				fmt.Print(grid[coord{r, c}])
			}
			fmt.Println()
		}
	}
}

func load() map[int]*tile {
	file, err := os.Open("./twenty-input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var current *tile
	tiles := make(map[int]*tile, 0)
	row := 0
	for scanner.Scan() {
		line := scanner.Text()
		if err != nil {
			log.Fatal(err)
		}
		if strings.HasPrefix(line, "Tile") {
			var id int
			fmt.Sscanf(line, "Tile %d:", &id)
			current = &tile{grid: make(map[coord]string, 0), id: id}
			row = 0
			tiles[id] = current
		} else if line != "" {
			for col, val := range line {
				current.grid[coord{row, col}] = string(val)
			}
			row++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return tiles
}

func (c coord) next(dim int) coord {
	if c.col+1 < dim {
		return coord{c.row, c.col + 1}
	}
	// wrap to next row, even if it might be out of bounds row-wise
	return coord{c.row + 1, 0}
}

func asString(grid map[coord]string, s side) string {
	var b strings.Builder
	for i := 0; i < tileDim; i++ {
		row := 0
		col := 0
		switch s {
		case TOP:
			col = i
		case BOTTOM:
			row = tileDim - 1
			col = i
		case LEFT:
			row = i
		case RIGHT:
			row = i
			col = tileDim - 1
		}
		b.WriteString(grid[coord{row, col}])
	}
	return b.String()
}

func (sol *solution) getVal(outer coord, inner coord, tiles map[int]*tile) string {
	return sol.getTileGrid(outer, tiles)[inner]
}

func (sol *solution) getTileGrid(outer coord, tiles map[int]*tile) map[coord]string {
	tp := sol.tileGrid[outer]
	t := tiles[tp.tileID]
	return t.applyCached(tp.p)
}

func arrange(sol *solution, solLoc coord, tiles map[int]*tile, cache map[corner][]tilePermutation,
	bottoms map[tilePermutation]string, rights map[tilePermutation]string) (*solution, bool) {
	if solLoc.row >= sol.dim {
		return sol, true
	}
	options := make([]tilePermutation, 0)

	if solLoc.row == 0 && solLoc.col == 0 {
		panic("Needs to have first tile filled in")
	} else {
		left, top := "", ""
		if solLoc.col > 0 {
			left = rights[sol.tileGrid[coord{solLoc.row, solLoc.col - 1}]]
		}
		if solLoc.row > 0 {
			top = bottoms[sol.tileGrid[coord{solLoc.row - 1, solLoc.col}]]
		}
		// all tile permutations matching this top left
		options, _ = cache[corner{top, left}]
	}
	// fmt.Println("exploring", len(options), "options at", solLoc, options)

	for _, tp := range options {
		_, used := sol.usedTiles[tp.tileID]
		if !used {
			// check this permutation of this grid item and see if it fits
			sol.tileGrid[solLoc] = tp
			sol.usedTiles[tp.tileID] = true
			// fmt.Println(solLoc, "placed", tp.tileID)
			newSol, valid := arrange(sol, solLoc.next(sol.dim), tiles, cache, bottoms, rights)
			if valid {
				return newSol, valid
			}
			delete(sol.usedTiles, tp.tileID)
			delete(sol.tileGrid, solLoc)
		}
	}
	return sol, false
}

func addToCache(c corner, tp tilePermutation, cache map[corner][]tilePermutation) {
	if _, exists := cache[c]; !exists {
		cache[c] = make([]tilePermutation, 0)
	}
	for _, val := range cache[c] {
		if val == tp {
			fmt.Println("would be dupe!", c, tp)
			return
		}
	}
	cache[c] = append(cache[c], tp)
}

var monster = [3]string{
	"                  # ",
	"#    ##    ##    ###",
	" #  #  #  #  #  #   ",
}

func checkMonsterAt(c coord, sol *solution, tiles map[int]*tile, monsterCells map[coord]bool) bool {
	cells := make([]coord, 0)
	for rowDelta, monsterCol := range monster {
		for colDelta, monsterVal := range monsterCol {
			if monsterVal == '#' {
				imgC := coord{c.row + rowDelta, c.col + colDelta}
				imgVal, ok := sol.getValInImage(imgC, tiles)
				if !ok || imgVal != "#" {
					return false
				}
				cells = append(cells, imgC)
			}
		}
	}
	// Mark this monster since we found one
	for _, cell := range cells {
		monsterCells[cell] = true
	}
	return true
}

func (sol *solution) imageFullDim() int {
	// Compensate for stripping the borders off the tiles
	return sol.imageTileDim() * sol.dim
}

func (sol *solution) imageTileDim() int {
	// Compensate for stripping the borders off the tiles
	return tileDim - 2
}

func (sol *solution) getValInImage(c coord, tiles map[int]*tile) (string, bool) {
	// Compensate for stripping the borders off the tiles
	if c.row < 0 || c.col < 0 || c.row >= sol.imageFullDim() || c.col >= sol.imageFullDim() {
		return "", false
	}
	outer := coord{c.row / sol.imageTileDim(), c.col / sol.imageTileDim()}
	inner := coord{(c.row % sol.imageTileDim()) + 1, (c.col % sol.imageTileDim()) + 1}
	// fmt.Println("converted", c, "to", outer, inner, "with tileDim", tileDim, sol.dim, sol.imageTileDim())
	// check the dims
	return sol.getVal(outer, inner, tiles), true
}

func findMonsters(sol *solution, tiles map[int]*tile) {
	monsterCells := make(map[coord]bool, 0)
	for r := 0; r < sol.imageFullDim(); r++ {
		for c := 0; c < sol.imageFullDim(); c++ {
			checkMonsterAt(coord{r, c}, sol, tiles, monsterCells)
		}
	}
	if len(monsterCells) > 0 {
		roughness := 0
		for r := 0; r < sol.imageFullDim(); r++ {
			for c := 0; c < sol.imageFullDim(); c++ {
				_, isMonster := monsterCells[coord{r, c}]
				if !isMonster {
					imageVal, ok := sol.getValInImage(coord{r, c}, tiles)
					if ok && imageVal == "#" {
						roughness++
					}
				}
			}
		}
		fmt.Println("ROUGHNESS is", roughness, sol.tileGrid[coord{0, 0}])
	}
}

func solve(ch chan int, tiles map[int]*tile, starting tilePermutation, cache map[corner][]tilePermutation,
	bottoms map[tilePermutation]string, rights map[tilePermutation]string) {
	dim := int(math.Sqrt(float64(len(tiles))))
	sol := &solution{make(map[coord]tilePermutation, 0), make(map[int]bool, 0), dim}
	sol.tileGrid[coord{0, 0}] = starting
	sol.usedTiles[starting.tileID] = true
	sol, ok := arrange(sol, coord{0, 1}, tiles, cache, bottoms, rights)
	if ok {
		product := (sol.tileGrid[coord{0, 0}].tileID * sol.tileGrid[coord{0, sol.dim - 1}].tileID * sol.tileGrid[coord{sol.dim - 1, 0}].tileID *
			sol.tileGrid[coord{sol.dim - 1, sol.dim - 1}].tileID)
		findMonsters(sol, tiles)
		ch <- product
	}
}

func main() {
	initPermutations()
	tiles := load()
	ch := make(chan int)
	// Build up a cache of all qualifying tiles for each string to
	// make the search faster
	cache := make(map[corner][]tilePermutation, 0)
	bottoms := make(map[tilePermutation]string)
	rights := make(map[tilePermutation]string)

	for _, t := range tiles {
		for _, permute := range allPermutes {
			grid := t.apply(permute)
			top := asString(grid, TOP)
			left := asString(grid, LEFT)
			tp := tilePermutation{t.id, permute}
			cachedPermuted[tp] = &grid
			// Add both this corner, as well as the top and left
			// independently so we can place it in the left col or top row
			addToCache(corner{top, left}, tp, cache)
			addToCache(corner{top, ""}, tp, cache)
			addToCache(corner{"", left}, tp, cache)
			bottoms[tp] = asString(grid, BOTTOM)
			rights[tp] = asString(grid, RIGHT)
		}
	}

	// Search across all possible starting locations, parallelized so we don't get
	// stuck down a rabbit hole of failure
	for _, t := range tiles {
		for _, permute := range allPermutes {
			go solve(ch, tiles, tilePermutation{t.id, permute}, cache, bottoms, rights)
		}
	}

	// the solution grid has 16 permutations, so wait for all of them to come back
	for i := 0; i < len(allPermutes); i++ {
		product := <-ch
		fmt.Println(i, "received product", product)
	}
	fmt.Println("tiles", len(tiles), "permutations", len(allPermutes), "cached permuted", len(cachedPermuted))

	// solve the found solution directly
	// solve(ch, tiles, tilePermutation{1951, permutation{3, true, true}}, cache, bottoms, rights)
	// print(tiles[1951].applyCached(permutation{3, true, true}), tileDim, true)
}
