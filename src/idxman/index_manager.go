package idxman

const order = 4 // We are implementing B+ Tree for n = 4.

// Required by Find
type DomainValue interface {
    IsLessThan(DomainValue) (bool)
    IsEqual(DomainValue) (bool)
}

// Required by Insert
type Key interface {
    GetValue() DomainValue
    IsLessThan(Key) (bool)
    IsEqual(Key) (bool)
} 

type node struct {
    isLeaf bool
    parent *node
    keyCnt int
    keys [] Key
    children [] *node
}

type idxMan struct {
    root *node
}

func NewIdxMan() *idxMan {
    im := new(idxMan)
    im.root = new(node)
    return im
}

func (self idxMan) Find(v DomainValue) (Key, bool) {
    return self.root.search(v)
}

func (self *idxMan) Insert(k Key) {
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

func (self *idxMan) Delete(k Key) {}

func (self node) findValueIndex(v DomainValue) (i int, found bool) {
    for i = 0; i < self.keyCnt; i++ {
        temp := self.keys[i]
        if ! temp.GetValue().IsLessThan(v) {
            return i, true
        }
    }
    return self.keyCnt, false
}

func (self node) search(v DomainValue) (Key, bool) {
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

func (self node) findLeafNode(v DomainValue) (l *node) {
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

func insertInLeaf(l *node, k Key) {
    i, _ := l.findValueIndex(k.GetValue())
    copy(l.keys[i + 1:l.keyCnt + 1], l.keys[i:l.keyCnt])
    l.keys[i] = k
    l.keyCnt++
}

func splitNode(l *node) (l1 *node, k1 Key) {
    // Create node l1
    l1 = new(node)
    l1.isLeaf = true
    l1.parent = l.parent
    l1.keyCnt = order / 2 + order % 2
    // Copy keys from l to l1
    copy(l1.keys, l.keys[0:order / 2])
    // Update l
    l.keyCnt = order / 2
    l1.children[order - 1] = l.children[order - 1]
    l.children[order - 1] = l1
    // return l1 and k1
    return l1, l1.keys[0]
}

func insertInParent(l *node, k1 Key, l1 *node) (r *node, newRoot bool) {
    // l is the root of the tree
    if l.parent == nil {
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
