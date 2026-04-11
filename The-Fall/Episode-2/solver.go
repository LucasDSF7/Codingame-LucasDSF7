package main

/*
Copy and paste this code in https://www.codingame.com/ide/puzzle/the-fall-episode-2
*/

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// {0, 0} is a invalid position
// In go its the deafult value for Room, so no need to write it
var Room = map[int]map[string][2]int{
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
	},
	5: {
		"TOP":  {0, 1},
		"LEFT": {1, 0},
	},
	6: {
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
		"TOP": {0, -1},
	},
	11: {
		"TOP": {0, 1},
	},
	12: {
		"RIGHT": {1, 0},
	},
	13: {
		"LEFT": {1, 0},
	},
}

var newPos = map[[2]int]string{
	{1, 0}:  "TOP",
	{0, 1}:  "LEFT",
	{0, -1}: "RIGHT",
	{0, 0}:  "X",
}

var RotateRight = map[int]int{
	2:  3,
	3:  2,
	4:  5,
	5:  4,
	6:  7,
	7:  8,
	8:  9,
	9:  6,
	10: 11,
	11: 12,
	12: 13,
	13: 10,
}

var RotateLeft = map[int]int{
	2:  3,
	3:  2,
	4:  5,
	5:  4,
	6:  9,
	7:  6,
	8:  7,
	9:  8,
	10: 13,
	11: 10,
	12: 11,
	13: 12,
}

var RotateTwice = map[int]int{
	6:  8,
	7:  9,
	8:  6,
	9:  7,
	10: 12,
	11: 13,
	12: 10,
	13: 11,
}

var Grid [][]int

// EX: the coordinate along the Y axis of the exit.
var EX int
var Actions = make([]Action, 0, 50)

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

type Entity struct {
	Pos string // X if its destroyed
	X   int
	Y   int
	XY  [2]int
}

// Modifies and return a copy of Entity
func (e Entity) Move() Entity {
	room := int(math.Abs(float64(Grid[e.X][e.Y])))
	xy := Room[room][e.Pos]
	e.Pos = newPos[xy]
	e.X, e.Y = e.X+xy[0], e.Y+xy[1]
	e.XY = [2]int{e.X, e.Y}
	if e.X >= len(Grid) || e.Y < 0 || e.Y >= len(Grid[0]) {
		e.Pos = "X"
	}
	return e
}

type Action struct {
	Move string
	X    int
	Y    int
}

func (a Action) Anwser() {
	if a.Move == "TWICE" {
		// Translate TWICE as simple LEFT. In DPS if Depth 0 and Action is TWICE, modify Grid to
		// Rotate LEFT once.
		a.Move = "LEFT"
	}
	if a.Move == "WAIT" {
		fmt.Println("WAIT")
	} else {
		fmt.Printf("%d %d %v\n", a.Y, a.X, a.Move)
	}
}

// Depth stores the depth of DPS.
// SpareMove is a counter of how many additional moves Indy can do. SpareMoves ups with WAIT command
// if SpareMove > 0, a tile can be rotated twice
type State struct {
	Depth     int
	SpareMove int
}

// Store possible move for a game turn. depth can be either 0, - 1, fallowing logic from State.
// spareMove can be 0, 1 or - 1.
type Move struct {
	x         int
	y         int
	move      string
	depth     int
	spareMove int
}

func getCantRotate(Indy Entity, Rocks []Entity) map[[2]int]bool {
	cant := make(map[[2]int]bool)
	cant[Indy.XY] = true
	for _, r := range Rocks {
		cant[r.XY] = true
	}
	return cant
}

func getMoves(newIndy Entity, newRocks []Entity, state State, cantRotate map[[2]int]bool) []Move {
	possibleMoves := make([]Move, 0, len(newRocks)*2+4)
	indyRoom := Grid[newIndy.X][newIndy.Y]
	for _, rock := range newRocks {
		if isInvalid := cantRotate[rock.XY]; Grid[rock.X][rock.Y] < 0 || isInvalid {
			continue
		}
		for _, move := range [2]string{"LEFT", "RIGHT"} {
			newMove := Move{rock.X, rock.Y, move, 0, 0}
			possibleMoves = append(possibleMoves, newMove)
		}
	}
	for _, move := range [4]string{"WAIT", "LEFT", "RIGHT", "TWICE"} {
		if state.SpareMove == 0 && move == "TWICE" {
			continue
		}
		if indyRoom < 0 && move != "WAIT" {
			continue
		}
		newMove := Move{newIndy.X, newIndy.Y, move, 0, 0}
		if move == "WAIT" {
			newMove.depth = -1
			newMove.spareMove = 1
		}
		if move == "TWICE" {
			newMove.spareMove = -1
		}
		possibleMoves = append(possibleMoves, newMove)
	}
	return possibleMoves
}

func Rotate(cmd string, room int) int {
	if cmd == "LEFT" {
		return RotateLeft[room]
	}
	if cmd == "RIGHT" {
		return RotateRight[room]
	}
	if cmd == "TWICE" {
		return RotateTwice[room]
	}
	return room
}

func MoveRocks(rs []Entity) []Entity {
	var newR []Entity
	for _, r := range rs {
		rock := r.Move()
		if rock.Pos == "X" {
			continue
		}
		newR = append(newR, rock)
	}
	return newR
}

func CheckCollisions(Indy Entity, Rocks []Entity) bool {
	for _, rock := range Rocks {
		if rock.Move().Pos == "X" {
			continue // A edge case where Indy next pos iquals Rock next pos.
		}
		if Indy.XY == rock.XY {
			return true
		}
	}
	return false
}

func Dps(Indy Entity, Rocks []Entity, state State) bool {
	if Indy.Y == EX && Indy.X == len(Grid)-1 {
		return true
	}
	newIndy := Indy.Move()
	newRocks := MoveRocks(Rocks)
	if CheckCollisions(Indy, Rocks) {
		return false
	}
	if newIndy.Pos == "X" {
		return false
	}
	cantRotate := getCantRotate(Indy, Rocks)
	moves := getMoves(newIndy, newRocks, state, cantRotate)
	for _, move := range moves {
		room := Grid[move.x][move.y]
		state.Depth += move.depth
		state.SpareMove += move.spareMove
		newRoom := Rotate(move.move, room)
		Grid[move.x][move.y] = newRoom
		if move.move != "WAIT" {
			Actions = append(Actions, Action{move.move, move.x, move.y})
		}
		if Dps(newIndy, newRocks, State{state.Depth + 1, state.SpareMove}) {
			if state.Depth != 0 {
				Grid[move.x][move.y] = room
			}
			if move.move == "TWICE" && state.Depth == 0 {
				Grid[move.x][move.y] = RotateLeft[room]
			}
			return true
		}
		state.Depth -= move.depth
		state.SpareMove -= move.spareMove
		if move.move != "WAIT" {
			Actions = Actions[:len(Actions)-1]
		}
		Grid[move.x][move.y] = room
	}
	return false
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1000000), 1000000)

	// W: number of columns.
	// H: number of rows.
	var W, H int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &W, &H)

	for i := 0; i < H; i++ {
		scanner.Scan()
		Line := scanner.Text()
		Grid = append(Grid, split_integers(Line))
	}
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &EX)
	for {
		var XI, YI int
		var POSI string
		scanner.Scan()
		fmt.Sscan(scanner.Text(), &XI, &YI, &POSI)
		Indy := Entity{POSI, YI, XI, [2]int{YI, XI}}

		// R: the number of rocks currently in the grid.
		var R int
		scanner.Scan()
		fmt.Sscan(scanner.Text(), &R)

		var Rocks []Entity

		for i := 0; i < R; i++ {
			var XR, YR int
			var POSR string
			scanner.Scan()
			fmt.Sscan(scanner.Text(), &XR, &YR, &POSR)
			Rocks = append(Rocks, Entity{POSR, YR, XR, [2]int{YR, XR}})
		}

		Actions = Actions[:0]
		Dps(Indy, Rocks, State{0, 0})
		if len(Actions) == 0 {
			fmt.Println("WAIT")
		} else {
			Actions[0].Anwser()
		}
	}
}
