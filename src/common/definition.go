package common

import (
	"fmt"
	"log"
	"os"
	"reflect"
)

// FAKE CONDITION
type Condition struct {
	nothing int
}

// the Talbe struct
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

// the Record struct
const (
	IntCol = iota
	StrCol
	FltCol

	DataDir string = "data"
)

type Comparable interface {
	EqualsTo(Comparable) bool
	LessThan(Comparable) bool
	GreaterThan(Comparable) bool
}

type CellValue interface {
	Comparable
	String() string
	Value() interface{}
}

type Record struct {
	Del    bool
	Values map[string]CellValue
}

// the Int type
type IntVal int

func (val1 IntVal) EqualsTo(val2 Comparable) bool {
	return int(val1) == int(reflect.ValueOf(val2).Int())
}

func (val1 IntVal) LessThan(val2 Comparable) bool {
	return int(val1) < int(reflect.ValueOf(val2).Int())
}

func (val1 IntVal) GreaterThan(val2 Comparable) bool {
	return int(val1) > int(reflect.ValueOf(val2).Int())
}

func (val IntVal) String() string {
	return fmt.Sprintf("%d", int(val))
}

func (val IntVal) Value() interface{} {
	return val
}

// the Char(n) type
type StrVal string

func (val1 StrVal) EqualsTo(val2 Comparable) bool {
	return string(val1) == reflect.ValueOf(val2).String()
}

func (val1 StrVal) LessThan(val2 Comparable) bool {
	return string(val1) < reflect.ValueOf(val2).String()
}

func (val1 StrVal) GreaterThan(val2 Comparable) bool {
	return string(val1) > reflect.ValueOf(val2).String()
}

func (val StrVal) String() string {
	return string(val)
}

func (val StrVal) Value() interface{} {
	return val
}

// the Float type
type FltVal float64

func (val1 FltVal) EqualsTo(val2 Comparable) bool {
	return float64(val1) == float64(reflect.ValueOf(val2).Float())
}

func (val1 FltVal) LessThan(val2 Comparable) bool {
	return float64(val1) < float64(reflect.ValueOf(val2).Float())
}

func (val1 FltVal) GreaterThan(val2 Comparable) bool {
	return float64(val1) > float64(reflect.ValueOf(val2).Float())
}

func (val FltVal) String() string {
	return fmt.Sprintf("%.2f", float64(val))
}

func (val FltVal) Value() interface{} {
	return val
}

// Logger
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
	OpLogger = log.New(op_log_file, "", log.Ldate|log.Ltime)

	err_log_file, err := os.OpenFile("logs/errors.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0600)
	if err != nil {
		panic(err)
	}
	ErrLogger = log.New(err_log_file, "", log.Ldate|log.Ltime|log.Lshortfile)

	os.MkdirAll(DataDir, 0700)
}
