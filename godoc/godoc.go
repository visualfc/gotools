// Copyright 2018 visualfc <visualfc@gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package godoc

import (
	"os/exec"

	"github.com/visualfc/gotools/pkg/command"
	"github.com/visualfc/goversion"
)

var Command = &command.Command{
	Run:       runDoc,
	UsageLine: "godoc [pkg]",
	Short:     "golang doc print",
	Long:      "golang document print",
}

var (
	GoVer112 = goversion.GoVersion{1, 12, -1, 0, 0, ""}
)

func runDoc(cmd *command.Command, args []string) error {
	if len(args) < 1 {
		return nil
	}
	ver, _ := goversion.Installed()
	gocmd, err := exec.LookPath("go")
	if err != nil {
		return err
	}
	var godoc_html bool
	godoc, err := exec.LookPath("godoc")
	if err == nil {
		godoc_html = true
	}
	if ver.AfterOrEqual(GoVer112) {
		godoc_html = false
	}
	var command *exec.Cmd
	if godoc_html {
		command = exec.Command(godoc, "-html", args[0])
	} else {
		command = exec.Command(gocmd, "doc", "-all", args[0])
	}
	command.Stdin = cmd.Stdin
	command.Stdout = cmd.Stdout
	command.Stderr = cmd.Stderr
	return command.Run()
}
