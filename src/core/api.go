package core

import (
	"../catman"
	"../common"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
)

func CreateTable(table common.Table) error {
	common.OpLogger.Printf("Creating table: %s\n", table.Name)
	exist_tables, err := catman.AllTables()
	for _, name := range exist_tables {
		if table.Name == name {
			return errors.New("Conflict with existing table.")
		}
	}
	// should use catman funcs
	// exist_tables, err := ListTables()
	// for t := range exist_tables {
	// 	if t.Name == table.Name {
	// 		return errors.New("Table name has been used.")
	// 	}
	// }
	// now create and save in catalog manager (or the catman just read dirs?)
	// need logging
	err = os.MkdirAll(common.DataDir+"/"+table.Name, 0700)
	if err != nil {
		common.ErrLogger.Println(err)
		return err
	}
	f, err := os.OpenFile("data.dbf", os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		common.ErrLogger.Println(err)
		return err
	}
	defer f.Close()

}

func DropTable(table_name string) error {

}

func CreateIndex(table_name string, index_name string, index_key string) error {

}

func DropIndex() {

}

func Insert() {
	// insert to end of db files
	// update index
}

func Select() {

}

func Delete() {
	// if has index then search
	// delete
}

func ListTables() ([]string, error) {

}

func ListIndex(table_name string) error {

}

func TableInfo(table_name string) common.Table {

}

// any other catalog manager api??
