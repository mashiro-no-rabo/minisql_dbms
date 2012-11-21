package idxman

import (
	"../catman"
	"../common"
	"../recman"
	"errors"
	"os"
)

func LinearSelectRange(tableName string, rangeConds []*rangeCondition, nonEQConds []*nonEQConditions) ([]int64, error) {
	common.OpLogger.Print("LinearSelectRange()")
	file, err := os.OpenFile(common.DataDir+"/"+tableName+"/data.dbf", os.O_RDONLY, 0600)
	if err != nil {
		common.OpLogger.Print("leave LinearSelectRange() with error")
		common.ErrLogger.Print("[LinearSelectRange]", err)
		return nil, err
	}
	defer file.Close()
	tabinfo, err := catman.TableInfo(tableName)
	if err != nil {
		common.OpLogger.Print("leave LinearSelectRange() with error")
		common.ErrLogger.Print("[LinearSelectRange]", err)
		return nil, err
	}
	resultIds := make([]int64, 0)
	records, recordIds := recman.ReadRecords(file, tabinfo)
	for i, record := range records {
		if testAllConditions(record, rangeConds, nonEQConds) {
			resultIds = append(resultIds, recordIds[i])
		}
	}
	common.OpLogger.Print("leave LinearSelectRange()")
	return resultIds, nil
}

func (self idxMan) SelectRange(rangeCond rangeCondition, nonEQConds nonEQConditions) ([]int64, error) {
	common.OpLogger.Print("SelectRange():\t", rangeCond)
	if rangeCond.left == nil {
		common.OpLogger.Print("leave SelectRange() with error")
		err := errors.New("B+ Tree should not handle range condition without left value.")
		common.ErrLogger.Print("[SelectRange]", err)
		return nil, err
	}
	l := self.root.findLeafNode(rangeCond.left)
	i, found := l.findKeyIndex(rangeCond.left)
	if !found {
		common.OpLogger.Print("leave SelectRange()")
		return nil, nil
	}
	if rangeCond.leftOp == OPEN_INTERVAL && l.keys[i].EqualsTo(rangeCond.left) {
		i = l.getNextKey(i)
	}
	resultIds := make([]int64, 0)
	for l != nil && rangeCond.containsCell(l.keys[i]) {
		if (nonEQConds == nil || nonEQConds.dontContainCell(l.keys[i])) {
			resultIds = append(resultIds, l.recordIds[i])
			i = l.getNextKey(i)
		}
	}

	common.OpLogger.Print("leave SelectRange()")
	return resultIds, nil
}

func (l *node) getNextKey(i int) int {
	i++
	if i == l.keyCnt() {
		l = l.children[0]
		i = 0
	}
	return i
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
	common.OpLogger.Print("findLeafNoe()\t", v)
	x := self
	for !x.isLeaf() {
		i := x.findChildIndex(v)
		x = x.children[i]
	}
	common.OpLogger.Print("leave findLeafNode()\t", x)
	return x
}
