//go:build go1.18
// +build go1.18

package astview

import (
	"go/ast"
	"go/types"
	"strings"
)

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
	case *ast.IndexExpr:
		return docBaseTypeName(t.X, showAll), false
	case *ast.IndexListExpr:
		return docBaseTypeName(t.X, showAll), false
	}
	return "", false
}

func typeName(ts *ast.TypeSpec, showTypeParams bool) string {
	if showTypeParams && ts.TypeParams != nil {
		var fs []string
		n := len(ts.TypeParams.List)
		for i := 0; i < n; i++ {
			f := ts.TypeParams.List[i]
			for _, name := range f.Names {
				fs = append(fs, name.String()+" "+types.ExprString(f.Type))
			}
		}
		return ts.Name.String() + "[" + strings.Join(fs, ",") + "]"
	}
	return ts.Name.String()
}
