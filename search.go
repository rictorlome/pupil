package main

import (
	"fmt"
	// "sort"
)

type Node struct {
	parent   *Node
	origin   Move
	children []Node
}

type Nodes []Node

func (n Nodes) Len() int {
	return len(n)
}
func (n Nodes) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}
func (n Nodes) Less(i, j int) bool {
	return int(move_src(n[i].origin)) < int(move_src(n[j].origin))
}

func (n Node) String() string {
	if n.parent == nil {
		return fmt.Sprintf("%v", n.origin)
	}
	return fmt.Sprintf("%v -> %v", n.parent, n.origin)
}

func build_tree(p *Position, parent *Node, move Move, depth_left int) Node {
	if 0 < depth_left {
		moves := (*p).generate_moves()
		children := make([]Node, len(moves))
		self := Node{parent, move, children}
		for i := 0; i < len(moves); i++ {
			(*p).do_move(moves[i], StateInfo{})
			children[i] = build_tree(p, &self, moves[i], depth_left-1)
			// sort.Sort(Nodes(children))
			(*p).undo_move(moves[i])
		}
		return self
	}
	return Node{parent, move, make([]Node, 0)}
}

func (n *Node) count_leaves() int {
	if len(n.children) == 0 {
		return 1
	}
	child_leaves := 0
	for _, child := range n.children {
		child_leaves += child.count_leaves()
	}
	return child_leaves
}

func divide(fen string, max_depth int, dividor int) {
	pos := parse_fen(fen)
	for i := 1; i <= max_depth; i++ {
		n := build_tree(&pos, nil, Move(0), i)
		if 2 < i {
			divide_tree("", n, 1, dividor, i)
		}
		fmt.Println(fmt.Sprintf("perft( %v)=          %--v", i, n.count_leaves()))
	}
}

func divide_tree(prefix string, n Node, cur_depth int, max_depth int, i int) {
	if max_depth < cur_depth {
		return
	}
	for i, child := range n.children {
		new_prefix := fmt.Sprintf("%v %v", prefix, child.origin)
		divide_tree(new_prefix, child, cur_depth+1, max_depth, i)
	}
	if 1 < cur_depth && cur_depth < i {
		fmt.Println(fmt.Sprintf("%v. %v moves =        %v", cur_depth, prefix, n.count_leaves()))
	}
}
