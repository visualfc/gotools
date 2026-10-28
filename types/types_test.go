package types

import (
	"fmt"
	"go/build"
	"os"
	"strings"
	"testing"
)

func TestTypes(t *testing.T) {
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

func TestImportHelperReturnsRecoveredPanics(t *testing.T) {
	w := NewPkgWalker(nil)
	w.SetOutput(os.Stdout, os.Stderr)

	_, _, err := w.ImportHelper("", "example.com/panic", "", DefaultPkgConfig(), nil)
	if err == nil || !strings.Contains(err.Error(), "panic while loading package") {
		t.Fatalf("ImportHelper error = %v, want recovered panic error", err)
	}
}

func test_error() {
	e := fmt.Errorf("error_test")
	_ = e.Error()
}

func TestError(t *testing.T) {
	w := NewPkgWalker(&build.Default)
	w.SetOutput(os.Stdout, os.Stderr)
	w.SetFindMode(&FindMode{Info: true, Doc: true, Define: true, Usage: true, UsageAll: true})
	conf := DefaultPkgConfig()
	dir, _ := os.Getwd()
	cursor := NewFileCursor(nil, dir, "types_test.go", 530)
	pkg, conf, err := w.Check(dir, conf, cursor)
	if err != nil {
		t.Fatalf("error %v\n", err)
	}
	w.LookupCursor(pkg, conf, cursor)
}

func TestOS(t *testing.T) {
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
