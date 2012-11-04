package idxman

import (
        "../common"
)

const order = 4 // We are implementing B+ Tree for n = 4.

type node struct {
    isLeaf bool
    parent *node
    keyCnt int
    keys [] KeyType
    // In leaf node, children[keyCnt] points to next leaf
    children [] *node
}

type idxMan struct {
    root *node
}

// func: NewIdxMan
// Creat an index manager and give it a empty root.
func NewIdxMan() *idxMan {
    im := new(idxMan)
    im.root = createNode()
    return im
}

func (self idxMan) Find(v ValueType) (KeyType, bool) {
    return self.root.search(v)
}

func (self *idxMan) Insert(k KeyType) {
    l := self.root.findLeafNode(k.GetValue())
    insertInLeaf(l, k)
    if l.isFull() {
        l1, k1 := splitNode(l)
        r, newRoot := insertInParent(l, k1, l1)
        if (newRoot) {
            self.root = r
        }
    }
}

func (self *idxMan) Delete(k KeyType) bool {
    l := self.root.findLeafNode(k.GetValue())
    i, found := l.findLeafNode(k.GetValue())
    if ! found {
        return false
    }
    root, newRoot := deleteEntry(l, i)
    if (newRoot) {
        self.root = root
    }
    return true
}

// func: createNode
// create an empty node and return its address
func createNode() *node {
    x := new(node)
    x.keys = make([] KeyType, 0, order)
    x.children = make([] KeyType, 0, order)
    return x
}

func (self node) findValueIndex(v ValueType) (i int, found bool) {
    for i = 0; i < self.keyCnt; i++ {
        temp := self.keys[i]
        if ! temp.GetValue().LessThan(v) {
            return i, true
        }
    }
    return self.keyCnt, false
}

func (self node) search(v ValueType) (KeyType, bool) {
    x := self
    for ! x.isLeaf {
        i, found := x.findValueIndex(v)
        temp := x.keys[i]
        if found && temp.GetValue().IsEqual(v) {
            return temp, true // Here we only return the first accepted key.
        }
        x = *x.children[i]
    }
    return nil, false
}

func (self node) findLeafNode(v ValueType) (l *node) {
    x := self
    for ! x.isLeaf {
        i, _ := x.findValueIndex(v)
        x = *x.children[i]
    }
    return &x
}

func (self node) isFull() bool {
    return self.keyCnt == order
}

func (self node) isRoot() bool {
    return self.parent == nil
}

func (self *node) deleteChild(i int) (k KeyType, c *node) {
    c = self.children[0]
    if (i == 0) {
        if self.isLeaf() {
           k = nil 
        } else {
            k = self.children[0].keys[0]
        }
    } else {
        k = self.keys[i - 1]
    }
    self.keys = append(self.keys[:i], self.keys[i + 1:]...)
    self.children = append(self.children[:i + 1], self.children[i + 2:]...)
    self.keyCnt--
    return k, c
}

func (self *node) insertChild(i int, k KeyType, c *node) {
    if (i == 0） {
        k1 = self.children[0].keys[0]
        self.keys = append([]KeyType{k1}, self.keys...)
        self.children = append([]*node{c}, self.children...}
    } else {
        self.keys = append(self.keys[:i], append([]KeyType{k}, self.keys[i:]...)...)
        self.children = append(self.children[:i - 1], append([]*node{c}, self.children[i - 1:]...)...)
    }
    self.keyCnt++
}

func insertInLeaf(l *node, k KeyType) {
    i, _ := l.findValueIndex(k.GetValue())
    l.keys = append(l.keys[:i + 1], append([]KeyType{k}, l,keys[i + 1:])...)
    l.keyCnt++
}

func splitNode(l *node) (l1 *node, k1 KeyType) {
    // Create node l1
    l1 = createNode()
    l1.isLeaf = true
    l1.parent = l.parent
    l1.keyCnt = order / 2 + order % 2
    // Copy keys from l to l1
    copy(l1.keys, l.keys[0:order / 2])
    // Update l
    l.keyCnt = order / 2
    l1.children[l1.keyCnt] = l.children[l.keyCnt]
    l.children[l.keyCnt] = l1 
    // return l1 and k1
    return l1, l1.keys[0]
}

func insertInParent(l *node, k1 KeyType, l1 *node) (r *node, newRoot bool) {
    // l is the root of the tree
    if l.isRoot() {
        temp := new(node)
        temp.keyCnt = 1
        temp.keys[0] = k1
        temp.children[0] = l
        temp.children[1] = l1
        return temp, true
    }
    
    p := l.parent
    p.keys[p.keyCnt] = k1
    p.children[p.keyCnt + 1] = l1 // Insert child link into right of key
    p.keyCnt++
    if p.isFull() {
        p1, k2 := splitNode(p)
        return insertInParent(p, k2, p1)
    }
    return nil, false
}

func deleteEntry(n *node, i int) （root int, newRoot bool） {
    // Delete k from l
    n.keys = append(n.keys[:i], n.keys[i + 1:])
    n.children = append(n.children[:i + 1], n.children[i + 2:])
    n.keyCnt--
    // if n is the root and n has only one remaining children
    if n.isRoot() && n.keyCnt == 1 {
        return n.children[0], true
    }
    // if n has too few values
    if n.keyCnt < order / 2 {
        p := n.parent
        i1 := p.findValueIndex(n.keys[0].GetValue())
        if i == p.keyCnt - 1 {
            n1 := p.children[i1 - 1]
            if (n.keyCnt + n1.keyCnt < order) {
                p.children[i1 - 1] = merge(n1, n)
            } else {
                // borrow entry from i1 - 1
               k, c := n1.delete(n1.keyCnt)
               n.insertChild(0, k, c)
            }
            deleteEntry(p, i1)
        } else {
            n1 := p.children[i1 + 1]
            if (n.keyCnt + n1,keyCnt < order) {
                p.children[i] = merge(n, n1)
            } else {
                // borrow an entry from i1 + 1
                k, c := n1.delete(0)
                n.insertChild(n.keyCnt, k, c)
            }
            deleteEntry(p, i1 + 1)
        }
        
    }
}