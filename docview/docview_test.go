package docview

import (
	"bufio"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPackageDocFiles(t *testing.T) {
	dir := filepath.Join("_testdata", "basic")
	fset := token.NewFileSet()
	pkg, err := parser.ParseDir(fset, dir, func(info os.FileInfo) bool {
		return strings.HasSuffix(info.Name(), ".go")
	}, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	astPkg := pkg["fixture"]
	if astPkg == nil {
		t.Fatal("fixture package not parsed")
	}
	doc := NewPackageDoc(astPkg, "fixture", true)
	checkDocWant(t, filepath.Join(dir, "want.txt"), doc)
}

func checkDocWant(t *testing.T, filename string, doc *PackageDoc) {
	t.Helper()
	f, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	byName := func(name string) *TypeDoc {
		for _, typ := range doc.Types {
			if typ.Type.Name.Name == name {
				return typ
			}
		}
		return nil
	}
	funcDoc := func(name string) *FuncDoc {
		for _, fn := range doc.Funcs {
			if fn.Name == name {
				return fn
			}
		}
		return nil
	}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		parts := strings.Split(sc.Text(), "|")
		if len(parts) < 2 {
			t.Fatalf("invalid want line %q", sc.Text())
		}
		var found bool
		switch parts[0] {
		case "package":
			found = doc.PackageName == parts[1]
		case "doc":
			found = strings.TrimSpace(doc.Doc) == parts[1]
		case "import":
			for _, imp := range doc.Imports {
				found = imp == parts[1] || found
			}
		case "type":
			typ := byName(parts[1])
			found = typ != nil && strings.TrimSpace(typ.Doc) == parts[2]
		case "method":
			typ := byName(parts[1])
			if typ != nil {
				for _, method := range typ.Methods {
					if method.Name == parts[2] && strings.TrimSpace(method.Doc) == parts[3] {
						found = true
					}
				}
			}
		case "factory":
			typ := byName(parts[1])
			if typ != nil {
				for _, fn := range typ.Funcs {
					if fn.Name == parts[2] && strings.TrimSpace(fn.Doc) == parts[3] {
						found = true
					}
				}
			}
		case "func":
			fn := funcDoc(parts[1])
			found = fn != nil && strings.TrimSpace(fn.Doc) == parts[2]
		default:
			t.Fatalf("unknown want category %q", parts[0])
		}
		if !found {
			t.Fatalf("missing doc entry %q", sc.Text())
		}
	}
	if err := sc.Err(); err != nil {
		t.Fatal(err)
	}
}
