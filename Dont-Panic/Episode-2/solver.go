package main

/*
Copy and paste this code at https://www.codingame.com/ide/puzzle/don't-panic-episode-2
*/

import (
	"fmt"
	"os"
	"strings"
)

var nbFloors, width, nbRounds, exitFloor, exitPos, nbTotalClones, nbAdditionalElevators, nbElevators, cloneFloor, generatorPos int

type Path struct {
	Floor                 int
	Pos                   int
	Direction             int
	Cost                  int
	NbAdditionalElevators int
	Moves                 map[int]int
}

func (p Path) GetMoves(d *Drive) []int {
	if d.Grid[p.Floor][p.Pos] == "^" {
		return []int{p.Pos}
	}
	moves := []int{p.Pos}
	leftPos, rightPos := 0, width-1
	upElevators := d.Elevators[p.Floor+1]
	elevators := d.Elevators[p.Floor]
	for _, elevator := range elevators {
		if elevator > leftPos && elevator < p.Pos {
			leftPos = elevator
		}
		if elevator < rightPos && elevator > p.Pos {
			rightPos = elevator
		}
	}
	if leftPos != 0 {
		moves = append(moves, leftPos)
	}
	if rightPos != width-1 {
		moves = append(moves, rightPos)
	}
	if exitPos > leftPos && exitPos < rightPos {
		moves = append(moves, exitPos)
	}
	for _, up := range upElevators {
		if up > leftPos && up < rightPos {
			moves = append(moves, up)
		}
	}
	return moves
}

func (p Path) NewPath(move int, d *Drive) Path {
	newMoves := make(map[int]int)
	for k, v := range p.Moves {
		newMoves[k] = v
	}
	newMoves[p.Floor] = move
	cost := (move - p.Pos) * p.Direction
	var addCost int
	var subElevator int
	if d.Grid[p.Floor][move] != "^" {
		addCost = 3
		subElevator = -1
	}
	if cost >= 0 {
		return Path{p.Floor + 1, move, p.Direction, p.Cost + cost + addCost, p.NbAdditionalElevators + subElevator, newMoves}
	}
	return Path{p.Floor + 1, move, p.Direction * -1, p.Cost + 3 + addCost - cost, p.NbAdditionalElevators + subElevator, newMoves}
}

type Drive struct {
	Grid      [][]string
	Rounds    int
	Elevators map[int][]int
	BestPath  Path
}

func (d Drive) String() string {
	var str strings.Builder
	for i := len(d.Grid) - 1; i >= 0; i-- {
		str.WriteString(strings.Join(d.Grid[i], "") + "\n")
	}
	fmt.Fprintf(&str, "Rounds %d\n", d.Rounds)
	return str.String()
}

func NewDrive(elevators map[int][]int) *Drive {
	grid := make([][]string, nbFloors)
	for i := 0; i < nbFloors; i++ {
		floor := make([]string, width)
		for j := 0; j < width; j++ {
			floor[j] = "_"
		}
		if positions, ok := elevators[i]; ok {
			for _, pos := range positions {
				floor[pos] = "^"
			}
		}
		if i == exitFloor {
			floor[exitPos] = "E"
		}
		if i == 0 {
			floor[generatorPos] = "G"
		}
		grid[i] = floor
	}
	return &Drive{grid, 0, elevators, Path{Cost: 1000}}
}

func (d *Drive) PruneDrive() {
	minMax := [2]int{exitPos, exitPos}
	for i := nbFloors - 1; i >= 0; i-- {
		if i > exitFloor {
			for j := range d.Grid[i] {
				d.Grid[i][j] = "X"
			}
			continue
		}
		newMin, newMax := 0, width-1
		if _, ok := d.Elevators[i]; !ok {
			continue
		}
		for _, pos := range d.Elevators[i] {
			if pos > newMin && pos < minMax[0] {
				newMin = pos
			}
			if pos < newMax && pos > minMax[1] {
				newMax = pos
			}
		}
		for j := 0; j <= newMin; j++ {
			if newMin == 0 {
				break
			}
			d.Grid[i][j] = "X"
		}
		for j := newMax; j < width; j++ {
			d.Grid[i][j] = "X"
		}
		minMax = [2]int{newMin, newMax}
	}
}

func (d *Drive) Dfs(path Path) {
	tile := d.Grid[path.Floor][path.Pos]
	if tile == "X" {
		return
	}
	if tile == "E" {
		if path.Cost < d.BestPath.Cost {
			d.BestPath = path
		}
		return
	}
	moves := path.GetMoves(d)
	for _, move := range moves {
		if d.Grid[path.Floor][move] != "^" && path.NbAdditionalElevators == 0 {
			continue
		}
		d.Dfs(path.NewPath(move, d))
	}
}

func (d *Drive) ApplyBestPath() {
	for floor := range d.Grid {
		move := d.BestPath.Moves[floor]
		for _, eleve := range d.Elevators[floor] {
			if move != eleve {
				d.Grid[floor][eleve] = "X"
			}
		}
		if d.Grid[floor][move] != "^" {
			d.Grid[floor][move] = "O"
		}
	}
}

type Clone struct {
	floor     int
	pos       int
	direction int // Init as 1 -> goes to the RIGHT
}

func PrintClones(clones []Clone) {
	var str strings.Builder
	for _, clone := range clones {
		fmt.Fprintf(&str, "Clone at %+v\n", clone)
	}
	fmt.Fprintln(os.Stderr, str.String())
}

func pruning(clone Clone, drive *Drive) bool {
	var floor []string
	if clone.direction == 1 {
		floor = drive.Grid[clone.floor][clone.pos:]
	} else {
		floor = drive.Grid[clone.floor][:clone.pos+1]
	}
	for _, v := range floor {
		if v == "^" || v == "E" || v == "B" || v == "O" {
			return false
		}
	}
	return true
}

func main() {
	fmt.Scan(&nbFloors, &width, &nbRounds, &exitFloor, &exitPos, &nbTotalClones, &nbAdditionalElevators, &nbElevators)
	elevators := make(map[int][]int)
	fmt.Fprintln(os.Stderr, nbFloors, width, nbRounds, exitFloor, exitPos, nbTotalClones, nbAdditionalElevators, nbElevators)

	for i := 0; i < nbElevators; i++ {
		var elevatorFloor, elevatorPos int
		fmt.Scan(&elevatorFloor, &elevatorPos)
		fmt.Fprintln(os.Stderr, elevatorFloor, elevatorPos)
		elevators[elevatorFloor] = append(elevators[elevatorFloor], elevatorPos)
	}
	firstTurn := true
	var drive *Drive
	for {
		var direction string
		fmt.Scan(&cloneFloor, &generatorPos, &direction)
		if firstTurn {
			fmt.Fprintln(os.Stderr, cloneFloor, generatorPos, direction)
			drive = NewDrive(elevators)
			fmt.Fprintln(os.Stderr, drive)
			drive.PruneDrive()
			fmt.Fprintln(os.Stderr, drive)
			path := Path{0, generatorPos, 1, 0, nbAdditionalElevators, make(map[int]int)}
			drive.Dfs(path)
			drive.ApplyBestPath()
			fmt.Fprintln(os.Stderr, drive)
			firstTurn = false
		}
		dir := 1
		if direction == "LEFT" {
			dir = -1
		}
		clone := Clone{cloneFloor, generatorPos, dir}
		if clone.floor == -1 {
			fmt.Println("WAIT")
			continue
		}
		if drive.Grid[clone.floor][clone.pos] == "O" {
			fmt.Println("ELEVATOR")
			drive.Grid[clone.floor][clone.pos] = "^"
			continue
		}
		if !pruning(clone, drive) {
			fmt.Println("WAIT")
		} else {
			fmt.Println("BLOCK")
		}
	}
}
