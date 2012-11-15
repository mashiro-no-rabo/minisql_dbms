package main

import (
    "reflect"
    "./idxman"
)

type IntVal int

func (val1 IntVal) GetValue() idxman.ValueType {
    return val1
}

func (val1 IntVal) EqualsTo(val2 idxman.ValueType) bool {
        return int(val1) == int(reflect.ValueOf(val2).Int())

}

func (val1 IntVal) LessThan(val2 idxman.ValueType) bool {
        return int(val1) < int(reflect.ValueOf(val2).Int())
}

func (val IntVal) ToString() string {
        return string(val)
}

func main() {
    idxMan := idxman.NewIdxMan()
    for i := 1; i <= 10; i++ {
        idxMan.Insert(IntVal(i))
    }
    idxMan.Print()
}