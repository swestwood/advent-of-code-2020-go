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

func main() {
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
