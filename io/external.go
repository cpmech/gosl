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

// Pipeline strings together the given exec.Cmd commands in a similar fashion
// to the Unix pipeline. Each command's standard output is connected to the
// standard input of the next command, and the output of the final command in
// the pipeline is returned, along with the collected standard error of all
// commands and the first error found (if any).
//
// by Kyle Lemons
//
// To provide input to the pipeline, assign an io.Reader to the first's Stdin.
func Pipeline(cmds ...*exec.Cmd) (pipeLineOutput, collectedStandardError []byte, pipeLineError error) {
	// Require at least one command
	if len(cmds) < 1 {
		return nil, nil, nil
	}

	// Collect the output from the command(s)
	var output bytes.Buffer
	var stderr bytes.Buffer

	last := len(cmds) - 1
	for i, cmd := range cmds[:last] {
		var err error
		// Connect each command's stdin to the previous command's stdout
		if cmds[i+1].Stdin, err = cmd.StdoutPipe(); err != nil {
			return nil, nil, err
		}
		// Connect each command's stderr to a buffer
		cmd.Stderr = &stderr
	}

	// Connect the output and error for the last command
	cmds[last].Stdout, cmds[last].Stderr = &output, &stderr

	// Start each command
	for _, cmd := range cmds {
		if err := cmd.Start(); err != nil {
			return output.Bytes(), stderr.Bytes(), err
		}
	}

	// Wait for each command to complete
	for _, cmd := range cmds {
		if err := cmd.Wait(); err != nil {
			return output.Bytes(), stderr.Bytes(), err
		}
	}

	// Return the pipeline output and the collected standard error
	return output.Bytes(), stderr.Bytes(), nil
}
