package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	lines, err := os.ReadFile("02.input")
	if err != nil {
		log.Fatalln(err)
	}

	x, y := 0, 0

	for _, line := range strings.Split(string(lines), "\n") {
		lineStripped := strings.TrimSpace(line)
		words := strings.Fields(lineStripped)
		if len(words) != 2 {
			log.Fatalln("Expected 2 words to be present per line")
		}

		value, err := strconv.Atoi(words[1])
		if err != nil {
			log.Fatalln("Expected the value to be numeric")
		}

		switch words[0] {
		case "down":
			y += value
		case "up":
			y -= value
		case "forward":
			x += value
		default:
			log.Fatalln(fmt.Sprintf("Unknown direction %s", words[0]))
		}
	}
	print(fmt.Sprintf("Eventual location %d %d\n", x, y))
	print(fmt.Sprintf("Eventual result %d\n", x*y))
}
