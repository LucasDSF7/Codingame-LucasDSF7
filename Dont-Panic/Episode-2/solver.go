package main

import (
	"fmt"
	"os"
	"strings"
)

/*
This solution creates a simulation of the Improbability Drive,
 thus the codingame loop is not necessary
*/

var nbFloors, width, nbRounds, exitFloor, exitPos, nbTotalClones, nbAdditionalElevators, nbElevators, cloneFloor, generatorPos int

type Path struct {
	Floor                 int
	Pos                   int
	Direction             int
	Cost                  int
	NbAdditionalElevators int
	Moves                 map[int]int
}

func (p *Path) GetMoves(d *Drive) []int {
	moves := []int{p.Pos}
	leftPos, rightPos := 0, width-1
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
	return moves
}

func (p Path) NewPath(move int, d *Drive) Path {
	p.Moves[p.Floor] = move
	cost := (move - p.Pos) * p.Direction
	var addCost int
	var subElevator int
	if d.Grid[p.Floor][move] != "^" {
		addCost = 3
		subElevator = -1
	}
	if cost > 0 {
		return Path{p.Floor + 1, move, p.Direction, p.Cost + cost + addCost, p.NbAdditionalElevators + subElevator, p.Moves}
	}
	return Path{p.Floor + 1, move, p.Direction * -1, p.Cost + 3 + addCost - cost, p.NbAdditionalElevators + subElevator, p.Moves}
}

type Drive struct {
	Grid      [][]string
	Rounds    int
	Commands  []string
	Elevators map[int][]int
}

func (d Drive) String() string {
	var str strings.Builder
	for i := len(d.Grid) - 1; i >= 0; i-- {
		str.WriteString(strings.Join(d.Grid[i], "") + "\n")
	}
	fmt.Fprintf(&str, "Rounds %d\nCommands:%v\n", d.Rounds, d.Commands)
	return str.String()
}

func NewDrive(elevators map[int][]int) Drive {
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
	return Drive{grid, 0, []string{}, elevators}
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
		fmt.Printf("%+v\n", path)
		return
	}
	moves := path.GetMoves(d)
	// fmt.Printf("%+v\nMoves: %v\n\n", path, moves)
	for _, move := range moves {
		if d.Grid[path.Floor][move] != "^" && path.NbAdditionalElevators == 0 {
			continue
		}
		d.Dfs(path.NewPath(move, d))
	}
}

type Clone struct {
	floor     int
	pos       int
	direction int  // Init as 1 -> goes to the RIGHT
	active    bool // When a Clone show up it cant move. Init as false
}

func (c Clone) Move(Grid [][]string) Clone {
	tile := Grid[c.floor][c.pos]
	if tile == "B" {
		c.direction *= -1
		c.pos += c.direction
		return c
	}
	if tile == "^" {
		c.floor += 1
		return c
	}
	c.pos += c.direction
	if c.pos >= 0 && c.pos <= len(Grid[0])-1 && Grid[c.floor][c.pos] == "B" {
		c.direction *= -1
		c.pos += c.direction * 2
	}
	return c
}

func NewClone(d *Drive) Clone {
	return Clone{0, generatorPos, 1, true}
}

func PrintClones(clones []Clone) {
	var str strings.Builder
	for _, clone := range clones {
		fmt.Fprintf(&str, "Clone at %+v\n", clone)
	}
	fmt.Fprintln(os.Stderr, str.String())
}

func Cmd(grid [][]string, clones []Clone, cmd string) []Clone {
	if len(clones) == 0 {
		return clones
	}
	newClones := make([]Clone, len(clones))
	leader := clones[0]
	if cmd == "BLOCK" && leader.active {
		grid[leader.floor][leader.pos] = "B"
	}
	for i, clone := range clones {
		newClones[i] = clone.Move(grid)
	}
	if cmd == "BLOCK" {
		newClones = newClones[1:]
	}
	return newClones
}

func getLeader(clones []Clone) Clone {
	if len(clones) > 0 {
		return clones[0]
	}
	return Clone{0, 0, 0, false}
}

func pruning(clone Clone, drive *Drive) bool {
	var floor []string
	if clone.direction == 1 {
		floor = drive.Grid[clone.floor][clone.pos:]
	} else {
		floor = drive.Grid[clone.floor][:clone.pos+1]
	}
	for _, v := range floor {
		if v == "^" || v == "E" || v == "B" {
			return false
		}
	}
	return true
}

func Dfs(drive *Drive, clones []Clone) bool {
	//PrintClones(clones)
	//fmt.Println(drive)
	leader := getLeader(clones)
	if leader.pos == exitPos && leader.floor == exitFloor && leader.direction != 0 {
		return true
	}
	if drive.Rounds%3 == 0 && len(clones) <= nbTotalClones {
		clones = append(clones, NewClone(drive))
	}
	tile := drive.Grid[leader.floor][leader.pos]
	drive.Rounds += 1
	if drive.Rounds > nbRounds {
		return false
	}
	for _, cmd := range [2]string{"WAIT", "BLOCK"} {
		if cmd == "WAIT" && leader.active && pruning(leader, drive) {
			continue
		}
		newClones := Cmd(drive.Grid, clones, cmd)
		newLeader := getLeader(newClones)
		if newLeader.pos < 0 || newLeader.pos > len(drive.Grid[0])-1 {
			continue
		}
		drive.Commands = append(drive.Commands, cmd)
		if Dfs(drive, newClones) {
			return true
		}
		drive.Commands = drive.Commands[:len(drive.Commands)-1]
	}
	drive.Rounds -= 1
	drive.Grid[leader.floor][leader.pos] = tile
	return false
}

func main() {
	fmt.Scan(&nbFloors, &width, &nbRounds, &exitFloor, &exitPos, &nbTotalClones, &nbAdditionalElevators, &nbElevators)
	// fmt.Fprintln(os.Stderr, nbFloors, width, nbRounds, exitFloor, exitPos, nbTotalClones, nbAdditionalElevators, nbElevators)
	elevators := make(map[int][]int)

	for i := 0; i < nbElevators; i++ {
		var elevatorFloor, elevatorPos int
		fmt.Scan(&elevatorFloor, &elevatorPos)
		// fmt.Fprintln(os.Stderr, elevatorFloor, elevatorPos)
		elevators[elevatorFloor] = append(elevators[elevatorFloor], elevatorPos)
	}

	var direction string
	fmt.Scan(&cloneFloor, &generatorPos, &direction) // Only geratorPos is used
	// fmt.Fprintln(os.Stderr, cloneFloor, generatorPos, direction)

	// fmt.Fprintln(os.Stderr, "Debug messages...")
	drive := NewDrive(elevators)
	fmt.Fprintln(os.Stderr, drive)
	drive.PruneDrive()
	fmt.Fprintln(os.Stderr, drive)
	path := Path{0, generatorPos, 1, 0, nbAdditionalElevators, make(map[int]int)}
	drive.Dfs(path)
	// Dfs(&drive, make([]Clone, 0, nbTotalClones))
	// for _, cmd := range drive.Commands {
	// 	fmt.Println(cmd)
	// }
}
