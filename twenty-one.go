package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

type food struct {
	ingredients []string
	allergens   []string
	foodID      int
}

func load() (map[int]food, map[string][]int) {
	file, err := os.Open("./twenty-one-input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	foods := make(map[int]food, 0)
	allergens := make(map[string][]int, 0)
	id := 0
	for scanner.Scan() {
		// Example: sqjhc fvjkl (contains soy)
		line := scanner.Text()
		if err != nil {
			log.Fatal(err)
		}
		parts := strings.Split(line, " (contains ")
		ingredientsStr := parts[0]
		allergensStr := strings.TrimSuffix(parts[1], ")")
		f := food{strings.Split(ingredientsStr, " "), strings.Split(allergensStr, ", "), id}
		foods[id] = f
		for _, a := range f.allergens {
			if _, exists := allergens[a]; !exists {
				allergens[a] = make([]int, 0)
			}
			allergens[a] = append(allergens[a], id)
		}
		id++
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return foods, allergens
}

func findShared(foodIDs []int, foods map[int]food) []string {
	// Returns a list of ingredients that all foods in foodIDs have in common
	// foodIDs must be unique, and ingredients per food must be unique
	appearances := make(map[string]int)
	for _, id := range foodIDs {
		f := foods[id]
		for _, i := range f.ingredients {
			appearances[i]++
		}
	}
	shared := make([]string, 0)
	for ingredient, count := range appearances {
		if count == len(foodIDs) {
			shared = append(shared, ingredient)
		}
	}
	return shared
}

func solve(possible map[string][]string, solvedAllergens map[string]string, solvedIngredients map[string]string, total int) bool {
	if len(solvedAllergens) == total {
		return true
	}

	options := make(map[string][]string, 0)
	for allergen, possibilities := range possible {
		if _, alreadyDone := solvedAllergens[allergen]; !alreadyDone {
			ingredientsLeft := make([]string, 0)
			for _, ingredient := range possibilities {
				if _, exists := solvedIngredients[ingredient]; !exists {
					ingredientsLeft = append(ingredientsLeft, ingredient)
				}
			}
			if len(ingredientsLeft) == 0 {
				// if there are no more valid possible for this allergen, this solution will never work
				return false
			}
			if len(ingredientsLeft) == 1 {
				// try making this assignment right away -- if it fails then there's no
				// way the solution we're on will succeed either since it is the only option.
				solvedAllergens[allergen] = ingredientsLeft[0]
				solvedIngredients[ingredientsLeft[0]] = allergen
				if solve(possible, solvedAllergens, solvedIngredients, total) {
					return true
				} else {
					delete(solvedAllergens, allergen)
					delete(solvedIngredients, ingredientsLeft[0])
					return false
				}
			} else {
				// track all the worse options in case it comes to that and there's no single paths left
				options[allergen] = ingredientsLeft
			}
		}
	}
	// find the option with the fewest paths and explore that one
	var fewest []string
	var fewestAllergen string
	for allergen, paths := range options {
		if fewest == nil || len(paths) < len(fewest) {
			fewest = paths
			fewestAllergen = allergen
		}
	}
	for _, option := range fewest {
		solvedAllergens[fewestAllergen] = option
		solvedIngredients[option] = fewestAllergen
		if solve(possible, solvedAllergens, solvedIngredients, total) {
			return true
		}
		delete(solvedAllergens, fewestAllergen)
		delete(solvedIngredients, option)
	}
	return false
}

func main() {
	foods, allergens := load()
	// fmt.Println(foods)
	// fmt.Println(allergens)
	possible := make(map[string][]string)
	solvedAllergens := make(map[string]string)
	solvedIngredients := make(map[string]string)
	for a, foodIDs := range allergens {
		possible[a] = make([]string, 0)
		possible[a] = append(possible[a], findShared(foodIDs, foods)...)
	}
	fmt.Println(possible)
	solve(possible, solvedAllergens, solvedIngredients, len(allergens))
	fmt.Println(solvedAllergens)
	count := 0
	for _, f := range foods {
		for _, i := range f.ingredients {
			if _, exists := solvedIngredients[i]; !exists {
				count++
			}
		}
	}
	fmt.Println("final count of inactive", count)
	dangerousA := make([]string, 0)
	for a := range solvedAllergens {
		dangerousA = append(dangerousA, a)
	}
	sort.Slice(dangerousA, func(i, j int) bool {
		return dangerousA[i] < dangerousA[j]
	})
	dangerousI := make([]string, 0)
	for _, a := range dangerousA {
		dangerousI = append(dangerousI, solvedAllergens[a])
	}
	fmt.Println(dangerousI)
	fmt.Println("Dangerous ingredients, alphabetically by allergen")
	fmt.Println(strings.Join(dangerousI, ","))
}
