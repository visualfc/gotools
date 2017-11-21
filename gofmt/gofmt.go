// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gofmt

import (
	"bytes"
	"fmt"
	"go/scanner"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/visualfc/gotools/command"
	"github.com/visualfc/gotools/godiff"
	"golang.org/x/tools/imports"
)

var Command = &command.Command{
	Run:       runGofmt,
	UsageLine: "gofmt [flags] [path ...]",
	Short:     "gofmt formats Go source.",
	Long:      `gofmt formats Go source`,
}

var (
	gofmtList         bool
	gofmtWrite        bool
	gofmtDiff         bool
	gofmtAllErrors    bool
	gofmtFixImports   bool
	gofmtSortImports  bool
	gofmtUseGodiffLib bool

	// layout control
	gofmtComments  bool
	gofmtTabWidth  int
	gofmtTabIndent bool
)

//func init
func init() {
	Command.Flag.BoolVar(&gofmtList, "l", false, "list files whose formatting differs from goimport's")
	Command.Flag.BoolVar(&gofmtWrite, "w", false, "write result to (source) file instead of stdout")
	Command.Flag.BoolVar(&gofmtDiff, "d", false, "display diffs instead of rewriting files")
	Command.Flag.BoolVar(&gofmtAllErrors, "e", false, "report all errors (not just the first 10 on different lines)")
	Command.Flag.BoolVar(&gofmtFixImports, "fiximports", false, "updates Go import lines, adding missing ones and removing unreferenced ones")
	Command.Flag.BoolVar(&gofmtSortImports, "sortimports", false, "sort Go import lines use goimports style")
	Command.Flag.BoolVar(&gofmtUseGodiffLib, "godiff", true, "diff use godiff library")

	// layout control
	Command.Flag.BoolVar(&gofmtComments, "comments", true, "print comments")
	Command.Flag.IntVar(&gofmtTabWidth, "tabwidth", 8, "tab width")
	Command.Flag.BoolVar(&gofmtTabIndent, "tabs", true, "indent with tabs")
}

var (
	fileSet  = token.NewFileSet() // per process FileSet
	exitCode = 0

	initModesOnce sync.Once // guards calling initModes
	//parserMode    parser.Mode
	//printerMode   printer.Mode
	options *imports.Options
)

func report(err error) {
	scanner.PrintError(os.Stderr, err)
	exitCode = 2
}

func runGofmt(cmd *command.Command, args []string) error {
	runtime.GOMAXPROCS(runtime.NumCPU())

	if gofmtTabWidth < 0 {
		fmt.Fprintf(os.Stderr, "negative tabwidth %d\n", gofmtTabWidth)
		exitCode = 2
		os.Exit(exitCode)
		return os.ErrInvalid
	}

	if gofmtFixImports {
		gofmtSortImports = true
	}

	options = &imports.Options{
		FormatOnly: !gofmtFixImports,
		TabWidth:   gofmtTabWidth,
		TabIndent:  gofmtTabIndent,
		Comments:   gofmtComments,
		AllErrors:  gofmtAllErrors,
		Fragment:   true,
	}

	if len(args) == 0 {
		if err := processFile("<standard input>", os.Stdin, os.Stdout, true); err != nil {
			report(err)
		}
	} else {
		for _, path := range args {
			switch dir, err := os.Stat(path); {
			case err != nil:
				report(err)
			case dir.IsDir():
				walkDir(path)
			default:
				if err := processFile(path, nil, os.Stdout, false); err != nil {
					report(err)
				}
			}
		}
	}
	os.Exit(exitCode)
	return nil
}

func isGoFile(f os.FileInfo) bool {
	// ignore non-Go files
	name := f.Name()
	return !f.IsDir() && !strings.HasPrefix(name, ".") && strings.HasSuffix(name, ".go")
}

func processFile(filename string, in io.Reader, out io.Writer, stdin bool) error {
	if in == nil {
		f, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer f.Close()
		in = f
	}

	src, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}

	res, err := imports.Process(filename, src, options)
	if err != nil {
		return err
	}

	if !bytes.Equal(src, res) {
		// formatting has changed
		if gofmtList {
			fmt.Fprintln(out, filename)
		}
		if gofmtWrite {
			err = ioutil.WriteFile(filename, res, 0)
			if err != nil {
				return err
			}
		}
		if gofmtDiff {
			if gofmtUseGodiffLib {
				data, err := godiff.UnifiedDiffString(string(src), string(res))
				if err != nil {
					return fmt.Errorf("computing diff: %s", err)
				}
				fmt.Printf("diff %s gofmt/%s\n", filename, filename)
				out.Write([]byte(data))
			} else {
				data, err := godiff.UnifiedDiffBytesByCmd(src, res)
				if err != nil {
					return fmt.Errorf("computing diff: %s", err)
				}
				fmt.Printf("diff %s gofmt/%s\n", filename, filename)
				out.Write(data)
			}
		}
	}

	if !gofmtList && !gofmtWrite && !gofmtDiff {
		_, err = out.Write(res)
	}

	return err
}

func visitFile(path string, f os.FileInfo, err error) error {
	if err == nil && isGoFile(f) {
		err = processFile(path, nil, os.Stdout, false)
	}
	if err != nil {
		report(err)
	}
	return nil
}

func walkDir(path string) {
	filepath.Walk(path, visitFile)
}
