package core

import (
	"../catman"
	"../common"
	"../idxman"
	"../recman"
	"encoding/binary"
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"strings"
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
	err = os.MkdirAll(tab_dir+"/index", 0700)
	if err != nil {
		common.ErrLogger.Printf("Cannot create index dir for %s, due to %s", table.Name, err)
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

	err = CreateIndex(table.Name, "default", table.PKey)

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

func CreateIndex(table_name string, index_name string, index_key int) error {
	filename := common.DataDir + "/" + table_name + "/index/" + index_name + "_" + strconv.Itoa(index_key) + ".idf"
	return idxman.NewIdxMan(filename, table_name, index_key)
}

func DropIndex(table_name string, index_name string) error {
	tidxs, err := catman.TableIndexes(table_name)
	if err != nil {
		return err
	}
	filename := ""
	for _, idxp := range tidxs {
		if index_name == strings.Split(idxp, "_")[0] {
			filename = common.DataDir + "/" + table_name + "/index/" + idxp
		}
	}
	if len(filename) == 0 {
		return errors.New("No index found")
	}
	return idxman.DestroyIdxMan(filename)
}

func Insert(table_name string, vals []common.CellValue) error {
	common.OpLogger.Printf("[Begin]Inserting record %v into %s\n", vals, table_name)

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
	offset, err := recman.Insert(dbf, tabinfo, vals)

	if err != nil {
		common.ErrLogger.Printf("recman.Insert error: %s, %v, %s", table_name, vals, err)
		return err
	}
	common.OpLogger.Println(offset)
	// tidxs, err := catman.TableIndexes(table_name)
	// tab_dir := common.DataDir + "/" + table_name
	// for _, idxp := range tidxs {
	// 	key, _ := strconv.Atoi(strings.Split(idxp, "_")[1])
	// 	idxman.Insert(tab_dir+"/index/"+idxp, vals[key], offset)
	// }

	common.OpLogger.Printf("[Done]Inserting record %v into %s\n", vals, table_name)
	return nil
}

func Select(table_name string, conds []common.Condition) ([]common.Record, error) {
	var recs []common.Record
	if checkExist(table_name) {
		dbf, err := os.OpenFile(common.DataDir+"/"+table_name+"/data.dbf", os.O_RDONLY|os.O_SYNC, 0600)
		if err != nil {
			common.ErrLogger.Println(err)
			return nil, err
		}
		defer dbf.Close()

		if len(conds) == 0 {
			common.OpLogger.Println("SelectAll")
			tabinfo, err := catman.TableInfo(table_name)
			if err != nil {
				return nil, err
			}
			recs, _ = recman.ReadRecords(dbf, tabinfo)
		} else {
			// offsets := idxman.Search(table, conds)
			// recs := idxman.SelectOffsets(table_name, offsets)
		}
	}
	return recs, nil
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
			if err != nil {
				return err
			}
			recs, offsets := recman.ReadRecords(dbf, tabinfo)
			recman.Delete(dbf, tabinfo, offsets)
			// tidxs, err := catman.TableIndexes(table_name)
			// tab_dir := common.DataDir + "/" + table_name

			common.OpLogger.Println(recs)
			// for _, idxp := range tidxs {
			// 	key, _ := strconv.Atoi(strings.Split(idxp, "_")[1])
			// 	for _, r := range recs {
			// 		idxman.Delete(tab_dir+"/index/"+idxp, r.Values[key])
			// 	}
			// }
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
		var vals []common.CellValue
		binary.Read(dbf, binary.LittleEndian, &del)
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
		result = append(result, *rec)
	}
	return result
}
