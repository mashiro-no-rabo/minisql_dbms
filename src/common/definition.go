package common

type Comparable interface {
	EqualsTo(Comparable) bool
	LessThan(Comparable) bool
}

type KeyType interface {
	GetValue() ValueType
	Comparable
}

type ValueType interface {
	Comparable
}
