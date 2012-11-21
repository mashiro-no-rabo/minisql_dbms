package idxman

import (
	"../common"
	"sort"
)

type rangeCondition struct {
	leftOp  int
	left    common.CellValue
	rightOp int
	right   common.CellValue
}

const (
	CLOSE_INTERVAL = iota
	OPEN_INTERVAL
)

func (rangeCond rangeCondition) containsRecord(record common.Record, colId int) bool {
	return rangeCond.containsCell(record.Values[colId])
}
func (rangeCond rangeCondition) containsCell(cell common.CellValue) bool {
	if rangeCond.left != nil && rangeCond.leftOp == CLOSE_INTERVAL && cell.LessThan(rangeCond.left) {
		return false
	}
	if rangeCond.left != nil && rangeCond.leftOp == OPEN_INTERVAL && rangeCond.left.GreaterThan(cell) {
		return false
	}
	if rangeCond.right != nil && rangeCond.rightOp == CLOSE_INTERVAL && cell.GreaterThan(rangeCond.right) {
		return false
	}
	if rangeCond.right != nil && rangeCond.rightOp == OPEN_INTERVAL && rangeCond.right.LessThan(cell) {
		return false
	}
	return true
}

type nonEQCondition common.CellValue

type nonEQConditions []*nonEQCondition

func (nonEQConds nonEQConditions) dontContainRecord(record common.Record, colId int) bool {
	return nonEQConds.dontContainCell(record.Values[colId])
}

func (nonEQConds nonEQConditions) dontContainCell(cell common.CellValue) bool {
	for _, nonEQCond := range nonEQConds {
		if (*nonEQCond).EqualsTo(cell) {
			return false
		}
	}
	return true
}

func testAllConditions(record common.Record, rangeConds []*rangeCondition, nonEQConds []*nonEQConditions) bool {
	for i, rangeCond := range rangeConds {
		if (rangeCond != nil && !rangeCond.containsRecord(record, i)) && 
			(nonEQConds != nil && !nonEQConds[i].dontContainRecord(record, i)) {
				return false
			}
	}
	return true
}

type int64Slice []int64

func (p int64Slice) Len() int           { return len(p) }
func (p int64Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p int64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p int64Slice) Sort()              { sort.Sort(p) }
