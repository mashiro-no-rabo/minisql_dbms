package core

import (
	"../catman"
	"../common"
	// "../idxman"
	"../recman"
	"encoding/json"
	"errors"
	"os"
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

func CreateTable(table common.Table) error {
	common.OpLogger.Printf("[Begin]Creating table: %v\n", table.Name, table)
	// defer a clean-up func

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

func DropIndex(table_name string, index_name string) error {
	return nil
}

func Insert(table_name string, rec common.Record) error {
	// insert to slot [planned]
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

	tabinfo, err := catman.TableInfo("test")
	offset, err := recman.Insert(dbf, tabinfo, rec)

	if err != nil {
		common.ErrLogger.Printf("recman.Insert error: %s, %v, %s", table_name, rec, err)
		return err
	}

	common.OpLogger.Println(offset)
	// tidxs, err := catman.TableIndexes(table_name)
	// for _, idxp := range tidxs {
	// 	tt := strings.Split(idxp, "_")
	// 	idxman.Insert(tab_dir+"/index/"+idxp, rec.Values[tt[1]], offset)
	// }

	common.OpLogger.Printf("[Done]Inserting record %v into %s\n", rec, table_name)
	return nil
}

func Select(table_name string, fields []string, conds []common.Condition) error {
	// how to implement0.0?
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
			err := recman.DeleteAll(dbf)
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
