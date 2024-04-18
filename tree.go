package main

import "sort"

type TreeNode interface {
	isLeaf() bool
	getWeight() uint32
	getVal() rune // hack - used for breaking ties when sorting
}

type LeafNode struct {
	char   rune
	weight uint32
}

func (l LeafNode) isLeaf() bool {
	return true
}

func (l LeafNode) getWeight() uint32 {
	return l.weight
}

func (l LeafNode) getVal() rune {
	return l.char
}

type InternalNode struct {
	left   TreeNode
	right  TreeNode
	weight uint32
}

func (i InternalNode) isLeaf() bool {
	return false
}

func (i InternalNode) getWeight() uint32 {
	return i.weight
}

func (i InternalNode) getVal() rune {
	return ' '
}

func buildTree(frequencies map[rune]uint32) TreeNode {
	var trees []TreeNode
	for char, frequency := range frequencies {
		trees = append(trees, LeafNode{char: char, weight: frequency})
	}

	for len(trees) > 1 {
		sort.Slice(trees, func(i, j int) bool {
			if trees[i].getWeight() == trees[j].getWeight() {
				return trees[i].getVal() < trees[j].getVal()
			}
			return trees[i].getWeight() < trees[j].getWeight()
		})

		left := trees[0]
		right := trees[1]

		newNode := InternalNode{
			left:   left,
			right:  right,
			weight: left.getWeight() + right.getWeight(),
		}

		if len(trees) > 2 {
			trees = trees[2:]
			trees = append(trees, newNode)
		} else {
			trees = []TreeNode{newNode}
		}

	}
	return trees[0]
}
