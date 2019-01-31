package types

import (
	"go/build"
	"os"
	"testing"
)

func _TestTypes(t *testing.T) {
	w := NewPkgWalker(&build.Default)
	w.SetOutput(os.Stdout, os.Stderr)
	w.SetFindMode(&FindMode{Info: true, Doc: true, Define: true})
	conf := DefaultPkgConfig()
	dir, _ := os.Getwd()
	cursor := NewFileCursor(nil, dir, "types_test.go", 126)
	pkg, conf, err := w.Check(dir, conf, cursor)
	if err != nil {
		t.Fatalf("error %v\n", err)
	}
	w.LookupCursor(pkg, conf, cursor)
}

func _TestOS(t *testing.T) {
	w := NewPkgWalker(&build.Default)

	w.SetOutput(os.Stdout, os.Stderr)
	w.SetFindMode(&FindMode{Info: true, Doc: true, Define: true})
	conf := DefaultPkgConfig()
	fn1 := func() {
		cursor := NewFileCursor(nil, "", "dir_windows.go", 235)
		pkg, conf, _ := w.Check("os", conf, cursor)
		w.LookupCursor(pkg, conf, cursor)
	}
	fn2 := func() {
		cursor := NewFileCursor(nil, "", "dir_unix.go", 1040)
		pkg, conf, _ := w.Check("os", conf, cursor)
		w.LookupCursor(pkg, conf, cursor)
	}
	for i := 0; i < 2; i++ {
		fn1()
		fn2()
	}
}
