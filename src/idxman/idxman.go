package idxman

import (
	"../common"
	"../catman"
	"../recman"
	"math"
	"os"
	"fmt"
	"io"
)

// CONST
const order = 4 // We are implementing B+ Tree for n = 4.
const maxNodeCnt = 10000
const maxRecordCnt = 100000

// STRUCT
type node struct {
	parent *node
	// Size == order, reserve one.
	keys []common.CellValue
	// if it is a leaf, children[0] should point to its right sibling
	// Size == order + 1, reserve one.
	children []*node
	recordIds []int64
	leaf     bool
}

type idxMan struct {
	root *node
}

// PUBLIC FUNCTION

func NewIdxMan(fileName string, tableName string, indexName string) error {
	common.OpLogger.Print("NewIdxMan(): file name ", fileName, ", table name ", tableName, ", index name ", indexName)
	
	im, err := NewIdxManInMemory(fileName, tableName, indexName)
	if err != nil {
		err = im.FlushToDisk(fileName)
	}
	
	if err != nil {
		common.OpLogger.Print("leave NewIdxMan() with error")
		common.ErrLogger.Println("[NewIdxMan]", err)
		return err
	}
	common.OpLogger.Print("leave NewIdxMan()")
	return nil
}

func NewIdxManInMemory(fileName string, tableName string, indexName string) (*idxMan, error) {
	common.OpLogger.Print("NewIdxManInMemory(): file name ", fileName, ", table name ", tableName, ", index name ", indexName)
	
	idxs, err := catman.TableIndexes(tableName)
	if err != nil && searchString(idxs, indexName) {
		common.OpLogger.Print("leave NewIdxManInMemory() with error")
		common.ErrLogger.Print("[NewIdxManInMemory]", err)
		return nil, err
	}
	
	file, err := os.OpenFile(common.DataDir+"/"+tableName+"/data.dbf", os.O_RDONLY, 0600)
	if err != nil {
		common.OpLogger.Print("leave NewIdxManInMemory() with error")
		common.ErrLogger.Print("[NewIdxManInMemory]", err)
		return nil, err
	}
	defer file.Close()
	tabinfo, err := catman.TableInfo(tableName)
	if err != nil {
		common.OpLogger.Print("leave NewIdxManInMemory() with error")
		common.ErrLogger.Print("[NewIdxManInMemory]", err)
		return nil, err
	}
	records, recordIds := recman.ReadRecords(file, tabinfo)
	
	im := NewEmptyIdxMan()
	for i, record := range records {
		im.Insert(record.Values[indexName], recordIds[i])
	}
	
	common.OpLogger.Print("leave NewIdxManInMemory()")
	return im, nil
}

func searchString(s []string, x string) bool {
	for _, y := range s {
		if x == y {
			return true
		}
	}
	return false
}

// Creat an index manager and give it a empty root.
func NewEmptyIdxMan() *idxMan {
	common.OpLogger.Print("NewEmptyIdxMan(): Create a empty B+ Tree!")
	
	im := new(idxMan)
	im.root = createNode(true)
	
	common.OpLogger.Print("leave NewEmptyIdxMan()")
	return im
}

func DestroyIdxMan(fileName string) error {
	common.OpLogger.Print("DestroyIdxMan\t", fileName)
	err := os.Remove(fileName)
	if err != nil {
		common.OpLogger.Print("leave DestroyIdxMan with error")
		common.ErrLogger.Print("[DestroyIdxMan]", err)
		return err
	}
	common.OpLogger.Print("leave DestroyIdxMan")
	return nil
}

// Disk file format(Each record corresponding to one line):
// no, pno, leaf, keyCnt, keys..., recordIds...
func (self *idxMan) FlushToDisk(fileName string) error {
	common.OpLogger.Print("FlushToDisk\t", fileName)
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {		
		common.OpLogger.Print("leave FlushToDisk with error")
		common.ErrLogger.Print("[FlushToDisk]", err)
		return err
	}
	defer file.Close()
	
	mapHelper := make(map[*node] int64)
	queue := make(chan *node, maxNodeCnt)
	queue <- self.root
	last := self.root
	var n *node
	var no int64
	var pno int64
	var ok bool
	no = 0
	for {
		n = <-queue
		mapHelper[n] = no
		pno, ok = mapHelper[n.parent]
		if !ok {
			pno = -1
		}
		fmt.Fprintf(file, "%d %d %t ", no, pno, n.leaf)
		no++
		fmt.Fprint(file, n.keyCnt())
		for i := 0; i < n.keyCnt(); i++ {
			fmt.Fprint(file, " ", n.keys[i])
		}
		if n.isLeaf() {
			for i := 0; i < n.keyCnt(); i++ {
				fmt.Fprint(file, " ", n.recordIds[i])
			}
		}
		if !n.isLeaf() {
			for _, child := range n.children {
				queue <- child
				last = child
			}
		}
		
		if last == n {
			break
		}
	}
	close(queue)
	common.OpLogger.Print("leave FlushToDisk\t")
	return nil
}

func ConstructFromDisk(fileName string) (*idxMan, error) {
	common.OpLogger.Print("ConstructFromDisk()")
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0600)
	if err != nil {
		common.OpLogger.Print("leave ConstructFromDisk with error")
		common.ErrLogger.Print("[ConstructFromDisk]", err)
		return nil, err
	}
	defer file.Close()
	
	mapHelper := make(map[int64] *node)
	var n *node
	var p *node
	var no int64
	var pno int64
	var leaf bool
	var keyCnt int
	fmt.Fscan(file, no, pno, leaf, keyCnt)
	im := new(idxMan)
	im.root = createNode(leaf)
	mapHelper[no] = im.root
	im.root.constructKeys(file, keyCnt)
	for {
		_, err := fmt.Fscan(file, no, pno, leaf, keyCnt)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		n = createNode(leaf)
		p, _ = mapHelper[pno]
		mapHelper[no] = n
		p.children = append(p.children, n)
		n.constructKeys(file, keyCnt)
	}
	common.OpLogger.Print("leave ConstructFromDisk")
	return im, nil
}

func (self *node) constructKeys(file *os.File, keyCnt int) {
	var key common.CellValue
	for i := 0; i < keyCnt; i++ {
		fmt.Fscan(file, key)
		self.keys = append(self.keys, key)
	}
	var recordId int64
	if self.isLeaf() {
		for i := 0; i < keyCnt; i++ {
			fmt.Fscan(file, recordId)
			self.recordIds = append(self.recordIds, recordId)
		}
	}
}

func Insert(fileName string, v common.CellValue, id int64) error {
	common.OpLogger.Print("Insert(): Insert a cell into the file")
	im, err := ConstructFromDisk(fileName)
	if err != nil {
		common.OpLogger.Print("leave Insert() with error")
		common.ErrLogger.Print("[Insert]", err)
		return err
	}
	im.Insert(v, id)
	im.FlushToDisk(fileName)
	common.OpLogger.Print("leave Insert()")
	return nil
}

func Delete(fileName string, v common.CellValue) (int64, bool, error) {
	common.OpLogger.Print("Delete(): Delete a node into the file")
	im, err := ConstructFromDisk(fileName)
	if err != nil {
		common.OpLogger.Print("leave Delete() with error")
		common.ErrLogger.Print("[Delete]", err)
		return 0, false, err
	}
	id, present := im.Delete(v)
	if present {
		common.OpLogger.Print("leave Delete(), no node deleted.")
		return 0, false, nil
	}
	im.FlushToDisk(fileName)
	common.OpLogger.Print("leave Delete()")
	
	return id, true, nil
}

func LinearSelectEqual(tableName string, keyName string, v common.CellValue) ([]int64, error) {
	file, err := os.OpenFile(common.DataDir+"/"+tableName+"/data.dbf", os.O_RDONLY, 0600)
	if err != nil {
		common.OpLogger.Print("leave LinearSelectEqual() with error")
		common.ErrLogger.Print("[LinearSelectEqual]", err)
		return nil, err
	}
	defer file.Close()
	tabinfo, err := catman.TableInfo(tableName)
	if err != nil {
		common.OpLogger.Print("leave LinearSelectEqual() with error")
		common.ErrLogger.Print("[LinearSelectEqual]", err)
		return nil, err
	}
	result := make([]int64, 0, maxRecordCnt)
	records, recordIds := recman.ReadRecords(file, tabinfo)
	for i, record := range records {
		if (record.Values[keyName].EqualsTo(v)) {
			result = append(result, recordIds[i])
		}
	}
	return result, nil
}

func LinearSelectRange(tableName string, keyName string, left common.CellValue, right common.CellValue) ([]int64, error) {
	file, err := os.OpenFile(common.DataDir+"/"+tableName+"/data.dbf", os.O_RDONLY, 0600)
	if err != nil {
		common.OpLogger.Print("leave LinearSelectEqual() with error")
		common.ErrLogger.Print("[LinearSelectEqual]", err)
		return nil, err
	}
	defer file.Close()
	tabinfo, err := catman.TableInfo(tableName)
	if err != nil {
		common.OpLogger.Print("leave LinearSelectEqual() with error")
		common.ErrLogger.Print("[LinearSelectEqual]", err)
		return nil, err
	}
	result := make([]int64, 0, maxRecordCnt)
	records, recordIds := recman.ReadRecords(file, tabinfo)
	for i, record := range records {
		if (!record.Values[keyName].LessThan(left) && !record.Values[keyName].GreaterThan(right)) {
			result = append(result, recordIds[i])
		}
	}
	return result, nil
}

// Find the first common.CellValue containing given v,
// return (nil, false) if nothing is found.
func (self idxMan) SelectEqual(v common.CellValue) []int64 {
	common.OpLogger.Print("SelectEqual():\t", v)
	l := self.root.findLeafNode(v)
	i, found := l.findKeyIndex(v)
	if found {
		common.OpLogger.Print("leave SelectRange()\t", l.keys[i])
		return []int64{l.recordIds[i]}
	}
	common.OpLogger.Print("leave SelectRange(), no record found.")
	return nil
}

func (self idxMan) SelectRange(left common.CellValue, right common.CellValue) []int64 {
	common.OpLogger.Print("SelectRange():\t", left, ", ", right)
	l := self.root.findLeafNode(left)
	i, found := l.findKeyIndex(left)
	if ! found {
		common.OpLogger.Print("leave SelectRange(), no record is found")
		return nil
	}
	result := make([]int64, 0, maxRecordCnt)
	for !l.keys[i].GreaterThan(right) {
		result = append(result, l.recordIds[i])
		i++
		if i == l.keyCnt() {
			l = l.children[0]
		}
		if l == nil {
			break
		}
	}
	
	common.OpLogger.Print("leave SelectRange()")
	return result
}

// Insert v into B+ Tree
func (self *idxMan) Insert(v common.CellValue, id int64) {
	common.OpLogger.Print("Insert(): Insert ", v)
	l := self.root.findLeafNode(v)
	l.insertKey(v, id)
	// If the l is full, split it.
	// Then insert two nodes l and l1 into their parent.
	// Update root if new root is created.
	if l.isFull() {
		common.OpLogger.Print("Split!")
		v1, l1 := l.splitNode()
		r, newRoot := l.insertInParent(v1, l1)
		if newRoot {
			self.root = r
		}
	}
	common.OpLogger.Print("leave Insert()")
}

// Delete the first common.CellValue containing given v,
// return (false) if nothing is deleted
func (self *idxMan) Delete(v common.CellValue) (int64, bool) {
	common.OpLogger.Print("Delete(): Delete ", v)
	// Find the corresponding leaf l
	l := self.root.findLeafNode(v)

	// Find the index of v in l
	i, found := l.findKeyIndex(v)
	if !found {
		common.OpLogger.Print("leave Delete(), no key deleted")
		return 0, false
	}
	_, id, _ := l.deleteKey(i)

	common.OpLogger.Print("Start checking children number")
	// Check node's children number back to root
	for n := l; ; n = n.parent {
		common.OpLogger.Print("n = ", n)
		// A leaf root node has between 0 and order - 1 values
		// A leaf node has between ceil((order - 1) / 2) and order - 1 values.
		if n.isLeaf() {
                        if n.isRoot() || n.keyCnt() >= int(math.Ceil(float64(order-1)/2.0)) {
				common.OpLogger.Print("Leaf is good.")
				break
			}
		} else {
		// A non-leaf root node has between 2 and order children
			if n.isRoot() {
				if n.childCnt() == 1 {
					self.root = n.children[0]
					self.root.parent = nil // root is root!
					common.OpLogger.Print("Root has only one child.")
				} else {
					common.OpLogger.Print("Non-Leaf Root is good.")
				}
				break
                        }
		// A node that is not a leaf or root has between ceil(order / 2) and order children.
			if n.keyCnt() + 1 >= int(math.Ceil(float64(order)/2.0)) {
				common.OpLogger.Print("Non-Leaf is good.")
				break
			}
		}
		p := n.parent
		i := p.findChildIndex(n.minKey())
		// n is ith child.
		if i == p.childCnt() - 1 {
			n1 := p.children[i-1]
			if n.keyCnt()+n1.keyCnt() < order {
				// Case0: merge n with its left brother
				common.OpLogger.Print("Case0: merge n with its left brother")
				common.OpLogger.Print(n1, n)
				n1.mergeRightNode(n, p.keys[i-1])
				p.deleteChild(i)
				common.OpLogger.Print(n1)
				common.OpLogger.Print("leave Case0")
			} else {
				// Case1: borrow a child from left brother
				common.OpLogger.Print("Case1: borrow a child from left brother")
				common.OpLogger.Print(n1, n)
				if n.isLeaf() {
					k, id, _ := n1.deleteKey(n1.keyCnt()-1)
					n.insertKey(k, id)
				} else {
					k, c, _ := n1.deleteChild(n1.childCnt()-1)
					n.insertChild(k, c)
				}
				common.OpLogger.Print(n1, n)
				common.OpLogger.Print("leave Case1")
			}
		} else {
			n1 := p.children[i+1]
			if n.keyCnt()+n1.keyCnt() < order {
				// Case2: merge n with its right brother
				common.OpLogger.Print("Case2: merge n with its right brother")
				common.OpLogger.Print(n, n1)
				n.mergeRightNode(n1, p.keys[i])
				p.deleteChild(i+1)
				common.OpLogger.Print(n)
				common.OpLogger.Print("leave Case2")
			} else {
				// Case3: borrow a child from right brother
				common.OpLogger.Print("Case3: borrow a child from right brother")
				common.OpLogger.Print(n, n1)
				if n.isLeaf() {
					k, id, _ := n1.deleteKey(0)
					n.insertKey(k, id)
				} else {
					k, c, _ := n1.deleteChild(0)
					n.insertChild(k, c)
				}
				common.OpLogger.Print(n, n1)
				common.OpLogger.Print("leave Case3")
			}
		}
	}
	common.OpLogger.Print("leave Delete()\t", id)	
	return id, true
}

// Print the whole index manager to logger
func (self *idxMan) Print() {
	c := make(chan *node, maxNodeCnt)
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
}

// create an empty node and return its address
func createNode(isLeaf bool) *node {
	common.OpLogger.Print("createNode()")
	x := new(node)
	x.keys = make([]common.CellValue, 0, order)
	x.leaf = isLeaf
	if isLeaf {
		x.children = make([]*node, 1, 1)
		x.recordIds = make([]int64, 0, order)
	} else {
		x.children = make([]*node, 0, order+1)
	}
	common.OpLogger.Print("leave createNode()")
	return x
}

// PRIVATE FUNCTION

// Split node and return the new node
func (self *node) splitNode() (common.CellValue, *node) {
	common.OpLogger.Print("splitNode()\t", self)
	// Create node
	// The new node is the right brother of self
	n := createNode(self.isLeaf())
	n.parent = self.parent
	var remainCnt int
	if n.isLeaf() {
		remainCnt = int(math.Ceil(float64(order-1)/2.0))
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
		n.keys = append(n.keys, self.keys[remainCnt + 1:]...)
		n.children = append(n.children, self.children[remainCnt + 1:]...)
		// Update child's parent
		for _, c := range n.children {
			c.parent = n
		}
		k = self.keys[remainCnt]
		self.keys = self.keys[:remainCnt]
		self.children = self.children[:remainCnt + 1]
	}

	common.OpLogger.Print("leave splitNode()", self, "\t", k, "\t", n)
	return k, n
}

func (self *node) mergeRightNode(rb *node, k common.CellValue) {
	common.OpLogger.Print("mergeRightNode()\t", self, ", ", rb, ", ", k)
	if self.isLeaf() {
		self.keys = append(self.keys, rb.keys...)
		self.recordIds = append(self.recordIds, rb.recordIds...)
		self.children[0] = rb.children[0]
	} else {
		self.keys = append(self.keys, append([]common.CellValue{k}, rb.keys...)...)
		self.children = append(self.children, rb.children...)
	}
	common.OpLogger.Print("leave mergeRightNode()", self)
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
	common.OpLogger.Print("minKey()")
	var n *node
	for n = &self; !n.isLeaf(); n = n.children[0] {
	}
	common.OpLogger.Print("leave minKey()\t", n.keys[0])
	return n.keys[0]
}

// return first index of key that is greater or equal to v
// return a bool value indicating an exact match found.
func (self node) findKeyIndex(v common.CellValue) (int, bool) {
	common.OpLogger.Print("findKeyIndex()\t", self, ", ", v)
	for i := 0; i < self.keyCnt(); i++ {
		if !self.keys[i].LessThan(v) {
			common.OpLogger.Print("leave findKeyIndex()\t", i)
			return i, self.keys[i].EqualsTo(v)
		}
	}
	common.OpLogger.Print("leave findKeyIndex()\t", self.keyCnt())
	return self.keyCnt(), false
}

// return first index of key that is greater than v.
// if no such key is found, return self.keyCnt（）
func (self node) findChildIndex(v common.CellValue) int {
	common.OpLogger.Print("findChildIndex()\t", self, ", ", v)
	for i := 0; i < self.keyCnt(); i++ {
		if self.keys[i].GreaterThan(v) {
			common.OpLogger.Print("leave findChildIndex()\t", i)
			return i
		}
	}
	common.OpLogger.Print("leave findChildIndex()\t", self.keyCnt())
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
func (self *node) insertKey(k common.CellValue, id int64) bool {
	common.OpLogger.Print("insertKey():\t", self, ",\t", k)
	// l should be a leaf and not full
	if (!self.isLeaf()) || self.isFull() {
		common.OpLogger.Print("Error!")
		common.ErrLogger.Print("Should be a leaf and not full!\t", self)
		common.ErrLogger.Println()
		return false
	}

	i, _ := self.findKeyIndex(k)
	self.keys = append(self.keys[:i], append([]common.CellValue{k}, self.keys[i:]...)...)
	self.recordIds = append(self.recordIds[:i], append([]int64{id}, self.recordIds[i:]...)...)

	common.OpLogger.Print("leave insertKey()\t", self)
	return true
}

// Insert c into non-leaf
func (self *node) insertChild(k common.CellValue, c *node) bool {
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
	if self.childCnt() == 0 {
		common.OpLogger.Print("First Child!")
		self.children = append(self.children, c)
	} else {
		i := self.findChildIndex(k)
		self.keys = append(self.keys[:i], append([]common.CellValue{k}, self.keys[i:]...)...)
		self.children = append(self.children[:i+1], append([]*node{c}, self.children[i+1:]...)...)
	}
	common.OpLogger.Print("leave insertChild()")
	return true
}

// Insert c into self's parent
// If new root is created, it will return (root, newRoot).
// Otherwise return (nil, false)
// @_@ It is not functional programming, keep it stupid.
func (self *node) insertInParent(k common.CellValue, c *node) (*node, bool) {
	common.OpLogger.Print("insertInParent()")
	// If l is the root of the tree, split it and create new root.
	if self.isRoot() {
		temp := createNode(false)
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
func (self *node) deleteKey(i int) (common.CellValue, int64, bool) {
	common.OpLogger.Print("deleteKey()\t", self, ", ", i)
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
	common.OpLogger.Print("leave deleteKey(), ", self)
	return k, id, true
}

// Delete ith child, return its minimum key and pointer
func (self *node) deleteChild(i int) (common.CellValue, *node, bool) {
	common.OpLogger.Print("deleteChild()\t", self, ", ", i)
	// Should be a non leaf
	if self.isLeaf() {
		common.OpLogger.Print("Error!")
		common.ErrLogger.Print("Should be a non leaf\t", self)
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

	common.OpLogger.Print("leave deleteChild()\t", k, ", ", c)
	return k, c, true
}
