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

	countLarger := 0
	var slidingWindow [3]int
	runningSum := 0
	for idx, line := range strings.Split(string(lines), "\n") {
		lineStripped := strings.TrimSpace(line)
		if value, err := strconv.Atoi(lineStripped); err == nil {
			previousValue := runningSum
			if idx > 2 {
				runningSum -= slidingWindow[idx%3]
			}
			slidingWindow[idx%3] = value
			runningSum += value
			if idx > 2 && runningSum > previousValue {
				countLarger += 1
			}
		}
	}
	print(fmt.Sprintf("Found %d values that were larger than the previous one\n", countLarger))
}
