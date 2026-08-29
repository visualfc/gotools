package gofmt

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"golang.org/x/tools/imports"
)

func TestFixImportsFile(t *testing.T) {
	testFixImportsCase(t, "fiximports")
	testFixImportsCase(t, "nonstd")
	testFixImportsCase(t, "nonstd_existing")
}

func testFixImportsCase(t *testing.T, name string) {
	t.Helper()
	src, err := os.ReadFile("_testdata/" + name + "/input.go")
	if err != nil {
		t.Fatal(err)
	}
	want, err := os.ReadFile("_testdata/" + name + "/want.txt")
	if err != nil {
		t.Fatal(err)
	}
	oldFix, oldOptions := gofmtFixImports, options
	defer func() { gofmtFixImports, options = oldFix, oldOptions }()
	gofmtFixImports = true
	options = &imports.Options{FormatOnly: false, TabWidth: 8, TabIndent: true, Comments: true, Fragment: true}
	var out bytes.Buffer
	if err := processFile("_testdata/"+name+"/input.go", bytes.NewReader(src), &out, false); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(out.Bytes(), want) {
		t.Fatalf("fiximports output mismatch\n got:\n%s\nwant:\n%s", out.Bytes(), want)
	}
	var again bytes.Buffer
	if err := processFile("_testdata/"+name+"/input.go", bytes.NewReader(out.Bytes()), &again, false); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(again.Bytes(), out.Bytes()) {
		t.Fatalf("fiximports is not idempotent\n first:\n%s\nsecond:\n%s", out.Bytes(), again.Bytes())
	}
	for _, line := range strings.Split(string(out.Bytes()), "\n") {
		if strings.TrimRight(line, " \t") != line {
			t.Fatalf("fiximports left trailing whitespace in %q", line)
		}
	}
}
