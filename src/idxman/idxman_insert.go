package idxman

import (
	"../common"
	"math"
)

// Insert v into B+ Tree
func (self *idxMan) Insert(v common.CellValue, id int64) {
	common.OpLogger.Print("[Insert]", v)
	defer common.OpLogger.Print("[leave Insert]")
	
	l := self.root.findLeafNode(v)
	l.insertKey(v, id)
	// If the l is full, split it.
	// Then insert two nodes l and l1 into their parent.
	// Update root if new root is created.
	if l.isFull() {
		common.OpLogger.Print("[Split!]")
		v1, l1 := l.splitNode()
		r, newRoot := l.insertInParent(v1, l1)
		if newRoot {
			self.root = r
		}
	}
}

// create an empty node and return its address
func createNode(isLeaf bool) *node {
	common.OpLogger.Print("[createNode]")
	defer common.OpLogger.Print("[leave createNode]")
	
	x := new(node)
	x.keys = make([]common.CellValue, 0, order)
	x.leaf = isLeaf
	if isLeaf {
		x.children = make([]*node, 1, 1)
		x.recordIds = make([]int64, 0, order)
	} else {
		x.children = make([]*node, 0, order+1)
	}
	return x
}

// PRIVATE FUNCTION

// Split node and return the new node
func (self *node) splitNode() (common.CellValue, *node) {
	common.OpLogger.Print("[splitNode()]", self)
	defer common.OpLogger.Print("[leave splitNode]")
	
	// Create node
	// The new node is the right brother of self
	n := createNode(self.isLeaf())
	n.parent = self.parent
	var remainCnt int
	if n.isLeaf() {
		remainCnt = int(math.Ceil(float64(order-1) / 2.0))
	} else {
		remainCnt = int(math.Ceil(float64(order)/2.0) - 1)
	}
	var k common.CellValue
	if self.isLeaf() {
		n.keys = append(n.keys, self.keys[remainCnt:]...)
		n.recordIds = append(n.recordIds, self.recordIds[remainCnt:]...)
		n.children = append(n.children, self.children[0])
		self.children[0] = n

		k = self.keys[remainCnt]
		self.keys = self.keys[:remainCnt]
		self.recordIds = self.recordIds[:remainCnt]
	} else {
		n.keys = append(n.keys, self.keys[remainCnt+1:]...)
		n.children = append(n.children, self.children[remainCnt+1:]...)
		// Update child's parent
		for _, c := range n.children {
			c.parent = n
		}
		k = self.keys[remainCnt]
		self.keys = self.keys[:remainCnt]
		self.children = self.children[:remainCnt+1]
	}

	return k, n
}


// Insert k into leaf
func (self *node) insertKey(k common.CellValue, id int64) bool {
	common.OpLogger.Print("[insertKey]", self, ",", k)
	defer common.OpLogger.Print("[leave insertKey]", self)
	
	// l should be a leaf and not full
	if (!self.isLeaf()) || self.isFull() {
		common.ErrLogger.Print("Should be a leaf and not full!\t", self)
		return false
	}

	i, _ := self.findKeyIndex(k)
	self.keys = append(self.keys[:i], append([]common.CellValue{k}, self.keys[i:]...)...)
	self.recordIds = append(self.recordIds[:i], append([]int64{id}, self.recordIds[i:]...)...)

	return true
}

// Insert c into non-leaf
func (self *node) insertChild(k common.CellValue, c *node) bool {
	common.OpLogger.Print("[insertChild]")
	defer common.OpLogger.Print("[leave insertChild]")
	
	// l should be a non-leaf and not full
	if self.isLeaf() || self.isFull() {
		common.ErrLogger.Print("Should be a non-leaf and not full ", self)
		return false
	}
	// update parent
	c.parent = self
	// case 1: self is empty
	if self.childCnt() == 0 {
		common.OpLogger.Print("[First Child!]")
		self.children = append(self.children, c)
	} else {
		i := self.findChildIndex(k)
		self.keys = append(self.keys[:i], append([]common.CellValue{k}, self.keys[i:]...)...)
		self.children = append(self.children[:i+1], append([]*node{c}, self.children[i+1:]...)...)
	}
	return true
}

// Insert c into self's parent
// If new root is created, it will return (root, newRoot).
// Otherwise return (nil, false)
// @_@ It is not functional programming, keep it stupid.
func (self *node) insertInParent(k common.CellValue, c *node) (*node, bool) {
	common.OpLogger.Print("[insertInParent]")
	defer common.OpLogger.Print("[leave insertInParent]")
	
	// If l is the root of the tree, split it and create new root.
	if self.isRoot() {
		temp := createNode(false)
		temp.insertChild(self.minKey(), self)
		temp.insertChild(k, c)
		common.OpLogger.Print("[new root]")
		return temp, true
	}

	p := self.parent
	p.insertChild(k, c)
	if p.isFull() {
		k1, p1 := p.splitNode()
		return p.insertInParent(k1, p1)
	}
	return nil, false
}
