//go:build go1.23
// +build go1.23

package stdlib

import "testing"

func TestGo123AndLaterPackages(t *testing.T) {
	for _, pkg := range []string{
		"cmp",
		"encoding/json/v2",
		"iter",
		"log/slog",
		"math/rand/v2",
		"slices",
		"testing/synctest",
		"unique",
		"weak",
	} {
		if !IsStdPkg(pkg) {
			t.Errorf("IsStdPkg(%q) = false", pkg)
		}
	}
}
