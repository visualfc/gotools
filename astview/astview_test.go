//go:build go1.18
// +build go1.18

package astview

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFiles(t *testing.T) {
	dir := filepath.Join("_testdata", "basic")
	files, err := filepath.Glob(filepath.Join(dir, "*.go"))
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("no testdata Go files")
	}
	var got bytes.Buffer
	AllFiles = nil
	astViewShowTypeParams = true
	defer func() { astViewShowTypeParams = false }()
	if err := PrintFilesTree(files, &got, true); err != nil {
		t.Fatalf("print files tree: %v", err)
	}
	want, err := os.ReadFile(filepath.Join(dir, "want.txt"))
	if err != nil {
		t.Fatal(err)
	}
	gotText := strings.ReplaceAll(strings.TrimSpace(got.String()), "\\", "/")
	wantText := strings.TrimSpace(strings.ReplaceAll(string(want), "\\", "/"))
	if gotText != wantText {
		t.Fatalf("file tree mismatch\n got:\n%s\nwant:\n%s", got.String(), want)
	}
}

func TestGenericSource(t *testing.T) {
	src := strings.NewReader("package p\n\ntype Pair[A, B any] struct { A A; B B }\n")
	view, err := NewFilePackageSource("stdin.go", src, true)
	if err != nil {
		t.Fatalf("parse generic source: %v", err)
	}
	var out strings.Builder
	view.PrintTree(&out)
	if !strings.Contains(out.String(), "Pair") {
		t.Fatalf("generic type missing from output: %s", out.String())
	}
}
