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
	Run:       runGotest,
	UsageLine: "gotest -f filename",
	Short:     "go test go filename",
	Long:      `go test go filename`,
}

var testFileName string

func init() {
	Command.Flag.StringVar(&testFileName, "f", "", "test go filename")
}

func runGotest(cmd *command.Command, args []string) error {
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
	testArgs = append(testArgs, "-v")
	testArgs = append(testArgs, testFiles...)

	command := exec.Command(gobin, testArgs...)
	command.Dir = pkg.Dir
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	return command.Run()
}
