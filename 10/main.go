package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

func loadInput() []string {
	bytes, err := os.ReadFile("input.txt")
	if err != nil {
		log.Fatalln(err)
	}
	lines := strings.Split(string(bytes), "\n")
	return lines
}

var CLOSING_SCORES = map[string]int{
	")": 3,
	"]": 57,
	"}": 1197,
	">": 25137,
}

var AUTOCOMPLETE_SCORES = map[string]int{
	")": 1,
	"]": 2,
	"}": 3,
	">": 4,
}

var CHUNK_SETS = map[string]string{
	"(": ")",
	"[": "]",
	"{": "}",
	"<": ">",
}

type StackFrame struct {
	Start       string
	ExpectedEnd string
}

func push(frame *StackFrame, stack []*StackFrame) []*StackFrame {
	return append(stack, frame)
}

func pop(stack []*StackFrame) ([]*StackFrame, *StackFrame, *StackFrame) {
	var lastFrame *StackFrame
	if len(stack) > 0 {
		lastFrame = stack[len(stack)-1]
	}
	stack = stack[:len(stack)-1]

	lastIdx := len(stack) - 1
	var newLastFrame *StackFrame
	if lastIdx >= 0 {
		newLastFrame = stack[lastIdx]
	}
	return stack, newLastFrame, lastFrame
}

func isStartingChar(c string) bool {
	_, found := CHUNK_SETS[c]
	return found
}

func openFrame(c string, curFrame *StackFrame, stack []*StackFrame) (*StackFrame, []*StackFrame) {
	expectedEnd, found := CHUNK_SETS[c]
	if !found {
		log.Fatalf("Closing character should have been handled before opening Frame")
	}
	curFrame = &StackFrame{
		Start:       c,
		ExpectedEnd: expectedEnd,
	}
	stack = push(curFrame, stack)
	return curFrame, stack
}

func determineRest(stack []*StackFrame) string {
	rest := ""
	var frame *StackFrame
	for {
		if len(stack) > 0 {
			stack, _, frame = pop(stack)
			rest += frame.ExpectedEnd
		} else {
			break
		}
	}

	return rest
}

func main() {
	lines := loadInput()
	score := 0
	completionScores := make([]uint64, 0)

LINES:
	for _, line := range lines {
		var curFrame *StackFrame
		stack := make([]*StackFrame, 0)
		for cI, c := range strings.Split(line, "") {
			if curFrame == nil {
				if isStartingChar(c) {
					curFrame, stack = openFrame(c, curFrame, stack)
				} else {
					log.Fatalf("Unexpected end")
				}
			} else {
				if isStartingChar(c) {
					curFrame, stack = openFrame(c, curFrame, stack)
				} else {
					if c != curFrame.ExpectedEnd {
						score += CLOSING_SCORES[c]
						print(fmt.Sprintf("Line: %s   %s expected %s got %s\n", line[:cI], line[cI:], curFrame.ExpectedEnd, c))
						continue LINES
					}
					stack, curFrame, _ = pop(stack)
				}
			}
		}
		rest := determineRest(stack)
		if len(rest) > 0 {
			lineAutoCompleteScore := uint64(0)
			print(fmt.Sprintf("Line: %s Complete by adding %s\n", line, rest))
			for _, c := range strings.Split(rest, "") {
				curCharScore, wasClosingChar := AUTOCOMPLETE_SCORES[c]
				if wasClosingChar {
					lineAutoCompleteScore *= 5
					lineAutoCompleteScore += uint64(curCharScore)
				}
			}
			print(fmt.Sprintf("Rest %s Got score %d\n", rest, score))
			completionScores = append(completionScores, uint64(lineAutoCompleteScore))
		}
	}
	print(fmt.Sprintf("Score %d\n", score))
	sort.Slice(completionScores, func(i, j int) bool {
		return completionScores[i] > completionScores[j]
	})
	middle := completionScores[len(completionScores)/2]
	print(fmt.Sprintf("AutoCompletionScore %d\n", middle))
}
