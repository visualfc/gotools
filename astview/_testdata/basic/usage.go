package fixture

import "fmt"

func (b Box[T]) Get() T { return b.Value }

func Use() string {
	b := Box[string]{Value: "ok"}
	return fmt.Sprint(Record{Name: b.Get()})
}
