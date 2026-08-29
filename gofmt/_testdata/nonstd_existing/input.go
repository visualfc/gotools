package fixture

import (
	"fmt"
	"github.com/visualfc/gotools/pkg/godiff"
	"github.com/visualfc/gotools/pkg/gomod"
)

func Use() string {
	_ = gomod.Module{}
	_ = pkgutil.IsVendorExperiment()
	fmt.Print("")
	return godiff.UnifiedDiffString("a\n", "b\n")
}
