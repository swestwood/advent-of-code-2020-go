package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func getBagsByContents() map[string]map[string]int {
	file, err := os.Open("./seven-input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	byContents := make(map[string]map[string]int)
	for scanner.Scan() {
		// Example: shiny gold bags contain 1 dark olive bag, 2 vibrant plum bags.
		line := scanner.Text()
		if err != nil {
			log.Fatal(err)
		}
		pieces := strings.Split(line, "contain")
		sourceColor := strings.TrimSuffix(pieces[0], " bags ")
		byContents[sourceColor] = make(map[string]int)
		contents := strings.Trim(pieces[1], ". ")
		contentPieces := strings.Split(contents, ", ")
		for _, phrase := range contentPieces {
			phrase = strings.TrimSuffix(phrase, " bag")
			phrase = strings.TrimSuffix(phrase, " bags")
			// No other is just an empty map
			if phrase != "no other" {
				phrasePieces := strings.SplitN(phrase, " ", 2)
				if len(phrasePieces) == 2 {
					numBags, err := strconv.Atoi(phrasePieces[0])
					if err != nil {
						log.Fatal(err)
					}
					byContents[sourceColor][phrasePieces[1]] = numBags
				}
			}

		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return byContents
}

// All the bags that this bag could go into
func getBagsByContainingBag(byContents map[string]map[string]int) map[string][]string {
	byContainingBag := make(map[string][]string)
	for bag, contents := range byContents {
		for innerBag := range contents {
			if _, ok := byContainingBag[innerBag]; !ok {
				byContainingBag[innerBag] = make([]string, 0)
			}
			byContainingBag[innerBag] = append(byContainingBag[innerBag], bag)
		}
	}
	return byContainingBag
}

func getAllPossibleContainers(targetBag string, byContainingBag map[string][]string,
	possibleContainers map[string]bool) {
	// Allowed bags are any bags that are in byContainingBag for targetBag, or any
	// of the containing bags in turn for those bags, on up. But, prevent cycles.
	if containers, ok := byContainingBag[targetBag]; ok {
		for _, container := range containers {
			if _, alreadyPresent := possibleContainers[container]; !alreadyPresent {
				possibleContainers[container] = true
				getAllPossibleContainers(container, byContainingBag, possibleContainers)
			}
		}
	}

}

func numBagsInside(targetBag string, byContents map[string]map[string]int) int {
	bagsInside := 0
	if contents, ok := byContents[targetBag]; ok {
		for innerBag, count := range contents {
			bagsInside += count + count*numBagsInside(innerBag, byContents)
			// fmt.Println("Adding", count, innerBag, "bags", "that's", bagsInside)
		}
	}
	return bagsInside
}

func main() {
	// Round 1
	byContents := getBagsByContents() // Example: bright white:map[shiny gold:1]
	byContainingBag := getBagsByContainingBag(byContents)
	possibleContainers := make(map[string]bool)
	targetBag := "shiny gold"
	getAllPossibleContainers(targetBag, byContainingBag, possibleContainers)
	fmt.Println(len(possibleContainers), "bags can eventually contain a", targetBag)

	// Round 2
	bagsInside := numBagsInside("shiny gold", byContents)
	fmt.Println("1", targetBag, "bag must contain", bagsInside, "other bags")
}
