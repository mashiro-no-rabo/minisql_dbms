package core

import (
	"../catman"
	"../common"
	"errors"
	"os"
)

func CreateTable(table common.Table) error {
	common.OpLogger.Printf("Creating table: %s\n", table.Name)
	// defer a clean-up func

	exist_tables, err := catman.AllTables()
	for _, name := range exist_tables {
		if table.Name == name {
			common.OpLogger.Printf("Cancel creating table: %s, conflict name.\n", table.Name)
			return errors.New("Conflict with existing table.")
		}
	}

	common.OpLogger.Printf("Creating folder structure for %s...\n", table.Name)
	tab_dir := common.DataDir + "/" + table.Name
	err = os.MkdirAll(tab_dir, 0700)
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

	common.OpLogger.Printf("Creating record file for %s...\n", table.Name)
	fd, err := os.OpenFile(tab_dir+"/data.dbf", os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		common.ErrLogger.Printf("Cannot create record file for %s, due to %s", table.Name, err)
		return err
	}
	defer fd.Close()

	return nil
}

func DropTable(table_name string) error {
	return nil
}

func CreateIndex(table_name string, index_name string, index_key string) error {
	return nil
}

func DropIndex() error {
	return nil
}

func Insert() error {
	// insert to end of db files
	// update index
	return nil
}

func Select() error {
	return nil
}

func Delete() error {
	// if has index then search
	// delete
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
