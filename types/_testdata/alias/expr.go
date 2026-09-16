package fixture

func UseBuilder(b Builder) {
	b.ChangeInterface()
	b.ChangeType()
}

func UseBuilderAgain(b Builder) {
	b.ChangeInterface()
}

func UseValue(v Value) {
	_ = v.Name()
	_ = v.X
}
