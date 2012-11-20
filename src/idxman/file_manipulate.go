package idxman

import (
	"os"
	"fmt"
	"io"
	"../common"
)

// Disk file format(Each record corresponding to one line):
// no, pno, leaf, keyCnt, keys..., recordIds...
func (self *idxMan) FlushToDisk(fileName string) error {
	common.OpLogger.Print("[FlushToDisk] fileName: ", fileName)
	defer common.OpLogger.Print("[leave FlushToDisk]")
	
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		common.ErrLogger.Print("[FlushToDisk]", err)
		return err
	}
	defer file.Close()

	fmt.Fprint(file, self.typ, " ", self.indexName, " ")

	nodeHelper := make(map[*node]int64)
	queue := make([]*node, 1)
	queue[0] = self.root
	var n *node
	var no int64  // offset of node
	var pno int64 // offset of node's parent
	var ok bool
	no = 0
	for i := 0; i < len(queue); i++ {
		n = queue[i]
		nodeHelper[n] = no
		pno, ok = nodeHelper[n.parent]
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
		} else {
			for _, child := range n.children {
				queue = append(queue, child)
			}
		}
	}
	return nil
}

func ConstructFromDisk(fileName string) (*idxMan, error) {
	common.OpLogger.Print("[ConstructFromDisk]", fileName)
	common.OpLogger.Print("[leave ConstructFromDisk]")
	
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0600)
	if err != nil {
		common.ErrLogger.Print("[ConstructFromDisk]", err)
		return nil, err
	}
	defer file.Close()

	im := new(idxMan)
	fmt.Fscan(file, &im.typ, &im.indexName)

	var n *node
	var p *node
	var no int64
	var pno int64
	var leaf bool
	var keyCnt int
	nodeHelper := make(map[int64]*node)
	fmt.Fscan(file, &no, &pno, &leaf, &keyCnt)
	im.root = createNode(leaf)
	nodeHelper[no] = im.root
	im.root.constructKeys(file, keyCnt, im.typ)
	for {
		_, err := fmt.Fscan(file, &no, &pno, &leaf, &keyCnt)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		n = createNode(leaf)
		p, _ = nodeHelper[pno]
		nodeHelper[no] = n
		p.children = append(p.children, n)
		n.constructKeys(file, keyCnt, im.typ)
	}
	return im, nil
}

func (self *node) constructKeys(file *os.File, keyCnt int, typ int) {
	var key common.CellValue
	var intval common.IntVal
	var fltval common.FltVal
	var strval common.StrVal
	for i := 0; i < keyCnt; i++ {
		switch typ {
		case common.IntCol:
			fmt.Fscan(file, &intval)
			key = intval
		case common.FltCol:
			fmt.Fscan(file, &fltval)
			key = fltval
		case common.StrCol:
			fmt.Fscan(file, &strval)
			key = strval
		}
		self.keys = append(self.keys, key)
	}
	var recordId int64
	if self.isLeaf() {
		for i := 0; i < keyCnt; i++ {
			fmt.Fscan(file, &recordId)
			self.recordIds = append(self.recordIds, recordId)
		}
	}
}