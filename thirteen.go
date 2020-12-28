package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	start := 1000677
	busStr := "29,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,41,x,x,x,x,x,x,x,x,x,661,x,x,x,x,x,x,x,x,x,x,x,x,13,17,x,x,x,x,x,x,x,x,23,x,x,x,x,x,x,x,521,x,x,x,x,x,37,x,x,x,x,x,x,x,x,x,x,x,x,19"
	buses := strings.Split(busStr, ",")
	fmt.Println(buses)
	earliest := start
	for {
		for _, bus := range buses {
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
