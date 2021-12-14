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

func loadInput() ([]string, []Substitution) {
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
	return polymer, substitutions
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

func mergeCounts(counts1 map[string]int, counts2 map[string]int) map[string]int {
	output := map[string]int{}
	for key, value := range counts2 {
		c1Value, found := counts1[key]
		if found {
			output[key] = c1Value + value
		} else {
			output[key] = value
		}
	}

	for key, value := range counts1 {
		_, found := counts2[key]
		if !found {
			output[key] = value
		}
	}
	return output
}

func getCounts(edge string, tranformations map[string][]string, substitutionMap map[string]string, cache map[string]map[string]int, depth int) map[string]int {
	t := substitutionMap[edge]
	cacheKey := fmt.Sprintf("%s%s%d", edge, t, depth)
	value, found := cache[cacheKey]
	if found {
		return value
	}

	if depth == 0 {
		counts := map[string]int{
			t: 1,
		}
		return counts
	}

	counts1 := getCounts(tranformations[edge][0], tranformations, substitutionMap, cache, depth-1)
	counts2 := getCounts(tranformations[edge][1], tranformations, substitutionMap, cache, depth-1)
	counts := mergeCounts(counts1, counts2)
	vl, found := counts[t]
	if found {
		counts[t] = vl + 1
	} else {
		counts[t] = 1
	}
	cache[cacheKey] = counts
	return counts
}

func challenge(polymer []string, substitutions []Substitution, iterations int) {
	tranformations := map[string][]string{}
	substitutionMap := map[string]string{}
	for _, substitution := range substitutions {
		possibilities := []string{
			string(substitution.from[0]) + substitution.insert,
			substitution.insert + string(substitution.from[1]),
		}

		tranformations[substitution.from] = possibilities
		substitutionMap[substitution.from] = substitution.insert
	}

	counts := map[string]int{}
	cache := map[string]map[string]int{}
	for i := 0; i < len(polymer)-1; i++ {
		edge := polymer[i] + polymer[i+1]
		countsI := getCounts(edge, tranformations, substitutionMap, cache, iterations-1)
		counts = mergeCounts(counts, countsI)
	}

	for _, s := range polymer {
		counts[s] += 1
	}

	printResult(counts)
}

func main() {
	polymer, substitutions := loadInput()
	challenge(polymer, substitutions, 10)
	challenge(polymer, substitutions, 40)
}
