// Copyright 2011-2015 visualfc <visualfc@gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package oracle

import (
	"fmt"
	"go/build"
	"os"
	"runtime"

	"github.com/visualfc/gotools/command"
	"github.com/visualfc/gotools/oracle/oracle"
)

//The mode argument determines the query to perform:

//	callees	  	show possible targets of selected function call
//	callers	  	show possible callers of selected function
//	callgraph 	show complete callgraph of program
//	callstack 	show path from callgraph root to selected function
//	describe  	describe selected syntax: definition, methods, etc
//	freevars  	show free variables of selection
//	implements	show 'implements' relation for selected type
//	peers     	show send/receive corresponding to selected channel op
//	referrers 	show all refs to entity denoted by selected identifier
//	what		show basic information about the selected syntax node

var Command = &command.Command{
	Run:       runOracle,
	UsageLine: "oracle",
	Short:     "golang oracle util",
	Long:      `golang oracle util.`,
}

var (
	oraclePos     string
	oracleReflect bool
)

func init() {
	Command.Flag.StringVar(&oraclePos, "pos", "", "filename:#offset")
	Command.Flag.BoolVar(&oracleReflect, "reflect", false, "Analyze reflection soundly (slow).")
}

func runOracle(cmd *command.Command, args []string) error {
	if len(args) < 2 {
		cmd.Usage()
		return os.ErrInvalid
	}
	if os.Getenv("GOMAXPROCS") == "" {
		n := runtime.NumCPU()
		if n < 4 {
			n = 4
		}
		runtime.GOMAXPROCS(n)
	}
	mode := args[0]
	args = args[1:]
	//	if args[0] == "." {
	//		pkgPath, err := os.Getwd()
	//		if err != nil {
	//			log.Fatalln(err)
	//		}
	//		pkg, err := build.Default.ImportDir(pkgPath, 0)
	//		if err != nil {
	//			log.Fatalln(err)
	//		}
	//		args = pkg.GoFiles
	//		//log.Println(pkg.ImportPath)
	//		if pkg.ImportPath != "." && pkg.ImportPath != "" {
	//			args = []string{pkg.ImportPath}
	//		}
	//	}
	query := oracle.Query{
		Mode:       mode,
		Pos:        oraclePos,
		Build:      &build.Default,
		Scope:      args,
		PTALog:     nil,
		Reflection: oracleReflect,
	}

	if err := oracle.Run(&query); err != nil {
		fmt.Fprintf(os.Stderr, "oracle: %s.\n", err)
		return err
	}

	if mode == "referrers" {
		ref := query.Serial().Referrers
		if ref != nil {
			fmt.Fprintln(os.Stdout, ref.Desc)
			fmt.Fprintln(os.Stdout, ref.ObjPos)
			for _, v := range ref.Refs {
				fmt.Fprintln(os.Stdout, v)
			}
		}
	} else {
		query.WriteTo(os.Stdout)
	}
	return nil
}
