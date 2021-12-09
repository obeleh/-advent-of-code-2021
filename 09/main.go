package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
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

type lowCoord struct {
	x   int
	y   int
	val int
}

func valIsLower(val int, xC int, yC int, xLen int, yLen int, heightMap [][]int) bool {
	if xC >= 0 && yC >= 0 && xC < xLen && yC < yLen {
		return val < heightMap[yC][xC]
	} else {
		// return true when out of bounds
		return true
	}
}

func main() {
	heightMap := loadInput()

	lowestIdxs := make([]lowCoord, 0)
	yLen := len(heightMap)
	xLen := len(heightMap[0])
	for yIdx, yVals := range heightMap {
		for xIdx, val := range yVals {
			if val == 9 {
				continue
			}

			// loop around, start top left go clockwise
			// skip to next field if current value is not lower than any of the surrounding

			// if !valIsLower(val, xIdx-1, yIdx-1, xLen, yLen, heightMap) {
			// 	continue // topleft
			// }
			if !valIsLower(val, xIdx, yIdx-1, xLen, yLen, heightMap) {
				continue // top
			}
			// if !valIsLower(val, xIdx+1, yIdx-1, xLen, yLen, heightMap) {
			// 	continue // topright
			// }
			if !valIsLower(val, xIdx+1, yIdx, xLen, yLen, heightMap) {
				continue // right
			}
			// if !valIsLower(val, xIdx+1, yIdx+1, xLen, yLen, heightMap) {
			// 	continue // bottomright
			// }
			if !valIsLower(val, xIdx, yIdx+1, xLen, yLen, heightMap) {
				continue // bottom
			}
			// if !valIsLower(val, xIdx-1, yIdx+1, xLen, yLen, heightMap) {
			// 	continue // bottomleft
			// }
			if !valIsLower(val, xIdx-1, yIdx, xLen, yLen, heightMap) {
				continue // left
			}

			// we checked all around, val is the lowest in the area
			lowestIdxs = append(lowestIdxs, lowCoord{
				x:   xIdx,
				y:   yIdx,
				val: val,
			})
		}
	}

	outputSum := 0
	for _, lC := range lowestIdxs {
		outputSum += lC.val + 1
	}

	print(fmt.Sprintf("lcSum %d\n", outputSum))
}

// 649 too high
