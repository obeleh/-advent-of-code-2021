package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type Substitution struct {
	from   string
	insert string
}

func loadInput() ([]*string, []Substitution) {
	bytes, err := os.ReadFile("input.txt")
	if err != nil {
		log.Fatalln(err)
	}
	lines := strings.Split(string(bytes), "\n")

	polymer := strings.Split(lines[0], "")
	substitutions := make([]Substitution, len(lines[2:]))
	for i, line := range lines[2:] {
		words := strings.Split(line, "->")
		from, insert := strings.TrimSpace(words[0]), strings.TrimSpace(words[1])
		substitutions[i] = Substitution{
			from:   from,
			insert: insert,
		}
	}

	polyPointer := make([]*string, len(polymer))
	for i := range polymer {
		polyPointer[i] = &polymer[i]
	}
	return polyPointer, substitutions
}

func runStep(polymer []*string, substitutions []Substitution) []*string {
	polymerStr := ""
	for _, s := range polymer {
		polymerStr += *s
	}

	for _, substitution := range substitutions {
		startIdx := 0
		for {
			pos := strings.Index(polymerStr[startIdx:], substitution.from)

			if pos == -1 {
				break
			}

			replacement := string(substitution.from[0]) + substitution.insert
			polymer[pos+startIdx] = &replacement
			startIdx += pos + 1
		}

	}

	newPolymer := []*string{}
	for _, s := range polymer {
		if len(*s) == 1 {
			newPolymer = append(newPolymer, s)
		} else if len(*s) == 2 {
			s0 := string((*s)[0])
			s1 := string((*s)[1])
			newPolymer = append(newPolymer, &s0, &s1)
		} else {
			log.Fatalln("Unexpected replacement lenght")
		}
	}

	return newPolymer
}

func countChars(polymer []*string) map[string]int {
	output := map[string]int{}
	for _, s := range polymer {
		key := *s
		_, found := output[key]
		if found {
			output[key] += 1
		} else {
			output[key] = 1
		}
	}
	return output
}

func printResult(counts map[string]int) {
	minVal := 99999999999999
	maxVal := 0
	for _, value := range counts {
		if value < minVal {
			minVal = value
		}

		if value > maxVal {
			maxVal = value
		}
	}

	print(fmt.Sprintf("Min: %d Max: %d Diff: %d\n", minVal, maxVal, maxVal-minVal))
}

func challenge1(polymer []*string, substitutions []Substitution) {

	for step := 0; step < 10; step++ {
		print(fmt.Sprintf("Step %d\n", step))
		polymer = runStep(polymer, substitutions)
	}

	counts := countChars(polymer)
	printResult(counts)
}

func getCounts(edge string, tranformations map[string][]string, substitutionMap map[string]string, counts *map[string]int, depth int) {
	(*counts)[substitutionMap[edge]] += 1
	if depth > 0 {
		getCounts(tranformations[edge][0], tranformations, substitutionMap, counts, depth-1)
		getCounts(tranformations[edge][1], tranformations, substitutionMap, counts, depth-1)
	}
}

func challenge2(polymer []*string, substitutions []Substitution) {
	tranformations := map[string][]string{}
	counts := map[string]int{}
	substitutionMap := map[string]string{}
	for _, s := range polymer {
		// make sure all _initial_ letters are in the counts map
		counts[*s] = 0
	}

	for _, substitution := range substitutions {
		possibilities := []string{
			string(substitution.from[0]) + substitution.insert,
			substitution.insert + string(substitution.from[1]),
		}
		// make sure all letters are in the counts map
		counts[string(substitution.from[0])] = 0
		counts[string(substitution.from[1])] = 0
		counts[substitution.insert] = 0

		tranformations[substitution.from] = possibilities
		substitutionMap[substitution.from] = substitution.insert
	}

	for i := 0; i < len(polymer)-1; i++ {
		edge := *polymer[i] + *polymer[i+1]
		getCounts(edge, tranformations, substitutionMap, &counts, 22-1)
	}

	for _, s := range polymer {
		counts[*s] += 1
	}

	printResult(counts)
}

func main() {
	polymer, substitutions := loadInput()
	// challenge1(polymer, substitutions)
	challenge2(polymer, substitutions)
}
