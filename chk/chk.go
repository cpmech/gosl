// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package chk contains functions for checking and testing computations
package chk

import (
	"fmt"
	"testing"
)

var (
	// AssertOn activates or deactivates asserts
	AssertOn = true

	// Verbose turn on verbose mode
	Verbose = false

	// ColorsOn turn on use of colours on console
	ColorsOn = true
)

// ET checks error in tests
func ET(tst *testing.T, err error) (failed bool) {
	if err != nil {
		tst.Errorf("%v\n", err)
		return true
	}
	return false
}

// EP checks error and panic
func EP(err error) {
	if err != nil {
		Panic("%v\n", err)
	}
}

// PanicSimple panicks without calling CallerInfo
func PanicSimple(msg string, prm ...interface{}) {
	panic(fmt.Sprintf(msg, prm...))
}

// Panic panicks
func Panic(msg string, prm ...interface{}) {
	CallerInfo(4)
	CallerInfo(3)
	CallerInfo(2)
	panic(fmt.Sprintf(msg, prm...))
}

// PanicErr panics if err != nil
func PanicErr(err error) {
	if err != nil {
		Panic("Error occurred:\n%v\n", err)
	}
}

// Err returns a new error
func Err(msg string, prm ...interface{}) error {
	return fmt.Errorf(msg, prm...)
}
