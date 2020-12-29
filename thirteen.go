package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Bus struct {
	num    int
	offset int
}

func roundOne() {
	start := 1000677
	busStr := "29,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,41,x,x,x,x,x,x,x,x,x,661,x,x,x,x,x,x,x,x,x,x,x,x,13,17,x,x,x,x,x,x,x,x,23,x,x,x,x,x,x,x,521,x,x,x,x,x,37,x,x,x,x,x,x,x,x,x,x,x,x,19"
	allBuses := strings.Split(busStr, ",")
	fmt.Println(allBuses)
	earliest := start
	// Round 1
	for {
		for _, bus := range allBuses {
			if string(bus) == "x" {
				continue
			}
			busNum, _ := strconv.Atoi(string(bus))
			fmt.Println(busNum)
			if earliest%busNum == 0 {
				fmt.Println("The earliest bus is", busNum, "at", earliest)
				fmt.Println("You waited", earliest-start)
				fmt.Println((earliest - start) * busNum)
				return
			}
		}
		earliest++
	}
}

func main() {
	busStr := "29,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,41,x,x,x,x,x,x,x,x,x,661,x,x,x,x,x,x,x,x,x,x,x,x,13,17,x,x,x,x,x,x,x,x,23,x,x,x,x,x,x,x,521,x,x,x,x,x,37,x,x,x,x,x,x,x,x,x,x,x,x,19"
	allBuses := strings.Split(busStr, ",")
	buses := make([]Bus, 0)

	// Round 2
	for i, bus := range allBuses {
		if string(bus) == "x" {
			continue
		}
		busNum, _ := strconv.Atoi(string(bus))
		buses = append(buses, Bus{busNum, i})
	}
	// [{29 0} {41 19} {661 29} {13 42} {17 43} {23 52} {521 60} {37 66} {19 79}]
	candidate := buses[0].num
	increment := buses[0].num
	for _, bus := range buses[1:] {
		// jump up in increments until we find the first number x that matches
		// (x + bus.offset) % bus.num == 0.
		for {
			if (candidate+bus.offset)%bus.num == 0 {
				break
			}
			candidate += increment
		}
		// then, figure out the next number where that happens and set increment to that so we can jump up by more next time
		next := candidate + increment
		for {
			if (next+bus.offset)%bus.num == 0 {
				break
			}
			next += increment
		}
		// the new amount we'll jump up by for the next bus number
		increment = next - candidate
	}
	fmt.Println(buses)
	fmt.Println("result is", candidate)
}
