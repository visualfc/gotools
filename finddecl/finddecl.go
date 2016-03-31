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
	printDecl(fset, decl)
	return nil
}

func printDecl(fset *token.FileSet, decl ast.Decl) {
	switch d := decl.(type) {
	case *ast.GenDecl:
		fmt.Println(d.Tok, fset.Position(d.Pos()).Line, fset.Position(d.End()).Line)
	case *ast.FuncDecl:
		fmt.Println("func", d.Name, fset.Position(d.Pos()).Line, fset.Position(d.End()).Line)
	}
}

func findDecl(fset *token.FileSet, file *ast.File, line int) ast.Decl {
	for _, decl := range file.Decls {
		if line >= fset.Position(decl.Pos()).Line && line <= fset.Position(decl.End()).Line {
			return decl
		}
	}
	return nil
}
