// Copyright 2011-2015 visualfc <visualfc@gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package godiff

import (
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/pmezard/go-difflib/difflib"
)

func UnifiedDiffLines(a []string, b []string) (string, error) {
	diff := difflib.UnifiedDiff{
		A:        a,
		B:        b,
		FromFile: "Original",
		ToFile:   "Current",
		Context:  3,
		Eol:      "\n",
	}
	return difflib.GetUnifiedDiffString(diff)
}

func UnifiedDiffString(a string, b string) (string, error) {
	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(a),
		B:        difflib.SplitLines(b),
		FromFile: "Original",
		ToFile:   "Current",
		Context:  3,
		Eol:      "\n",
	}
	return difflib.GetUnifiedDiffString(diff)
}

func UnifiedDiffBytesByCmd(b1, b2 []byte) (data []byte, err error) {
	f1, err := ioutil.TempFile("", "godiff")
	if err != nil {
		return
	}
	defer os.Remove(f1.Name())
	defer f1.Close()

	f2, err := ioutil.TempFile("", "godiff")
	if err != nil {
		return
	}
	defer os.Remove(f2.Name())
	defer f2.Close()

	f1.Write(b1)
	f2.Write(b2)

	data, err = exec.Command("diff", "-u", f1.Name(), f2.Name()).CombinedOutput()
	if len(data) > 0 {
		// diff exits with a non-zero status when the files don't match.
		// Ignore that failure as long as we get output.
		err = nil
	}
	return
}
