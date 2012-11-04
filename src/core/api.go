package core

import (
	"../catman"
	"../common"
	"errors"
	"fmt"
	"os"
)

func CreateTable(table common.Table) error {
	exist_tables, err := ListTables()
	for t := range exist_tables {
		if t.Name == table.Name {
			return errors.New("Table name has been used.")
		}
	}
	// now create and save in catalog manager (or the catman just read dirs?)
	err = os.MkdirAll("data/"+table.Name, 0600)
	if err != nil {
		return err
	}
}

func DropTable(table_name string) error {

}

func CreateIndex(table_name string, index_name string) error {

}

func DropIndex() {

}

func Insert() {
	// insert to end
	// update index
}

func Select() {

}

func Delete() {
	// if has index then search
	// delete
}

func ListTables() ([]common.Table, error) {

}

// any other catalog manager api??
