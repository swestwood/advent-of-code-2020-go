package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func load() []string {
	file, err := os.Open("./eighteen-input.txt")
	equations := make([]string, 0)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		// Example:
		line := scanner.Text()
		if err != nil {
			log.Fatal(err)
		}
		equations = append(equations, line)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return equations
}

const (
	ADD      = 1
	MULTIPLY = 2
)

type parens struct {
	start    int
	end      int
	solution int
}

func solve(equation string) int {
	// find all groups of parens and solve them inside-out (for each end paren that is found)
	groups := make([]parens, 0)
	solved := make(map[int]parens, 0) // indexed by index of start paren
	for i, ch := range equation {
		if ch == '(' {
			groups = append(groups, parens{start: i, end: -1, solution: -1})
		} else if ch == ')' {
			// assign it to the last parens missing a closing paren, and solve this group.
			// because we're processing end parens as we find them, this group is solvable
			// (it either contains parens that are already solved, or contains no parens)
			for j := len(groups) - 1; j >= 0; j-- {
				if groups[j].end == -1 {
					groups[j].end = i
					// +1 to strip out parens, end is non-inclusive
					// solution := solvePartialRoundOne(equation, groups[j].start+1, groups[j].end, solved)
					solution := solvePartialRoundTwo(equation, groups[j].start+1, groups[j].end, solved)
					// fmt.Println(equation[groups[j].start:groups[j].end+1], "solution", solution)
					groups[j].solution = solution
					solved[groups[j].start] = groups[j]
					break
				}
			}
		}
	}
	// the final solve
	return solvePartialRoundTwo(equation, 0, len(equation), solved)
}

func doOperation(operation int, num1 int, num2 int) int {
	if operation == ADD {
		return num1 + num2
	}
	return num1 * num2
}

func solvePartialRoundOne(equation string, start int, end int, solved map[int]parens) int {
	// Solve an equation where all paren indices have corresponding values in the solved map,
	// with already-computed solutions. This can just add and multiply as it goes along encountering
	// operators. When a paren is found, look up the solution for that paren group in solved,
	// then advance past the end of that paren group.
	operation := ADD
	result := 0
	numInProgress := ""
	for i := start; i < end; i++ {
		ch := string(equation[i])

		// see if it's a number
		_, err := strconv.Atoi(ch)
		if err == nil {
			// build up the num in progress
			// note - this is actually overkill for round 1 bc everything is single digits, but oh well
			numInProgress += ch
			endOfNum := false
			// do the operation if this is the last digit in this number
			if i+1 >= end {
				endOfNum = true
			} else {
				_, err2 := strconv.Atoi(string(equation[i+1]))
				if err2 != nil {
					endOfNum = true
				}
			}
			if endOfNum {
				fullNum, _ := strconv.Atoi(numInProgress)
				result = doOperation(operation, result, fullNum)
				numInProgress = "" // reset for the next one
			}
		} else {
			switch ch {
			case "+":
				operation = ADD
			case "*":
				operation = MULTIPLY
			case "(":
				parens, ok := solved[i]
				if !ok {
					panic("paren not actually solved")
				}
				num := parens.solution
				result = doOperation(operation, result, num)
				i = parens.end // advance to the end
			default:
				// hopefully just spaces
				if ch != " " {
					panic("did not expect" + ch)
				}
			}
		}
	}
	return result
}

func solvePartialRoundTwo(equation string, start int, end int, solved map[int]parens) int {
	// Solve an equation where all paren indices have corresponding values in the solved map,
	// with already-computed solutions
	numInProgress := ""
	// Create a flat equation without parens, like: 1 + 2 * 3 + 4 * 5 + 6
	flatEquation := ""

	for i := start; i < end; i++ {
		ch := string(equation[i])

		// see if it's a number
		_, err := strconv.Atoi(ch)
		if err == nil {
			// build up the num in progress
			numInProgress += ch
			endOfNum := false
			// do the operation if this is the last digit in this number
			if i+1 >= end {
				endOfNum = true
			} else {
				_, err2 := strconv.Atoi(string(equation[i+1]))
				if err2 != nil {
					endOfNum = true
				}
			}
			if endOfNum {
				flatEquation += numInProgress
				numInProgress = "" // reset for the next one
			}
		} else {
			switch ch {
			case "(":
				parens, ok := solved[i]
				if !ok {
					fmt.Println(solved)
					fmt.Println(i)
					panic("paren not actually solved")
				}
				flatEquation += strconv.Itoa(parens.solution)
				i = parens.end // advance to the end
			default:
				// hopefully just spaces
				if ch != " " && ch != "+" && ch != "*" {
					panic("did not expect" + ch)
				}
				flatEquation += ch
			}
		}
	}
	return solveFlat(flatEquation)
}

func solveFlat(flat string) int {
	// Solves equations without parens like 6 + 9 * 8 + 6, which
	// is easy because the addition parts can be solved first, then
	// multiplied together
	toMultiply := strings.Split(flat, " * ")
	result := 1
	for _, s := range toMultiply {
		toAdd := strings.Split(s, " + ")
		sum := 0
		for _, a := range toAdd {
			num, err := strconv.Atoi(a)
			if err != nil {
				panic(err)
			}
			sum += num
		}
		result *= sum
	}
	return result
}

func main() {
	equations := load()
	totalSum := 0

	for _, equation := range equations {
		// fmt.Println(equation)
		result := solve(equation)
		// fmt.Println("result", result)
		totalSum += result
	}
	fmt.Println("total sum", totalSum)
}
