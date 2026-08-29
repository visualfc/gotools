package fixture

func (b Box[T]) Get() T { return b.Value } // cursor:Get

func Identity[T comparable](v T) T { return v }
