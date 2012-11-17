package recman

import (
	"../common"
	"bufio"
	"encoding/binary"
	"io"
	"os"
	"reflect"
)

func Insert(dbf *os.File, tab *common.Table, vals []common.CellValue) (int64, error) {
	dbs, err := dbf.Stat()
	if err != nil {
		return -1, err
	}
	offset := dbs.Size()
	w := bufio.NewWriter(dbf)
	binary.Write(w, binary.LittleEndian, uint8(1))

	for i, val := range vals {
		switch tab.Columns[i].Type {
		case common.IntCol:
			binary.Write(w, binary.LittleEndian, reflect.ValueOf(val.Value()).Int())
		case common.StrCol:
			w.Write([]byte(reflect.ValueOf(val.Value()).String()))
		case common.FltCol:
			binary.Write(w, binary.LittleEndian, reflect.ValueOf(val.Value()).Float())
		}
	}
	err = w.Flush()
	if err != nil {
		return -1, nil
	}
	return offset, nil
}

func DeleteAll(dbf *os.File, tab *common.Table) error {
	var valsSize int64
	valsSize = 0
	for _, c := range tab.Columns {
		switch c.Type {
		case common.IntCol:
			valsSize += 8
		case common.StrCol:
			valsSize += c.Length
		case common.FltCol:
			valsSize += 8
		}
	}
	if stat, _ := dbf.Stat(); stat.Size() == 0 {
		return nil
	}
	for {
		binary.Write(dbf, binary.LittleEndian, uint8(0))
		if _, err := dbf.Seek(valsSize+1, os.SEEK_CUR); err == io.EOF {
			break
		}
	}

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

func ReadRecords(dbf *os.File, tab *common.Table) ([]common.Record, []int64) {
	var del uint8
	var valsSize int64
	valsSize = 0
	for _, c := range tab.Columns {
		switch c.Type {
		case common.IntCol:
			valsSize += 8
		case common.StrCol:
			valsSize += c.Length
		case common.FltCol:
			valsSize += 8
		}
	}

	var recs []common.Record
	var offsets []int64
	var now_ofst int64
	now_ofst = 0
	var intval common.IntVal
	var fltval common.FltVal
	for {
		if err := binary.Read(dbf, binary.LittleEndian, &del); err == io.EOF {
			break
		}
		if del == 1 {
			rec := new(common.Record)
			var vals []common.CellValue
			for _, col := range tab.Columns {
				switch col.Type {
				case common.IntCol:
					binary.Read(dbf, binary.LittleEndian, &intval)
					vals = append(vals, intval)
				case common.FltCol:
					binary.Read(dbf, binary.LittleEndian, &fltval)
					vals = append(vals, fltval)
				case common.StrCol:
					raw_bytes := make([]byte, col.Length)
					dbf.Read(raw_bytes)
					vals = append(vals, common.StrVal(raw_bytes))
				}
			}
			rec.Values = vals
			rec.Del = false
			recs = append(recs, *rec)
			offsets = append(offsets, now_ofst)
		} else {
			dbf.Seek(valsSize, os.SEEK_CUR)
		}
		now_ofst += valsSize + 1
	}
	return recs, offsets
}
