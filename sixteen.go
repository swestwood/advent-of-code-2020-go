package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type rule struct {
	name        string
	firstStart  int
	firstEnd    int
	secondStart int
	secondEnd   int
}

const (
	RULES  = 1
	YOURS  = 2
	NEARBY = 3
)

func invalidSum(nums []int, rules []rule) (int, bool) {
	sum := 0
	valid := true
	for _, num := range nums {
		ok := false
		for _, r := range rules {
			if num >= r.firstStart && num <= r.firstEnd ||
				num >= r.secondStart && num <= r.secondEnd {
				ok = true
				break
			}
		}
		if !ok {
			sum += num
			valid = false
		}
	}
	return sum, valid
}
func parseTicket(line string) []int {
	fields := strings.Split(line, ",")
	nums := make([]int, 0)
	for _, f := range fields {
		num, err := strconv.Atoi(f)
		if err != nil {
			panic(err)
		}
		nums = append(nums, num)
	}
	return nums
}

func scanTickets() ([][]int, []rule, []int) {
	file, err := os.Open("./sixteen-input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	section := RULES
	rules := make([]rule, 0)
	invalid := 0
	validByPosition := make([][]int, 0)
	yours := make([]int, 0)

	for scanner.Scan() {
		line := scanner.Text()
		if err != nil {
			log.Fatal(err)
		}
		if line == "" {
			section++
		} else if line == "your ticket:" {
			section = YOURS
		} else if line == "nearby tickets:" {
			section = NEARBY
		} else if section == RULES {
			parts := strings.Split(line, ": ")
			r := rule{name: parts[0]}
			_, err := fmt.Sscanf(parts[1], "%d-%d or %d-%d", &r.firstStart, &r.firstEnd, &r.secondStart, &r.secondEnd)
			if err != nil {
				panic(err)
			}
			rules = append(rules, r)
		} else if section == YOURS {
			yours = parseTicket(line)
		} else if section == NEARBY {
			nums := parseTicket(line)
			sum, isValid := invalidSum(nums, rules)
			invalid += sum // for part 1
			if isValid {
				// make a corresponding array to store all values at this index across valid tickets
				for i, num := range nums {
					if len(validByPosition) <= i {
						validByPosition = append(validByPosition, make([]int, 0))
					}
					validByPosition[i] = append(validByPosition[i], num)
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("num invalid", invalid)
	return validByPosition, rules, yours
}

func getProposal(mappings []mapping, i int, proposal map[int]string,
	rulesLeft map[string]bool) (bool, map[int]string) {
	if i == len(mappings) {
		return true, proposal
	}
	// check if anything in the future has no options to bail out early
	for j := i; j < len(mappings); j++ {
		m := mappings[j]
		anyLeft := false
		for _, o := range m.options {
			if rulesLeft[o] {
				anyLeft = true
				break
			}
		}
		if !anyLeft {
			return false, proposal
		}
	}
	// try out each rule
	for _, o := range mappings[i].options {
		if !rulesLeft[o] {
			continue
		}
		proposal[mappings[i].position] = o
		rulesLeft[o] = false
		success, proposal := getProposal(mappings, i+1, proposal, rulesLeft)
		if success {
			return true, proposal
		}
		delete(proposal, mappings[i].position)
		rulesLeft[o] = true
	}
	// tried all the rules and none worked
	return false, proposal
}

type mapping struct {
	options  []string
	position int
}

func main() {
	validByPosition, rules, yours := scanTickets()
	proposal := make(map[int]string, 0)
	mappings := make([]mapping, 0)
	for pos, toMatch := range validByPosition {
		m := mapping{make([]string, 0), pos}
		for _, r := range rules {
			isMatch := true
			for _, num := range toMatch {
				if (num < r.firstStart || num > r.firstEnd) &&
					(num < r.secondStart || num > r.secondEnd) {
					isMatch = false
					break
				}
			}
			if isMatch {
				m.options = append(m.options, r.name)
			}
		}
		mappings = append(mappings, m)
	}
	// sort array to have positions with fewest choices first
	sort.Slice(mappings, func(i, j int) bool {
		return len(mappings[i].options) < len(mappings[j].options)
	})
	// initially all rules are available
	rulesLeft := make(map[string]bool, 0)
	for _, r := range rules {
		rulesLeft[r.name] = true
	}

	fmt.Println("There are", len(rules), "rules and", len(validByPosition), "positions")
	fmt.Println(yours)
	// Map the fields
	success, proposal := getProposal(mappings, 0, proposal, rulesLeft)
	fmt.Println(success, proposal)
	if !success {
		panic("welp")
	}
	final := 1
	for pos, name := range proposal {
		if strings.HasPrefix(name, "departure") {
			final *= yours[pos]
		}
	}
	fmt.Println("Final product", final)

}
