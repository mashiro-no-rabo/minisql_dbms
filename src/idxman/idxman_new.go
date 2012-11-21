package idxman

import (
	"../catman"
	"../common"
	"../recman"
	"os"
	"strconv"
	"strings"
)

func NewInMemory(fileName string, tableName string, indexName int) (*idxMan, error) {
	common.OpLogger.Print("[NewInMemory]", fileName, ",", tableName, ",", indexName)
	defer common.OpLogger.Print("[leave NewInMemory]")

	idxs, err := catman.TableIndexes(tableName)
	if err != nil && searchString(idxs, indexName) {
		common.ErrLogger.Print("[NewInMemory]", err)
		return nil, err
	}

	file, err := os.OpenFile(common.DataDir+"/"+tableName+"/data.dbf", os.O_RDONLY, 0600)
	if err != nil {
		common.ErrLogger.Print("[NewInMemory]", err)
		return nil, err
	}
	defer file.Close()

	tabinfo, err := catman.TableInfo(tableName)
	if err != nil {
		common.ErrLogger.Print("[NewInMemory]", err)
		return nil, err
	}
	records, recordIds := recman.ReadRecords(file, tabinfo)

	im := NewEmpty(tabinfo.Columns[indexName].Type, indexName)
	for i, record := range records {
		im.Insert(record.Values[indexName], recordIds[i])
	}

	return im, nil
}

func searchString(s []string, x int) bool {
	for _, y := range s {
		if testx, _ := strconv.Atoi(strings.Split(y, "_")[1]); x == testx {
			return true
		}
	}
	return false
}

// Creat an index manager and give it a empty root.
func NewEmpty(idx_typ int, indexName int) *idxMan {
	common.OpLogger.Print("[NewEmpty]", idx_typ, ",", indexName)
	defer common.OpLogger.Print("[leave NewEmpty]")

	im := new(idxMan)
	im.root = createNode(true)
	im.typ = idx_typ
	im.indexName = indexName

	return im
}
