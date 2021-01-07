package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

type node struct {
	val  int
	next *node
}

type player struct {
	start *node
	end   *node
	len   int
}

func (p *player) append(val int) {
	n := node{val, nil}
	if p.start == nil {
		p.start = &n
	}
	p.end.next = &n
	p.end = &n
	p.len++
}

func (p *player) draw() int {
	if p.len == 0 {
		panic("length 0 for player")
	}
	top := p.start.val
	p.start = p.start.next
	p.len--
	return top
}

func (n node) String() string {
	str := fmt.Sprintf("%d", n.val)
	if n.next != nil {
		str += "->" + (*n.next).String()
	}
	return str
}

func (p *player) score() int {
	s := 0
	multiplier := p.len
	for n := p.start; n != nil; n = n.next {
		s += n.val * multiplier
		multiplier--
	}
	return s
}

func (p *player) copy(len int) *player {
	if len < 1 || len > p.len {
		panic("length makes no sense < 1")
	}
	newStart := &node{p.start.val, nil}
	ptr := newStart
	oldPtr := p.start
	count := 1
	for count < len {
		if oldPtr.next == nil {
			panic("can't copy this far")
		}
		ptr.next = &node{oldPtr.next.val, nil}
		ptr = ptr.next
		oldPtr = oldPtr.next
		count++
	}
	return &player{newStart, ptr, len}
}

func (p *player) String() string {
	str := fmt.Sprintf("Player: %d cards: ", p.len)
	if p.start != nil {
		str += (*p.start).String()
	}
	return str
}

func load() (*player, *player) {
	file, err := os.Open("./twenty-two-input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var player1, player2 *player
	isPlayer1 := true
	for scanner.Scan() {
		// Example:
		line := scanner.Text()
		if err != nil {
			log.Fatal(err)
		}
		if line == "" || line == "Player 1:" {
			continue
		} else if line == "Player 2:" {
			isPlayer1 = false
		} else {
			val, err := strconv.Atoi(line)
			if err != nil {
				panic(err)
			}
			if isPlayer1 {
				if player1 == nil {
					n := node{val, nil}
					player1 = &player{&n, &n, 1}
				} else {
					player1.append(val)
				}
			} else {
				if player2 == nil {
					n := node{val, nil}
					player2 = &player{&n, &n, 1}
				} else {
					player2.append(val)
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return player1, player2
}

func roundKey(p1 *player, p2 *player) string {
	str1 := "none"
	str2 := "none"
	if p1.start != nil {
		str1 = p1.start.String()
	}
	if p2.start != nil {
		str2 = p2.start.String()
	}
	return str1 + "--" + str2
}
func playRound2(p1 *player, p2 *player, prevRounds map[string]bool) bool {
	// Returns true if p1 won, false if p2 won
	for {
		if p1.len == 0 {
			return false
		}
		if p2.len == 0 {
			return true
		}
		_, alreadyPlayed := prevRounds[roundKey(p1, p2)]
		if alreadyPlayed {
			return true
		}
		// fmt.Println(">>>")
		prevRounds[roundKey(p1, p2)] = true
		// fmt.Println("Player 1's deck:", p1)
		// fmt.Println("Player 2's deck:", p2)
		v1 := p1.draw()
		v2 := p2.draw()
		// fmt.Println("Player 1 plays:", v1)
		// fmt.Println("Player 2 plays:", v2)
		p1Won := v1 > v2 // regular combat
		if p1.len >= v1 && p2.len >= v2 {
			// recursive combat
			new1 := p1.copy(v1)
			new2 := p2.copy(v2)
			newRounds := make(map[string]bool, 0)
			// fmt.Println("-- entering new game --")
			p1Won = playRound2(new1, new2, newRounds)
			// fmt.Println("-- done with new game --")
		}
		if p1Won {
			// fmt.Println("Player 1 won!")
			p1.append(v1)
			p1.append(v2)
		} else {
			// fmt.Println("Player 2 won!")
			p2.append(v2)
			p2.append(v1)
		}
	}
}

func playRound1(p1 *player, p2 *player) {
	for {
		if p1.len == 0 || p2.len == 0 {
			break
		}
		v1 := p1.draw()
		v2 := p2.draw()
		if v1 > v2 {
			p1.append(v1)
			p1.append(v2)
		} else {
			p2.append(v2)
			p2.append(v1)
		}
	}
}

func main() {
	p1, p2 := load()
	fmt.Println(p1)
	fmt.Println(p2)
	playRound2(p1, p2, make(map[string]bool, 0))
	fmt.Println("Player 1:", p1, p1.score())
	fmt.Println("Player 2:", p2, p2.score())
}
