package types

import (
	"go/build"
	"log"
	"testing"
)

func TestTypes(t *testing.T) {
	w := NewPkgWalker(&build.Default)
	conf := DefaultPkgConfig()
	conf.Cursor = NewFileCursor(nil, "types.go", 100)
	w.Check(".", conf)
	for _, v := range w.Imported {
		log.Println(v.Name(), v.Path())
	}
}
