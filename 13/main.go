package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type coord struct {
	x int
	y int
}

type fold struct {
	dir byte
	pos int
}

func loadInput() ([][]bool, []fold) {
	bytes, err := os.ReadFile("input.txt")
	if err != nil {
		log.Fatalln(err)
	}
	lines := strings.Split(string(bytes), "\n")
	coords := []coord{}
	maxX := 0
	maxY := 0
	folds := []fold{}
	for _, line := range lines {
		if strings.Contains(line, ",") {
			parts := strings.Split(line, ",")
			x, err := strconv.Atoi(parts[0])
			if err != nil {
				log.Fatal("Failed parsing X")
			}
			y, err := strconv.Atoi(parts[1])
			if err != nil {
				log.Fatal("Failed parsing Y")
			}
			if x > maxX {
				maxX = x
			}
			if y > maxY {
				maxY = y
			}
			coords = append(coords, coord{
				x: x,
				y: y,
			})
		} else if strings.Contains(line, "=") {
			val, err := strconv.Atoi(string(line[13:]))
			if err != nil {
				log.Fatal("Failed parsing fold")
			}
			folds = append(folds, fold{
				dir: line[11],
				pos: val,
			})
		}
	}

	paper := make([][]bool, maxY+1)
	for y := 0; y < maxY+1; y++ {
		paper[y] = make([]bool, maxX+1)
	}

	for _, coord := range coords {
		paper[coord.y][coord.x] = true
	}

	return paper, folds
}

func printPaper(paper [][]bool) {
	count := 0
	for _, xValues := range paper {
		for _, bit := range xValues {
			if bit {
				print("#")
				count++
			} else {
				print(".")
			}
		}
		print("\n")
	}
	print("\n")
	print(fmt.Sprintf("Got %d dots \n", count))
}

func foldPaper(paper [][]bool, f fold) [][]bool {
	var newPaper [][]bool
	if f.dir == 'y' {
		newHeight := f.pos - 1
		newPaper = paper[:newHeight+1]

		for y, xVals := range paper[newHeight+2:] {
			for x, val := range xVals {
				if val {
					newPaper[newHeight-y][x] = val
				}
			}
		}
	} else {
		newPaper = paper
		newWidth := f.pos - 1
		for y, xVals := range paper {
			for x, val := range xVals[newWidth+2:] {
				if val {
					newPaper[y][newWidth-x] = val
				}
			}
			newPaper[y] = paper[y][:newWidth+1]
		}

	}
	return newPaper
}

func main() {
	paper, folds := loadInput()
	for _, fold := range folds {
		paper = foldPaper(paper, fold)
		printPaper(paper)
	}
}
