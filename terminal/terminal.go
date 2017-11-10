// Copyright 2011-2017 visualfc <visualfc@gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package terminal

import (
	"os"
	"os/exec"

	"github.com/visualfc/gotools/command"
)

var Command = &command.Command{
	Run:         runTerminal,
	UsageLine:   "terminal [program_name arguments...]",
	Short:       "terminal [program]",
	Long:        `terminal program and arguments`,
	CustomFlags: true,
}

func runTerminal(cmd *command.Command, args []string) (err error) {
	var c *exec.Cmd
	if len(args) >= 1 {
		var carg []string
		if len(args) >= 2 {
			carg = append(carg, args[1:]...)
		}
		c = exec.Command(args[0], carg...)
	} else {
		shellCmd, shellArgs := GetShell()
		c = exec.Command(shellCmd, shellArgs...)
	}
	if c == nil {
		return os.ErrInvalid
	}
	err = Execute(c)

	return
}
