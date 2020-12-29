package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type command struct {
	mask  string
	index int
	value int
}

func loadProgram() []command {
	file, err := os.Open("./fourteen-input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	commands := make([]command, 0)
	for scanner.Scan() {
		// Example: mask = XXXXXXXXXXXXXXXXXXXXXXXXXXXXX1XXXX0X
		//          mem[8] = 11
		line := scanner.Text()
		if err != nil {
			log.Fatal(err)
		}
		parts := strings.Split(line, " = ")
		if parts[0] == "mask" {
			commands = append(commands, command{mask: parts[1]})
		} else {
			indexStr := parts[0][4:]
			index, _ := strconv.Atoi(indexStr[:len(indexStr)-1])
			value, _ := strconv.Atoi(parts[1])
			commands = append(commands, command{index: index, value: value})
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return commands
}

func roundOne() {
	commands := loadProgram()
	memory := make(map[int]int)
	mask := ""
	for _, command := range commands {
		if len(command.mask) > 0 {
			mask = command.mask
		} else {
			value := command.value
			for i, maskRune := range mask {
				if maskRune != 'X' {
					maskBit, _ := strconv.Atoi(string(maskRune))
					if maskBit == 0 {
						// Clear this bit
						value &= ^(1 << (35 - i))
					} else {
						value |= (1 << (35 - i))
					}
				}
			}
			memory[command.index] = value
		}
	}
	sum := 0
	for _, contents := range memory {
		sum += contents
	}
	fmt.Println("sum of memory is", sum)
}

// Round 2
type Mask struct {
	root     int
	floating []int // a list of all indexes that are floating, where 0 is rightmost
}

func setBit(num int, bitI int, bit int) int {
	if bit == 0 {
		// Clear this bit
		num &= ^(1 << bitI)
	} else {
		num |= (1 << bitI)
	}
	return num
}

func main() {
	commands := loadProgram()
	memory := make(map[int]int)
	mask := Mask{0, make([]int, 0)}
	for _, command := range commands {
		if len(command.mask) > 0 {
			// Parse the mask into a Mask struct with root and floating
			floating := make([]int, 0)
			root := 0
			for i, maskRune := range command.mask {
				bitI := 35 - i
				if maskRune == 'X' {
					floating = append(floating, bitI)
				} else {
					maskBit, _ := strconv.Atoi(string(maskRune))
					if maskBit == 1 {
						// Set this bit to 1, all other bits are 0 by default
						root |= (1 << bitI)
					}
				}
			}
			mask = Mask{root, floating}
		} else {
			indexes := make([]int, 0)
			// overwrite all 1s in the mask in the index
			indexes = append(indexes, command.index|mask.root)
			for i := 0; i < len(mask.floating); i++ {
				// Expand all current saved indexes by flipping the bit to 1 and to 0
				newIndexes := make([]int, 0)
				for _, num := range indexes {
					// One of these is already in the list and one isn't, but just
					// make a new list to avoid modifying indexes while iterating
					newIndexes = append(newIndexes, setBit(num, mask.floating[i], 0))
					newIndexes = append(newIndexes, setBit(num, mask.floating[i], 1))
				}
				indexes = newIndexes
			}
			// write the value to all masked indices
			for _, index := range indexes {
				memory[index] = command.value
			}
		}
	}
	sum := 0
	for _, contents := range memory {
		sum += contents
	}
	fmt.Println("sum of memory is", sum)
}
