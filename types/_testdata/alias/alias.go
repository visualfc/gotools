package fixture

type aBuilder struct{}

// Builder is a pointer alias, matching types such as llgo/ssa.Builder = *aBuilder.
type Builder = *aBuilder

func (b Builder) ChangeInterface() {}

func (b Builder) ChangeType() {}

type aValue struct {
	X int
}

type Value = aValue

func (v Value) Name() string { return "v" }
