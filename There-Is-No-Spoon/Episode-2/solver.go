/*
https://www.conceptispuzzles.com/index.aspx?uri=puzzle%2Fhashi%2Ftechniques
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

var Commands []string
var Visited = make(map[string]struct{}, 10000)

type Move struct {
	From   *Node
	To     *Node
	Amount int
}

type Node struct {
	Y               int
	X               int
	InitAmount      int
	Amount          int
	Neighbors       []*Node
	NumberNeighbors int
	Bridges         map[*Node]int
}

func (n *Node) String() string {
	str := strings.Builder{}
	fmt.Fprintf(&str, "Y: %d X: %d Links: %d Nº Neighbors: %v Non-zero Neighbors: %d\n", n.X, n.Y, n.Amount, len(n.Neighbors), n.NumberNeighbors)
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

func Link(move Move) {
	n := move.From
	node := move.To
	Commands = append(Commands, fmt.Sprintf("%d %d %d %d %d", n.X, n.Y, node.X, node.Y, move.Amount))
	node.Amount -= move.Amount
	n.Amount -= move.Amount
	node.Bridges[n] += move.Amount
	n.Bridges[node] += move.Amount
	if node.Amount == 0 {
		for _, neighbor := range node.Neighbors {
			neighbor.NumberNeighbors -= 1
		}
	}
	if n.Amount == 0 {
		for _, neighbor := range n.Neighbors {
			neighbor.NumberNeighbors -= 1
		}
	}
}

func UnLink(moves []Move) {
	for _, move := range moves {
		node := move.From
		n := move.To
		if node.Amount == 0 {
			for _, neighbor := range node.Neighbors {
				neighbor.NumberNeighbors += 1
			}
		}
		if n.Amount == 0 {
			for _, neighbor := range n.Neighbors {
				neighbor.NumberNeighbors += 1
			}
		}
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
	}
	Commands = Commands[:len(Commands)-len(moves)]
}

func TrivialSolutions(nodes []*Node) []Move {
	var moves []Move
	for {
		links := false
		for _, node := range nodes {
			if node.Amount == 0 || node.NumberNeighbors == 0 {
				continue
			}
			if node.Amount/node.NumberNeighbors == 2 {
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
			if (node.Amount == 1 || node.Amount == 2) && node.NumberNeighbors == 1 {
				for _, neighbor := range node.Neighbors {
					if neighbor.Amount < node.Amount || node.Bridges[neighbor]+node.Amount > 2 || neighbor.Amount == 0 {
						continue
					}
					move := Move{node, neighbor, node.Amount}
					Link(move)
					moves = append(moves, move)
					links = true
				}
				continue
			}
			if node.Amount == 3 && node.NumberNeighbors == 2 {
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
			if node.Amount == 5 && node.NumberNeighbors == 3 {
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
			if node.Amount == 7 && node.NumberNeighbors == 4 {
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
	for neighbor := range node.Bridges {
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

func DFS(nodes []*Node, turn int) bool {
	if CheckSolution(nodes) {
		return CheckGroup(nodes[0], make(map[*Node]struct{}), len(nodes))
	}
	stateKey := GetState(nodes)
	if _, seen := Visited[stateKey]; seen {
		return false
	}
	for _, node := range nodes {
		for _, neighbor := range node.Neighbors {
			for _, amount := range [1]int{1} {
				if node.Amount < amount || neighbor.Amount < amount {
					continue
				}
				if node.InitAmount == amount && neighbor.InitAmount == amount && len(node.Neighbors) > 1 {
					continue
				}
				if node.Bridges[neighbor]+amount > 2 {
					continue
				}
				move := Move{node, neighbor, amount}
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
				if DFS(nodes, turn+1) {
					return true
				}
				UnLink(moves)
			}
		}
	}
	Visited[stateKey] = struct{}{}
	return false
}

func main() {
	var Nodes []*Node
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
			newNode := &Node{y, x, amount, amount, []*Node{}, 0, map[*Node]int{}}
			Nodes = append(Nodes, newNode)
			Ynodes[y] = append(Ynodes[y], newNode)
			Xnodes[x] = append(Xnodes[x], newNode)
			// Ynodes and Xnodes are sorted by default.
		}
	}
	start := time.Now()
	FindNeighbors(Xnodes)
	FindNeighbors(Ynodes)
	for _, node := range Nodes {
		node.NumberNeighbors = len(node.Neighbors)
	}
	//TrivialSolutions(Nodes)
	fmt.Fprintln(os.Stderr, time.Since(start))
	sort.Slice(Nodes, func(i, j int) bool {
		return Nodes[i].Amount < Nodes[j].Amount
	})
	fmt.Fprintln(os.Stderr, Nodes)
	DFS(Nodes, 1)
	fmt.Fprintln(os.Stderr, time.Since(start))
	fmt.Fprintln(os.Stderr, Nodes)
	for _, cmd := range Commands {
		fmt.Println(cmd)
	}
	// fmt.Fprintln(os.Stderr, "Debug messages...")
	fmt.Fprintln(os.Stderr, "{", strings.Join(Commands, "\",\""), "}")
}
