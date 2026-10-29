package fixture

func broken() {
	var _ int = "type error"
}
