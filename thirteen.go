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
	rawBuses := strings.Split(busStr, ",")
	fmt.Println(rawBuses)
	earliest := start
	// Round 1
	for {
		for _, bus := range rawBuses {
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
	rawBuses := strings.Split(busStr, ",")
	fmt.Println(rawBuses)
	buses := make([]Bus, 0)

	// Round 2
	for i, bus := range rawBuses {
		if string(bus) == "x" {
			continue
		}
		busNum, _ := strconv.Atoi(string(bus))
		fmt.Println(busNum)
		buses = append(buses, Bus{busNum, i})
	}
	fmt.Println(buses)

}
