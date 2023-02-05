//go:build !go1.18
// +build !go1.18

package astview

import "go/ast"

func docBaseTypeName(typ ast.Expr, showAll bool) string {
	name, _ := recvTypeName(typ, showAll)
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

func typeName(ts *ast.TypeSpec, showTypeParams bool) string {
	return ts.Name.String()
}

func funcName(d *ast.FuncDecl, showTypeParams bool) string {
	return d.Name.String()
}
