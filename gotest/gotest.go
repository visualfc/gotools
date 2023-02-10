// Copyright 2011-2015 visualfc <visualfc@gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gotest

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/visualfc/gotools/pkg/command"
	"golang.org/x/tools/go/buildutil"
)

var Command = &command.Command{
	Run:         runGotest,
	UsageLine:   "gotest -f filename [build/test flags]",
	Short:       "go test go filename",
	Long:        `go test go filename`,
	CustomFlags: true,
}

var testFileName string
var testFileArgs string

func init() {
	//Command.Flag.StringVar(&testFileName, "f", "", "test go filename")
	ApplyBuildTags()
}

func runGotest(cmd *command.Command, args []string) error {
	index := -1
	for n, arg := range args {
		if arg == "-f" {
			index = n
			break
		}
	}
	if index >= 0 && index < len(args) {
		testFileName = args[index+1]
		var r []string
		r = append(r, args[0:index]...)
		r = append(r, args[index+2:]...)
		args = r
	}

	if testFileName == "" {
		cmd.Usage()
		return os.ErrInvalid
	}
	if !strings.HasSuffix(testFileName, "_test.go") {
		fmt.Println("The test filename must xxx_test.go")
		return os.ErrInvalid
	}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, testFileName, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var fnList []string
	for _, decl := range f.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			name := fn.Name.Name
			if strings.HasPrefix(name, "Test") || strings.HasPrefix(name, "Benchmark") {
				fnList = append(fnList, name)
			}
		}
	}
	if len(fnList) == 0 {
		return fmt.Errorf("testing: warning: no tests to run")
	}

	if filepath.IsAbs(testFileName) {
		dir, _ := filepath.Split(testFileName)
		os.Chdir(dir)
	}

	gobin, err := exec.LookPath("go")
	if err != nil {
		return fmt.Errorf("error lookup go: %v", err)
	}

	var testArgs []string
	testArgs = append(testArgs, "test")
	if len(args) > 0 {
		testArgs = append(testArgs, args...)
	}
	testArgs = append(testArgs, "-run", fmt.Sprintf("^(%v)$", strings.Join(fnList, "|")))

	command := exec.Command(gobin, testArgs...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	return command.Run()
}

func ApplyBuildTags() {
	nexttag := false
	for _, arg := range os.Args[1:] {
		if nexttag {
			var tags buildutil.TagsFlag
			tags.Set(arg)

			build.Default.BuildTags = tags
			nexttag = false
			continue
		}
		if arg == "-tags" {
			nexttag = true
		}
	}
}
