package recman

import (
	"../common"
	"bufio"
	"encoding/binary"
	"os"
	"reflect"
	"sort"
)

func Insert(dbf *os.File, tab *common.Table, rec common.Record) (int64, error) {
	dbs, err := dbf.Stat()
	if err != nil {
		return -1, err
	}
	offset := dbs.Size()
	// enc := json.NewEncoder(dbf)
	// enc.Encode(rec)
	w := bufio.NewWriter(dbf)
	binary.Write(w, binary.LittleEndian, uint8(1))

	var keys []string
	for k, _ := range tab.Columns {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		switch tab.Columns[k].Type {
		case common.IntCol:
			binary.Write(w, binary.LittleEndian, reflect.ValueOf(rec.Values[k].Value()).Int())
		case common.StrCol:
			w.Write([]byte(reflect.ValueOf(rec.Values[k].Value()).String()))
		case common.FltCol:
			binary.Write(w, binary.LittleEndian, reflect.ValueOf(rec.Values[k].Value()).Float())
		}
	}
	err = w.Flush()
	if err != nil {
		return -1, nil
	}
	return offset, nil
}

func DeleteAll(dbf *os.File) error {
	return nil
}

func Delete(dbf *os.File, tab *common.Table, offsets []int64) error {
	for num, ofst := range offsets {
		common.OpLogger.Printf("Deleting no.%d record of table %s at offset %d", num, tab.Name, ofst)
		dbf.Seek(ofst, os.SEEK_SET)
		binary.Write(dbf, binary.LittleEndian, uint8(0))
	}

	return nil
}
