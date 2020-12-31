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

type node struct {
	val      string
	children []*node
}

func (n *node) String() string {
	if len(n.children) == 0 {
		return n.val
	}
	result := n.val + "["
	for _, c := range n.children {
		result += c.String() + " "
	}
	return result + "]"
}

type rule struct {
	id int
	// only one of val or options should be present
	val     string
	options [][]int
}

func getOrMakeNext(roots []*node, val string) []*node {
	// Build children nodes with val for all roots if not already present
	// Return a list of all these children (existing and new)
	frontier := make([]*node, 0)

	for _, r := range roots {
		found := false
		for _, c := range r.children {
			if c.val == val {
				frontier = append(frontier, c)
				found = true
			}
		}
		if !found {
			child := node{val, make([]*node, 0)}
			r.children = append(r.children, &child)
			frontier = append(frontier, &child)
		}
	}
	return frontier
}

func addRuleAt(index int, subindex int, roots []*node, rules map[int]rule) []*node {
	// Build out the grammar tree.
	// roots is all the paths that are currently being build out for this rule
	r := rules[index]
	if r.val != "" {
		// add this final value to all paths currently being built out
		result := getOrMakeNext(roots, r.val)
		// fmt.Println("final rule #", r.id, "with", r.val, "- roots", roots, "- frontier", result)
		return result
	}
	if subindex != -1 {
		// work through the rules sequentially in this array to reach the
		// FINAL frontier for this sequence of this rule
		seq := rules[index].options[subindex]
		nextFrontier := make([]*node, 0)
		nextFrontier = append(nextFrontier, roots...)
		for _, seqRuleI := range seq {
			seqRule := rules[seqRuleI]
			nextFrontier = addRuleAt(seqRule.id, -1, nextFrontier, rules)
			// fmt.Println("done seq rule #", r.id, "seq", seq, "at", seqRuleI, "with", "- roots", roots, "- nextfrontier", nextFrontier)
		}
		// fmt.Println("adding", seq, roots)
		return nextFrontier
	}
	// combine the frontiers of the OR of each option in this rule
	frontier := make([]*node, 0)
	// fmt.Println("started rule", r.id)
	for i := range r.options {
		// add each branch subindex
		f := addRuleAt(index, i, roots, rules)
		frontier = append(frontier, f...)
	}
	return frontier
}

func build(rules map[int]rule, startRuleI int) *node {
	root := &node{"^", make([]*node, 0)}
	roots := []*node{root}
	endNodes := addRuleAt(startRuleI, -1, roots, rules)
	// add a terminal character
	getOrMakeNext(endNodes, "#")
	return root
}

func load() (map[int]rule, []string) {
	file, err := os.Open("./nineteen-input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	rules := make(map[int]rule, 0)
	messages := make([]string, 0)
	isRule := true
	for scanner.Scan() {
		line := scanner.Text()
		if err != nil {
			log.Fatal(err)
		}
		if line == "" {
			isRule = false
		} else if isRule {
			parts := strings.Split(line, ": ")
			id, err := strconv.Atoi(parts[0])
			if err != nil {
				panic(err)
			}
			r := rule{id: id}
			if parts[1][0] == '"' {
				r.val = parts[1][1 : len(parts[1])-1]
			} else {
				r.options = make([][]int, 0)
				opts := strings.Split(parts[1], " | ")
				for _, opt := range opts {
					nums := strings.Split(opt, " ")
					optArr := make([]int, 0)
					for _, val := range nums {
						num, err := strconv.Atoi(val)
						if err != nil {
							panic(err)
						}
						optArr = append(optArr, num)
					}
					r.options = append(r.options, optArr)
				}
			}
			rules[id] = r
		} else {
			messages = append(messages, line)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return rules, messages
}

type loopedMatch struct {
	endIndex int
	numLoops int
}

func loopedPrefixes(str string, index int, root *node, rules map[int]rule, loopNum int) []loopedMatch {
	prefixes := make([]loopedMatch, 0)
	originalRoot := root
	for i := index; i < len(str); i++ {
		val := string(str[i])
		matchedChild := false
		for _, child := range root.children {
			if child.val == val {
				root = child
				matchedChild = true
			}
			if child.val == "#" {
				// even if there are upcoming substrings, if we hit at # that means
				// this substring is valid (where i is the noninclusive end index)
				prefixes = append(prefixes, loopedMatch{i, loopNum})
				// fmt.Println((loopNum + 1), "-", "found prefix", i, str[index:i])
				// also find all future looped options starting from here onward
				loopingFromHere := loopedPrefixes(str, i, originalRoot, rules, loopNum+1)
				prefixes = append(prefixes, loopingFromHere...)
			}
		}
		if !matchedChild {
			break
		}
	}
	return prefixes
}

func matches(str string, index int, root *node, rules map[int]rule) bool {
	// 8 is one or more 42s in a row
	// 11 is one or more 42s in a row followed by one or more 31s in a row, with
	// the same number of 42s and 31s.
	// 0 is 8 then 11
	// For example: [42]^j[42]^k[31]^k for j, k >= 1 to match 0
	// as long as #42s is >=2, and #31s < #42s, there is some j, k matching.
	// First figure out all the ways to get k 42s matching, then for each option,
	// find the max number of matching 31s reaching to the end of the string.
	// Then see if the repetitions for each are okay.
	root42 := build(rules, 42)
	root31 := build(rules, 31)
	match42s := loopedPrefixes(str, 0, root42, rules, 1)
	// put options with most loops first
	sort.Slice(match42s, func(i, j int) bool {
		return match42s[i].numLoops > match42s[j].numLoops
	})

	for _, m42 := range match42s {
		if m42.numLoops < 2 {
			continue
		}
		match31s := loopedPrefixes(str, m42.endIndex, root31, rules, 1)
		for _, m31 := range match31s {
			// minus 1 because of the trailing # -- did this reach the end of the str
			if m31.endIndex == len(str)-1 && m31.numLoops < m42.numLoops {
				// fmt.Println("\nSUCCESS FOUND for")
				// fmt.Println(str)
				// fmt.Println(str[:m42.endIndex], m42.numLoops, m42.endIndex)
				// fmt.Println(str[m42.endIndex:m31.endIndex], m31.numLoops, m31.endIndex)
				return true
			}
		}
	}
	return false
}

func matchesRound1(m string, root *node, rules map[int]rule) bool {
	if root.val != "^" {
		panic("huh")
	}
	for _, v := range m {
		val := string(v)
		found := false
		for _, child := range root.children {
			if child.val == val {
				root = child
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func main() {
	rules, messages := load()
	root := build(rules, 0)
	// fmt.Println(root)
	count := 0
	for _, m := range messages {
		// if matchesRound1(m+"#", root, rules) {
		if matches(m+"#", 0, root, rules) {
			count++
		}
	}
	fmt.Println(count, "messages of ", len(messages), "match")
}
