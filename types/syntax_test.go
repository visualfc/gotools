//go:build go1.18
// +build go1.18

package types

import (
	"bytes"
	"go/ast"
	"go/build"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestFiles(t *testing.T) {
	for _, name := range []string{"syntax", "generic"} {
		t.Run(name, func(t *testing.T) { testFiles(t, name) })
	}
}

func testFiles(t *testing.T, name string) {
	dir, err := filepath.Abs(filepath.Join("_testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	w := NewPkgWalker(&build.Default)
	var out bytes.Buffer
	w.SetOutput(&out, &out)
	conf := DefaultPkgConfig()
	pkg, conf, err := w.Check(dir, conf, nil)
	if err != nil {
		t.Fatalf("%s package check: %v", name, err)
	}
	if pkg == nil || pkg.Name() != "fixture" {
		t.Fatalf("checked package = %v", pkg)
	}
	if len(conf.Files) == 0 {
		t.Fatal("parsed files = 0")
	}
	if len(conf.Info.Defs) == 0 || len(conf.Info.Uses) == 0 {
		t.Fatalf("type info too small: defs=%d uses=%d", len(conf.Info.Defs), len(conf.Info.Uses))
	}
	if len(conf.Info.Types) == 0 || len(conf.Info.Scopes) == 0 {
		t.Fatalf("missing expression type/scope information")
	}
	checkWant(t, filepath.Join(dir, "want.txt"), w.FileSet, conf)
}

func TestBuiltinInfoMap(t *testing.T) {
	for _, name := range types.Universe.Names() {
		if _, ok := types.Universe.Lookup(name).(*types.Builtin); !ok {
			continue
		}
		if _, ok := builtinInfoMap[name]; !ok {
			t.Errorf("builtinInfoMap missing %q", name)
		}
	}
}

func checkWant(t *testing.T, filename string, fset *token.FileSet, conf *PkgConfig) {
	t.Helper()
	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		line = strings.TrimSpace(line)
		parts := strings.Split(line, "|")
		if len(parts) < 2 || len(parts) > 4 {
			t.Fatalf("invalid want line %q", line)
		}
		category := parts[0]
		switch category {
		case "def":
			category = "defs"
		case "use":
			category = "uses"
		case "instance":
			category = "instances"
		}
		file, lineNo, colNo := parseWantPosition(t, parts[1])
		wantPos := token.Pos(0)
		var name string
		if len(parts) >= 3 {
			name = parts[2]
		}
		match := func(pos token.Pos) bool {
			p := fset.Position(pos)
			return filepath.Base(p.Filename) == file && p.Line == lineNo && p.Column == colNo
		}
		switch category {
		case "defs":
			for ident := range conf.Info.Defs {
				if ident.Name == name && match(ident.Pos()) {
					wantPos = ident.Pos()
					break
				}
			}
		case "uses":
			for ident := range conf.Info.Uses {
				if ident.Name == name && match(ident.Pos()) {
					wantPos = ident.Pos()
					break
				}
			}
		case "types":
			for expr := range conf.Info.Types {
				if match(expr.Pos()) {
					wantPos = expr.Pos()
					break
				}
			}
		case "info":
			for expr, typeAndValue := range conf.Info.Types {
				if match(expr.Pos()) && (len(parts) < 3 || types.TypeString(typeAndValue.Type, nil) == parts[2]) {
					wantPos = expr.Pos()
					break
				}
			}
		case "selections":
			for expr := range conf.Info.Selections {
				if expr.Sel.Name == name && match(expr.Sel.Pos()) {
					wantPos = expr.Sel.Pos()
					break
				}
			}
		case "instances":
			for ident := range conf.Info.Instances {
				instance := conf.Info.Instances[ident]
				argsMatch := len(parts) < 4 || typeArgsString(instance.TypeArgs) == parts[3]
				if ident.Name == name && match(ident.Pos()) && argsMatch {
					wantPos = ident.Pos()
					break
				}
			}
		case "builtins":
			for _, fileAST := range conf.Files {
				ast.Inspect(fileAST, func(node ast.Node) bool {
					ident, ok := node.(*ast.Ident)
					if ok && ident.Name == name && match(ident.Pos()) {
						if _, builtin := types.Universe.Lookup(name).(*types.Builtin); builtin {
							wantPos = ident.Pos()
						}
						return false
					}
					return true
				})
			}
		default:
			t.Fatalf("unknown want category %q", category)
		}
		if wantPos == token.NoPos {
			t.Fatalf("missing %s entry %s at %s:%d:%d", parts[0], name, file, lineNo, colNo)
		}
	}
}

func typeArgsString(args *types.TypeList) string {
	if args == nil {
		return "[]"
	}
	items := make([]string, args.Len())
	for i := range items {
		items[i] = types.TypeString(args.At(i), nil)
	}
	return "[" + strings.Join(items, ", ") + "]"
}

func parseWantPosition(t *testing.T, value string) (string, int, int) {
	t.Helper()
	parts := strings.Split(value, ":")
	if len(parts) != 3 {
		t.Fatalf("invalid want position %q", value)
	}
	lineNo, err := strconv.Atoi(parts[1])
	if err != nil {
		t.Fatal(err)
	}
	colNo, err := strconv.Atoi(parts[2])
	if err != nil {
		t.Fatal(err)
	}
	return parts[0], lineNo, colNo
}
