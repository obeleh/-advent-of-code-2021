package main

import "fmt"

type targetArea struct {
	fromX int
	toX   int
	fromY int
	toY   int
}

type coord struct {
	x int
	y int
}

var PUZZLE_INPUT = targetArea{
	fromX: 265,
	toX:   287,
	fromY: -58,
	toY:   -103,
}
var EXAMPLE_INPUT = targetArea{
	fromX: 20,
	toX:   30,
	fromY: -5,
	toY:   -10,
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

const STATUS_PRE_T = 0
const STATUS_IN_T = 1
const STATUS_POST_T = 2

func (t targetArea) getStatus(x int, y int) (int, int) {
	var xStatus int
	var yStatus int
	if x > t.toX {
		xStatus = STATUS_POST_T
	} else if x < t.fromX {
		xStatus = STATUS_PRE_T
	} else {
		xStatus = STATUS_IN_T
	}

	if y > t.fromY {
		yStatus = STATUS_PRE_T
	} else if y < t.toY {
		yStatus = STATUS_POST_T
	} else {
		yStatus = STATUS_IN_T
	}

	return xStatus, yStatus
}

func calculateHit(xVel int, yVel int, t targetArea) (bool, *coord, int) {
	x := 0
	y := 0

	maxY := y

	for {
		if y > maxY {
			maxY = y
		}
		xStatus, yStatus := t.getStatus(x, y)
		if xStatus == STATUS_POST_T || yStatus == STATUS_POST_T {
			return false, &coord{
				x: x,
				y: y,
			}, maxY
		}
		if xStatus == STATUS_IN_T && yStatus == STATUS_IN_T {
			return true, &coord{
				x: x,
				y: y,
			}, maxY
		}

		x += xVel
		y += yVel

		xVel = max(xVel-1, 0)
		yVel -= 1
	}
}

const FINDING_FIRST_HIT = 0

type hit struct {
	xV int
	yV int
	c  *coord
}

func main() {
	t := PUZZLE_INPUT

	hits := []hit{}
	maxY := 0
	loopCnt := 600
	for xV := -loopCnt; xV < loopCnt; xV++ {
		for yV := -loopCnt; yV < loopCnt; yV++ {
			//print(fmt.Sprintf("x:%d y:%d\n", xV, yV))
			isHit, c, curMaxY := calculateHit(xV, yV, t)
			if isHit {
				hits = append(hits, hit{
					xV: xV,
					yV: yV,
					c:  c,
				})
				if curMaxY > maxY {
					maxY = curMaxY
				}
			}
		}
	}

	print(fmt.Sprintf("MaxY: %d HitCount:%d\n", maxY, len(hits)))
}
