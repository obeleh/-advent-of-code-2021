package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type octopus struct {
	level          int
	lastFlashRound int
	x              int
	y              int
}

type coord struct {
	x int
	y int
}

var OFFSETCOORDS = [...]coord{
	{x: -1, y: -1},
	{x: 0, y: -1},
	{x: +1, y: -1},
	{x: +1, y: 0},
	{x: +1, y: +1},
	{x: 0, y: +1},
	{x: -1, y: +1},
	{x: -1, y: 0},
}

func loadInput() [][]*octopus {
	bytes, err := os.ReadFile("input.txt")
	if err != nil {
		log.Fatalln(err)
	}
	lines := strings.Split(string(bytes), "\n")
	output := make([][]*octopus, len(lines))
	for idx, line := range lines {
		numberStrs := strings.Split(line, "")
		octoPussies := make([]*octopus, len(numberStrs))
		for nIdx, nrStr := range numberStrs {
			value, err := strconv.Atoi(nrStr)
			if err != nil {
				log.Fatal("Failed loading numbers")
			}
			octoPussy := &octopus{
				level:          value,
				x:              nIdx,
				y:              idx,
				lastFlashRound: -1,
			}
			octoPussies[nIdx] = octoPussy
		}
		output[idx] = octoPussies
	}
	return output
}

func getSurroundingCoords(x int, y int, xLen int, yLen int) []coord {
	output := make([]coord, 0)
	for _, coordOffset := range OFFSETCOORDS {
		if x+coordOffset.x >= 0 && y+coordOffset.y >= 0 && x+coordOffset.x < xLen && y+coordOffset.y < yLen {
			output = append(output, coord{
				x: x + coordOffset.x,
				y: y + coordOffset.y,
			})
		}
	}
	return output
}

func getAdjacentOcupussies(octoPussiesMap [][]*octopus, x int, y int) []*octopus {
	yLen := len(octoPussiesMap)
	xLen := len(octoPussiesMap[0])
	octoPussySlice := make([]*octopus, 0)
	for _, c := range getSurroundingCoords(x, y, xLen, yLen) {
		octoPussySlice = append(octoPussySlice, octoPussiesMap[c.y][c.x])
	}
	return octoPussySlice
}

func playRoundOnOctopussy(octoPussy *octopus, roundNr int, octoPussiesMap [][]*octopus) int {
	// keep level at 0 and no longer flash
	if octoPussy.lastFlashRound == roundNr {
		return 0
	}
	flashCount := 0
	octoPussy.level += 1
	if octoPussy.level > 9 {
		//flash
		flashCount++
		octoPussy.lastFlashRound = roundNr
		octoPussy.level = 0
		for _, adjacentOctoPussy := range getAdjacentOcupussies(octoPussiesMap, octoPussy.x, octoPussy.y) {
			flashCount += playRoundOnOctopussy(adjacentOctoPussy, roundNr, octoPussiesMap)
		}
	}
	return flashCount
}

func playRound(octoPussiesMap [][]*octopus, roundNr int) int {
	flashCount := 0
	yLen := len(octoPussiesMap)
	xLen := len(octoPussiesMap[0])
	for y := 0; y < yLen; y++ {
		for x := 0; x < xLen; x++ {
			flashCount += playRoundOnOctopussy(octoPussiesMap[y][x], roundNr, octoPussiesMap)
		}
	}
	return flashCount
}

func printMap(octoPussiesMap [][]*octopus) {
	yLen := len(octoPussiesMap)
	xLen := len(octoPussiesMap[0])
	for y := 0; y < yLen; y++ {
		for x := 0; x < xLen; x++ {
			print(octoPussiesMap[y][x].level)
		}
		print("\n")
	}
	print("\n")
}

func main() {
	octoPussiesMap := loadInput()

	mapSize := len(octoPussiesMap) * len(octoPussiesMap[0])

	flashCount := 0
	for round := 0; round < 999; round++ {
		roundFlashCount := playRound(octoPussiesMap, round)
		flashCount += roundFlashCount
		printMap(octoPussiesMap)
		if round == 99 {
			print(fmt.Sprintf("flashCount %d\n", flashCount))
		}

		if roundFlashCount == mapSize {
			print(fmt.Sprintf("All flash round %d\n", round+1))
			break
		}
	}

}
