package idxman

// CONST
const order = 4 // We are implementing B+ Tree for n = 4.

// INTERFACE
type ValueType interface {
	EqualsTo(ValueType) bool
	LessThan(ValueType) bool
	Val() (int, float64, string)
}

type KeyType interface {
	GetValue() ValueType
}

// STRUCT
type node struct {
	parent *node
	keyCnt int
	keys   []KeyType
	// if it is a leaf, children[0] should point to its right sibling
	children []*node
	leaf     bool
}

type idxMan struct {
	root *node
}

// PUBLIC FUNCTION

// Creat an index manager and give it a empty root.
func NewIdxMan() *idxMan {
	im := new(idxMan)
	im.root = createNode()
	return im
}

// Find the first KeyType containing given v,
// return (nil, false) if nothing is found.
func (self idxMan) Find(v ValueType) (KeyType, bool) {
	l := self.root.findLeafNode(v)
	i, found := l.findValueIndex(v)
	temp := l.keys[i]
	if found && temp.GetValue().EqualsTo(v) {
		return temp, true
	}
	return nil, false
}

// Insert k into B+ Tree
func (self *idxMan) Insert(k KeyType) {
	l := self.root.findLeafNode(k.GetValue())
	l.insertKey(k)
	// If the l is full, split it.
	// Then insert two nodes l and l1 into their parent.
	// Update root if new root is created.
	if l.isFull() {
		l1 := l.splitNode()
		r, newRoot := l.insertInParent(l1)
		if newRoot {
			self.root = r
		}
	}
}

// Delete the first KeyType containing given v,
// return (false) if nothing is deleted
func (self *idxMan) Delete(v ValueType) (KeyType, bool) {
	// Find the corresponding leaf l
	l := self.root.findLeafNode(v)

	// Find the index of v in l
	i, found := l.findValueIndex(v)
	if !found {
		return nil, false
	}
	k, _ := l.deleteKey(i)

	// Check node's children number back to root
	for n := l; n.keyCnt < order/2; n = n.parent {
		if n.isRoot() && n.keyCnt == 1 {
			self.root = n.children[0]
			break
		}
		p := n.parent
		i, _ := p.findValueIndex(n.minKey().GetValue())
		if i == p.keyCnt-1 {
			n1 := p.children[i-1]
			if n.keyCnt+n1.keyCnt < order {
				// Case0: merge n with its left brother
				n1.mergeRightNode(n, p.keys[i-1])
				p.deleteChild(i)
			} else {
				// Case1: borrow a child from left brother
				if n.isLeaf() {
					k, _ := n1.deleteKey(n1.keyCnt - 1)
					n.insertKey(k)
				} else {
					c, _ := n1.deleteChild(n1.keyCnt - 1)
					n.insertChild(c)
				}
			}
		} else {
			n1 := p.children[i+1]
			if n.keyCnt+n1.keyCnt < order {
				// Case2: merge n with its right brother
				n.mergeRightNode(n1, p.keys[i])
				p.deleteChild(i + 1)
			} else {
				// Case3: borrow a child from right brother
				if n.isLeaf() {
					k, _ := n1.deleteKey(0)
					n.insertKey(k)
				} else {
					c, _ := n1.deleteChild(0)
					n.insertChild(c)
				}
			}
		}
	}

	return k, true
}

// create an empty node and return its address
func createNode() *node {
	x := new(node)
	x.keys = make([]KeyType, 0, order)
	x.children = make([]*node, 0, order)
	return x
}

// PRIVATE FUNCTION

// Split node and return the new node
func (self *node) splitNode() *node {
	// Create node
	// The new node is the right brother of self
	n := createNode()
	n.leaf = self.isLeaf()
	n.parent = self.parent
	n.keyCnt = order/2 + order%2
	if self.isLeaf() {
		n.keys = append(n.keys, self.keys[order/2:]...)
		n.children[0] = self.children[0]
		self.children[0] = n
	} else {
		n.keys = append(n.keys, self.keys[order/2:]...)
		n.children = append(n.children, self.children[order/2:]...)
	}
	// Update self
	self.keyCnt = order / 2
	return n
}

func (self *node) mergeRightNode(rb *node, k KeyType) {
	self.keys = append(self.keys, append([]KeyType{k}, rb.keys...)...)
	self.children = append(self.children, rb.children...)
	self.keyCnt += rb.keyCnt
}

func (self node) isFull() bool {
	return self.keyCnt == order
}

func (self node) isRoot() bool {
	return self.parent == nil
}

func (self node) isLeaf() bool {
	return self.leaf
}

// func: minKey
func (self node) minKey() KeyType {
	var n *node
	for n = &self; !n.isLeaf(); n = n.children[0] {
	}
	return n.keys[0]
}

// return first index of Child having key greater or
// equal than v. 0 <= i <= keyCnt - 1
// In a leaf, the ith Child is ith Key
// In a non-leaf, the ith Child's key is (i - 1)th Key,
// and 0th Child has no key in the Keys.
func (self node) findValueIndex(v ValueType) (i int, found bool) {
	for i = 0; i < self.keyCnt; i++ {
		if !self.keys[i].GetValue().LessThan(v) {
			return i, true
		}
	}
	return self.keyCnt, false
}

// Find the leaf node contains (or should contain) v.
func (self node) findLeafNode(v ValueType) *node {
	x := &self
	for !x.isLeaf() {
		i, _ := x.findValueIndex(v)
		x = x.children[i]
	}
	return x
}

// Insert k into leaf
func (self *node) insertKey(k KeyType) bool {
	// l should be a leaf and not full
	if (!self.isLeaf()) || self.isFull() {
		return false
	}

	i, _ := self.findValueIndex(k.GetValue())
	self.keys = append(self.keys[:i], append([]KeyType{k}, self.keys[i:]...)...)
	self.keyCnt++
	return true
}

// Insert c into non-leaf
func (self *node) insertChild(c *node) bool {
	// l should be a non-leaf and not full
	if self.isLeaf() || self.isFull() {
		return false
	}

	c.parent = self

	k := c.minKey()
	i, _ := self.findValueIndex(k.GetValue())
	if i == 0 {
		k = self.minKey()
		self.keys = append([]KeyType{k}, self.keys...)
		self.children = append([]*node{c}, self.children...)
	} else {
		self.keys = append(self.keys[:i-1], append([]KeyType{k}, self.keys[i-1:]...)...)
		self.children = append(self.children[:i], append([]*node{c}, self.children[i:]...)...)
	}
	self.keyCnt++
	return true
}

// Insert c into self's parent
// If new root is created, it will return (root, newRoot).
// Otherwise return (nil, false)
// @_@ It is not functional programming, keep it stupid.
func (self *node) insertInParent(c *node) (*node, bool) {
	// If l is the root of the tree, split it and create new root.
	if self.isRoot() {
		temp := new(node)
		temp.insertChild(self)
		temp.insertChild(c)
		return temp, true
	}

	p := self.parent
	p.insertChild(c)
	if p.isFull() {
		p1 := p.splitNode()
		return p.insertInParent(p1)
	}
	return nil, false
}

// Delete ith key.
func (self *node) deleteKey(i int) (KeyType, bool) {
	if !self.isLeaf() {
		return nil, false
	}
	k := self.keys[i]
	self.keys = append(self.keys[:i], self.keys[i+1:]...)
	return k, true
}

// Delete ith child, return its minimum key and pointer
func (self *node) deleteChild(i int) (*node, bool) {
	if self.isLeaf() {
		return nil, false
	}

	c := self.children[i]

	self.keys = append(self.keys[:i-1], self.keys[i:]...)
	self.children = append(self.children[:i], self.children[i+1:]...)
	self.keyCnt--

	return c, true
}
