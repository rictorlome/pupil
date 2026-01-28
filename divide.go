package main

import (
	"fmt"
)

type Node struct {
	parent   *Node
	origin   Move
	children []Node
}

func buildTree(p *Position, parent *Node, move Move, depthLeft int, c chan Node) {
	c <- buildTreeRecursive(p, parent, move, depthLeft)
}

func buildTreeParallel(p *Position, depthLeft int) Node {
	moves := (*p).generateMoves()
	movesLen := len(moves)
	children := make([]Node, movesLen)
	c := make(chan Node)
	self := Node{nil, Move(0), make([]Node, movesLen)}
	for i := 0; i < movesLen; i++ {
		duped := (*p).dup()
		duped.doMove(moves[i], &StateInfo{})
		go buildTree(&duped, &self, moves[i], depthLeft-1, c)
	}
	for i := 0; i < movesLen; i++ {
		children[i] = <-c
	}
	self.children = children
	return self
}

func buildTreeRecursive(p *Position, parent *Node, move Move, depthLeft int) Node {
	self := Node{parent, move, make([]Node, 0)}
	if 0 < depthLeft {
		moves := (*p).generateMoves()
		children := make([]Node, len(moves))
		for i := 0; i < len(moves); i++ {
			(*p).doMove(moves[i], &StateInfo{})
			children[i] = buildTreeRecursive(p, &self, moves[i], depthLeft-1)
			(*p).undoMove(moves[i])
		}
		self.children = children
	}
	return self
}

func (n *Node) countLeaves() int {
	if len(n.children) == 0 {
		return 1
	}
	childLeaves := 0
	for _, child := range n.children {
		childLeaves += child.countLeaves()
	}
	return childLeaves
}

func divide(fen string, maxDepth int, dividor int) {
	pos := parseFen(fen)
	for i := 1; i <= maxDepth; i++ {
		n := buildTreeRecursive(&pos, nil, Move(0), i)
		if 2 < i {
			divideTree("", n, 1, dividor, i)
		}
		fmt.Println(fmt.Sprintf("perft( %v)=          %--v", i, n.countLeaves()))
	}
}

func divideTree(prefix string, n Node, curDepth int, maxDepth int, i int) {
	if maxDepth < curDepth {
		return
	}
	for i, child := range n.children {
		newPrefix := fmt.Sprintf("%v %v", prefix, child.origin)
		divideTree(newPrefix, child, curDepth+1, maxDepth, i)
	}
	if 1 < curDepth && curDepth < i {
		fmt.Println(fmt.Sprintf("%v. %v moves =        %v", curDepth, prefix, n.countLeaves()))
	}
}

func (n Node) String() string {
	if n.parent == nil {
		return fmt.Sprintf("%v", n.origin)
	}
	return fmt.Sprintf("%v -> %v", n.parent, n.origin)
}
