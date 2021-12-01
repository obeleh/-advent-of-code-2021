package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	lines, err := os.ReadFile("01.input")
	if err != nil {
		log.Fatalln(err)
	}

	var previousValue *int
	previousValue = nil
	countLarger := 0
	for _, line := range strings.Split(string(lines), "\n") {
		lineStripped := strings.TrimSpace(line)
		if value, err := strconv.Atoi(lineStripped); err == nil {
			if previousValue != nil && value > *previousValue {
				countLarger += 1
			}
			previousValue = &value
		}
	}
	print(fmt.Sprintf("Found %d values that were larger than the previous one\n", countLarger))
}
