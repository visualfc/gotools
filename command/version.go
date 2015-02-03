// Copyright 2011-2015 visualfc <visualfc@gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package command

import (
	"fmt"
	"runtime"
)

func init() {
	Register(cmdVersion)
}

var AppVersion string = "1.0"

var cmdVersion = &Command{
	Run:       runVersion,
	UsageLine: "version",
	Short:     "print tool version",
	Long:      `Version prints the version.`,
}

func runVersion(cmd *Command, args []string) {
	if len(args) != 0 {
		cmd.Usage()
	}

	fmt.Printf("%s version %s [%s %s/%s]\n", AppName, AppVersion, runtime.Version(), runtime.GOOS, runtime.GOARCH)
}
