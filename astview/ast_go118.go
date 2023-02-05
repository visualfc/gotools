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

func typeName(d *ast.TypeSpec, showTypeParams bool) string {
	if showTypeParams && d.TypeParams != nil {
		tparams := d.TypeParams
		var fs []string
		n := len(tparams.List)
		for i := 0; i < n; i++ {
			f := tparams.List[i]
			for _, name := range f.Names {
				fs = append(fs, name.String()+" "+types.ExprString(f.Type))
			}
		}
		return d.Name.String() + "[" + strings.Join(fs, ", ") + "]"
	}
	return d.Name.String()
}

func funcName(d *ast.FuncDecl, showTypeParams bool) string {
	if showTypeParams && d.Type.TypeParams != nil {
		tparams := d.Type.TypeParams
		var fs []string
		n := len(tparams.List)
		for i := 0; i < n; i++ {
			f := tparams.List[i]
			for _, name := range f.Names {
				fs = append(fs, name.String()+" "+types.ExprString(f.Type))
			}
		}
		return d.Name.String() + "[" + strings.Join(fs, ", ") + "]"
	}
	return d.Name.String()
}
