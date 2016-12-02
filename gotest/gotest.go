// Copyright 2011-2015 visualfc <visualfc@gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gotest

import (
	"fmt"
	"go/build"
	"os"
	"os/exec"
	"strings"

	"github.com/visualfc/gotools/command"
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

//func init() {
//	Command.Flag.StringVar(&testFileName, "f", "", "test go filename")
//}

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

	pkg, err := build.ImportDir(".", 0)
	if err != nil {
		fmt.Println("import dir error", err)
		return err
	}

	var testFiles []string

	for _, file := range pkg.XTestGoFiles {
		if file == testFileName {
			testFiles = append(testFiles, file)
			break
		}
	}
	for _, file := range pkg.TestGoFiles {
		if file == testFileName {
			testFiles = append(testFiles, pkg.GoFiles...)
			testFiles = append(testFiles, file)
			break
		}
	}

	gobin, err := exec.LookPath("go")
	if err != nil {
		fmt.Println("error look go", err)
		return err
	}

	var testArgs []string
	testArgs = append(testArgs, "test")
	if len(args) > 0 {
		testArgs = append(testArgs, args...)
	}
	testArgs = append(testArgs, testFiles...)

	command := exec.Command(gobin, testArgs...)
	command.Dir = pkg.Dir
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	return command.Run()
}
