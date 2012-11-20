package idxman

import (
	"../common"
)

const order = 4 // We are implementing B+ Tree for n = 4.

// B+ Tree Node Definition
type node struct {
	parent *node
	// Size == order, reserve one.
	keys []common.CellValue
	// if it is a leaf, children[0] should point to its right sibling
	// Size == order + 1, reserve one.
	children  []*node
	recordIds []int64
	leaf      bool
}

func (self node) isFull() bool {
	return self.keyCnt() == order
}

func (self node) isRoot() bool {
	return self.parent == nil
}

func (self node) isLeaf() bool {
	return self.leaf
}

func (self node) keyCnt() int {
	return len(self.keys)
}

func (self node) childCnt() int {
	return len(self.children)
}

func (self node) minKey() common.CellValue {
	common.OpLogger.Print("[minKey]")
	defer common.OpLogger.Print("[leave minKey]")
	
	var n *node
	for n = &self; !n.isLeaf(); n = n.children[0] {
	}
	return n.keys[0]
}