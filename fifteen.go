package main

import (
	"fmt"
	"strconv"
	"strings"
)

func roundOne() {
	input := "9,12,1,4,17,0,18"
	parts := strings.Split(input, ",")
	nums := make([]int, 0)
	for _, part := range parts {
		num, _ := strconv.Atoi(part)
		nums = append(nums, num)
	}
	last := nums[len(nums)-1]
	for {
		next := 0
		for i, num := range nums {
			if num == last && i < len(nums)-1 {
				// it was spoken before
				next = len(nums) - 1 - i
			}
		}
		nums = append(nums, next)
		last = next
		if len(nums) == 2020 {
			fmt.Println("the 2020th is", nums[2019])
			break
		}
	}
}

func main() {
	input := "9,12,1,4,17,0,18"
	parts := strings.Split(input, ",")
	nums := make(map[int]int, 0)
	last := -1
	desired := 30000000
	// iteration = number of numbers seen (including last)
	iteration := 0
	for i, part := range parts {
		num, err := strconv.Atoi(part)
		if err != nil {
			panic("can't parse")
		}
		fmt.Println(num)

		nums[num] = i
		last = num
		iteration++
	}
	delete(nums, last) // the most recent one should never be in the map
	fmt.Println(nums)
	for {
		next := 0
		lastI, exists := nums[last]
		if exists {
			// it was spoken before
			next = iteration - lastI - 1
		}
		// fmt.Println("last", last, "| exists", exists, "| lastI", lastI, "| iteration", iteration, "| next", next)
		// fill this in retroactively since we're one step behind doing it
		nums[last] = iteration - 1
		last = next
		iteration++
		if iteration == desired {
			fmt.Println(desired, " is", next)
			break
		}
	}
}
