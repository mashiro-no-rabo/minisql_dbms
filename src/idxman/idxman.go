package idxman

import (
	"../common"
	"math"
)

// CONST
const order = 4 // We are implementing B+ Tree for n = 4.

type KeyType interface {
	GetValue() common.CellValue
}

// STRUCT
type node struct {
	parent *node
	// Size == order, reserve one.
	keys []KeyType
	// if it is a leaf, children[0] should point to its right sibling
	// Size == order + 1, reserve one.
	children []*node
	leaf     bool
}

type idxMan struct {
	root *node
}

// PUBLIC FUNCTION

// Creat an index manager and give it a empty root.
func NewIdxMan() *idxMan {
	common.OpLogger.Print("NewIdxMan(): New Index Manager!")
	im := new(idxMan)
	im.root = createNode()
	im.root.leaf = true
	im.root.children = append(im.root.children, nil)
	common.OpLogger.Print("The root is ", im.root)
	common.OpLogger.Print("leave NewIdxMan()")
	common.OpLogger.Println()
	return im
}

// Find the first KeyType containing given v,
// return (nil, false) if nothing is found.
func (self idxMan) Find(v common.CellValue) (KeyType, bool) {
	common.OpLogger.Print("Find():\t", v)
	l := self.root.findLeafNode(v)
	i, found := l.findKeyIndex(v)
	if found {
		common.OpLogger.Print("leave Find()\t", l,keys[i])
		return l.keys[i], true
	}
	common.OpLogger.Print("leave Find(), no record found.")
	return nil, false
}

// Insert k into B+ Tree
func (self *idxMan) Insert(k KeyType) {
	common.OpLogger.Print("Insert(): Insert ", k.GetValue())
	l := self.root.findLeafNode(k.GetValue())
	l.insertKey(k)
	// If the l is full, split it.
	// Then insert two nodes l and l1 into their parent.
	// Update root if new root is created.
	if l.isFull() {
		common.OpLogger.Print("Split!")
		k1, l1 := l.splitNode()
		r, newRoot := l.insertInParent(k1, l1)
		if newRoot {
			self.root = r
		}
	}
	common.OpLogger.Print("leave Insert()")
	common.OpLogger.Println()
}

// Delete the first KeyType containing given v,
// return (false) if nothing is deleted
func (self *idxMan) Delete(v common.CellValue) (KeyType, bool) {
	common.OpLogger.Print("Delete(): Delete ", v)
	// Find the corresponding leaf l
	l := self.root.findLeafNode(v)

	// Find the index of v in l
	i, found := l.findKeyIndex(v)
	if !found {
		common.OpLogger.Print("leave Delete(), no key deleted")
		return nil, false
	}
	k, _ := l.deleteKey(i)

	// Check node's children number back to root
	for n := l; ; n = n.parent {
		// A leaf root node has between 0 and order - 1 values
		// A leaf node has between ceil((order - 1) / 2) and order - 1 values.
		if n.isLeaf() {
                        if n.isRoot() || n.keyCnt() >= int(math.Ceil(float64(order-1)/2.0)) {
				break
			}
		} else {
		// A non-leaf root node has between 2 and order children
			if n.isRoot() && n.keyCnt() == 0 {
				self.root = n.children[0]
				break
                        }
		// A node that is not a leaf or root has between ceil(order / 2) and order children.
			if n.keyCnt() + 1 >= int(math.Ceil(float64(order)/2.0)) {
				break
			}
		}
		p := n.parent
		i := p.findChildIndex(n.minKey().GetValue())
		if i == p.keyCnt() {
			n1 := p.children[i-1]
			if n.keyCnt()+n1.keyCnt()+1 < order {
				// Case0: merge n with its left brother
				n1.mergeRightNode(n, p.keys[i-1])
				p.deleteChild(i)
			} else {
				// Case1: borrow a child from left brother
				if n.isLeaf() {
					k, _ := n1.deleteKey(n1.keyCnt() - 1)
					n.insertKey(k)
				} else {
					k, c, _ := n1.deleteChild(n1.keyCnt() - 1)
					n.insertChild(k, c)
				}
			}
		} else {
			n1 := p.children[i+1]
			if n.keyCnt()+n1.keyCnt()+1 < order {
				// Case2: merge n with its right brother
				n.mergeRightNode(n1, p.keys[i])
				p.deleteChild(i + 1)
			} else {
				// Case3: borrow a child from right brother
				if n.isLeaf() {
					k, _ := n1.deleteKey(0)
					n.insertKey(k)
				} else {
					k, c, _ := n1.deleteChild(0)
					n.insertChild(k, c)
				}
			}
		}
	}
	common.OpLogger.Print("leave Delete()\t", k)
	return k, true
}

// create an empty node and return its address
func createNode() *node {
	common.OpLogger.Print("createNode()")
	x := new(node)
	x.keys = make([]KeyType, 0, order)
	x.children = make([]*node, 0, order+1)
	common.OpLogger.Print("leave createNode()")
	return x
}

// PRIVATE FUNCTION

// Split node and return the new node
func (self *node) splitNode() (KeyType, *node) {
	common.OpLogger.Print("splitNode()\t", self)
	// Create node
	// The new node is the right brother of self
	n := createNode()
	n.leaf = self.isLeaf()
	n.parent = self.parent
	var remainCnt int
	if n.isLeaf() {
		remainCnt = int(math.Ceil(float64(order-1)/2.0))
	} else {
		remainCnt = int(math.Ceil(float64(order)/2.0) - 1)
	}
	var k KeyType
	if self.isLeaf() {
		n.keys = append(n.keys, self.keys[remainCnt:]...)
		n.children = append(n.children, self.children[0])
		self.children[0] = n

		k = self.keys[remainCnt]
		self.keys = self.keys[:remainCnt]
	} else {
		n.keys = append(n.keys, self.keys[remainCnt + 1:]...)
		n.children = append(n.children, self.children[remainCnt + 1:]...)

		k = self.keys[remainCnt]
		self.keys = self.keys[:remainCnt]
		self.children = self.children[:remainCnt + 1]
	}

	common.OpLogger.Print("leave splitNode()", self, "\t", k.GetValue(), "\t", n)
	return k, n
}

func (self *node) mergeRightNode(rb *node, k KeyType) {
	self.keys = append(self.keys, append([]KeyType{k}, rb.keys...)...)
	self.children = append(self.children, rb.children...)
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

func (self node) ChildCnt() int {
	return len(self.children)
}

func (self node) minKey() KeyType {
	common.OpLogger.Print("minKey()")
	var n *node
	for n = &self; !n.isLeaf(); n = n.children[0] {
	}
	common.OpLogger.Print("leave minKey()", n.keys[0])
	return n.keys[0]
}

// return first index of key that is greater or equal to v
// return a bool value indicating an exact match found.
func (self node) findKeyIndex(v common.CellValue) (int, bool) {
	for i := 0; i < self.keyCnt(); i++ {
		if !self.keys[i].GetValue().LessThan(v) {
			return i, self.keys[i].GetValue().EqualsTo(v)
		}
	}
	return self.keyCnt(), false
}

// return first index of key that is greater or equal than v.
func (self node) findChildIndex(v common.CellValue) int {
	for i := 0; i < self.keyCnt(); i++ {
		if !self.keys[i].GetValue().LessThan(v) {
			return i
		}
	}
	return self.keyCnt()
}

func (self *node) findLeafNode(v common.CellValue) *node {
	common.OpLogger.Print("search()\t", v)
	x := self
	for !x.isLeaf() {
		i := x.findChildIndex(v)
		x = x.children[i]
	}
	common.OpLogger.Print("leave search()\t", x)
	return x
}

// Insert k into leaf
func (self *node) insertKey(k KeyType) bool {
	common.OpLogger.Print("insertKey():\t", self, ",\t", k.GetValue())
	// l should be a leaf and not full
	if (!self.isLeaf()) || self.isFull() {
		common.OpLogger.Print("Error!")
		common.ErrLogger.Print("Should be a leaf and not full!\t", self)
		common.ErrLogger.Println()
		return false
	}

	i, _ := self.findKeyIndex(k.GetValue())
	self.keys = append(self.keys[:i], append([]KeyType{k}, self.keys[i:]...)...)

	common.OpLogger.Print("leave insertKey()\t", self)
	return true
}

// Insert c into non-leaf
func (self *node) insertChild(k KeyType, c *node) bool {
	common.OpLogger.Print("insertChild()")
	// l should be a non-leaf and not full
	if self.isLeaf() || self.isFull() {
		common.OpLogger.Print("Error!")
		common.ErrLogger.Print("Should be a non-leaf and not full\t", self)
		return false
	}
	// update parent
	c.parent = self
	// case 1: self is empty
	if self.ChildCnt() == 0 {
		common.OpLogger.Print("First Child!")
		self.children = append(self.children, c)
	} else {
		i := self.findChildIndex(k.GetValue())
		self.keys = append(self.keys[:i], append([]KeyType{k}, self.keys[i:]...)...)
		self.children = append(self.children[:i+1], append([]*node{c}, self.children[i+1:]...)...)
	}
	common.OpLogger.Print("leave insertChild()")
	return true
}

// Insert c into self's parent
// If new root is created, it will return (root, newRoot).
// Otherwise return (nil, false)
// @_@ It is not functional programming, keep it stupid.
func (self *node) insertInParent(k KeyType, c *node) (*node, bool) {
	common.OpLogger.Print("insertInParent()")
	// If l is the root of the tree, split it and create new root.
	if self.isRoot() {
		temp := createNode()
		temp.insertChild(self.minKey(), self)
		temp.insertChild(k, c)
		common.OpLogger.Print("leave insertInParent() with new root")
		return temp, true
	}

	p := self.parent
	p.insertChild(k, c)
	if p.isFull() {
		k1, p1 := p.splitNode()
		return p.insertInParent(k1, p1)
	}
	common.OpLogger.Print("leave insertInParent()")
	return nil, false
}

// Delete ith key.
func (self *node) deleteKey(i int) (KeyType, bool) {
	common.OpLogger.Print("deleteKey()\t", self, ", ", i)
	// Should be a leaf
	if !self.isLeaf() {
		common.OpLogger.Print("Error!")
		common.ErrLogger.Print("Should be a leaf\t", self)
		return nil, false
	}
	k := self.keys[i]
	self.keys = append(self.keys[:i], self.keys[i+1:]...)
	common.OpLogger.Print("leave deleteKey(), ", self)
	return k, true
}

// Delete ith child, return its minimum key and pointer
func (self *node) deleteChild(i int) (KeyType, *node, bool) {
	common.OpLogger.Print("deleteChild()\t", self, ", ", i)
	// Should be a non leaf
	if self.isLeaf() {
		common.OpLogger.Print("Error!")
		common.ErrLogger.Print("Should be a non leaf\t", self)
		return nil, nil, false
	}

	c := self.children[i]

	var k KeyType
	if i == 0 {
		k = self.minKey()
	} else {
		k = self.keys[i-1]
	}

	self.keys = append(self.keys[:i-1], self.keys[i:]...)
	self.children = append(self.children[:i], self.children[i+1:]...)

	common.OpLogger.Print("
	return k, c, true
}

// Test Functions
func (self *idxMan) Print() {
	c := make(chan *node, 1000)
	c <- self.root
	last := self.root
	for {
		n := <-c

		common.OpLogger.Printf("[%p]\t%v", n, n)

		if !n.isLeaf() {
			for _, child := range n.children {
				c <- child
				last = child
			}
		}

		if last == n {
			break
		}
	}
	close(c)
	common.OpLogger.Println()
}
