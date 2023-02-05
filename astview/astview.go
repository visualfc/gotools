// Copyright 2011-2015 visualfc <visualfc@gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package astview

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/visualfc/gotools/pkg/command"
	"github.com/visualfc/gotools/pkg/pkgutil"
)

var Command = &command.Command{
	Run:       runAstView,
	UsageLine: "astview [-stdin] files...",
	Short:     "print go files astview",
	Long:      `print go files astview`,
}

var (
	astViewStdin          bool
	astViewShowEndPos     bool
	astViewShowTodo       bool
	astViewOutline        bool
	astViewShowTypeParams bool
	astViewSep            string
)

func init() {
	Command.Flag.BoolVar(&astViewStdin, "stdin", false, "input from stdin")
	Command.Flag.BoolVar(&astViewShowEndPos, "end", false, "show decl end pos")
	Command.Flag.BoolVar(&astViewShowTodo, "todo", false, "show todo list")
	Command.Flag.BoolVar(&astViewShowTypeParams, "tp", false, "show typeparams")
	Command.Flag.BoolVar(&astViewOutline, "outline", false, "set outline mode")
	Command.Flag.StringVar(&astViewSep, "sep", ",", "set output seperator")
}

func runAstView(cmd *command.Command, args []string) error {
	if len(args) == 0 {
		cmd.Usage()
		return os.ErrInvalid
	}
	if astViewStdin {
		view, err := NewFilePackageSource(args[0], cmd.Stdin, true)
		if err != nil {
			return err
		}
		view.PrintTree(cmd.Stdout)
	} else {
		if len(args) == 1 && astViewOutline {
			err := PrintFileOutline(args[0], cmd.Stdout, astViewSep, true)
			if err != nil {
				return err
			}
		} else {
			err := PrintFilesTree(args, cmd.Stdout, true)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

const (
	tag_package        = "p"
	tag_imports_folder = "+m"
	tag_import         = "mm"
	tag_type           = "t"
	tag_struct         = "s"
	tag_interface      = "i"
	tag_value          = "v"
	tag_const          = "c"
	tag_func           = "f"
	tag_value_folder   = "+v"
	tag_const_folder   = "+c"
	tag_func_folder    = "+f"
	tag_factor_folder  = "+tf"
	tag_type_method    = "tm"
	tag_type_factor    = "tf"
	tag_type_value     = "tv"
	tag_todo           = "b"
	tag_todo_folder    = "+b"
)

type PackageView struct {
	fset *token.FileSet
	pdoc *PackageDoc
	pkg  *ast.Package
	expr bool
}

var AllFiles []string

func (p *PackageView) posFileIndex(pos token.Position) int {
	var index = -1
	for i := 0; i < len(AllFiles); i++ {
		if AllFiles[i] == pos.Filename {
			index = i
			break
		}
	}
	if index == -1 {
		AllFiles = append(AllFiles, pos.Filename)
		index = len(AllFiles) - 1
	}
	return index
}

func (p *PackageView) posText(node ast.Node) string {
	pos := p.fset.Position(node.Pos())
	index := p.posFileIndex(pos)
	if astViewShowEndPos {
		end := p.fset.Position(node.End())
		return fmt.Sprintf("%d:%d:%d:%d:%d", index, pos.Line, pos.Column, end.Line, end.Column)
	}
	return fmt.Sprintf("%d:%d:%d", index, pos.Line, pos.Column)
}

func NewFilePackage(filename string) (*PackageView, error) {
	p := new(PackageView)
	p.fset = token.NewFileSet()
	file, err := parser.ParseFile(p.fset, filename, nil, parser.AllErrors)
	if file == nil {
		return nil, err
	}
	m := make(map[string]*ast.File)
	m[filename] = file
	pkg, err := ast.NewPackage(p.fset, m, nil, nil)
	if err != nil {
		return nil, err
	}
	p.pkg = pkg
	p.pdoc = NewPackageDoc(pkg, pkg.Name, true)
	return p, nil
}

func NewPackageView(pkg *ast.Package, fset *token.FileSet, expr bool) (*PackageView, error) {
	p := new(PackageView)
	p.fset = fset
	p.pkg = pkg
	p.pdoc = NewPackageDoc(pkg, pkg.Name, true)
	p.expr = expr
	return p, nil
}

func ParseFiles(fset *token.FileSet, filenames []string, mode parser.Mode) (pkgs map[string]*ast.Package, pkgsfiles []string, first error) {
	pkgs = make(map[string]*ast.Package)
	for _, filename := range filenames {
		if src, err := parser.ParseFile(fset, filename, nil, mode); src != nil {
			name := src.Name.Name
			pkg, found := pkgs[name]
			if !found {
				pkg = &ast.Package{
					Name:  name,
					Files: make(map[string]*ast.File),
				}
				pkgs[name] = pkg
			}
			pkg.Files[filename] = src
			pkgsfiles = append(pkgsfiles, filename)
		} else {
			first = err
			return
		}
	}
	return
}

func PrintFilesTree(filenames []string, w io.Writer, expr bool) error {
	fset := token.NewFileSet()
	mode := parser.AllErrors
	if astViewShowTodo {
		mode |= parser.ParseComments
	}
	pkgs, pkgsfiles, err := ParseFiles(fset, filenames, mode)
	if err != nil {
		return err
	}
	AllFiles = pkgsfiles
	for i := 0; i < len(AllFiles); i++ {
		fmt.Fprintf(w, "@%s\n", AllFiles[i])
	}
	for _, pkg := range pkgs {
		view, err := NewPackageView(pkg, fset, expr)
		if err != nil {
			return err
		}
		view.PrintTree(w)
	}
	return nil
}

func NewFilePackageSource(filename string, f io.Reader, expr bool) (*PackageView, error) {
	src, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	p := new(PackageView)
	p.fset = token.NewFileSet()
	p.expr = expr
	mode := parser.AllErrors
	if astViewShowTodo {
		mode |= parser.ParseComments
	}
	file, err := parser.ParseFile(p.fset, filename, src, mode)
	if err != nil {
		return nil, err
	}
	m := make(map[string]*ast.File)
	m[filename] = file
	pkg, err := ast.NewPackage(p.fset, m, nil, nil)
	if err != nil {
		return nil, err
	}

	p.pdoc = NewPackageDoc(pkg, pkg.Name, true)
	return p, nil
}

func (p *PackageView) out0(w io.Writer, level int, text ...string) {
	fmt.Fprintf(w, "%v%s%s\n", level, astViewSep, strings.Join(text, astViewSep))
}
func (p *PackageView) out1(w io.Writer, level int, pos ast.Node, text ...string) {
	fmt.Fprintf(w, "%v%s%s%s%s\n", level, astViewSep, strings.Join(text, astViewSep), astViewSep, p.posText(pos))
}
func (p *PackageView) out2(w io.Writer, level int, pos ast.Node, expr ast.Expr, text ...string) {
	fmt.Fprintf(w, "%v%s%s%s%s@%v\n", level, astViewSep, strings.Join(text, astViewSep), astViewSep, p.posText(pos), types.ExprString(expr))
}
func (p *PackageView) out2s(w io.Writer, level int, pos ast.Node, expr string, text ...string) {
	fmt.Fprintf(w, "%v%s%s%s%s@%v\n", level, astViewSep, strings.Join(text, astViewSep), astViewSep, p.posText(pos), expr)
}

func (p *PackageView) printFuncsHelper(w io.Writer, funcs []*FuncDoc, level int, tag string, tag_folder string) {
	for _, f := range funcs {
		if p.expr {
			p.out2(w, level, f.Decl, f.Decl.Type, tag, f.Name)
		} else {
			p.out1(w, level, f.Decl, tag, f.Name)
		}
	}
}

func (p *PackageView) PrintVars(w io.Writer, vars []*ValueDoc, level int, tag string, tag_folder string) {
	if len(tag_folder) > 0 && len(vars) > 0 {
		if tag_folder == tag_value_folder {
			p.out0(w, level, tag_folder, "Variables")
		} else if tag_folder == tag_const_folder {
			p.out0(w, level, tag_folder, "Constants")
		}
		level++
	}
	for _, v := range vars {
		if v.Decl == nil {
			continue
		}
		for _, s := range v.Decl.Specs {
			if m, ok := s.(*ast.ValueSpec); ok {
				for i := 0; i < len(m.Names); i++ {
					if p.expr && m.Type != nil {
						p.out2(w, level, m, m.Type, tag, m.Names[i].String())
					} else {
						p.out1(w, level, m, tag, m.Names[i].String())
					}
				}
			}
		}
	}
}
func (p *PackageView) PrintTypes(w io.Writer, types []*TypeDoc, level int) {
	for _, d := range types {
		if d.Decl == nil {
			continue
		}
		typespec := d.Decl.Specs[0].(*ast.TypeSpec)
		var tag = tag_type
		if _, ok := typespec.Type.(*ast.InterfaceType); ok {
			tag = tag_interface
		} else if _, ok := typespec.Type.(*ast.StructType); ok {
			tag = tag_struct
		}
		p.out1(w, level, d.Decl, tag, typeName(typespec, astViewShowTypeParams))
		p.printFuncsHelper(w, d.Funcs, level+1, tag_type_factor, "")
		p.printFuncsHelper(w, d.Methods, level+1, tag_type_method, "")
		p.PrintTypeFields(w, d.Decl, level+1)
		//p.PrintVars(w, d.Consts, level+1, tag_const, "")
		//p.PrintVars(w, d.Vars, level+1, tag_value, "")
	}
}

func (p *PackageView) PrintTypeFields(w io.Writer, decl *ast.GenDecl, level int) {
	spec, ok := decl.Specs[0].(*ast.TypeSpec)
	if ok == false {
		return
	}
	switch d := spec.Type.(type) {
	case *ast.StructType:
		for _, list := range d.Fields.List {
			if list.Names == nil {
				continue
			}
			for _, m := range list.Names {
				if list.Type != nil {
					p.out2(w, level, m, list.Type, tag_type_value, m.Name)
				} else {
					p.out1(w, level, m, tag_type_value, m.Name)
				}
			}
		}
	case *ast.InterfaceType:
		for _, list := range d.Methods.List {
			if list.Names == nil {
				continue
			}
			for _, m := range list.Names {
				p.out1(w, level, m, tag_type_method, m.Name)
			}
		}
	}
}

func (p *PackageView) PrintHeader(w io.Writer, level int) {
	p.out0(w, level, tag_package, p.pdoc.PackageName)
}

func (p *PackageView) PrintImports(w io.Writer, level int, tag, tag_folder string) {
	if tag_folder != "" && len(p.pdoc.Imports) > 0 {
		p.out0(w, level, tag_folder, "Imports")
		level++
	}
	var parentPkg *pkgutil.Package
	if pkgutil.IsVendorExperiment() && p.pkg != nil {
		for filename, _ := range p.pkg.Files {
			if !filepath.IsAbs(filename) {
				name, err := filepath.Abs(filename)
				if err == nil {
					filename = name
				}
			}
			parentPkg = pkgutil.ImportFile(filename)
			break
		}
	}
	for _, name := range p.pdoc.Imports {
		vname := "\"" + name + "\""
		var ps []string
		for _, file := range p.pkg.Files {
			for _, v := range file.Imports {
				if v.Path.Value == vname {
					ps = append(ps, p.posText(v))
				}
			}
		}
		if parentPkg != nil {
			name, _ = pkgutil.VendoredImportPath(parentPkg, name)
		}
		p.out0(w, level, tag, name, strings.Join(ps, ";"))
	}
}

func (p *PackageView) PrintFuncs(w io.Writer, level int, tag_folder string) {
	hasFolder := false
	if len(p.pdoc.Funcs) > 0 || len(p.pdoc.Factorys) > 0 {
		hasFolder = true
	}
	if !hasFolder {
		return
	}
	if len(tag_folder) > 0 {
		p.out0(w, level, tag_folder, "Functions")
		level++
	}
	p.printFuncsHelper(w, p.pdoc.Factorys, level, tag_type_factor, tag_func_folder)
	p.printFuncsHelper(w, p.pdoc.Funcs, level, tag_func, tag_func_folder)
}

func (p *PackageView) PrintTodos(w io.Writer, level int, tag, tag_folder string) {
	hasFolder := false
	if len(p.pdoc.Todos) > 0 {
		hasFolder = true
	}
	if !hasFolder {
		return
	}
	if len(tag_folder) > 0 {
		p.out0(w, level, tag_folder, "TodoList")
		level++
	}
	for _, todo := range p.pdoc.Todos {
		c := todo.Comments.List[0]
		p.out2s(w, level, c, todo.Text, tag, todo.Tag)
	}
}

func (p *PackageView) PrintPackage(w io.Writer, level int) {
	p.PrintHeader(w, level)
	level++
	p.PrintImports(w, level, tag_import, tag_imports_folder)
	p.PrintVars(w, p.pdoc.Vars, level, tag_value, tag_value_folder)
	p.PrintVars(w, p.pdoc.Consts, level, tag_const, tag_const_folder)
	p.PrintFuncs(w, level, tag_func_folder)
	p.PrintTypes(w, p.pdoc.Types, level)
	p.PrintTodos(w, level, tag_todo, tag_todo_folder)
}

// level,tag,pos@info
func (p *PackageView) PrintTree(w io.Writer) {
	p.PrintPackage(w, 0)
}

// level,tag,pos@info
func PrintFileOutline(filename string, w io.Writer, sep string, showexpr bool) error {
	fset := token.NewFileSet()
	mode := parser.AllErrors
	if astViewShowTodo {
		mode |= parser.ParseComments
	}
	f, err := parser.ParseFile(fset, filename, nil, mode)
	if err != nil {
		return err
	}
	posText := func(node ast.Node) string {
		pos := fset.Position(node.Pos())
		if astViewShowEndPos {
			end := fset.Position(node.End())
			return fmt.Sprintf("%d:%d:%d:%d:%d", 0, pos.Line, pos.Column, end.Line, end.Column)
		}
		return fmt.Sprintf("%d:%d:%d", 0, pos.Line, pos.Column)
	}

	out0 := func(level int, text ...string) {
		fmt.Fprintf(w, "%v%s%s\n", level, sep, strings.Join(text, sep))
	}
	out1 := func(level int, pos ast.Node, text ...string) {
		fmt.Fprintf(w, "%v%s%s%s%s\n", level, sep, strings.Join(text, sep), sep, posText(pos))
	}
	out2 := func(level int, pos ast.Node, expr ast.Expr, text ...string) {
		fmt.Fprintf(w, "%v%s%s%s%s@%v\n", level, sep, strings.Join(text, sep), sep, posText(pos), types.ExprString(expr))
	}
	out2s := func(level int, pos ast.Node, expr string, text ...string) {
		fmt.Fprintf(w, "%v%s%s%s%s@%v\n", level, sep, strings.Join(text, sep), sep, posText(pos), expr)
	}

	fmt.Fprintf(w, "@%s\n", filename)
	level := 0
	out1(level, f.Name, tag_package, f.Name.Name)
	// level++
	if len(f.Imports) > 0 {
		sort.Slice(f.Imports, func(i, j int) bool {
			return f.Imports[i].Pos() < f.Imports[j].Pos()
		})
		out0(level, tag_imports_folder, "Imports")
		level++
		for _, imp := range f.Imports {
			path, _ := strconv.Unquote(imp.Path.Value)
			out1(level, imp, tag_import, path)
		}
		level--
	}

	sort.Slice(f.Decls, func(i, j int) bool {
		return f.Decls[i].Pos() < f.Decls[j].Pos()
	})
	for _, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			switch d.Tok {
			case token.IMPORT:
			case token.TYPE:
				for _, spec := range d.Specs {
					ts := spec.(*ast.TypeSpec)
					switch t := ts.Type.(type) {
					case *ast.StructType:
						out2(level, ts, t, tag_struct, typeName(ts, astViewShowTypeParams))
						n := len(t.Fields.List)
						if n > 0 {
							level++
							for i := 0; i < n; i++ {
								f := t.Fields.List[i]
								for _, name := range f.Names {
									out2(level, name, f.Type, tag_type_value, name.String())
								}
							}
							level--
						}
					case *ast.InterfaceType:
						out2(level, ts, t, tag_interface, typeName(ts, astViewShowTypeParams))
						n := len(t.Methods.List)
						if n > 0 {
							level++
							for i := 0; i < n; i++ {
								f := t.Methods.List[i]
								for _, name := range f.Names {
									out2(level, name, f.Type, tag_type_method, name.String())
								}
							}
							level--
						}
					default:
						out2(level, ts, t, tag_type, typeName(ts, astViewShowTypeParams))
					}
				}
			case token.CONST:
				for _, spec := range d.Specs {
					vs := spec.(*ast.ValueSpec)
					for i, name := range vs.Names {
						if vs.Values == nil {
							out1(level, name, tag_const, name.String())
						} else {
							out2(level, name, vs.Values[i], tag_const, name.String())
						}
					}
				}
			case token.VAR:
				for _, spec := range d.Specs {
					vs := spec.(*ast.ValueSpec)
					for _, name := range vs.Names {
						if vs.Type == nil {
							out1(level, name, tag_value, name.String())
						} else {
							out2(level, name, vs.Type, tag_value, name.String())
						}
					}
				}
			}
		case *ast.FuncDecl:
			if d.Recv != nil {
				name, star := recvTypeName(d.Recv.List[0].Type, true)
				if star {
					name = "*" + name
				}
				if astViewShowTypeParams {
					out2(level, d, d.Type, tag_func, types.ExprString(d.Recv.List[0].Type))
				} else {
					out2(level, d, d.Type, tag_func, "("+name+")."+d.Name.String())
				}
			} else {
				out2(level, d, d.Type, tag_func, d.Name.String())
			}
		}
	}

	if astViewShowTodo {
		var todoList []*TodoDoc
		for _, c := range f.Comments {
			text := c.List[0].Text
			if m := todo_markers.FindStringSubmatchIndex(text); m != nil {
				todoList = append(todoList, &TodoDoc{text[m[2]:m[3]], text[m[2]:], c})
			}
		}
		if len(todoList) > 0 {
			out0(level, tag_todo_folder, "TodoList")
			level++
			for _, todo := range todoList {
				c := todo.Comments.List[0]
				out2s(level, c, todo.Text, tag_todo, todo.Tag)
			}
			level--
		}
	}
	return nil
}
