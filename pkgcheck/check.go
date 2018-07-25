package pkgcheck

import (
	"fmt"
	"os"

	"github.com/visualfc/gotools/pkg/pkgutil"

	"github.com/visualfc/gotools/pkg/gomod"

	"github.com/visualfc/gotools/pkg/command"
)

var Command = &command.Command{
	Run:       runCheck,
	UsageLine: "pkgcheck -pkg -w .",
	Short:     "pkg check utils",
	Long:      "check pkg mod or vendor path",
}

var (
	flagCheckPkg string
	flagCheckDir string = "."
)

func init() {
	Command.Flag.StringVar(&flagCheckPkg, "pkg", "", "check pkg name")
	Command.Flag.StringVar(&flagCheckDir, "w", "", "work path")
}

func runCheck(cmd *command.Command, args []string) error {
	if flagCheckPkg == "" {
		cmd.Usage()
		return os.ErrInvalid
	}
	if flagCheckDir == "" || flagCheckDir == "." {
		flagCheckDir, _ = os.Getwd()
	}
	modList := gomod.LooupModList(flagCheckDir)
	// check mod, check vendor
	if modList != nil {
		m, path, _ := modList.LookupModule(flagCheckPkg)
		if m != nil {
			fmt.Printf("%s,mod\n", path)
			return nil
		}
	} else {
		pkg := pkgutil.ImportDir(flagCheckDir)
		if pkg != nil {
			found, _ := pkgutil.VendoredImportPath(pkg, flagCheckPkg)
			if found != "" {
				fmt.Printf("%s,vendor\n", found)
				return nil
			}
		}
	}
	fmt.Printf("%s,pkg\n", flagCheckPkg)
	return nil
}
