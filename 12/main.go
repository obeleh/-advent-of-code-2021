package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/thoas/go-funk"
)

func mapAppend(key string, value string, edgeMap map[string][]string) {
	_, exists := edgeMap[key]
	if exists {
		edgeMap[key] = append(edgeMap[key], value)
	} else {
		edgeMap[key] = []string{value}
	}
}

func loadInput() map[string][]string {
	bytes, err := os.ReadFile("input.txt")
	if err != nil {
		log.Fatalln(err)
	}
	lines := strings.Split(string(bytes), "\n")
	edgeMap := make(map[string][]string)
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		parts := strings.Split(line, "-")
		// add both directions
		mapAppend(parts[0], parts[1], edgeMap)
		if parts[0] != "start" {
			mapAppend(parts[1], parts[0], edgeMap)
		}
	}
	return edgeMap
}

func findRoutes1(curNode string, curRoute []string, edgeMap map[string][]string) [][]string {
	output := make([][]string, 0)
	curRoute = append(curRoute, curNode)
	for _, nextNode := range edgeMap[curNode] {
		if nextNode == strings.ToLower(nextNode) && funk.Contains(curRoute, nextNode) {
			// skip small caverns
			continue
		}

		nextRoute := append([]string{}, curRoute...)
		if nextNode == "end" {
			// finish up
			nextRoute = append(nextRoute, "end")
			output = append(output, nextRoute)
			continue
		}

		output = append(output, findRoutes1(nextNode, nextRoute, edgeMap)...)
	}
	return output
}

func smallNodeOkForRoute(node string, route []string) bool {
	if node == "start" {
		return false
	}
	if !funk.Contains(route, node) {
		return true
	}

	counts := make(map[string]int)
	for _, routeNode := range route {
		if routeNode != strings.ToLower(routeNode) {
			continue
		}
		_, found := counts[routeNode]
		if found {
			return false
		} else {
			counts[routeNode] = 1
		}
	}
	return true
}

func findRoutes2(curNode string, curRoute []string, edgeMap map[string][]string) [][]string {
	output := make([][]string, 0)
	curRoute = append(curRoute, curNode)
	for _, nextNode := range edgeMap[curNode] {
		if nextNode == strings.ToLower(nextNode) && !smallNodeOkForRoute(nextNode, curRoute) {
			// skip small caverns
			continue
		}

		nextRoute := append([]string{}, curRoute...)
		if nextNode == "end" {
			// finish up
			nextRoute = append(nextRoute, "end")
			output = append(output, nextRoute)
			continue
		}

		output = append(output, findRoutes2(nextNode, nextRoute, edgeMap)...)
	}
	return output
}

func main() {
	edgeMap := loadInput()
	startNode := "start"
	startRoute := []string{}
	routes := findRoutes1(startNode, startRoute, edgeMap)
	for _, route := range routes {
		fmt.Println(strings.Join(route, ","))
	}
	print(fmt.Sprintf("Found %d paths\n\n", len(routes)))

	routes = findRoutes2(startNode, startRoute, edgeMap)
	for _, route := range routes {
		fmt.Println(strings.Join(route, ","))
	}
	print(fmt.Sprintf("Found %d paths\n\n", len(routes)))
}

// 118234 too high
