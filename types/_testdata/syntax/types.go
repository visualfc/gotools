package fixture

import "errors"

const (
	First = iota
	Second
)

type Alias = string
type Number int

type Reader interface {
	Read([]byte) (int, error)
}

type Embedded interface {
	Reader
	Close() error
}

type Record struct {
	Name   Alias
	Values []Number
	Next   *Record
}

func (r Record) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, errors.New("empty")
	}
	return len(p), nil
}

func (r *Record) Close() error {
	r.Next = nil
	return nil
}

var Global = map[string]Number{}
