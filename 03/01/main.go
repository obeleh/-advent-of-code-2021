package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const NUMBITS = 12

func main() {

	lines, err := os.ReadFile("input.txt")
	if err != nil {
		log.Fatalln(err)
	}

	masks := make([]int, NUMBITS)
	counts := make([]int, NUMBITS)
	for shift := 0; shift < NUMBITS; shift++ {
		masks[shift] = 1 << shift
		counts[shift] = 0
	}

	lineCount := 0
	for _, line := range strings.Split(string(lines), "\n") {
		lineStripped := strings.TrimSpace(line)
		words := strings.Fields(lineStripped)
		if len(words) != 1 {
			log.Fatalln("Expected 1 words to be present per line")
		}

		// NUMBITS+1 because ParseInt parses into signed ints
		value64, err := strconv.ParseInt(words[0], 2, NUMBITS+1)
		if err != nil {
			log.Fatalln(fmt.Errorf("Expected the value to be 12 bits,  %v", err))
		}
		value := int(value64)
		for bit := 0; bit < NUMBITS; bit++ {
			if masks[bit]&value == masks[bit] {
				counts[bit] += 1
			}
		}
		lineCount += 1
	}

	halfLineCount := lineCount / 2

	gamma := 0
	for bit := 0; bit < NUMBITS; bit++ {
		if counts[bit] > halfLineCount {
			gamma += 1 << bit
		}
	}
	xorMask := (1 << NUMBITS) - 1
	epsilon := gamma ^ xorMask
	print(fmt.Sprintf("gamma:%d epsilon:%d multiplied:%d\n", gamma, epsilon, gamma*epsilon))
}
