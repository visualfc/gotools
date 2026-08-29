package fixture

// Box stores a value of any type.
type Box[T any] struct {
	Value T
}

type Reader interface {
	Read() string
}

type Record struct {
	Name string
}
