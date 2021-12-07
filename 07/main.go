package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

func loadInput() ([]int, int, int, int, int) {
	bytes, err := os.ReadFile("input.txt")
	if err != nil {
		log.Fatalln(err)
	}
	strValues := strings.Split(string(bytes), ",")

	values := make([]int, len(strValues))
	var value int
	sum := 0
	count := 0
	min := 99999
	max := 0
	for idx, strValue := range strValues {
		value, err = strconv.Atoi(strValue)
		if err != nil {
			log.Fatalf("Failed loading input, %v", err)
		}
		sum += value
		count += 1
		if value < min {
			min = value
		}
		if value > max {
			max = value
		}
		values[idx] = value
	}
	return values, sum, count, min, max
}

func challenge1(values []int, avg int) {
	fuel := 0
	for _, value := range values {
		fuel += int(math.Abs(float64(avg - value)))
	}

	var prevFuel int
	for {
		prevFuel = fuel
		fuel = 0
		avg -= 1
		for _, value := range values {
			fuel += int(math.Abs(float64(avg - value)))
		}
		if fuel >= prevFuel {
			break
		}
	}

	print(fmt.Sprintf("avg:%d fuel:%d prevFuel:%d smaller:%v\n", avg, fuel, prevFuel, fuel < prevFuel))
}

func challenge2(values []int, avg int, min int, max int) {
	maxDiff := max - min
	distanceCosts := make([]int, maxDiff)
	runningTotal := 0
	for i := 1; i < maxDiff; i++ {
		// not using first array position, how wasteful!
		runningTotal += i
		distanceCosts[i] = runningTotal
	}

	var distance int
	fuel := 0
	for _, value := range values {
		distance = int(math.Abs(float64(avg - value)))
		fuel += distanceCosts[distance]
	}

	var prevFuel int
	for {
		prevFuel = fuel
		fuel = 0
		avg -= 1
		for _, value := range values {
			distance = int(math.Abs(float64(avg - value)))
			fuel += distanceCosts[distance]
		}
		if fuel >= prevFuel {
			break
		}
	}

	print(fmt.Sprintf("avg:%d fuel:%d prevFuel:%d smaller:%v\n", avg, fuel, prevFuel, fuel < prevFuel))
}

func main() {
	values, sum, count, min, max := loadInput()
	avg := sum / count

	challenge1(values, avg)
	challenge2(values, avg, min, max)
}

//avg:320 fuel:343443 prevFuel:343441 smaller:false
//avg:472 fuel:98926146 prevFuel:98925151 smaller:false
