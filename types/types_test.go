package types

import (
	"go/build"
	"os"
	"testing"

	"github.com/visualfc/gotools/pkg/command"
)

func testCommand() *command.Command {
	cmd := &command.Command{}
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd
}

func TestTypes(t *testing.T) {
	typesFindInfo = true
	typesFindDoc = true
	typesFindDef = true
	w := NewPkgWalker(&build.Default)
	w.cmd = testCommand()
	conf := DefaultPkgConfig()
	conf.Cursor = NewFileCursor(nil, "types_test.go", 325)
	pkg, err := w.Check(".", conf)
	if err != nil {
		t.Fatalf("error %v\n", err)
	}
	w.LookupCursor(pkg, conf, conf.Cursor)
}
