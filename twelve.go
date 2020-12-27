package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

type Command struct {
	direction string
	amount    int
}

func loadCommands() []Command {
	file, err := os.Open("./twelve-input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	commands := make([]Command, 0)
	for scanner.Scan() {
		// Example: #F10
		line := scanner.Text()
		if err != nil {
			log.Fatal(err)
		}
		amount, _ := strconv.Atoi(line[1:])
		commands = append(commands, Command{string(line[0]), amount})
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return commands
}

var orientations = [4]string{"E", "S", "W", "N"}

func roundOne(commands []Command) {
	x, y, orientation := 0, 0, 0
	for _, command := range commands {
		direction := command.direction
		if direction == "R" || direction == "L" {
			change := command.amount / 90
			if direction == "L" {
				change = command.amount / 90 * -1
			}
			// + len(orientations) so we cycle backwards properly without going negative
			orientation = (orientation + len(orientations) + change) % len(orientations)
			continue
		}
		if direction == "F" {
			direction = orientations[orientation]
		}
		switch direction {
		case "N":
			y += command.amount
		case "S":
			y -= command.amount
		case "E":
			x += command.amount
		case "W":
			x -= command.amount
		}
	}
	fmt.Println("Ended up at", x, y, "facing", orientations[orientation])
	fmt.Println("Total manhattan distance:", manhattan(x, y))
}

func manhattan(x int, y int) (manhattan int) {
	if x > 0 {
		manhattan += x
	} else {
		manhattan -= x
	}
	if y > 0 {
		manhattan += y
	} else {
		manhattan -= y
	}
	return
}

type position struct {
	x int
	y int
}

func roundTwo(commands []Command) {
	ship := position{0, 0}
	waypoint := position{10, 1}
	for _, command := range commands {
		direction := command.direction
		if direction == "F" {
			ship.x += waypoint.x * command.amount
			ship.y += waypoint.y * command.amount
		} else if direction == "L" || direction == "R" {
			if command.amount == 180 {
				// Rotate around
				waypoint.x *= -1
				waypoint.y *= -1
			} else if (direction == "L" && command.amount == 90) || (direction == "R" && command.amount == 270) {
				// Rotate left
				waypoint.x, waypoint.y = -waypoint.y, waypoint.x
			} else if (direction == "R" && command.amount == 90) || (direction == "L" && command.amount == 270) {
				// Rotate right
				waypoint.x, waypoint.y = waypoint.y, -waypoint.x
			} else {
				panic("Need to handle this after all")
			}
		} else {
			switch direction {
			case "N":
				waypoint.y += command.amount
			case "S":
				waypoint.y -= command.amount
			case "E":
				waypoint.x += command.amount
			case "W":
				waypoint.x -= command.amount
			}
		}
	}
	fmt.Println("Waypoint up at", waypoint.x, waypoint.y)
	fmt.Println("Ship ended up at", ship.x, ship.y)
	fmt.Println("Total manhattan distance:", manhattan(ship.x, ship.y))
}

func main() {
	commands := loadCommands()
	roundTwo(commands)
}
