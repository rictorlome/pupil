package main

import (
	"fmt"
)

type Node struct {
	origin   Move
	children []Node
}

func build_tree(p *Position, move Move, depth_left int) Node {
	var children []Node
	if 0 < depth_left {
		moves := (*p).generate_moves()
		for _, move := range moves {
			(*p).do_move(move, StateInfo{})
			children = append(children, build_tree(p, move, depth_left-1))
			(*p).undo_move(move)
		}
	}
	return Node{move, children}
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
		n := build_tree(&pos, Move(0), i)
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
