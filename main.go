// Copyright 2011-2017 visualfc <visualfc@gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/visualfc/gotools/astview"
	"github.com/visualfc/gotools/command"
	"github.com/visualfc/gotools/docview"
	"github.com/visualfc/gotools/finddecl"
	"github.com/visualfc/gotools/finddoc"
	"github.com/visualfc/gotools/goapi"
	"github.com/visualfc/gotools/gofmt"
	"github.com/visualfc/gotools/gopresent"
	"github.com/visualfc/gotools/gotest"
	"github.com/visualfc/gotools/jsonfmt"
	"github.com/visualfc/gotools/oracle"
	"github.com/visualfc/gotools/pkgs"
	"github.com/visualfc/gotools/runcmd"
	"github.com/visualfc/gotools/terminal"
	"github.com/visualfc/gotools/types"
)

func init() {
	command.Register(types.Command)
	command.Register(jsonfmt.Command)
	command.Register(finddoc.Command)
	command.Register(runcmd.Command)
	command.Register(docview.Command)
	command.Register(astview.Command)
	command.Register(gopresent.Command)
	command.Register(goapi.Command)
	command.Register(pkgs.Command)
	command.Register(oracle.Command)
	command.Register(gofmt.Command)
	command.Register(gotest.Command)
	command.Register(finddecl.Command)
	command.Register(terminal.Command)
}

func main() {
	command.AppName = "gotools"
	command.AppVersion = "1.9"
	command.AppInfo = "Go tools for liteide."
	command.Main()
}
