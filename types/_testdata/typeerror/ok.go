package fixture

type Config struct {
	Wasm bool
}

func Foo() {}

func Bar() {
	Foo()
	_ = Config{Wasm: true}
}
