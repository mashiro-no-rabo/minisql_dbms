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
		return err
	}

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

// any other catalog manager api??
