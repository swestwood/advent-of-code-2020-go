package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Instruction struct {
	operation string
	argument  int
}

func loadInstructions() []Instruction {
	file, err := os.Open("./eight-input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	instructions := make([]Instruction, 0)
	for scanner.Scan() {
		// Example: acc +1
		line := scanner.Text()
		if err != nil {
			log.Fatal(err)
		}
		pieces := strings.Fields(line)
		argument, _ := strconv.Atoi(pieces[1])
		instructions = append(instructions, Instruction{pieces[0], argument})
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return instructions
}

func execute(instructions []Instruction) (int, bool) {
	accumulator := 0
	currentLine := 0
	linesSeen := make(map[int]bool)
	didFinish := true
	for {
		if currentLine < 0 || currentLine > len(instructions) {
			// Illegal stuff
			didFinish = false
			break
		}
		if currentLine == len(instructions) {
			// Finished the program
			break
		}
		instruction := instructions[currentLine]
		if _, ok := linesSeen[currentLine]; ok {
			didFinish = false
			break
		}
		linesSeen[currentLine] = true
		switch instruction.operation {
		case "acc":
			accumulator += instruction.argument
			currentLine++
		case "jmp":
			currentLine += instruction.argument
		case "nop":
			currentLine++
		default:
			fmt.Println("Something went wrong, unrecognized operation", instruction)
		}
	}
	return accumulator, didFinish
}

func main() {
	// Round 1
	instructions := loadInstructions()
	initialValue, completed := execute(instructions)
	fmt.Println("The bootcode value is", initialValue, completed)

	// Round 2
	for i := range instructions {
		switch instructions[i].operation {
		case "jmp":
			instructions[i].operation = "nop"
			// fmt.Println("changed jmp", instructions[i].argument, "to", instructions[i].operation, instructions[i].argument)
			value, didFinish := execute(instructions)
			if didFinish {
				fmt.Println("The program finished successfully with value", value)
				return
			}
			instructions[i].operation = "jmp" // put it back
		case "nop":
			instructions[i].operation = "jmp"
			value, didFinish := execute(instructions)
			if didFinish {
				fmt.Println("The program finished successfully with value", value)
				return
			}
			instructions[i].operation = "nop" // put it back
		}
	}

}
