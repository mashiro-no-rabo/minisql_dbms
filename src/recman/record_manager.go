package recman

import (
	"../common"
	"bufio"
	"encoding/binary"
	"io"
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

func ReadRecords(dbf *os.File, tab *common.Table) []common.Record {
	// r := bufio.NewReader(dbf)
	var del uint8
	var valsSize int64
	valsSize = 0
	var keys []string
	for k, c := range tab.Columns {
		keys = append(keys, k)
		switch c.Type {
		case common.IntCol:
			valsSize += 8
		case common.StrCol:
			valsSize += c.Length
		case common.FltCol:
			valsSize += 8
		}
	}
	sort.Strings(keys)
	common.OpLogger.Println(valsSize)

	var recs []common.Record
	var intval common.IntVal
	var fltval common.FltVal
	for {
		if err := binary.Read(dbf, binary.LittleEndian, &del); err == io.EOF {
			break
		}
		if del == 1 {
			rec := new(common.Record)
			vals := make(map[string]common.CellValue)
			for _, k := range keys {
				switch tab.Columns[k].Type {
				case common.IntCol:
					binary.Read(dbf, binary.LittleEndian, &intval)
					vals[k] = intval
				case common.FltCol:
					binary.Read(dbf, binary.LittleEndian, &fltval)
					vals[k] = fltval
				case common.StrCol:
					raw_bytes := make([]byte, tab.Columns[k].Length)
					dbf.Read(raw_bytes)
					vals[k] = common.StrVal(raw_bytes)
				}
			}
			rec.Values = vals
			rec.Del = false
			recs = append(recs, *rec)
		} else {
			dbf.Seek(valsSize, os.SEEK_CUR)
		}
	}
	return recs
}
