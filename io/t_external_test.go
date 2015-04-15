// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package io

import (
	"os/exec"
	"testing"

	"github.com/cpmech/gosl/chk"
)

func Test_runcmd01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("runcmd01")

	ver := false
	Pforan("running 'ls -la'\n")
	out, err := RunCmd(ver, "ls", "-la")
	Pfblue2("\noutput:\n%v\n", out)
	if err != nil {
		tst.Errorf("error: %v\n", err)
		return
	}
}

func Test_pipe01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("runcmd01")

	Pforan("running pipe\n")

	// find $DIR -type f # Find all files
	dir := "."
	find := exec.Command("find", dir, "-type", "f")

	// | grep -v '/[._]' # Ignore hidden/temporary files
	egrep := exec.Command("egrep", "-v", `/[._]`)

	// | sort -t. -k2 # Sort by file extension
	sort := exec.Command("sort", "-t.", "-k2")

	output, stderr, err := Pipeline(find, egrep, sort)
	Pfblue2("\noutput:\n%v\n", string(output))
	Pfcyan("stderr:\n%v\n", string(stderr))
	if err != nil {
		tst.Errorf("error: %v\n", err)
		return
	}
}
