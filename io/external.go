// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package io

import (
	"bytes"
	"os/exec"

	"github.com/cpmech/gosl/chk"
)

// RunCmd runs external command
func RunCmd(verbose bool, cmd string, args ...string) (string, error) {
	exe := exec.Command(cmd, args...)
	var out, serr bytes.Buffer
	exe.Stdout = &out
	exe.Stderr = &serr
	err := exe.Run()
	if err != nil {
		Pf("\n")
		return Sf("%s\n%s\n", out.String(), serr.String()), chk.Err("%q command failed:", cmd)
	}
	if verbose {
		Pf("\n")
		Pf("%s", out.String())
	}
	return out.String(), nil
}

// CopyFileOver copies file (Linux only with cp), overwriting if it exists already
func CopyFileOver(destination, source string) (err error) {
	_, err = RunCmd(false, "cp", "-rf", source, destination)
	return
}
