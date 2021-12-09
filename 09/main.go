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
	x   int
	y   int
	val int
}
type bassin struct {
	Coord
	surroundingValues []Coord
}

var OFFSETCOORDS = [...]Coord{
	{x: 0, y: -1},
	{x: +1, y: 0},
	{x: 0, y: +1},
	{x: -1, y: +0},
}

func getSurroundingCoords(x int, y int, heightMap [][]int, yLen int, xLen int) []Coord {
	output := make([]Coord, 0)
	for _, coordOffset := range OFFSETCOORDS {
		if x+coordOffset.x >= 0 && y+coordOffset.y >= 0 && x+coordOffset.x < xLen && y+coordOffset.y < yLen {
			output = append(output, Coord{
				x:   x + coordOffset.x,
				y:   y + coordOffset.y,
				val: heightMap[y+coordOffset.y][x+coordOffset.x],
			})
		}
	}
	return output
}

func challenge1(heightMap [][]int, yLen int, xLen int) []*bassin {
	bassins := make([]*bassin, 0)
	for yIdx, yVals := range heightMap {
	VAL_LOOP:
		for xIdx, val := range yVals {
			if val == 9 {
				continue
			}

			bassinCoords := make([]Coord, 0)
			surroundingCoords := getSurroundingCoords(xIdx, yIdx, heightMap, yLen, xLen)
			for _, coord := range surroundingCoords {
				if val >= coord.val {
					// skip to next field if current value is not lower than any of the surrounding
					continue VAL_LOOP
				}

				if coord.val < 9 {
					bassinCoords = append(bassinCoords, coord)
				}
			}

			// we checked all around, val is the lowest in the area
			b := bassin{
				Coord: Coord{
					x:   xIdx,
					y:   yIdx,
					val: val,
				},
				surroundingValues: bassinCoords,
			}
			bassins = append(bassins, &b)
		}
	}

	outputSum := 0
	for _, lC := range bassins {
		outputSum += lC.val + 1
	}

	print(fmt.Sprintf("lcSum %d\n", outputSum))
	return bassins
}

func loadBassinSurroundings(heightMap [][]int, yLen int, xLen int, bassin *bassin) {
	nextBatch := make([]Coord, len(bassin.surroundingValues))
	copy(nextBatch, bassin.surroundingValues)
	for {
		candidates := make([]Coord, 0)
		for _, coord := range nextBatch {
			surroundingCoords := getSurroundingCoords(coord.x, coord.y, heightMap, yLen, xLen)
			for _, surCoord := range surroundingCoords {
				if surCoord.val > coord.val && surCoord.val < 9 {
					if funk.Contains(bassin.surroundingValues, surCoord) {
						continue
					}
					candidates = append(candidates, surCoord)
				}
			}
		}
		if len(candidates) == 0 {
			break
		}
		bassin.surroundingValues = append(bassin.surroundingValues, candidates...)

		nextBatch = candidates
	}
	bassin.surroundingValues = funk.Uniq(bassin.surroundingValues).([]Coord)
}

func printBassin(curBassin bassin, xLen int, yLen int) {
	heightMap := make([][]int, yLen)
	for y := 0; y < yLen; y++ {
		heightMap[y] = make([]int, xLen)
	}
	heightMap[curBassin.y][curBassin.x] = curBassin.val
	for _, coord := range curBassin.surroundingValues {
		heightMap[coord.y][coord.x] = coord.val
	}
	print("\n")
	for y := 0; y < yLen; y++ {
		for x := 0; x < xLen; x++ {
			if heightMap[y][x] == 0 {
				print(" ")
			} else {
				print(heightMap[y][x])
			}
		}
		print("\n")
	}
	print("\n")
}

func challenge2(heightMap [][]int, yLen int, xLen int, bassins []*bassin) {
	for _, bassin := range bassins {
		loadBassinSurroundings(heightMap, yLen, xLen, bassin)
	}

	sort.Slice(bassins, func(i int, j int) bool {
		return len(bassins[i].surroundingValues) > len(bassins[j].surroundingValues)
	})

	curMultiple := 1
	for i := 0; i < 3; i++ {
		curBassin := bassins[i]
		curMultiple *= (len(curBassin.surroundingValues) + 1)
		printBassin(*curBassin, xLen, yLen)
	}
	print(fmt.Sprintf("Multiple: %d\n", curMultiple))
}

func main() {
	heightMap := loadInput()

	yLen := len(heightMap)
	xLen := len(heightMap[0])
	bassins := challenge1(heightMap, yLen, xLen)
	challenge2(heightMap, yLen, xLen, bassins)
}
