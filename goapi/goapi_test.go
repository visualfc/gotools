//go:build go1.23
// +build go1.23

package goapi

import (
	"bytes"
	"go/build"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenericConstraintPackage(t *testing.T) {
	w := NewWalker()
	w.context = &build.Default
	w.WalkPackage("cmp")
	if p := w.findPackage("cmp"); p == nil {
		t.Fatal("cmp package was not loaded")
	}
}

func TestGenericConstraintFeatures(t *testing.T) {
	w := NewWalker()
	w.context = &build.Default
	w.wantedPkg["cmp"] = true
	w.WalkPackage("cmp")
	features := strings.Join(w.Features(""), "\n")
	for _, want := range []string{
		"pkg cmp, func Compare(T, T) int",
		"pkg cmp, func Less(T, T) bool",
		"pkg cmp, type Ordered interface {}",
	} {
		if !strings.Contains(features, want) {
			t.Errorf("features do not contain %q:\n%s", want, features)
		}
	}
}

func TestCompareAPI(t *testing.T) {
	tests := []struct {
		name     string
		features []string
		required []string
		optional []string
		except   []string
		allowNew bool
		ok       bool
		want     string
	}{
		{
			name:     "compatible",
			features: []string{"pkg p, func F()"},
			required: []string{"pkg p, func F()"},
			ok:       true,
		},
		{
			name:     "missing required",
			required: []string{"pkg p, func F()"},
			want:     "-pkg p, func F()\n",
		},
		{
			name:     "optional addition",
			features: []string{"pkg p, func F()"},
			optional: []string{"pkg p, func F()"},
			allowNew: false,
			ok:       true,
		},
		{
			name:     "exception",
			required: []string{"pkg p, func F()"},
			except:   []string{"pkg p, func F()"},
			want:     "~pkg p, func F()\n",
			ok:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got bytes.Buffer
			if ok := compareAPI(&got, tt.features, tt.required, tt.optional, tt.except, tt.allowNew); ok != tt.ok {
				t.Fatalf("compareAPI ok = %v, want %v", ok, tt.ok)
			}
			if got.String() != tt.want {
				t.Fatalf("compareAPI output = %q, want %q", got.String(), tt.want)
			}
		})
	}
}

func TestContextAndFeatureHelpers(t *testing.T) {
	c := parseContext("linux-amd64-cgo")
	if c.GOOS != "linux" || c.GOARCH != "amd64" || !c.CgoEnabled {
		t.Fatalf("parseContext returned %+v", c)
	}
	if got := contextName(c); got != "linux-amd64-cgo" {
		t.Fatalf("contextName = %q", got)
	}
	if got := featureWithoutContext("func F (linux-amd64), x"); got != "func F, x" {
		t.Fatalf("featureWithoutContext = %q", got)
	}
}

func TestCursorInfo(t *testing.T) {
	dir := t.TempDir()
	src := "package sample\n\nimport \"fmt\"\n\ntype Item struct { Name string }\n\nfunc (i Item) String() string { return i.Name }\n\nfunc Hello() { value := Item{Name: \"hi\"}; fmt.Println(value.Name) }\n"
	if err := os.WriteFile(filepath.Join(dir, "sample.go"), []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		text string
		kind Kind
		sym  string
	}{
		{name: "package", text: "sample", kind: KindPackage, sym: "sample"},
		{name: "import", text: "fmt", kind: KindImport, sym: "fmt"},
		{name: "type", text: "Item", kind: KindStruct, sym: "Item"},
		{name: "field", text: "Name string", kind: KindField, sym: "Item.Name"},
		{name: "method", text: "String()", kind: KindMethod, sym: "Item.String"},
		{name: "function", text: "Hello", kind: KindFunc, sym: "Hello"},
		{name: "local", text: "value :=", kind: KindVar, sym: "value"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos := strings.Index(src, tt.text) + 1
			w := NewWalker()
			w.context = &build.Default
			w.cursorInfo = &CursorInfo{pkg: "sample", file: "sample.go", pos: tokenPos(pos)}
			w.WalkPackage(dir)
			if w.cursorInfo.info == nil {
				t.Fatalf("no cursor info for %q", tt.text)
			}
			if w.cursorInfo.info.Kind != tt.kind || w.cursorInfo.info.Name != tt.sym {
				t.Fatalf("cursor info = %s %q, want %s %q", w.cursorInfo.info.Kind, w.cursorInfo.info.Name, tt.kind, tt.sym)
			}
		})
	}
}

func TestCursorInfoFromStandardInput(t *testing.T) {
	src := "package iter\n\ntype Seq[V any] func(yield func(V) bool)\n"
	w := NewWalker()
	w.context = &build.Default
	w.cursorInfo = &CursorInfo{
		pkg:  "iter",
		file: "iter.go",
		pos:  tokenPos(strings.Index(src, "Seq") + 1),
		src:  []byte(src),
		std:  true,
	}
	w.WalkPackage("iter")
	if w.cursorInfo.info == nil {
		t.Fatal("no cursor info from standard input")
	}
	if w.cursorInfo.info.Kind != KindType || w.cursorInfo.info.Name != "Seq" {
		t.Fatalf("cursor info = %s %q, want type Seq", w.cursorInfo.info.Kind, w.cursorInfo.info.Name)
	}
}

func TestGenericMethodFeaturesAndCursor(t *testing.T) {
	dir := t.TempDir()
	src := "package generic\n\ntype Box[T any] struct { Value T }\n\nfunc (b Box[T]) Get() T { return b.Value }\n\nfunc Identity[T comparable](v T) T { return v }\n"
	if err := os.WriteFile(filepath.Join(dir, "generic.go"), []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}

	w := NewWalker()
	w.context = &build.Default
	w.wantedPkg["generic"] = true
	w.WalkPackage(dir)
	features := strings.Join(w.Features(""), "\n")
	for _, want := range []string{
		"pkg generic, func Identity(T) T",
		"pkg generic, method (Box[T]) Get() T",
	} {
		if !strings.Contains(features, want) {
			t.Errorf("generic features do not contain %q:\n%s", want, features)
		}
	}

	pos := strings.Index(src, "Get") + 1
	w = NewWalker()
	w.context = &build.Default
	w.cursorInfo = &CursorInfo{pkg: "generic", file: "generic.go", pos: tokenPos(pos)}
	w.WalkPackage(dir)
	if w.cursorInfo.info == nil {
		t.Fatal("no cursor info for generic method")
	}
	if w.cursorInfo.info.Kind != KindMethod || w.cursorInfo.info.Name != "Box.Get" {
		t.Fatalf("generic method cursor info = %s %q", w.cursorInfo.info.Kind, w.cursorInfo.info.Name)
	}
}

func TestCursorInfoAcrossFileSet(t *testing.T) {
	dir := t.TempDir()
	files := map[string]string{
		"types.go":   "package split\n\ntype Box[T any] struct { Value T }\n",
		"methods.go": "package split\n\nfunc (b Box[T]) Get() T {\n\treturn b.Value\n}\n",
	}
	for name, src := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(src), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	src := files["methods.go"]
	pos := strings.Index(src, "Get") + 1
	w := NewWalker()
	w.context = &build.Default
	w.cursorInfo = &CursorInfo{pkg: "split", file: "methods.go", pos: tokenPos(pos)}
	w.WalkPackage(dir)
	info := w.cursorInfo.info
	if info == nil {
		t.Fatal("no cursor info in second file")
	}
	if info.Kind != KindMethod || info.Name != "Box.Get" {
		t.Fatalf("cursor info = %s %q, want method Box.Get", info.Kind, info.Name)
	}
	if info.T == nil {
		t.Fatal("cursor info has no AST position")
	}
	file := w.fset.File(info.T.Pos())
	if file == nil || filepath.Base(file.Name()) != "methods.go" {
		t.Fatalf("AST position file = %v, want methods.go", file)
	}
	if got := w.fset.Position(info.T.Pos()); got.Line != 3 {
		t.Fatalf("method position = %v, want line 3", got)
	}
}

func TestFiles(t *testing.T) {
	testFiles(t, "generic")
}

func testFiles(t *testing.T, name string) {
	t.Helper()
	dir, err := filepath.Abs(filepath.Join("_testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	w := NewWalker()
	w.context = &build.Default
	w.wantedPkg["fixture"] = true
	w.WalkPackage(dir)

	wantBytes, err := os.ReadFile(filepath.Join(dir, "want.txt"))
	if err != nil {
		t.Fatal(err)
	}
	want := strings.TrimSpace(strings.ReplaceAll(string(wantBytes), "\r\n", "\n"))
	got := strings.TrimSpace(strings.Join(w.Features(""), "\n"))
	if got != want {
		t.Fatalf("API features mismatch:\n got:\n%s\nwant:\n%s", got, want)
	}

	src, err := os.ReadFile(filepath.Join(dir, "methods.go"))
	if err != nil {
		t.Fatal(err)
	}
	marker := "// cursor:Get"
	markerPos := strings.Index(string(src), marker)
	if markerPos < 0 {
		t.Fatalf("missing cursor marker %q", marker)
	}
	w = NewWalker()
	w.context = &build.Default
	w.cursorInfo = &CursorInfo{
		pkg:  "fixture",
		file: "methods.go",
		pos:  tokenPos(strings.Index(string(src[:markerPos]), "Get") + 1),
	}
	w.WalkPackage(dir)
	if w.cursorInfo.info == nil || w.cursorInfo.info.Kind != KindMethod || w.cursorInfo.info.Name != "Box.Get" {
		t.Fatalf("cursor info = %#v, want method Box.Get", w.cursorInfo.info)
	}
}

// tokenPos keeps cursor offsets in the same representation used by runApi.
func tokenPos(pos int) token.Pos { return token.Pos(pos) }
