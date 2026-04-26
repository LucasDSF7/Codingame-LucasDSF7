/*
Hashi puzzle tips for constructing TrivialSolutions function:
https://www.conceptispuzzles.com/index.aspx?uri=puzzle%2Fhashi%2Ftechniques

Copy and paste this code in https://www.codingame.com/ide/puzzle/there-is-no-spoon-episode-2 IDE to solve the puzzle.
*/

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

var Nodes []*Node
var Commands []string
var Visited = make(map[string]struct{}, 10000)

type Move struct {
	From   *Node
	To     *Node
	Amount int
}

type Node struct {
	Y          int
	X          int
	InitAmount int
	Amount     int
	Neighbors  []*Node
	Bridges    map[*Node]int
}

func (n *Node) String() string {
	str := strings.Builder{}
	fmt.Fprintf(&str, "Y: %d X: %d Links: %d Nº Neighbors: %v Non-zero Neighbors: %d\n", n.X, n.Y, n.Amount, len(n.Neighbors), n.NumberNeighbors())
	return str.String()
}

func FindNeighbors(nodesMap map[int][]*Node) {
	for _, v := range nodesMap {
		for i := 0; i < len(v)-1; i++ {
			node1 := v[i]
			node2 := v[i+1]
			node1.Neighbors = append(node1.Neighbors, node2)
			node2.Neighbors = append(node2.Neighbors, node1)
		}
	}
}

func (n *Node) NumberNeighbors() int {
	counter := 0
	for _, node := range n.Neighbors {
		if node.Amount > 0 && node.Bridges[n] < 2 {
			counter += 1
		}
	}
	return counter
}

func (n *Node) CanCompleteBridges() bool {
	counter := 0
	for _, node := range n.Neighbors {
		if node.Amount > 0 && node.Bridges[n] < 2 {
			counter += 2 - node.Bridges[n]
		}
	}
	return n.Amount == counter
}

func Link(move Move) {
	n := move.From
	node := move.To
	Commands = append(Commands, fmt.Sprintf("%d %d %d %d %d", n.X, n.Y, node.X, node.Y, move.Amount))
	node.Amount -= move.Amount
	n.Amount -= move.Amount
	node.Bridges[n] += move.Amount
	n.Bridges[node] += move.Amount
	if n.Y == node.Y {
		BarrierY(move, 1)
	} else {
		BarrierX(move, 1)
	}
}

func UnLink(moves []Move) {
	for _, move := range moves {
		node := move.From
		n := move.To
		node.Amount += move.Amount
		n.Amount += move.Amount
		node.Bridges[n] -= move.Amount
		n.Bridges[node] -= move.Amount
		if node.Bridges[n] == 0 {
			delete(node.Bridges, n)
		}
		if n.Bridges[node] == 0 {
			delete(n.Bridges, node)
		}
		if n.Y == node.Y {
			BarrierY(move, -1)
		} else {
			BarrierX(move, -1)
		}
	}
	Commands = Commands[:len(Commands)-len(moves)]
}

// BarrierY ocour when a bidge is placed along Y axes. If do == 1 it places the barrier. If -1 it removes.
func BarrierY(move Move, do int) {
	for _, node := range Nodes {
		if (node.X < move.From.X && node.X > move.To.X) || (node.X > move.From.X && node.X < move.To.X) {
			for _, neighbor := range node.Neighbors {
				if (node.Y < move.From.Y && neighbor.Y > move.From.Y) || (node.Y > move.From.Y && neighbor.Y < move.From.Y) {
					node.Bridges[neighbor] += 2 * do
					neighbor.Bridges[node] += 2 * do
				}
			}
		}
	}
}

// BarrierX ocour when a bidge is placed along Y axes. If do == 1 it places the barrier. If -1 it removes.
func BarrierX(move Move, do int) {
	for _, node := range Nodes {
		if (node.Y < move.From.Y && node.Y > move.To.Y) || (node.Y > move.From.Y && node.Y < move.To.Y) {
			for _, neighbor := range node.Neighbors {
				if (node.X < move.From.X && neighbor.X > move.From.X) || (node.X > move.From.X && neighbor.X < move.From.X) {
					node.Bridges[neighbor] += 2 * do
					neighbor.Bridges[node] += 2 * do
				}
			}
		}
	}
}

func TrivialSolutions(nodes []*Node) []Move {
	var moves []Move
	for {
		links := false
		for _, node := range nodes {
			if node.Amount == 0 || node.NumberNeighbors() == 0 {
				continue
			}
			if node.Amount/node.NumberNeighbors() == 2 {
				for _, neighbor := range node.Neighbors {
					if neighbor.Amount < 2 || node.Bridges[neighbor] > 0 || neighbor.Amount == 0 {
						continue
					}
					move := Move{node, neighbor, 2}
					Link(move)
					moves = append(moves, move)
					links = true
				}
				continue
			}
			if (node.Amount == 1 || node.Amount == 2) && node.NumberNeighbors() == 1 {
				for _, neighbor := range node.Neighbors {
					if neighbor.Amount < node.Amount || node.Bridges[neighbor]+node.Amount > 2 || neighbor.Amount == 0 || node.Amount == 0 {
						continue
					}
					move := Move{node, neighbor, node.Amount}
					Link(move)
					moves = append(moves, move)
					links = true
				}
				continue
			}
			if node.Amount == 3 && node.NumberNeighbors() == 2 {
				for _, neighbor := range node.Neighbors {
					if neighbor.Amount == 0 || node.Bridges[neighbor] > 1 {
						continue
					}
					move := Move{node, neighbor, 1}
					moves = append(moves, move)
					Link(move)
					links = true
				}
				continue
			}
			if node.Amount == 5 && node.NumberNeighbors() == 3 {
				for _, neighbor := range node.Neighbors {
					if neighbor.Amount == 0 || node.Bridges[neighbor] > 1 {
						continue
					}
					move := Move{node, neighbor, 1}
					moves = append(moves, move)
					Link(move)
					links = true
				}
				continue
			}
			if node.Amount == 7 && node.NumberNeighbors() == 4 {
				for _, neighbor := range node.Neighbors {
					if neighbor.Amount == 0 || node.Bridges[neighbor] > 1 {
						continue
					}
					move := Move{node, neighbor, 1}
					moves = append(moves, move)
					Link(move)
					links = true
				}
				continue
			}
			if node.CanCompleteBridges() {
				for _, neighbor := range node.Neighbors {
					if neighbor.Amount == 0 || node.Bridges[neighbor] >= 2 {
						continue
					}
					move := Move{node, neighbor, 1}
					moves = append(moves, move)
					Link(move)
					links = true
				}
				continue
			}
		}
		if !links {
			break
		}
	}
	return moves
}

func GetState(nodes []*Node) string {
	str := strings.Builder{}
	for _, node := range nodes {
		fmt.Fprintf(&str, "%d%d%d", node.X, node.Y, node.Amount)
		for _, neighbor := range node.Neighbors {
			fmt.Fprintf(&str, "%d%d%d", neighbor.X, neighbor.Y, node.Bridges[neighbor])
		}
	}
	return str.String()
}

func CheckGroup(node *Node, memo map[*Node]struct{}, target int) bool {
	if node.Amount > 0 || len(memo) == target { // it can still branch out
		return true
	}
	for neighbor, n := range node.Bridges {
		if n > 2 {
			continue
		}
		if _, visited := memo[neighbor]; !visited {
			memo[neighbor] = struct{}{}
			if CheckGroup(neighbor, memo, target) {
				return true
			}
		}
	}
	return false
}

func CheckSolution(nodes []*Node) bool {
	solveds := 0
	for _, node := range nodes {
		if node.Amount == 0 {
			solveds += 1
		}
	}
	return solveds == len(nodes)
}

func DFS(nodes []*Node) bool {
	if CheckSolution(nodes) {
		return CheckGroup(nodes[0], make(map[*Node]struct{}), len(nodes))
	}
	stateKey := GetState(nodes)
	if _, seen := Visited[stateKey]; seen {
		return false
	}
	for _, node := range nodes {
		for _, neighbor := range node.Neighbors {
			if node.Amount < 1 || neighbor.Amount < 1 || node.Amount > 1 {
				continue
			}
			if node.InitAmount == 2 && neighbor.InitAmount == 2 && node.Bridges[neighbor] == 1 {
				continue
			}
			if node.Bridges[neighbor]+1 > 2 {
				continue
			}
			move := Move{node, neighbor, 1}
			Link(move)
			moves := append(TrivialSolutions(nodes), move)
			if node.Amount < 0 || neighbor.Amount < 0 {
				UnLink(moves)
				continue
			}
			if node.Amount == 0 && !CheckGroup(node, make(map[*Node]struct{}), len(nodes)) {
				UnLink(moves)
				continue
			}
			if DFS(nodes) {
				return true
			}
			UnLink(moves)
		}
	}
	Visited[stateKey] = struct{}{}
	return false
}

func main() {
	Ynodes := make(map[int][]*Node)
	Xnodes := make(map[int][]*Node)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1000000), 1000000)
	// width: the number of cells on the X axis
	var width int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &width)
	fmt.Fprintln(os.Stderr, width)

	// height: the number of cells on the Y axis
	var height int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &height)
	fmt.Fprintln(os.Stderr, height)

	for y := 0; y < height; y++ {
		scanner.Scan()
		line := scanner.Text()
		fmt.Fprintln(os.Stderr, line)
		for x, node := range line {
			if node == '.' {
				continue
			}
			amount, _ := strconv.Atoi(string(node))
			newNode := &Node{y, x, amount, amount, []*Node{}, map[*Node]int{}}
			Nodes = append(Nodes, newNode)
			Ynodes[y] = append(Ynodes[y], newNode)
			Xnodes[x] = append(Xnodes[x], newNode)
			// Ynodes and Xnodes are sorted by default.
		}
	}
	start := time.Now()
	FindNeighbors(Xnodes)
	FindNeighbors(Ynodes)
	TrivialSolutions(Nodes)
	fmt.Fprintln(os.Stderr, time.Since(start))
	sort.Slice(Nodes, func(i, j int) bool {
		return Nodes[i].Amount < Nodes[j].Amount
	})
	fmt.Fprintln(os.Stderr, Nodes)
	DFS(Nodes)
	fmt.Fprintln(os.Stderr, time.Since(start))
	fmt.Fprintln(os.Stderr, Nodes)
	for _, cmd := range Commands {
		fmt.Println(cmd)
	}
	// fmt.Fprintln(os.Stderr, "Debug messages...")
	fmt.Fprintln(os.Stderr, "{", strings.Join(Commands, "\",\""), "}")
}
