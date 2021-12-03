package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/thoas/go-funk"
)

const NUM_BITS = 12
const LINECOUNT = 1000

func loadInput() []int {
	bytes, err := os.ReadFile("03.input")
	if err != nil {
		log.Fatalln(err)
	}
	values := make([]int, LINECOUNT)
	for idx, line := range strings.Split(string(bytes), "\n") {
		lineStripped := strings.TrimSpace(line)
		words := strings.Fields(lineStripped)
		if len(words) != 1 {
			log.Fatalln("Expected 1 words to be present per line")
		}

		// NUM_BITS+1 because ParseInt parses into signed ints
		value64, err := strconv.ParseInt(words[0], 2, NUM_BITS+1)
		if err != nil {
			log.Fatalln(fmt.Errorf("expected the value to be 12 bits,  %v", err))
		}
		values[idx] = int(value64)
	}
	return values
}

func getValuesWithBitSet(mask int, values []int, more bool) []int {
	count := 0
	halfCount := float32(len(values)) / 2.0
	for _, value := range values {
		if mask&value == mask {
			count += 1
		}
	}

	moreOnes := float32(count) >= halfCount
	weWantOnes := moreOnes == more

	return funk.FilterInt(values, func(value int) bool {
		bitIsSet := mask&value == mask
		return bitIsSet == weWantOnes
	})
}

func filterDownValues(masks []int, values []int, more bool) int {
	workingList := make([]int, LINECOUNT)
	copy(workingList, values)
	for i := NUM_BITS - 1; i >= 0; i-- {
		workingList = getValuesWithBitSet(masks[i], workingList, more)
		if len(workingList) == 1 {
			break
		}
	}
	return workingList[0]
}

func main() {

	masks := make([]int, NUM_BITS)

	for shift := 0; shift < NUM_BITS; shift++ {
		masks[shift] = 1 << shift
	}

	values := loadInput()
	oxigenValue := filterDownValues(masks, values, true)
	scrubberValue := filterDownValues(masks, values, false)
	print(fmt.Sprintf("oxigen:%d scrubber:%d multiplied:%d\n", oxigenValue, scrubberValue, oxigenValue*scrubberValue))
}

// oxigen:1679 scrubber:3648 multiplied:6124992
