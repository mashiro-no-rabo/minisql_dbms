package main

import (
	"./common"
	"./idxman"
	"fmt"
)

type IntKey struct {
	i common.IntVal
}

func (key IntKey) GetValue() common.CellValue {
	return key.i
}

func main() {
	idxMan := idxman.NewIdxMan()
	for i := 10; i >= 1; i-- {
		idxMan.Insert(IntKey{common.IntVal(i)})
		idxMan.Print()
	}
	k, found := idxMan.Find(common.IntVal(6))
	if found {
		fmt.Println(k.GetValue())
	} else {
		fmt.Println("Not found: ", 6)
	}
	k, found = idxMan.Delete(common.IntVal(6))
	if found {
		fmt.Println("Delete Succeed!")
	} else {
		fmt.Println("Delete failed!")
	}
	idxMan.Print()
}
