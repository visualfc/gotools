package types

import (
	"go/build"
	"os"
	"testing"
)

func TestTypes(t *testing.T) {
	w := NewPkgWalker(&build.Default)
	w.SetOutput(os.Stdout, os.Stderr)
	w.SetFindMode(&FindMode{Info: true, Doc: true, Define: true})
	conf := DefaultPkgConfig()
	dir, _ := os.Getwd()
	cursor := NewFileCursor(nil, dir, "types_test.go", 138)
	pkg, conf, err := w.Check(".", conf)
	if err != nil {
		t.Fatalf("error %v\n", err)
	}
	w.LookupCursor(pkg, conf, cursor)
	//	pkg, err = w.Check(".", conf)
	w.LookupCursor(pkg, conf, cursor)
}
