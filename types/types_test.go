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
	conf.Cursor = NewFileCursor(nil, "types_test.go", 138)
	pkg, err := w.Check(".", conf)
	if err != nil {
		t.Fatalf("error %v\n", err)
	}
	w.LookupCursor(pkg, conf, conf.Cursor)
}
