package fixture

func (b Box[T]) Get() T { return b.Value }

func Identity[T comparable](v T) T { return v }

func Map[T, U any](values []T, fn func(T) U) []U {
	out := make([]U, len(values))
	for i, value := range values {
		out[i] = fn(value)
	}
	return out
}

func Sum[T Numeric](a, b T) T { return a + b }

func Use() int { // cursor:Use
	b := Box[int]{Value: 1}
	p := Pair[string, int]{First: "x", Second: b.Get()}
	values := Slice[int]{p.Second}
	result := Map(values, func(v int) int { return Identity(v) })
	return int(Sum(p.Second, len(result)))
}
