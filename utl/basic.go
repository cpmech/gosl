// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// utl (utilities) implements auxiliary functions for a number of tasks
// (printing, parsing, handling files and directories, etc.),
// including for making the transition from NumPy to Go easier.
package utl

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"testing"
)

var (
	UseColors = true
	Tsilent   = false // silent during testing, except with errors
)

// Sf wraps Sprintf
func Sf(msg string, prm ...interface{}) string {
	return fmt.Sprintf(msg, prm...)
}

// Ff wraps Fprintf
func Ff(b *bytes.Buffer, msg string, prm ...interface{}) {
	fmt.Fprintf(b, msg, prm...)
}

// Err returns a new error
func Err(msg string, prm ...interface{}) error {
	if UseColors {
		return errors.New(Sf("[1;35m"+msg+"[0m", prm...))
	}
	return errors.New(Sf(msg, prm...))
}

// PanicSimple panicks without calling CallerInfo
func PanicSimple(msg string, prm ...interface{}) {
	if UseColors {
		panic(Sf("[1;35m"+msg+"[0m", prm...))
	}
	panic(Sf(msg, prm...))
}

// Panic panicks
func Panic(msg string, prm ...interface{}) {
	CallerInfo(4)
	CallerInfo(3)
	CallerInfo(2)
	if UseColors {
		panic(Sf("[1;35m"+msg+"[0m", prm...))
	}
	panic(Sf(msg, prm...))
}

// Catch catches a panic
func Catch(tst *testing.T) {
	if err := recover(); err != nil {
		tst.Error("[1;31mSome error has happened:[0m\n", err)
	}
}

// PrintTime prints simulation time
func PrintTime(t float64) {
	if Tsilent {
		return
	}
	fmt.Printf("%13.8f\b\b\b\b\b\b\b\b\b\b\b\b\b", t)
}

// PrintTime prints simulation time (long format)
func PrintTimeLong(t float64) {
	if Tsilent {
		return
	}
	fmt.Printf("%20.8f\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b", t)
}

// Digits returns the nubmer of digits
func Digits(maxint int) (ndigits int, format string) {
	ndigits = int(math.Log10(float64(maxint))) + 1
	format = Sf("%%%dd", ndigits)
	return
}

// Expon returns the exponent
func Expon(val float64) (ndigits int) {
	if val == 0.0 {
		return
	}
	ndigits = int(math.Log10(math.Abs(val)))
	return
}

// Atob converts string to bool
func Atob(val string) (bres bool) {
	if strings.ToLower(val) == "true" {
		return true
	}
	if strings.ToLower(val) == "false" {
		return false
	}
	res, err := strconv.Atoi(val)
	if err != nil {
		Panic("utl.Atob: cannot parse string representing integer: %s", val)
	}
	if res != 0 {
		bres = true
	}
	return
}

// Atoi converts string to integer
func Atoi(val string) (res int) {
	res, err := strconv.Atoi(val)
	if err != nil {
		Panic("utl.Atoi: cannot parse string representing integer number: %s", val)
	}
	return
}

// Atof converts string to float64
func Atof(val string) (res float64) {
	res, err := strconv.ParseFloat(val, 64)
	if err != nil {
		Panic("utl.Atof: cannot parse string representing float number: %s", val)
	}
	return
}

// Itob converts from integer to bool
//  Note: only zero returns false
//        anything else returns true
func Itob(val int) bool {
	if val == 0 {
		return false
	}
	return true
}

// Btoi converts flag to interger
//  Note: true  => 1
//        false => 0
func Btoi(flag bool) int {
	if flag {
		return 1
	}
	return 0
}

// Btoa converts flag to string
//  Note: true  => "true"
//        false => "false"
func Btoa(flag bool) string {
	if flag {
		return "true"
	}
	return "false"
}

// PrintFormat commands ---------------------------------------------------------

// low intensity

func Pf(msg string, prm ...interface{}) {
	if Tsilent {
		return
	}
	fmt.Printf(msg, prm...)
}
func Pfcyan(msg string, prm ...interface{}) {
	if Tsilent {
		return
	}
	if UseColors {
		fmt.Printf("[0;36m"+msg+"[0m", prm...)
	} else {
		fmt.Printf(msg, prm...)
	}
}
func Pfyel(msg string, prm ...interface{}) {
	if Tsilent {
		return
	}
	if UseColors {
		fmt.Printf("[0;33m"+msg+"[0m", prm...)
	} else {
		fmt.Printf(msg, prm...)
	}
}
func Pfred(msg string, prm ...interface{}) {
	if Tsilent {
		return
	}
	if UseColors {
		fmt.Printf("[0;31m"+msg+"[0m", prm...)
	} else {
		fmt.Printf(msg, prm...)
	}
}
func Pfgreen(msg string, prm ...interface{}) {
	if Tsilent {
		return
	}
	if UseColors {
		fmt.Printf("[0;32m"+msg+"[0m", prm...)
	} else {
		fmt.Printf(msg, prm...)
	}
}
func Pfblue(msg string, prm ...interface{}) {
	if Tsilent {
		return
	}
	if UseColors {
		fmt.Printf("[0;34m"+msg+"[0m", prm...)
	} else {
		fmt.Printf(msg, prm...)
	}
}
func Pfmag(msg string, prm ...interface{}) {
	if Tsilent {
		return
	}
	if UseColors {
		fmt.Printf("[0;35m"+msg+"[0m", prm...)
	} else {
		fmt.Printf(msg, prm...)
	}
}
func Pflmag(msg string, prm ...interface{}) {
	if Tsilent {
		return
	}
	if UseColors {
		fmt.Printf("[0;95m"+msg+"[0m", prm...)
	} else {
		fmt.Printf(msg, prm...)
	}
}
func Pfpink(msg string, prm ...interface{}) {
	if Tsilent {
		return
	}
	if UseColors {
		fmt.Printf("[38;5;205m"+msg+"[0m", prm...)
	} else {
		fmt.Printf(msg, prm...)
	}
}
func Pfdgreen(msg string, prm ...interface{}) {
	if Tsilent {
		return
	}
	if UseColors {
		fmt.Printf("[38;5;22m"+msg+"[0m", prm...)
	} else {
		fmt.Printf(msg, prm...)
	}
}
func Pfgreen2(msg string, prm ...interface{}) {
	if Tsilent {
		return
	}
	if UseColors {
		fmt.Printf("[38;5;2m"+msg+"[0m", prm...)
	} else {
		fmt.Printf(msg, prm...)
	}
}
func Pfpurple(msg string, prm ...interface{}) {
	if Tsilent {
		return
	}
	if UseColors {
		fmt.Printf("[38;5;55m"+msg+"[0m", prm...)
	} else {
		fmt.Printf(msg, prm...)
	}
}
func Pfcyan2(msg string, prm ...interface{}) {
	if Tsilent {
		return
	}
	if UseColors {
		fmt.Printf("[38;5;50m"+msg+"[0m", prm...)
	} else {
		fmt.Printf(msg, prm...)
	}
}
func Pfdyel(msg string, prm ...interface{}) {
	if Tsilent {
		return
	}
	if UseColors {
		fmt.Printf("[38;5;58m"+msg+"[0m", prm...)
	} else {
		fmt.Printf(msg, prm...)
	}
}
func Pfdyel2(msg string, prm ...interface{}) {
	if Tsilent {
		return
	}
	if UseColors {
		fmt.Printf("[38;5;94m"+msg+"[0m", prm...)
	} else {
		fmt.Printf(msg, prm...)
	}
}
func Pfgrey(msg string, prm ...interface{}) {
	if Tsilent {
		return
	}
	if UseColors {
		fmt.Printf("[38;5;59m"+msg+"[0m", prm...)
	} else {
		fmt.Printf(msg, prm...)
	}
}
func Pfblue2(msg string, prm ...interface{}) {
	if Tsilent {
		return
	}
	if UseColors {
		fmt.Printf("[38;5;69m"+msg+"[0m", prm...)
	} else {
		fmt.Printf(msg, prm...)
	}
}
func Pfgrey2(msg string, prm ...interface{}) {
	if Tsilent {
		return
	}
	if UseColors {
		fmt.Printf("[38;5;60m"+msg+"[0m", prm...)
	} else {
		fmt.Printf(msg, prm...)
	}
}
func Pforan(msg string, prm ...interface{}) {
	if Tsilent {
		return
	}
	if UseColors {
		fmt.Printf("[38;5;202m"+msg+"[0m", prm...)
	} else {
		fmt.Printf(msg, prm...)
	}
}

// high intensity

func PfCyan(msg string, prm ...interface{}) {
	if Tsilent {
		return
	}
	if UseColors {
		fmt.Printf("[1;36m"+msg+"[0m", prm...)
	} else {
		fmt.Printf(msg, prm...)
	}
}
func PfYel(msg string, prm ...interface{}) {
	if Tsilent {
		return
	}
	if UseColors {
		fmt.Printf("[1;33m"+msg+"[0m", prm...)
	} else {
		fmt.Printf(msg, prm...)
	}
}
func PfRed(msg string, prm ...interface{}) {
	if Tsilent {
		return
	}
	if UseColors {
		fmt.Printf("[1;31m"+msg+"[0m", prm...)
	} else {
		fmt.Printf(msg, prm...)
	}
}
func PfGreen(msg string, prm ...interface{}) {
	if Tsilent {
		return
	}
	if UseColors {
		fmt.Printf("[1;32m"+msg+"[0m", prm...)
	} else {
		fmt.Printf(msg, prm...)
	}
}
func PfBlue(msg string, prm ...interface{}) {
	if Tsilent {
		return
	}
	if UseColors {
		fmt.Printf("[1;34m"+msg+"[0m", prm...)
	} else {
		fmt.Printf(msg, prm...)
	}
}
func PfMag(msg string, prm ...interface{}) {
	if Tsilent {
		return
	}
	if UseColors {
		fmt.Printf("[1;35m"+msg+"[0m", prm...)
	} else {
		fmt.Printf(msg, prm...)
	}
}
func PfWhite(msg string, prm ...interface{}) {
	if Tsilent {
		return
	}
	if UseColors {
		fmt.Printf("[1;37m"+msg+"[0m", prm...)
	} else {
		fmt.Printf(msg, prm...)
	}
}

// remove format ----------------------------------------------------------------

// UnColor removes console characters used to show colors
func UnColor(msg string) string {
	if len(msg) < 7 {
		Panic("utl.UnColor: string must have at least 7 characters (colored)")
	}
	return strings.Trim(msg, "")[6:]
}
