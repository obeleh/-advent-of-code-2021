package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/thoas/go-funk"
)

func loadInput() [][]int {
	bytes, err := os.ReadFile("input.txt")
	if err != nil {
		log.Fatalln(err)
	}
	lines := strings.Split(string(bytes), "\n")
	input := make([][]int, len(lines))
	for idx, line := range lines {
		numberStrs := strings.Split(line, "")
		numbersSlice := make([]int, len(numberStrs))
		for nIdx, nrStr := range numberStrs {
			value, err := strconv.Atoi(nrStr)
			if err != nil {
				log.Fatal("Failed loading numbers")
			}
			numbersSlice[nIdx] = value
		}
		input[idx] = numbersSlice
	}
	return input
}

type Coord struct {
	x int
	y int
}
type lowCoord struct {
	Coord
	val               int
	surroundingValues []Coord
}

var OFFSETCOORDS = [...]Coord{
	{x: -1, y: -1},
	{x: 0, y: -1},
	{x: +1, y: -1},
	{x: +1, y: 0},
	{x: +1, y: +1},
	{x: 0, y: +1},
	{x: -1, y: +1},
	{x: -1, y: +0},
}

func valIsLower(val int, xC int, yC int, xLen int, yLen int, heightMap [][]int) bool {
	if xC >= 0 && yC >= 0 && xC < xLen && yC < yLen {
		return val < heightMap[yC][xC]
	} else {
		// return true when out of bounds
		return true
	}
}

func challenge1(heightMap [][]int, yLen int, xLen int) []lowCoord {
	lowestIdxs := make([]lowCoord, 0)
	for yIdx, yVals := range heightMap {
		for xIdx, val := range yVals {
			if val == 9 {
				continue
			}

			// loop around, start top left go clockwise
			// skip to next field if current value is not lower than any of the surrounding

			if !valIsLower(val, xIdx, yIdx-1, xLen, yLen, heightMap) {
				continue // top
			}
			if !valIsLower(val, xIdx+1, yIdx, xLen, yLen, heightMap) {
				continue // right
			}
			if !valIsLower(val, xIdx, yIdx+1, xLen, yLen, heightMap) {
				continue // bottom
			}
			if !valIsLower(val, xIdx-1, yIdx, xLen, yLen, heightMap) {
				continue // left
			}

			// we checked all around, val is the lowest in the area
			lowestIdxs = append(lowestIdxs, lowCoord{
				Coord: Coord{
					x: xIdx,
					y: yIdx,
				},
				val: val,
			})
		}
	}

	outputSum := 0
	for _, lC := range lowestIdxs {
		outputSum += lC.val + 1
	}

	print(fmt.Sprintf("lcSum %d\n", outputSum))
	return lowestIdxs
}

func getSurroundingCoords(x int, y int, heightMap [][]int, yLen int, xLen int) []Coord {
	output := make([]Coord, 0)
	for _, coordOffset := range OFFSETCOORDS {
		if x+coordOffset.x >= 0 && y+coordOffset.y >= 0 && x+coordOffset.x < xLen && y+coordOffset.y < yLen {
			output = append(output, Coord{
				x: x + coordOffset.x,
				y: y + coordOffset.y,
			})
		}
	}
	return output
}

func challenge2(heightMap [][]int, yLen int, xLen int, lowcoords []lowCoord) {
	lowestIdxs := make([]lowCoord, 0)
	alreadyDiscoveredCoords := make([]Coord, 0)
	for yIdx, yVals := range heightMap {
	VAL_LOOP:
		for xIdx, val := range yVals {
			if val == 9 {
				continue
			}

			surroundingCoords := getSurroundingCoords(xIdx, yIdx, heightMap, yLen, xLen)
			for _, coord := range surroundingCoords {
				if val >= heightMap[coord.y][coord.x] {
					// skip to next field if current value is not lower than any of the surrounding
					continue VAL_LOOP
				}
			}
			centerVal := val

			// we've found the lowest point, now we find all steps up

			nextStep := make([]Coord, len(surroundingCoords))
			copy(nextStep, surroundingCoords)
			rejects := make([]Coord, 0)
		WHILE_LOOP:
			for { // while new steps up are being found
				candidates := make([]Coord, 0)
				copy(rejects, surroundingCoords)
				for _, coord := range nextStep {
					val = heightMap[coord.y][coord.x]
					newSurroundingCoords := getSurroundingCoords(coord.x, coord.y, heightMap, yLen, xLen)
					candidates = append(candidates, newSurroundingCoords...)
					for _, newCoord := range newSurroundingCoords {
						newCoordVal := heightMap[newCoord.y][newCoord.x]
						if val >= newCoordVal || newCoordVal == 9 {
							rejects = append(rejects, newCoord)
						}
					}
				}
				candidates = funk.Subtract(funk.Uniq(candidates), surroundingCoords).([]Coord)
				rejects = funk.Uniq(rejects).([]Coord)
				nextStep = funk.Subtract(candidates, rejects).([]Coord)
				found := len(nextStep) > 0
				if found {
					surroundingCoords = append(surroundingCoords, nextStep...)
				} else {
					// break condition
					uniqSurrounding := funk.Uniq(surroundingCoords).([]Coord)
					lc := lowCoord{
						Coord: Coord{
							x: xIdx,
							y: yIdx,
						},
						val:               centerVal,
						surroundingValues: uniqSurrounding,
					}
					duplicates := funk.Intersect(uniqSurrounding, alreadyDiscoveredCoords).([]Coord)
					if len(duplicates) > 0 {
						for idx, lowestIdx := range lowestIdxs {
							intersection := funk.Intersect(lowestIdx.surroundingValues, uniqSurrounding).([]Coord)
							if len(intersection) > 0 {
								//xreplace old found lc only if the new one is bigger
								if len(lowestIdx.surroundingValues) > len(uniqSurrounding) {
									lowestIdxs[idx] = lc
								}
								break WHILE_LOOP
							}
						}
					}
					alreadyDiscoveredCoords = append(alreadyDiscoveredCoords, uniqSurrounding...)
					lowestIdxs = append(lowestIdxs, lc)
					break
				}
			}
		}
	}

	print(fmt.Sprintf("LowestIndexes %d", len(lowestIdxs)))
	sort.Slice(lowestIdxs, func(i int, j int) bool {
		return len(lowestIdxs[i].surroundingValues) > len(lowestIdxs[j].surroundingValues)
	})
	curMultiple := 1
	for i := 0; i < 3; i++ {
		curMultiple *= (len(lowestIdxs[i].surroundingValues) + 1)
	}
	print(fmt.Sprintf("Multiple: %d\n", curMultiple))
}

func main() {
	heightMap := loadInput()

	yLen := len(heightMap)
	xLen := len(heightMap[0])
	lowCoords := challenge1(heightMap, yLen, xLen)
	challenge2(heightMap, yLen, xLen, lowCoords)
}

// 85008 too low
