package common

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
	ColInt = itoa
	ColString
	ColFloat
)

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
