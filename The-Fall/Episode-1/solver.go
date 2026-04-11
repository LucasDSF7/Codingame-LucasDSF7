package main

/**
Copy and paste this code in https://www.codingame.com/ide/puzzle/the-fall-episode-1
**/

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var Room = map[int]map[string][2]int{
	0: {
		"": {-1, -1},
	},
	1: {
		"TOP":   {1, 0},
		"LEFT":  {1, 0},
		"RIGHT": {1, 0},
	},
	2: {
		"LEFT":  {0, 1},
		"RIGHT": {0, -1},
	},
	3: {
		"TOP": {1, 0},
	},
	4: {
		"TOP":   {0, -1},
		"RIGHT": {1, 0},
		"LEFT":  {-1, -1},
	},
	5: {
		"TOP":   {0, 1},
		"LEFT":  {1, 0},
		"RIGHT": {-1, -1},
	},
	6: {
		"TOP":   {-1, -1},
		"LEFT":  {0, 1},
		"RIGHT": {0, -1},
	},
	7: {
		"TOP":   {1, 0},
		"RIGHT": {1, 0},
	},
	8: {
		"LEFT":  {1, 0},
		"RIGHT": {1, 0},
	},
	9: {
		"LEFT": {1, 0},
		"TOP":  {1, 0},
	},
	10: {
		"TOP":  {0, -1},
		"LEFT": {-1, -1},
	},
	11: {
		"TOP":   {0, 1},
		"RIGHT": {-1, -1},
	},
	12: {
		"RIGHT": {1, 0},
	},
	13: {
		"LEFT": {1, 0},
	},
}

func split_integers(str string) []int {
	str_split := strings.Split(str, " ")
	var int_array []int
	for _, s := range str_split {
		integer, err := strconv.Atoi(s)
		_ = err
		int_array = append(int_array, integer)
	}
	return int_array
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1000000), 1000000)

	// W: number of columns.
	// H: number of rows.
	var W, H int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &W, &H)

	var Grid [][]int

	for i := 0; i < H; i++ {
		scanner.Scan()
		Line := scanner.Text()
		fmt.Fprintln(os.Stderr, Line)
		Grid = append(Grid, split_integers(Line))
	}
	// EX: the coordinate along the X axis of the exit (not useful for this first mission, but must be read).
	var EX int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &EX)
	for {
		var XI, YI int // INVERTED
		var POS string
		scanner.Scan()
		fmt.Sscan(scanner.Text(), &XI, &YI, &POS)
		fmt.Fprintln(os.Stderr, XI, YI, POS)
		fmt.Fprintln(os.Stderr, Grid)
		RoomType := Grid[YI][XI]
		move := Room[RoomType][POS]
		NewX, NewY := XI+move[1], YI+move[0]

		// fmt.Fprintln(os.Stderr, "Debug messages...")

		// One line containing the X Y coordinates of the room in which you believe Indy will be on the next turn.
		fmt.Println(NewX, NewY)
	}
}
