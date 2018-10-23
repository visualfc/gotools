package pkgcheck

import (
	"fmt"
	"go/build"
	"os"
	"path/filepath"

	"github.com/visualfc/fastmod"
	"github.com/visualfc/gotools/pkg/command"
	"github.com/visualfc/gotools/pkg/pkgutil"
)

var Command = &command.Command{
	Run:       runCheck,
	UsageLine: "pkgcheck [-pkg | -name] -w .",
	Short:     "pkg check utils",
	Long:      "check pkg mod or vendor path",
}

var (
	flagCheckPkg  string
	flagCheckDir  string = "."
	flagCheckName bool
)

func init() {
	Command.Flag.StringVar(&flagCheckPkg, "pkg", "", "check pkg name")
	Command.Flag.BoolVar(&flagCheckName, "name", false, "check module name")
	Command.Flag.StringVar(&flagCheckDir, "w", "", "work path")
}

func runCheck(cmd *command.Command, args []string) error {
	if flagCheckPkg == "" && !flagCheckName {
		cmd.Usage()
		return os.ErrInvalid
	}
	if flagCheckDir == "" || flagCheckDir == "." {
		flagCheckDir, _ = os.Getwd()
	}
	mods := fastmod.NewModuleList(&build.Default)
	mod, _ := mods.LoadModule(flagCheckDir)
	if flagCheckName {
		if mod != nil {
			fmt.Println(mod.Path())
		} else {
			_, fname := filepath.Split(flagCheckDir)
			fmt.Println(fname)
		}
		return nil
	}
	// check mod, check vendor
	if mod != nil {
		_, dir, _ := mod.Lookup(flagCheckPkg)
		if dir != "" {
			fmt.Printf("%s,mod\n", dir)
			return nil
		}
	} else {
		pkg := pkgutil.ImportDir(flagCheckDir)
		if pkg != nil {
			found, _ := pkgutil.VendoredImportPath(pkg, flagCheckPkg)
			if found != "" && found != flagCheckPkg {
				fmt.Printf("%s,vendor\n", found)
				return nil
			}
		}
	}
	fmt.Printf("%s,pkg\n", flagCheckPkg)
	return nil
}
