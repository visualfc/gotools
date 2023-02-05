//go:build !go1.18
// +build !go1.18

package astview

import "go/ast"

func docBaseTypeName(typ ast.Expr, showAll bool) string {
	name, _ := recvTypeName(typ)
	return name
}

func recvTypeName(typ ast.Expr, showAll bool) (string, bool) {
	switch t := typ.(type) {
	case *ast.Ident:
		// if the type is not exported, the effect to
		// a client is as if there were no type name
		if showAll || t.IsExported() {
			return t.Name, false
		}
	case *ast.StarExpr:
		return docBaseTypeName(t.X, showAll), true
	}
	return "", false
}
