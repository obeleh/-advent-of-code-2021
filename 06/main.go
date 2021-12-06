package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const NUMDAYS int = 256

type fishGeneration struct {
	count int
	age   int
}

func loadInput() []*fishGeneration {
	bytes, err := os.ReadFile("input.txt")
	if err != nil {
		log.Fatalln(err)
	}
	strValues := strings.Split(string(bytes), ",")

	fish := []*fishGeneration{
		&fishGeneration{count: 0, age: 0},
		&fishGeneration{count: 0, age: 1},
		&fishGeneration{count: 0, age: 2},
		&fishGeneration{count: 0, age: 3},
		&fishGeneration{count: 0, age: 4},
		&fishGeneration{count: 0, age: 5},
		&fishGeneration{count: 0, age: 6},
		&fishGeneration{count: 0, age: 7},
		&fishGeneration{count: 0, age: 8},
	}
	var value int
	for _, strValue := range strValues {
		value, err = strconv.Atoi(strValue)
		if err != nil {
			log.Fatalf("Failed loading input, %v", err)
		}
		fish[value].count += 1
	}
	return fish
}

func main() {
	fishies := loadInput()

	age6Idx := 6
	for day := 0; day < NUMDAYS; day++ {
		newFish := 0
		for _, gen := range fishies {
			gen.age -= 1
			if gen.age < 0 {
				gen.age = 6
				newFish = gen.count
			}
		}

		age6Idx += 1
		if age6Idx > 6 {
			age6Idx = 0
		}
		fishies[age6Idx].count += fishies[7].count
		fishies[7].count = fishies[8].count
		fishies[7].age = 7
		fishies[8].count = newFish
		fishies[8].age = 8
	}

	total := 0
	for _, gen := range fishies {
		total += gen.count
	}

	print(fmt.Sprintf("%d fish", total))
}
