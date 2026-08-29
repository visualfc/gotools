package fixture

type Integer interface {
	~int | ~int64
}

type Pair[A, B any] struct {
	First  A
	Second B
}

type List[T any] []T

type Alias[T any] = List[T]

type Combiner[T, U any] interface {
	Combine(T, U) T
}

func Identity[T any](v T) T { return v }

func Map[T, U any](xs []T, f func(T) U) []U {
	out := make([]U, len(xs))
	for i, x := range xs {
		out[i] = f(x)
	}
	return out
}

func (p Pair[A, B]) FirstValue() A { return p.First }

var _ = Pair[string, int]{First: "x", Second: 1}
