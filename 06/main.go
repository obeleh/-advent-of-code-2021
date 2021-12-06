package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const NUMDAYS int = 256

func loadInput() []int {
	bytes, err := os.ReadFile("input.txt")
	if err != nil {
		log.Fatalln(err)
	}
	strValues := strings.Split(string(bytes), ",")

	values := make([]int, len(strValues))
	for idx, strValue := range strValues {
		values[idx], err = strconv.Atoi(strValue)
		if err != nil {
			log.Fatalf("Failed loading input, %v", err)
		}
	}
	return values
}

func main() {
	fish := loadInput()
	fishyFish := make([]*[]int, NUMDAYS+1)
	fishyFish[0] = &fish
	var dayArr []int
	for day := 0; day < NUMDAYS; day++ {
		newFish := make([]int, 0)
		for dayArrIdx := 0; dayArrIdx < day; dayArrIdx++ {
			dayArr = *fishyFish[dayArrIdx]
			for fishIdx := 0; fishIdx < len(dayArr); fishIdx++ {
				dayArr[fishIdx] -= 1
				if dayArr[fishIdx] < 0 {
					dayArr[fishIdx] = 6
					newFish = append(newFish, 8)
				}
			}
		}
		fishyFish[day+1] = &newFish
		print(day)
		print("\n")
	}

	cnt := 0
	for _, curArr := range fishyFish {
		cnt += len(*curArr)
	}
	print(fmt.Sprintf("%d fish", len(fish)))
}

// found 5124 intersections
// found 19771 intersections
