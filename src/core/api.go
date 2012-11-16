package core

import (
	"../catman"
	"../common"
	"../recman"
	"encoding/json"
	"errors"
	"os"
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

func DropIndex() error {
	return nil
}

func Insert(table_name string, rec common.Record) error {
	// insert to slot [planned]
	common.OpLogger.Printf("[Begin]Inserting record %v into %s\n", rec, table_name)

	if !checkExist(table_name) {
		common.OpLogger.Printf("[Cancel]Inserting record, table %s not exist.\n", table_name)
		return errors.New("Can't find target table")
	}

	tab_dir := common.DataDir + "/" + table_name
	dbf, err := os.OpenFile(tab_dir+"/data.dbf", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		common.ErrLogger.Printf("Cannot open record file for %s, due to %s", table_name, err)
		return err
	}
	defer dbf.Close()

	recman.Insert(dbf, rec)
	// update index

	common.OpLogger.Printf("[Done]Inserting record %v into %s\n", rec, table_name)
	return nil
}

func Select() error {
	// how to implement0.0?
	return nil
}

func Delete(table_name string, pkvals ...common.CellValue) error {
	// if has index then search
	// else raw delete 
	// and update index
	return nil
}

func ListTables() ([]string, error) {
	return catman.AllTables()
}

func ListIndex(table_name string) error {
	return nil
}

// func TableInfo(table_name string) common.Table {

// }

// any other catalog manager api??
