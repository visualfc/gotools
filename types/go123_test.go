//go:build go1.23
// +build go1.23

package types

import (
	"go/build"
	"testing"
)

func TestGo123StdlibTypeCheck(t *testing.T) {
	w := NewPkgWalker(&build.Default)
	conf := DefaultPkgConfig()
	if _, _, err := w.Check("iter", conf, nil); err != nil {
		t.Fatalf("type checking iter: %v", err)
	}
}
