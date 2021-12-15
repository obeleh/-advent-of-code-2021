package main

import (
	"container/heap"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func expandInputMap(input [][]int) [][]int {
	factor := 5
	output := make([][]int, len(input)*factor)

	for yT := 0; yT < len(input)*factor; yT++ {
		output[yT] = make([]int, len(input[0])*factor)
	}

	for yT := 0; yT < factor; yT++ {
		for xT := 0; xT < factor; xT++ {
			for y := 0; y < len(input); y++ {
				for x := 0; x < len(input[0]); x++ {
					value := (input[y][x] + yT + xT)
					yPos := yT*len(input) + y
					xPos := xT*len(input[0]) + x
					if value > 9 {
						value = value - 9
					}
					output[yPos][xPos] = value
				}
			}
		}
	}

	return output
}

func loadInput(expandInput bool) [][]int {
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
	if expandInput {
		return expandInputMap(input)
	}
	return input
}

type coord struct {
	x   int
	y   int
	val int
}

var OFFSETCOORDS = [...]coord{
	{x: +1, y: 0},
	{x: 0, y: +1},
	{x: 0, y: -1},
	{x: -1, y: +0},
}

func getSurroundingCoords(x int, y int, m [][]int) []coord {
	yLen := len(m)
	xLen := len(m[0])
	output := make([]coord, 0)
	for _, coordOffset := range OFFSETCOORDS {
		if x+coordOffset.x >= 0 && y+coordOffset.y >= 0 && x+coordOffset.x < xLen && y+coordOffset.y < yLen {
			output = append(output, coord{
				x:   x + coordOffset.x,
				y:   y + coordOffset.y,
				val: m[y+coordOffset.y][x+coordOffset.x],
			})
		}
	}
	return output
}

func addToRiskMap(x int, y int, surCoord coord, dirRisks map[int]map[int][]coord) {
	xVals, found := dirRisks[y]
	if !found {
		xVals = map[int][]coord{}
		dirRisks[y] = xVals
	}

	coords, found := xVals[x]
	if !found {
		coords = []coord{}
	}

	xVals[x] = append(coords, surCoord)
}

func loadDirectionRisks(m [][]int) map[int]map[int][]coord {
	dirRisks := map[int]map[int][]coord{}
	yLen := len(m)
	xLen := len(m[0])
	for y := 0; y < yLen; y++ {
		for x := 0; x < xLen; x++ {
			for _, surCoord := range getSurroundingCoords(x, y, m) {
				addToRiskMap(x, y, surCoord, dirRisks)
			}
		}
	}

	return dirRisks
}

/*func findPath(x int, y int, m [][]int, hist []coord, score int) ([]coord, int) {
}*/

func pushPq(pq *PriorityQueue, c coord, coordRiskScore int) {
	item := Item{
		value:    c,
		priority: -coordRiskScore,
	}
	heap.Push(pq, &item)
	//pq.Push(&item)
}

func popPq(pq *PriorityQueue) coord {
	item := heap.Pop(pq).(*Item)
	return item.value
}

type risk struct {
	score     int
	prevCoord *coord
}

func findPath(m [][]int, directionRisks map[int]map[int][]coord) {
	pq := make(PriorityQueue, 0)
	riskSoFar := map[coord]risk{}
	destinationY := len(m) - 1
	destinationX := len(m[0]) - 1

	curCoord := coord{
		x:   0,
		y:   0,
		val: m[0][0],
	}

	riskSoFar[curCoord] = risk{
		score: 0,
	}
	for _, nextCoord := range directionRisks[curCoord.y][curCoord.x] {
		riskSoFar[nextCoord] = risk{
			score:     nextCoord.val,
			prevCoord: &curCoord,
		}
		pushPq(&pq, nextCoord, nextCoord.val)
	}

	var lowestRiskPath *risk

	for {
		if pq.Len() == 0 {
			break
		}

		curCoord = popPq(&pq)

		if curCoord.y == destinationY && curCoord.x == destinationX {
			risk := riskSoFar[curCoord]
			lowestRiskPath = &risk
			break
		}

		// loop over next available coords
		for _, nextCoord := range directionRisks[curCoord.y][curCoord.x] {
			_, found := riskSoFar[nextCoord]
			// keep risk score along the path
			nextRiskScore := riskSoFar[curCoord].score + nextCoord.val
			if !found || nextRiskScore < riskSoFar[nextCoord].score {
				if found && riskSoFar[nextCoord].prevCoord == &curCoord {
					continue
				}
				riskSoFar[nextCoord] = risk{
					score:     nextRiskScore,
					prevCoord: &curCoord,
				}

				pushPq(&pq, nextCoord, nextRiskScore)
			}
		}
	}

	print(fmt.Sprintf("Here, %d", lowestRiskPath.score))
}

func challenge1() {
	m := loadInput(false)
	directionRisks := loadDirectionRisks(m)
	findPath(m, directionRisks)
}

func challenge2() {
	m := loadInput(true)
	directionRisks := loadDirectionRisks(m)
	findPath(m, directionRisks)
}

func main() {
	challenge1()
	challenge2()
}
