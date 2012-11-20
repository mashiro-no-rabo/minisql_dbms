package idxman

import (
	"../catman"
	"../common"
	"os"
	"sort"
)

// Index Manager Definition
type idxMan struct {
	root *node
	typ  int
	indexName int
}

/*
 * There are Five main Funcs:
 * New: create a new index manager
 * Destroy: destroy an index manager
 * Insert: insert a key into index manager
 * Delete: delete a key from index manager
 * Select: select records using index manager
 * 
 * AND one Func for test:
 * Print: print the whole tree to logger
 */

func New(fileName string, tableName string, indexName int) error {
	common.OpLogger.Print("[NewIdxMan]", fileName, ",", tableName, ",", indexName)
	defer common.OpLogger.Print("[leave NewIdxMan]")

	im, err := NewInMemory(fileName, tableName, indexName)
	if err != nil {
		common.ErrLogger.Println("[NewIdxMan]", err)
		return err
	}

	im.FlushToDisk(fileName)
	if err != nil {
		common.ErrLogger.Println("[NewIdxMan]", err)
		return err
	}
	
	return nil
}

func Destroy(fileName string) error {
	common.OpLogger.Print("[DestroyIdxMan] fileName: ", fileName)
	defer common.OpLogger.Print("[leave DestroyIdxMan]")
	
	err := os.Remove(fileName)
	if err != nil {
		common.ErrLogger.Print("[DestroyIdxMan]", err)
		return err
	}
	
	return nil
}

func Insert(fileName string, v common.CellValue, id int64) error {
	common.OpLogger.Print("[Insert]", fileName, ",", v, ",", id)
	defer common.OpLogger.Print("[leave Insert]")
	im, err := ConstructFromDisk(fileName)
	if err != nil {
		common.ErrLogger.Print("[Insert]", err)
		return err
	}
	
	im.Insert(v, id)
	
	im.FlushToDisk(fileName)
	
	return nil
}

func Delete(fileName string, v common.CellValue) (int64, bool, error) {
	common.OpLogger.Print("[Delete]", fileName, ",", v)
	defer common.OpLogger.Print("[leave Delete]")
	
	im, err := ConstructFromDisk(fileName)
	if err != nil {
		common.ErrLogger.Print("[Delete]", err)
		return 0, false, err
	}
	
	id, present := im.Delete(v)
	if present {
		common.OpLogger.Print("[No node deleted.]")
		return 0, false, nil
	}
	
	im.FlushToDisk(fileName)

	return id, true, nil
}

func Select(tableName string, conditions []common.Condition, indexFile string) ([]int64, error) {
	common.OpLogger.Print("[Select]", tableName, ",", conditions, ",", indexFile)
	defer common.OpLogger.Print("[leave Select]")

	table, err := catman.TableInfo(tableName) 
	if err != nil {
		common.ErrLogger.Print("[Select]", err)
		return nil, err
	}
	
	colNameHelper := make(map[string]int)
	for i, col := range table.Columns { 
		colNameHelper[col.Name] = i
	} 

	rangeConds := make([]*rangeCondition, len(table.Columns))
	nonEQConds := make([]*nonEQConditions, len(table.Columns))
	for _, cond := range conditions {
		colId := colNameHelper[cond.ColName]
		if cond.Op == common.OP_NEQ {
			if nonEQConds[colId] == nil {
				*nonEQConds[colId] = make([]*nonEQCondition, 0, len(conditions))
			}
			nonEQCond := new(nonEQCondition)
			*nonEQCond = cond.Value()
			*nonEQConds[colId] = append(*nonEQConds[colId], nonEQCond)
		} else
		if rangeConds[colId] == nil {
			rangeCond := new(rangeCondition)
			switch cond.Op {
				case common.OP_EQ:
					rangeCond.leftOp = CLOSE_INTERVAL
					rangeCond.left = cond.Value()
					rangeCond.rightOp = CLOSE_INTERVAL
					rangeCond.right = cond.Value()
				case common.OP_LT:
					rangeCond.rightOp = OPEN_INTERVAL
					rangeCond.right = cond.Value()
				case common.OP_GT:
					rangeCond.leftOp = OPEN_INTERVAL
					rangeCond.left = cond.Value()
				case common.OP_LEQ:
					rangeCond.rightOp = CLOSE_INTERVAL
					rangeCond.right = cond.Value()
				case common.OP_GEQ:
					rangeCond.leftOp = CLOSE_INTERVAL
					rangeCond.left = cond.Value()
			}
			rangeConds[colId] = rangeCond
		} else {
			rangeCond := rangeConds[colId]
			switch cond.Op {
				case common.OP_EQ:
					if cond.Value().GreaterThan(rangeCond.left) {
						rangeCond.leftOp = CLOSE_INTERVAL
						rangeCond.left = cond.Value()
					}
					if cond.Value().LessThan(rangeCond.right) {
						rangeCond.rightOp = CLOSE_INTERVAL
						rangeCond.right = cond.Value()
					}
				case common.OP_LT:
					if ! cond.Value().GreaterThan(rangeCond.right) {
						rangeCond.rightOp = OPEN_INTERVAL
						rangeCond.right = cond.Value()
					}
				case common.OP_GT:
					if ! cond.Value().LessThan(rangeCond.left) {
						rangeCond.leftOp = OPEN_INTERVAL
						rangeCond.left = cond.Value()
					}
				case common.OP_LEQ:
					if cond.Value().LessThan(rangeCond.right) {
						rangeCond.rightOp = CLOSE_INTERVAL
						rangeCond.right = cond.Value()
					}
				case common.OP_GEQ:
					if cond.Value().GreaterThan(rangeCond.left) {
						rangeCond.leftOp = CLOSE_INTERVAL
						rangeCond.left = cond.Value()
					}
			}
		}
	}

	im, err := ConstructFromDisk(indexFile)
	if err != nil {
		common.ErrLogger.Print("[Select]", err)
		return nil, err
	}

	var resultIds []int64
	firstCond := true
	for i, rangeCond := range rangeConds {
		if rangeCond == nil {
			continue
		}
		
		var tempResult int64Slice
		if i == im.indexName && rangeCond.left != nil {
			tempResult, err = im.SelectRange(*rangeCond, *nonEQConds[i])
			if err != nil {
				common.ErrLogger.Print("[Select]", err)
				return nil, err
			}
		} else {
			tempResult, err = LinearSelectRange(table.Name, *rangeCond, *nonEQConds[i], i)
			if err != nil {
				common.ErrLogger.Print("[Select]", err)
				return nil, err
			}
		}
		
		if firstCond {
			tempResult.Sort()
			resultIds = tempResult
			firstCond = false
		} else {
			for _, recordId := range tempResult {
				i := sort.Search(len(resultIds), func(i int) bool {
					return resultIds[i] >= recordId
				})
				if !(i < len(resultIds) && resultIds[i] == recordId) {
					resultIds = append(resultIds[:i], resultIds[i + 1:]...)
				}
			}
		}
	}
	return resultIds, nil
}

// Print the whole index manager to logger
func (self *idxMan) Print() {
	common.OpLogger.Print("[Print]")
	defer common.OpLogger.Print("[leave Print]")
	
	nodeList := make([]*node, 1)
	nodeList[0] = self.root
	for i := 0; i < len(nodeList); i++ {
		n := nodeList[i]

		common.OpLogger.Printf("[%p]\t%v", n, n)

		if !n.isLeaf() {
			for _, child := range n.children {
				nodeList = append(nodeList, child)
			}
		}
	}
}