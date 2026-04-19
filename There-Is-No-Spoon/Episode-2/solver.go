package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

var Commands []string

type Node struct {
	Y         int
	X         int
	Amount    int
	Neighbors []*Node
}

func (n *Node) String() string {
	return fmt.Sprintf("Y: %d X: %d Links: %d Nº Neighbors: %v\n", n.X, n.Y, n.Amount, len(n.Neighbors))
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

func TrivialSolutions(nodes []*Node) {
	for _, node := range nodes {
		if len(node.Neighbors) != 1 {
			continue
		}
		neighbor := node.Neighbors[0]
		Commands = append(Commands, fmt.Sprintf("%d %d %d %d %d", node.X, node.Y, neighbor.X, neighbor.Y, node.Amount))
		neighbor.Amount -= node.Amount
		node.Amount = 0
	}
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
			newNode := &Node{y, x, amount, []*Node{}}
			Nodes = append(Nodes, newNode)
			Ynodes[y] = append(Ynodes[y], newNode)
			Xnodes[x] = append(Xnodes[x], newNode)
			// Ynodes and Xnodes are sorted by default.
		}
	}
	FindNeighbors(Xnodes)
	FindNeighbors(Ynodes)
	fmt.Println(Nodes)
	TrivialSolutions(Nodes)
	fmt.Println(Nodes)
	for _, cmd := range Commands {
		fmt.Println(cmd)
	}
	// fmt.Fprintln(os.Stderr, "Debug messages...")
}
