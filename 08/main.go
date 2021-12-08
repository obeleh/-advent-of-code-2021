package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/thoas/go-funk"
)

type digit struct {
	chars string
	digit int
}

type registrationLine struct {
	example []*digit
	output  []*digit
}

func wordsToDigits(words []string) []*digit {
	return funk.Map(words, func(word string) *digit {
		// sort chars
		s := strings.Split(word, "")
		sort.Strings(s)

		d := digit{
			chars: strings.Join(s, ""),
			digit: -1,
		}
		return &d
	}).([]*digit)
}

func segmentDelta(word1 string, word2 string) (int, int) {
	w1 := strings.Split(word1, "")
	w2 := strings.Split(word2, "")

	d1, d2 := funk.DifferenceString(w1, w2)
	return len(d1), len(d2)
}

func loadInput() []*registrationLine {
	bytes, err := os.ReadFile("input.txt")
	if err != nil {
		log.Fatalln(err)
	}
	lines := strings.Split(string(bytes), "\n")

	output := make([]*registrationLine, len(lines))
	for idx, line := range lines {
		parts := strings.Split(line, "|")
		rl := registrationLine{
			example: wordsToDigits(strings.Fields(parts[0])),
			output:  wordsToDigits(strings.Fields(parts[1])),
		}
		output[idx] = &rl
	}
	return output
}

func decodeDigit(d *digit) bool {
	switch len(d.chars) {
	case 2:
		d.digit = 1
		return true
	case 4:
		d.digit = 4
		return true
	case 3:
		d.digit = 7
		return true
	case 7:
		d.digit = 8
		return true
	}
	return false
}

func RemoveIndex(s []*digit, index int) []*digit {
	return append(s[:index], s[index+1:]...)
}

func showPatternDifs() {
	// Taking an 8 going clockwise starting topleft and with ending the horizontal bar in the middle as last
	digits := []*digit{
		&digit{chars: "012345", digit: 0},
		&digit{chars: "12", digit: 1},
		&digit{chars: "01346", digit: 2},
		&digit{chars: "01236", digit: 3},
		&digit{chars: "1256", digit: 4},
		&digit{chars: "02356", digit: 5},
		&digit{chars: "023456", digit: 6},
		&digit{chars: "012", digit: 7},
		&digit{chars: "0123456", digit: 8},
		&digit{chars: "012356", digit: 9},
	}

	for _, digit := range digits {
		for _, compareDigit := range digits {
			d1, d2 := segmentDelta(digit.chars, compareDigit.chars)
			print(fmt.Sprintf("%d -> %d [%d,%d]\n", digit.digit, compareDigit.digit, d1, d2))
		}
		print("\n")
	}
}

func deductDigit(digitMap map[int]*digit, unFound []*digit, fromDigit int, foundDigit int, d1r int, d2r int) []*digit {
	found := false
	var newUnFound []*digit
	for idx, uf := range unFound {
		d1, d2 := segmentDelta(digitMap[fromDigit].chars, uf.chars)
		if d1 == d1r && d2 == d2r {
			digitMap[foundDigit] = uf
			uf.digit = foundDigit
			if found {
				log.Fatal("Got duplicate")
			}
			found = true
			newUnFound = RemoveIndex(unFound, idx)
		}
	}
	if !found {
		log.Fatal("Could not find idx")
	}
	return newUnFound
}

func guessDigits(rl *registrationLine) int {
	// Because the digits 1, 4, 7, and 8 each use a unique number of segments
	count := 0
	unFound := make([]*digit, 0)
	digitMap := make(map[int]*digit)
	for _, d := range rl.example {
		found := decodeDigit(d)
		if !found {
			unFound = append(unFound, d)
		} else {
			digitMap[d.digit] = d
		}
	}

	for _, d := range rl.output {
		found := decodeDigit(d)
		if found {
			count += 1
		}
	}

	unFound = deductDigit(digitMap, unFound, 1, 3, 0, 3) // 1 -> 3 [0,3]
	unFound = deductDigit(digitMap, unFound, 1, 6, 1, 5) // 1 -> 6 [1,5]
	unFound = deductDigit(digitMap, unFound, 4, 2, 2, 3) // 4 -> 2 [2,3]
	unFound = deductDigit(digitMap, unFound, 4, 9, 0, 2) // 4 -> 9 [0,2]
	unFound = deductDigit(digitMap, unFound, 2, 5, 2, 2) // 2 -> 5 [2,2]
	unFound = deductDigit(digitMap, unFound, 5, 0, 1, 2) // 2 -> 5 [1,2]

	if len(unFound) != 0 {
		log.Fatalf("Still matches left???")
	}

	return count
}

func decodeOutput(rl *registrationLine) int {
	output := 0
	for _, d := range rl.output {
		if d.digit == -1 {
			// lookup from the examples
			for _, ex := range rl.example {
				if ex.chars == d.chars {
					output = output*10 + ex.digit
					break
				}
			}
		} else {
			output = output*10 + d.digit
		}

	}
	return output
}

func main() {
	showPatternDifs()
	rls := loadInput()

	challenge1 := 0
	challenge2 := 0
	for _, rl := range rls {
		challenge1 += guessDigits(rl)
		challenge2 += decodeOutput(rl)
	}
	print(fmt.Sprintf("challenge1 %d\n", challenge1))
	print(fmt.Sprintf("challenge2 %d\n", challenge2))
}
