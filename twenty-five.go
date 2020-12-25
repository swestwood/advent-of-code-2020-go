package main

import "fmt"

func loop(value int, subject int) int {
	value = value * subject
	value = value % 20201227
	return value
}

func main() {
	// card public key = 7 looped X times
	// door public key = 7 looped Y times
	loopCount := 0
	value := 1
	subject := 7
	for {
		value = loop(value, subject)
		loopCount++
		// door public 4126658
		// card public 10604480
		if value == 10604480 {
			fmt.Println("Loop count", loopCount)
			break
		}
	}
	// door loop 9709101
	// card loop 1568743
	doorPublic := 4126658
	cardLoop := 1568743
	loopCount = 0
	value = 1
	for {
		value = loop(value, doorPublic)
		loopCount++
		if loopCount == cardLoop {
			fmt.Println("Value count", value)
			return
		}
	}
}
