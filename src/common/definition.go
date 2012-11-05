package common

import (
	"log"
	"os"
)

type Comparable interface {
	EqualsTo(Comparable) bool
	LessThan(Comparable) bool
}

// move if only for index manager
// type KeyType interface {
// 	GetValue() ValueType
// 	Comparable
// }

// type ValueType interface {
// 	Comparable
// }

const (
	ColInt = iota
	ColString
	ColFloat

	DataDir string = "data"
)

var (
	OpLogger  *log.Logger
	ErrLogger *log.Logger
)

func init() {
	os.MkdirAll("logs", 0700)
	op_log_file, err := os.OpenFile("logs/ops.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0600)
	if err != nil {
		panic(err)
	}
	OpLogger = log.New(op_log_file, "- ", log.Ldate|log.Ltime)

	err_log_file, err := os.OpenFile("logs/errors.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0600)
	if err != nil {
		panic(err)
	}
	ErrLogger = log.New(err_log_file, "!!> ", log.Ldate|log.Ltime|log.Lshortfile)
}

type CellValue interface {
	Comparable
	ToString() string
}

type Table struct {
	Name    string
	Columns map[string]Column
	PKey    string
}

type Column struct {
	Type   int
	Unique bool
	Length int
}
