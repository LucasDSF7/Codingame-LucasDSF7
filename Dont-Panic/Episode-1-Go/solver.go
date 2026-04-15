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

type Drive struct {
	Grid         [][]string
	Rounds       int
	Commands     []string
	GeneratorPos int
	ExitFloor    int
	ExitPos      int
}

func (d Drive) String() string {
	var str strings.Builder
	for i := len(d.Grid) - 1; i >= 0; i-- {
		str.WriteString(strings.Join(d.Grid[i], "") + "\n")
	}
	fmt.Fprintf(&str, "Rounds %d\nCommands:%v\n", d.Rounds, d.Commands)
	return str.String()
}

func NewDrive(elevators map[int]int, minWidth int, maxWidth int) Drive {
	grid := make([][]string, nbFloors)
	for i := 0; i < nbFloors; i++ {
		floor := make([]string, maxWidth)
		for j := 0; j < maxWidth; j++ {
			floor[j] = "_"
		}
		if pos, ok := elevators[i]; ok {
			floor[pos] = "^"
		}
		if i == exitFloor {
			floor[exitPos] = "E"
		}
		if i == 0 {
			floor[generatorPos] = "G"
		}
		grid[i] = floor[minWidth:]
	}
	return Drive{grid, 0, []string{}, generatorPos - minWidth, exitFloor, exitPos - minWidth}
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
	return Clone{0, d.GeneratorPos, 1, true}
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
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
	if leader.pos == drive.ExitPos && leader.floor == drive.ExitFloor && leader.direction != 0 {
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
	fmt.Fprintln(os.Stderr, nbFloors, width, nbRounds, exitFloor, exitPos, nbTotalClones, nbAdditionalElevators, nbElevators)
	minWidth, maxWidth := exitPos, 0
	elevators := make(map[int]int)

	for i := 0; i < nbElevators; i++ {
		var elevatorFloor, elevatorPos int
		fmt.Scan(&elevatorFloor, &elevatorPos)
		fmt.Fprintln(os.Stderr, elevatorFloor, elevatorPos)
		elevators[elevatorFloor] = elevatorPos
		minWidth, maxWidth = min(minWidth, elevatorPos), max(maxWidth, elevatorPos+1)
	}

	var direction string
	fmt.Scan(&cloneFloor, &generatorPos, &direction) // Only geratorPos is used
	fmt.Fprintln(os.Stderr, cloneFloor, generatorPos, direction)
	minWidth, maxWidth = min(minWidth, generatorPos), max(maxWidth, generatorPos+1)

	// Implement clone generator to replace codingame loop
	// fmt.Fprintln(os.Stderr, "Debug messages...")
	drive := NewDrive(elevators, minWidth, maxWidth)
	Dfs(&drive, make([]Clone, 0, nbTotalClones))
	for _, cmd := range drive.Commands {
		fmt.Println(cmd)
	}
}
