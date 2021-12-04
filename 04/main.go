package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/thoas/go-funk"
)

func loadInput() []string {
	bytes, err := os.ReadFile("input.txt")
	if err != nil {
		log.Fatalln(err)
	}
	return strings.Split(string(bytes), "\n")
}

func loadDraws(line string) []int {
	words := strings.Split(line, ",")
	return wordsToIntarray(words)
}

type board struct {
	numbers []int
	marks   []bool
	dim     int
	bingoed bool
}

func wordsToIntarray(words []string) []int {
	return funk.Map(words, func(word string) int {
		value, err := strconv.Atoi(strings.TrimSpace(word))
		if err != nil {
			log.Fatalf("Failed parsing numbers %v", err)
		}
		return value
	}).([]int)
}

func lineToInts(line string) []int {
	words := strings.Fields(line)
	return wordsToIntarray(words)
}

func loadBoardNumbers(numbers []int, newLine string) []int {
	return append(numbers, lineToInts(newLine)...)
}

func NewBoard(numbers []int) *board {
	if len(numbers) == 0 {
		log.Fatal("Expected a slice of numbers")
	}
	marks := make([]bool, len(numbers))
	board := board{
		numbers: numbers,
		marks:   marks,
		dim:     int(math.Sqrt(float64(len(numbers)))),
	}
	return &board
}

func bingo(b *board) bool {
	// check rows
	for r := 0; r < b.dim; r++ {
		markCnt := 0
		for c := 0; c < b.dim; c++ {
			if b.marks[r*b.dim+c] {
				markCnt += 1
			} else {
				break
			}
		}
		if markCnt == b.dim {
			b.bingoed = true
			return true
		}
	}
	// check columns
	for c := 0; c < b.dim; c++ {
		markCnt := 0
		for r := 0; r < b.dim; r++ {
			if b.marks[r*b.dim+c] {
				markCnt += 1
			} else {
				break
			}
		}
		if markCnt == b.dim {
			b.bingoed = true
			return true
		}
	}
	return false
}

func addDraw(b *board, draw int) {
	for r := 0; r < b.dim; r++ {
		for c := 0; c < b.dim; c++ {
			if b.numbers[r*b.dim+c] == draw {
				b.marks[r*b.dim+c] = true
				return
			}
		}
	}
}

func getUnMarkedNumbers(b *board) []int {
	unMarkedNumbers := make([]int, 0)
	for r := 0; r < b.dim; r++ {
		for c := 0; c < b.dim; c++ {
			if !b.marks[r*b.dim+c] {
				unMarkedNumbers = append(unMarkedNumbers, b.numbers[r*b.dim+c])
			}
		}
	}
	return unMarkedNumbers
}

func printOutcome(draw int, board *board, drawIdx int) {
	marks := getUnMarkedNumbers(board)
	sum := funk.SumInt(marks)
	print(fmt.Sprintf("bingo on draw %d with number %d unmarked %v %d\n", drawIdx, draw, marks, sum*draw))
}

func challenge1(draws []int, boards []*board) {
	for drawIdx, draw := range draws {
		for _, board := range boards {
			addDraw(board, draw)
			if bingo(board) {
				printOutcome(draw, board, drawIdx)
				return
			}
		}
	}
}

func challenge2(draws []int, boards []*board) {
	for drawIdx, draw := range draws {
		for _, board := range boards {
			if board.bingoed {
				continue
			}
			addDraw(board, draw)
			bingo(board)
		}

		if len(boards) == 1 && boards[0].bingoed {
			printOutcome(draw, boards[0], drawIdx)
			return
		}

		boards = funk.Filter(boards, func(board *board) bool {
			return !board.bingoed
		}).([]*board)
	}
}

func main() {
	lines := loadInput()
	draws := loadDraws(lines[0])
	boards := make([]*board, 0)

	curNumbers := make([]int, 0)
	for _, line := range lines[2:] {
		if len(line) > 0 {
			curNumbers = loadBoardNumbers(curNumbers, line)
		} else {
			board := NewBoard(curNumbers)
			boards = append(boards, board)
			curNumbers = make([]int, 0)
		}
	}

	challenge1(draws, boards)
	challenge2(draws, boards)
}
