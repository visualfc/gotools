// Copyright 2011-2018 visualfc <visualfc@gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package debugflags

import (
	"runtime"

	"github.com/go-delve/delve/pkg/goversion"
	"github.com/visualfc/gotools/pkg/command"
)

var Command = &command.Command{
	Run:         runDebugFlags,
	UsageLine:   "debugflags",
	Short:       "print go debug flags",
	Long:        `print go debug flags`,
	CustomFlags: true,
}

func runDebugFlags(cmd *command.Command, args []string) error {
	var buildFlagsDefault string
	if runtime.GOOS == "windows" {
		ver, _ := goversion.Installed()
		if ver.Major > 0 && !ver.AfterOrEqual(goversion.GoVersion{1, 9, -1, 0, 0, ""}) {
			// Work-around for https://github.com/golang/go/issues/13154
			buildFlagsDefault = "-ldflags='-linkmode internal'"
		}
	}
	flags := debugflags()
	if buildFlagsDefault != "" {
		flags += " " + buildFlagsDefault
	}
	cmd.Println(flags)
	return nil
}

// copy from github.com/derekparker/delve/cmd/dlv/cmds/commands.go
func debugflags() string {
	// after go1.9 building with -gcflags='-N -l' and -a simultaneously works.
	// after go1.10 specifying -a is unnecessary because of the new caching strategy, but we should pass -gcflags=all=-N -l to have it applied to all packages
	// see https://github.com/golang/go/commit/5993251c015dfa1e905bdf44bdb41572387edf90

	ver, _ := goversion.Installed()
	var flags string
	switch {
	case ver.Major < 0 || ver.AfterOrEqual(goversion.GoVersion{1, 10, -1, 0, 0, ""}):
		flags = "-gcflags=\"all=-N -l\""
	case ver.AfterOrEqual(goversion.GoVersion{1, 9, -1, 0, 0, ""}):
		flags = "-gcflags=\"-N -l\" -a"
	default:
		flags = "-gcflags=\"-N -l\" -a"
	}
	return flags
}
