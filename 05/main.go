package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func loadInput() []string {
	bytes, err := os.ReadFile("input.txt")
	if err != nil {
		log.Fatalln(err)
	}
	return strings.Split(string(bytes), "\n")
}

func NewVent(line string) (*vent, error) {
	// 0,9 -> 5,9
	words := strings.Split(line, " -> ")
	fromCoords := strings.Split(words[0], ",")
	toCoords := strings.Split(words[1], ",")
	fromX, err := strconv.Atoi(fromCoords[0])
	if err != nil {
		return nil, err
	}
	fromY, err := strconv.Atoi(fromCoords[1])
	if err != nil {
		return nil, err
	}
	toX, err := strconv.Atoi(toCoords[0])
	if err != nil {
		return nil, err
	}
	toY, err := strconv.Atoi(toCoords[1])
	if err != nil {
		return nil, err
	}
	v := vent{
		fromX: fromX,
		fromY: fromY,
		toX:   toX,
		toY:   toY,
	}
	return &v, nil
}

type vent struct {
	fromX int
	fromY int
	toX   int
	toY   int
}

func loadVents(lines []string) (int, int, []*vent, error) {
	maxX := 0
	maxY := 0
	vents := make([]*vent, 0)
	for _, line := range lines {
		vent, err := NewVent(line)
		if err != nil {
			return 0, 0, nil, err
		}
		vents = append(vents, vent)

		// get board dims so we can size arrays correctly
		// just saved out a loop :P
		maxX = maxInt(maxInt(maxX, vent.toX), vent.fromX)
		maxY = maxInt(maxInt(maxY, vent.toY), vent.fromY)
	}
	return maxX, maxY, vents, nil
}

func maxInt(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func initiateVentSums(maxX int, maxY int) [][]int {
	ventSums := make([][]int, maxX+1)
	for x := range ventSums {
		ventSums[x] = make([]int, maxY+1)
	}
	return ventSums
}

func drawVentLines(vents []*vent, ventSums [][]int, onlyConsiderHvLines bool) {
	for _, vent := range vents {
		// x1 = x2 or y1 = y2.
		isDiagonal := !(vent.fromX == vent.toX || vent.fromY == vent.toY)
		if onlyConsiderHvLines && isDiagonal {
			continue
		}

		var xStep, yStep int
		if vent.fromX < vent.toX {
			xStep = 1
		} else if vent.fromX > vent.toX {
			xStep = -1
		} else {
			xStep = 0
		}
		if vent.fromY < vent.toY {
			yStep = 1
		} else if vent.fromY > vent.toY {
			yStep = -1
		} else {
			yStep = 0
		}

		x := vent.fromX
		y := vent.fromY

		for {
			ventSums[x][y] += 1
			if x == vent.toX && y == vent.toY {
				break
			}
			x += xStep
			y += yStep
		}
	}
}

func countIntersections(vents []*vent, ventSums [][]int, minScore int) {
	width := len(ventSums)
	height := len(ventSums[0])
	intersectionCount := 0
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if ventSums[x][y] >= minScore {
				intersectionCount += 1
			}
		}
	}
	print(fmt.Sprintf("found %d intersections\n", intersectionCount))
}

func main() {
	lines := loadInput()
	maxX, maxY, vents, err := loadVents(lines)
	if err != nil {
		log.Fatalf("Failed loading vents %v", err)
	}
	ventSums := initiateVentSums(maxX, maxY)
	drawVentLines(vents, ventSums, true)
	countIntersections(vents, ventSums, 2)
	ventSums = initiateVentSums(maxX, maxY)
	drawVentLines(vents, ventSums, false)
	countIntersections(vents, ventSums, 2)
}

// found 5124 intersections
// found 19771 intersections
