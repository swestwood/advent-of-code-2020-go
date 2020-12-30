package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type world struct {
	grid      [][][][]string
	zeroIndex int
}

const kSize = 40        // starting size, set manually to encompass the max final size after 6 cycles
const kStartOffset = 15 // room to grow on each side

func (w world) String() string {
	str := ""
	g := w.grid
	for i := 0; i < len(g); i++ {
		for j := 0; j < len(g[0]); j++ {
			// note that w here is y in our grid, not sure if it matches the sample anymore in 4D
			str += fmt.Sprintf("z=%d, y=%d\n", i-w.zeroIndex, j-w.zeroIndex)
			for k := 0; k < len(g[0][0]); k++ {
				str += strings.Join(g[i][j][k], "") + "\n"
			}
		}
		str += "\n"
	}
	return str
}

func makeWorld(z, y, x, u, zero int) world {
	// [z][y][x][u]
	g := make([][][][]string, 0)
	for i := 0; i < z; i++ {
		g = append(g, make([][][]string, 0))
		for j := 0; j < y; j++ {
			g[i] = append(g[i], make([][]string, 0))
			for k := 0; k < x; k++ {
				g[i][j] = append(g[i][j], make([]string, 0))
				for h := 0; h < u; h++ {
					g[i][j][k] = append(g[i][j][k], ".")
				}
			}
		}
	}
	return world{g, zero}
}

func load() world {
	file, err := os.Open("./seventeen-input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	w := makeWorld(kSize, kSize, kSize, kSize, kStartOffset)
	j, h := 0, 0
	for scanner.Scan() {
		line := scanner.Text()
		// Example .#.
		if err != nil {
			log.Fatal(err)
		}
		for k, ch := range line {
			if string(ch) == "#" {
				w.grid[0+kStartOffset][j+kStartOffset][k+kStartOffset][h+kStartOffset] = "#"
			}
		}
		j++
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return w
}

func isActive(i, j, k, h int, w world) bool {
	// range protected check
	g := w.grid
	if i < 0 || j < 0 || k < 0 || h < 0 ||
		i >= len(g) || j >= len(g[0]) || k >= len(g[0][0]) || h >= len(g[0][0][0]) {
		return false
	}
	return g[i][j][k][h] == "#"
}

func becomesActive(i, j, k, h int, w1 world) bool {
	neighbors := 0
	for di := -1; di <= 1; di++ {
		for dj := -1; dj <= 1; dj++ {
			for dk := -1; dk <= 1; dk++ {
				for dh := -1; dh <= 1; dh++ {
					if di == 0 && dj == 0 && dk == 0 && dh == 0 {
						continue
					}
					if isActive(i+di, j+dj, k+dk, h+dh, w1) {
						neighbors++
					}
				}
			}
		}
	}
	if w1.grid[i][j][k][h] == "#" {
		return neighbors == 2 || neighbors == 3
	}
	return neighbors == 3
}

func cycle(w1 world) world {
	w2 := makeWorld(len(w1.grid), len(w1.grid[0]), len(w1.grid[0][0]), len(w1.grid[0][0][0]), w1.zeroIndex)
	g2 := w2.grid
	for i := 0; i < len(g2); i++ {
		for j := 0; j < len(g2[0]); j++ {
			for k := 0; k < len(g2[0][0]); k++ {
				for h := 0; h < len(g2[0][0]); h++ {
					if becomesActive(i, j, k, h, w1) {
						g2[i][j][k][h] = "#"
					}
				}
			}
		}
	}
	return w2
}

func main() {
	w := load()
	for round := 1; round <= 6; round++ {
		// fmt.Println("\n\n\n======== ROUND ", round)
		w = cycle(w)
		// fmt.Println(w)
	}

	// count active
	g := w.grid
	active := 0
	for i := 0; i < len(g); i++ {
		for j := 0; j < len(g[0]); j++ {
			for k := 0; k < len(g[0][0]); k++ {
				for h := 0; h < len(g[0][0][0]); h++ {
					if isActive(i, j, k, h, w) {
						active++
					}
				}
			}
		}
	}
	fmt.Println("num active:", active)
}
