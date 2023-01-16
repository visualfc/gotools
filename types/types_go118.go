//go:build go1.18
// +build go1.18

package types

import (
	"go/ast"
	"go/types"
)

const enableTypeParams = true

func DefaultPkgConfig() *PkgConfig {
	conf := &PkgConfig{IgnoreFuncBodies: false, AllowBinary: true, WithTestFiles: true}
	conf.Info = &types.Info{
		Uses:       make(map[*ast.Ident]types.Object),
		Defs:       make(map[*ast.Ident]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Scopes:     make(map[ast.Node]*types.Scope),
		Implicits:  make(map[ast.Node]types.Object),
		Instances:  make(map[*ast.Ident]types.Instance),
	}
	conf.XInfo = &types.Info{
		Uses:       make(map[*ast.Ident]types.Object),
		Defs:       make(map[*ast.Ident]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Scopes:     make(map[ast.Node]*types.Scope),
		Implicits:  make(map[ast.Node]types.Object),
		Instances:  make(map[*ast.Ident]types.Instance),
	}
	return conf
}

func sameNamed(n1, n2 *types.Named) bool {
	return n1 != nil && n2 != nil && n1.Origin().String() == n2.Origin().String()
}
