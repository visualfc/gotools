package fixture

type Box[T any] struct {
	Value T
}

type Pair[A, B any] struct {
	First  A
	Second B
}

type Numeric interface {
	~int | ~int64
}

type Slice[T any] []T
