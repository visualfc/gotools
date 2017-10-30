// Copyright 2011-2015 visualfc <visualfc@gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package finddecl

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"

	"github.com/visualfc/gotools/command"
)

var Command = &command.Command{
	Run:       runFindDecl,
	UsageLine: "finddecl",
	Short:     "golang finddecl util",
	Long:      `golang finddecl util.`,
}
var (
	filePath string
	fileLine int
)

func init() {
	Command.Flag.StringVar(&filePath, "file", "", "file path")
	Command.Flag.IntVar(&fileLine, "line", -1, "file line")
}

func runFindDecl(cmd *command.Command, args []string) error {
	if len(filePath) == 0 || fileLine == -1 {
		cmd.Usage()
		return os.ErrInvalid
	}
	if !filepath.IsAbs(filePath) {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}
		filePath = filepath.Join(dir, filePath)
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, 0)
	if err != nil {
		return err
	}
	decl := findDecl(fset, f, fileLine)
	if decl == nil {
		fmt.Println("-")
		return errors.New("error find decl")
	}
	printDecl(fset, decl, fileLine)
	return nil
}

type Info struct {
	Type      string
	Name      string
	BeginLine int
	EndLine   int
}

func printDecl(fset *token.FileSet, decl ast.Decl, line int) {
	var tag string
	var name string

	tag = "-"
	name = "-"
	switch d := decl.(type) {
	case *ast.GenDecl:
		switch d.Tok {
		case token.IMPORT:
			tag = "import"
			if len(d.Specs) > 0 {
				if ts := d.Specs[0].(*ast.ImportSpec); ts != nil {
					name, _ = strconv.Unquote(ts.Path.Value)
				}
			}
		case token.TYPE:
			tag = "type"
			if len(d.Specs) > 0 {
				if ts := d.Specs[0].(*ast.TypeSpec); ts != nil {
					name = ts.Name.Name
					switch ts.Type.(type) {
					case *ast.StructType:
						tag = "struct"
					case *ast.InterfaceType:
						tag = "interface"
					default:
						tag = "type"
					}
				}
			}
		case token.VAR, token.CONST:
			tag = d.Tok.String()
			var testName string
			for _, ds := range d.Specs {
				if ts := ds.(*ast.ValueSpec); ts != nil {
					name = ts.Names[0].Name
					for _, n := range ts.Names {
						if line >= fset.Position(n.Pos()).Line && line <= fset.Position(n.End()).Line {
							testName = n.Name
							break
						}
					}
				}
			}
			if testName != "" {
				name = testName
			}
		default:
			tag = d.Tok.String()
		}
	case *ast.FuncDecl:
		tag = "func"
		name = d.Name.Name
	}
	fmt.Println(tag, name, fset.Position(decl.Pos()).Line, fset.Position(decl.End()).Line)
}

func findDecl(fset *token.FileSet, file *ast.File, line int) ast.Decl {
	for _, decl := range file.Decls {
		if line >= fset.Position(decl.Pos()).Line && line <= fset.Position(decl.End()).Line {
			return decl
		}
	}
	return nil
}
