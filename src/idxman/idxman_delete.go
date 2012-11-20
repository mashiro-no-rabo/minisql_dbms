package idxman

import (
	"../common"
	"math"
)

// Delete the first common.CellValue containing given v,
// return (false) if nothing is deleted
func (self *idxMan) Delete(v common.CellValue) (int64, bool) {
	common.OpLogger.Print("[Delete]", v)
	defer common.OpLogger.Print("[leave Delete]")
	// Find the corresponding leaf l
	l := self.root.findLeafNode(v)

	// Find the index of v in l
	i, found := l.findKeyIndex(v)
	if !found {
		common.OpLogger.Print("[No key deleted.]")
		return 0, false
	}
	_, id, _ := l.deleteKey(i)

	common.OpLogger.Print("[Start checking children number]")
	// Check node's children number back to root
	for n := l; ; n = n.parent {
		common.OpLogger.Print("n = ", n)
		// A leaf root node has between 0 and order - 1 values
		// A leaf node has between ceil((order - 1) / 2) and order - 1 values.
		if n.isLeaf() {
			if n.isRoot() || n.keyCnt() >= int(math.Ceil(float64(order-1)/2.0)) {
				common.OpLogger.Print("[Leaf is good.]")
				break
			}
		} else {
			// A non-leaf root node has between 2 and order children
			if n.isRoot() {
				if n.childCnt() == 1 {
					self.root = n.children[0]
					self.root.parent = nil // root is root!
					common.OpLogger.Print("[Root has only one child.]")
				} else {
					common.OpLogger.Print("[Non-Leaf Root is good.]")
				}
				break
			}
			// A node that is not a leaf or root has between ceil(order / 2) and order children.
			if n.keyCnt()+1 >= int(math.Ceil(float64(order)/2.0)) {
				common.OpLogger.Print("[Non-Leaf is good.]")
				break
			}
		}
		p := n.parent
		i := p.findChildIndex(n.minKey())
		// n is ith child.
		if i == p.childCnt()-1 {
			n1 := p.children[i-1]
			if n.keyCnt()+n1.keyCnt() < order {
				// Case0: merge n with its left brother
				common.OpLogger.Print("[Case0: merge n with its left brother]")
				common.OpLogger.Print(n1, n)
				n1.mergeRightNode(n, p.keys[i-1])
				p.deleteChild(i)
				common.OpLogger.Print(n1)
				common.OpLogger.Print("[leave Case0]")
			} else {
				// Case1: borrow a child from left brother
				common.OpLogger.Print("[Case1: borrow a child from left brother]")
				common.OpLogger.Print(n1, n)
				if n.isLeaf() {
					k, id, _ := n1.deleteKey(n1.keyCnt() - 1)
					n.insertKey(k, id)
				} else {
					k, c, _ := n1.deleteChild(n1.childCnt() - 1)
					n.insertChild(k, c)
				}
				common.OpLogger.Print(n1, n)
				common.OpLogger.Print("[leave Case1]")
			}
		} else {
			n1 := p.children[i+1]
			if n.keyCnt()+n1.keyCnt() < order {
				// Case2: merge n with its right brother
				common.OpLogger.Print("[Case2: merge n with its right brother]")
				common.OpLogger.Print(n, n1)
				n.mergeRightNode(n1, p.keys[i])
				p.deleteChild(i + 1)
				common.OpLogger.Print(n)
				common.OpLogger.Print("[leave Case2]")
			} else {
				// Case3: borrow a child from right brother
				common.OpLogger.Print("[Case3: borrow a child from right brother]")
				common.OpLogger.Print(n, n1)
				if n.isLeaf() {
					k, id, _ := n1.deleteKey(0)
					n.insertKey(k, id)
				} else {
					k, c, _ := n1.deleteChild(0)
					n.insertChild(k, c)
				}
				common.OpLogger.Print(n, n1)
				common.OpLogger.Print("[leave Case3]")
			}
		}
	}
	return id, true
}

func (self *node) mergeRightNode(rb *node, k common.CellValue) {
	common.OpLogger.Print("[mergeRightNode]", self, ",", rb, ",", k)
	defer common.OpLogger.Print("[leave mergeRightNode]", self)
	
	if self.isLeaf() {
		self.keys = append(self.keys, rb.keys...)
		self.recordIds = append(self.recordIds, rb.recordIds...)
		self.children[0] = rb.children[0]
	} else {
		self.keys = append(self.keys, append([]common.CellValue{k}, rb.keys...)...)
		self.children = append(self.children, rb.children...)
	}
}

// Delete ith key.
func (self *node) deleteKey(i int) (common.CellValue, int64, bool) {
	common.OpLogger.Print("[deleteKey]", self, ",", i)
	defer common.OpLogger.Print("[leave deleteKey]", self)
	
	// Should be a leaf
	if !self.isLeaf() {
		common.OpLogger.Print("Error!")
		common.ErrLogger.Print("Should be a leaf\t", self)
		return nil, 0, false
	}
	k := self.keys[i]
	id := self.recordIds[i]
	self.keys = append(self.keys[:i], self.keys[i+1:]...)
	self.recordIds = append(self.recordIds[:i], self.recordIds[i+1:]...)
	return k, id, true
}

// Delete ith child, return its minimum key and pointer
func (self *node) deleteChild(i int) (common.CellValue, *node, bool) {
	common.OpLogger.Print("[deleteChild]", self, ",", i)
	defer common.OpLogger.Print("[leave deleteChild]")
	
	// Should be a non leaf
	if self.isLeaf() {
		common.ErrLogger.Print("Should be a non leaf ", self)
		return nil, nil, false
	}

	c := self.children[i]

	var k common.CellValue
	if i == 0 {
		k = self.minKey()
	} else {
		k = self.keys[i-1]
	}

	self.keys = append(self.keys[:i-1], self.keys[i:]...)
	self.children = append(self.children[:i], self.children[i+1:]...)

	return k, c, true
}
