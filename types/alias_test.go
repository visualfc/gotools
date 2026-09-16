package types

import (
	"bytes"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"go/types"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"
)

func TestParserMethodPointerAlias(t *testing.T) {
	const src = `package p
type aBuilder struct{}
type Builder = *aBuilder
func (b Builder) ChangeInterface() {}
func use(b Builder) { b.ChangeInterface() }
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "p.go", src, 0)
	if err != nil {
		t.Fatal(err)
	}
	info := &types.Info{
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
	}
	if _, err := new(types.Config).Check("p", fset, []*ast.File{f}, info); err != nil {
		t.Fatal(err)
	}

	var method types.Object
	for id, obj := range info.Uses {
		if id.Name == "ChangeInterface" && obj != nil {
			method = obj
			break
		}
	}
	if method == nil {
		t.Fatal("ChangeInterface use not found")
	}
	recv := method.Type().(*types.Signature).Recv().Type()
	if _, ok := recv.(*types.Alias); !ok {
		t.Fatalf("receiver type = %T, want *types.Alias", recv)
	}

	named, name, ok := parserMethod(method)
	if !ok {
		t.Fatalf("parserMethod failed for alias receiver %T %s", recv, recv)
	}
	if named.Obj().Name() != "aBuilder" || name != "ChangeInterface" {
		t.Fatalf("parserMethod = %s.%s, want aBuilder.ChangeInterface", named.Obj().Name(), name)
	}
}

func TestAliasFindUsages(t *testing.T) {
	dir, err := filepath.Abs(filepath.Join("_testdata", "alias"))
	if err != nil {
		t.Fatal(err)
	}
	for _, ident := range []string{"ChangeInterface", "ChangeType", "Name", "X"} {
		t.Run(ident, func(t *testing.T) {
			refs := identRefs(t, dir, ident)
			if len(refs) < 2 {
				t.Fatalf("testdata needs a definition and at least one use for %s, got %v", ident, refs)
			}
			want := make([]string, len(refs))
			for i, ref := range refs {
				want[i] = ref.key
			}
			sort.Strings(want)
			for _, ref := range refs {
				got := lookupUsagePositions(t, dir, ref.file, ref.offset)
				if !sameStringSet(got, want) {
					t.Errorf("cursor %s:%d\n got %v\nwant %v", ref.file, ref.offset, got, want)
				}
			}
		})
	}
}

type identRef struct {
	file   string
	offset int
	key    string
}

func identRefs(t *testing.T, dir, name string) []identRef {
	t.Helper()
	entries, err := filepath.Glob(filepath.Join(dir, "*.go"))
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(entries)
	var refs []identRef
	for _, filename := range entries {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, filename, nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		ast.Inspect(f, func(node ast.Node) bool {
			id, ok := node.(*ast.Ident)
			if !ok || id.Name != name {
				return true
			}
			pos := fset.Position(id.Pos())
			file := filepath.Base(filename)
			refs = append(refs, identRef{
				file:   file,
				offset: pos.Offset,
				key:    posKey(file, pos.Line, pos.Column),
			})
			return true
		})
	}
	return refs
}

func lookupUsagePositions(t *testing.T, dir, file string, offset int) []string {
	t.Helper()
	var buf bytes.Buffer
	w := NewPkgWalker(&build.Default)
	w.SetOutput(&buf, &buf)
	w.SetFindMode(&FindMode{Usage: true})
	conf := DefaultPkgConfig()
	cursor := NewFileCursor(nil, dir, file, offset)
	pkg, conf, err := w.Check(dir, conf, cursor)
	if err != nil {
		t.Fatalf("check %s:%d: %v", file, offset, err)
	}
	if err := w.LookupCursor(pkg, conf, cursor); err != nil {
		t.Fatalf("lookup %s:%d: %v\n%s", file, offset, err, buf.String())
	}
	var got []string
	for _, line := range strings.Split(strings.TrimSpace(buf.String()), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fileName, lineNo, col, ok := parsePositionLine(line)
		if !ok {
			t.Fatalf("unexpected output %q", line)
		}
		got = append(got, posKey(fileName, lineNo, col))
	}
	sort.Strings(got)
	return uniqueStrings(got)
}

func parsePositionLine(s string) (file string, line, col int, ok bool) {
	colIdx := strings.LastIndex(s, ":")
	if colIdx < 0 {
		return
	}
	lineIdx := strings.LastIndex(s[:colIdx], ":")
	if lineIdx < 0 {
		return
	}
	var err error
	line, err = strconv.Atoi(s[lineIdx+1 : colIdx])
	if err != nil {
		return
	}
	col, err = strconv.Atoi(s[colIdx+1:])
	if err != nil {
		return
	}
	return filepath.Base(s[:lineIdx]), line, col, true
}

func posKey(file string, line, col int) string {
	return filepath.Base(file) + ":" + strconv.Itoa(line) + ":" + strconv.Itoa(col)
}

func sameStringSet(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func uniqueStrings(in []string) []string {
	if len(in) == 0 {
		return in
	}
	out := in[:0]
	var last string
	for i, s := range in {
		if i == 0 || s != last {
			out = append(out, s)
			last = s
		}
	}
	return out
}
