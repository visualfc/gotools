package types

import (
	"bytes"
	"go/build"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/visualfc/gotools/pkg/command"
)

func TestRunTypesFindUsageDespiteTypeErrors(t *testing.T) {
	dir, err := filepath.Abs(filepath.Join("_testdata", "typeerror"))
	if err != nil {
		t.Fatal(err)
	}
	assertCheckTypeErrorKeepsPackage(t, dir)

	for _, ident := range []string{"Foo", "Wasm"} {
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
				got, stdout, stderr := runTypesUsagePositions(t, dir, ref.file, ref.offset)
				if stderr != "" {
					t.Fatalf("cursor %s:%d stderr = %q", ref.file, ref.offset, stderr)
				}
				if strings.Contains(stdout, "error loading package") {
					t.Fatalf("cursor %s:%d stdout = %q", ref.file, ref.offset, stdout)
				}
				if !sameStringSet(got, want) {
					t.Errorf("cursor %s:%d\n got %v\nwant %v\nstdout=%s", ref.file, ref.offset, got, want, stdout)
				}
			}
		})
	}
}

func TestRunTypesFindDefDespiteTypeErrors(t *testing.T) {
	dir, err := filepath.Abs(filepath.Join("_testdata", "typeerror"))
	if err != nil {
		t.Fatal(err)
	}
	assertCheckTypeErrorKeepsPackage(t, dir)

	refs := identRefs(t, dir, "Foo")
	if len(refs) < 2 {
		t.Fatalf("testdata needs a definition and at least one use for Foo, got %v", refs)
	}
	def := refs[0]
	for _, ref := range refs {
		stdout, stderr := runTypesOutput(t, dir, ref.file, ref.offset, false, true, true)
		if stderr != "" {
			t.Fatalf("cursor %s:%d stderr = %q", ref.file, ref.offset, stderr)
		}
		if !strings.Contains(stdout, def.key) && !strings.Contains(stdout, def.file+":") {
			t.Fatalf("cursor %s:%d stdout = %q, want definition %s", ref.file, ref.offset, stdout, def.key)
		}
	}
}

func assertCheckTypeErrorKeepsPackage(t *testing.T, dir string) {
	t.Helper()
	w := NewPkgWalker(&build.Default)
	w.SetOutput(&bytes.Buffer{}, &bytes.Buffer{})
	pkg, _, err := w.Check(dir, DefaultPkgConfig(), nil)
	if err == nil {
		t.Fatal("Check error = nil; typeerror testdata must keep a type error")
	}
	if pkg == nil {
		t.Fatalf("Check package = nil, err = %v", err)
	}
	if !strings.Contains(err.Error(), "cannot use") {
		t.Fatalf("Check error = %v, want a type error", err)
	}
}

func runTypesUsagePositions(t *testing.T, dir, file string, offset int) (got []string, stdout, stderr string) {
	t.Helper()
	stdout, stderr = runTypesOutput(t, dir, file, offset, true, true, false)
	for _, line := range strings.Split(stdout, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fileName, lineNo, col, ok := parsePositionLine(line)
		if !ok {
			continue
		}
		got = append(got, posKey(fileName, lineNo, col))
	}
	sort.Strings(got)
	return uniqueStrings(got), stdout, stderr
}

func runTypesOutput(t *testing.T, dir, file string, offset int, use, info, def bool) (stdout, stderr string) {
	t.Helper()
	oldUse, oldInfo, oldDef, oldPos, oldStdin, oldVerbose := typesFindUse, typesFindInfo, typesFindDef, typesFilePos, typesFileStdin, typesVerbose
	defer func() {
		typesFindUse, typesFindInfo, typesFindDef, typesFilePos, typesFileStdin, typesVerbose = oldUse, oldInfo, oldDef, oldPos, oldStdin, oldVerbose
	}()
	typesFindUse = use
	typesFindInfo = info
	typesFindDef = def
	typesFileStdin = false
	typesVerbose = false
	typesFilePos = file + ":" + strconv.Itoa(offset)

	var outBuf, errBuf bytes.Buffer
	cmd := &command.Command{Stdout: &outBuf, Stderr: &errBuf}
	if err := runTypes(cmd, []string{dir}); err != nil {
		t.Fatalf("runTypes %s:%d: %v\nstderr=%s\nstdout=%s", file, offset, err, errBuf.String(), outBuf.String())
	}
	return outBuf.String(), errBuf.String()
}
