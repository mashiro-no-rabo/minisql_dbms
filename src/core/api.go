package core

import (
	"../catman"
	"../common"
	// "../idxman"
	"../recman"
	"encoding/binary"
	"encoding/json"
	"errors"
	"os"
	"sort"
	// "strings"
)

func checkExist(table_name string) bool {
	exist_tables, err := catman.AllTables()
	if err != nil {
		return false
	}
	found := false
	for _, name := range exist_tables {
		if table_name == name {
			found = true
			break
		}
	}
	return found
}

func CreateTable(table *common.Table) error {
	common.OpLogger.Printf("[Begin]Creating table: %v\n", table.Name, table)

	if checkExist(table.Name) {
		common.OpLogger.Printf("[Cancel]Creating table: %s, conflict table name.\n", table.Name)
		return errors.New("Conflict with existing table name.")
	}

	common.OpLogger.Printf("Creating folder structure for %s...\n", table.Name)
	tab_dir := common.DataDir + "/" + table.Name
	err := os.MkdirAll(tab_dir, 0700)
	if err != nil {
		common.ErrLogger.Printf("Cannot create table dir for %s, due to %s", table.Name, err)
		return err
	}

	common.OpLogger.Printf("Saving schema for %s...\n", table.Name)
	fs, err := os.OpenFile(tab_dir+"/schema.dbf", os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		common.ErrLogger.Printf("Cannot create schema file for %s, due to %s", table.Name, err)
		return err
	}
	defer fs.Close()
	enc := json.NewEncoder(fs)
	enc.Encode(table)

	common.OpLogger.Printf("Creating record file for %s\n", table.Name)
	fd, err := os.OpenFile(tab_dir+"/data.dbf", os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		common.ErrLogger.Printf("Cannot create record file for %s, due to %s", table.Name, err)
		return err
	}
	defer fd.Close()

	// create index?

	common.OpLogger.Printf("[Done]Created table: %s\n", table.Name)
	return nil
}

func DropTable(table_name string) error {
	common.OpLogger.Printf("[Begin]Dropping table: %s\n", table_name)

	if !checkExist(table_name) {
		common.OpLogger.Printf("[Cancel]Dropping table, table %s not exist.\n", table_name)
		return errors.New("Can't find target table")
	}

	tab_dir := common.DataDir + "/" + table_name
	err := os.RemoveAll(tab_dir)
	if err != nil {
		common.ErrLogger.Printf("Cannot delete folder of %s, due to %s", table_name, err)
		return err
	}
	common.OpLogger.Printf("[Done]Dropped table: %s\n", table_name)
	return nil
}

func CreateIndex(table_name string, index_name string, index_key string) error {
	return nil
}

func DropIndex(index_name string) error {
	return nil
}

func Insert(table_name string, vals []common.CellValue) error {
	common.OpLogger.Printf("[Begin]Inserting record %v into %s\n", rec, table_name)

	if !checkExist(table_name) {
		common.OpLogger.Printf("[Cancel]Inserting record, table %s not exist.\n", table_name)
		return errors.New("Can't find target table")
	}

	dbf, err := os.OpenFile(common.DataDir+"/"+table_name+"/data.dbf", os.O_APPEND|os.O_WRONLY|os.O_SYNC, 0600)
	if err != nil {
		return err
	}
	defer dbf.Close()

	tabinfo, err := catman.TableInfo(table_name)
	var keys []string
	for k, _ := range tabinfo.Columns {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	mapvals := make(map[string]CellValue)
	for i, k := range keys {
		mapvals[k] = CellValue[i]
	}
	rec := common.Record{Del: false, Values: mapvals}

	offset, err := recman.Insert(dbf, tabinfo, rec)

	if err != nil {
		common.ErrLogger.Printf("recman.Insert error: %s, %v, %s", table_name, rec, err)
		return err
	}

	common.OpLogger.Println(offset)
	tidxs, err := catman.TableIndexes(table_name)
	for _, idxp := range tidxs {
		tt := strings.Split(idxp, "_")
		idxman.Insert(tab_dir+"/index/"+idxp, rec.Values[tt[1]], offset)
	}

	common.OpLogger.Printf("[Done]Inserting record %v into %s\n", rec, table_name)
	return nil
}

func Select(table_name string, conds []common.Condition) error {

	return nil
}

func Delete(table_name string, conds []common.Condition) error {
	if checkExist(table_name) {
		dbf, err := os.OpenFile(common.DataDir+"/"+table_name+"/data.dbf", os.O_RDWR|os.O_SYNC, 0600)

		if err != nil {
			common.ErrLogger.Println(err)
			return err
		}
		defer dbf.Close()
		if len(conds) == 0 {
			common.OpLogger.Println("DeleteAll")
			tabinfo, err := catman.TableInfo(table_name)
			err = recman.DeleteAll(dbf, tabinfo)
			if err != nil {
				return err
			}
		} else {
			// offsets := idxman.Search(table_name, conds)
			// err := recman.Delete(dbf, offsets)
			// if err != nil {
			// 	return err
			// }
		}
	} else {
		common.OpLogger.Printf("[Cancel]Table %s not exist", table_name)
		return errors.New("Table not exist.")
	}

	// and update index
	return nil
}

func SelectOffsets(table_name string, offsets []int64) []common.Record {
	tab, err := catman.TableInfo(table_name)
	var keys []string
	for k, _ := range tab.Columns {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	dbf, err := os.OpenFile(common.DataDir+"/"+table_name+"/data.dbf", os.O_RDONLY, 0600)
	if err != nil {
		return nil
	}
	var del uint8
	var intval common.IntVal
	var fltval common.FltVal
	var result []common.Record
	for _, ofst := range offsets {
		dbf.Seek(ofst, os.SEEK_SET)
		rec := new(common.Record)
		vals := make(map[string]common.CellValue)
		binary.Read(dbf, binary.LittleEndian, &del)
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
		result = append(result, *rec)
	}
	return result
}
