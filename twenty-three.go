package main

import (
	"fmt"
	"strconv"
)

func findTargetIRound1(candidates string, dest int) int {
	for {
		for i, val := range candidates {
			valInt, _ := strconv.Atoi(string(val))
			if valInt == dest {
				return i
			}
		}
		dest--
		if dest < 1 {
			dest = 9 // reset to the max
		}
	}
}

// pick up 3, remove them from consideration, wrap if needed
// find closest destination to find one tha wasn't picked up, wrapping if needed
// insert the picked up ones to the right of the destination cup
// pick a new current cup, clockwise of the current cup

func round1() {
	cups := "186524973"
	currI := 0
	for x := 0; x < 100; x++ {
		currVal := string(cups[currI])
		// rotate the string so it starts at currI + 1 in index 0 and excludes curr
		rotated := cups[currI+1:]
		pickup := rotated[:3]
		remaining := rotated[3:]
		currLabel, _ := strconv.Atoi(currVal)
		destTargetI := findTargetIRound1(remaining, currLabel-1)
		reinserted := remaining[:destTargetI+1] + pickup + remaining[destTargetI+1:]
		fmt.Println("-- move", x, "--")
		fmt.Println("cups:", cups)
		fmt.Println("pickup:", pickup)
		fmt.Println("destination:", remaining[destTargetI:destTargetI+1])
		cups = reinserted + currVal
	}
	fmt.Println("Final order:", cups)
	fmt.Println("Done")
}

type node struct {
	val  int
	next *node
	prev *node
}

func (n *node) append(val int) *node {
	next := node{val, nil, n}
	n.next = &next
	return n.next
}

func (n *node) forward(distance int) *node {
	result := n
	for i := 0; i < distance; i++ {
		result = result.next
	}
	return result
}

func (n *node) backward(distance int) *node {
	result := n
	for i := 0; i < distance; i++ {
		result = result.prev
	}
	return result
}

const maxInclusive = 1000000
const numIterations = 10000000

func findTargetNode(start *node, dest int, blocklist *node, valMap map[int]*node) *node {
	if dest < 1 {
		dest = maxInclusive
	}

	for {
		// all numbers are present, so we can figure out what to look for
		// just based on the blocklist
		found := false
		for n := blocklist; n != nil; n = n.next {
			if n.val == dest {
				found = true
				break
			}
		}
		if !found {
			break
		}
		dest--
		if dest < 1 {
			dest = maxInclusive // reset to the max
		}
	}
	return valMap[dest]
}

func (n *node) String() string {
	str := fmt.Sprintf("%d", n.val)
	m := n.next
	// since it is a circle, assume unique values
	for m != nil && m.val != n.val {
		str += fmt.Sprintf("->%d", m.val)
		m = m.next
	}
	return str
}

func main() {
	start := &node{1, nil, nil} // first digit here
	// Track which nodes we have for which values for quick lookup
	valMap := make(map[int]*node)
	valMap[1] = start
	end := start
	for _, d := range "86524973" { // remaining input digits here
		val, _ := strconv.Atoi(string(d))
		end = end.append(val)
		valMap[val] = end
	}

	for i := 10; i <= maxInclusive; i++ {
		end = end.append(i)
		valMap[i] = end
	}
	// make it a circle
	end.next = start
	start.prev = end

	curr := start
	for x := 0; x < numIterations; x++ {
		pickup := curr.forward(1)
		remaining := curr.forward(4)
		remaining.prev = nil       // splice out pickup from remaining
		curr.forward(3).next = nil // splice out pickup from remaining
		pickup.prev = nil
		curr.next = nil // splice out pickup from curr, don't add remaining yet
		dest := findTargetNode(curr.prev, curr.val-1, pickup, valMap)
		// splice in pickup
		oldDestNext := dest.next
		dest.next = pickup
		pickup.prev = dest
		pickup.forward(2).next = oldDestNext
		oldDestNext.prev = pickup.forward(2)
		curr.next = remaining
		remaining.prev = curr

		curr = curr.next
		if x%100000 == 0 {
			fmt.Println("iteration", x)
		}
	}
	fmt.Println("Done")
	n := curr
	for {
		if n.val == 1 {
			fmt.Println("Found 1")
			once := n.forward(1).val
			twice := n.forward(2).val
			fmt.Println("Clockwise once:", once)
			fmt.Println("Clockwise twice:", twice)
			fmt.Println("product:", once*twice)
			break
		}
		n = n.next
	}
}
