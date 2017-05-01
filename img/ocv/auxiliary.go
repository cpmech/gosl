// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ocv

/*
#include "auxiliary.h"
*/
import "C"

import (
	"errors"
	"unsafe"
)

// makeArrayChar allocates C array of chars
// Note: make sure to call free() to deallocate memory
func makeArrayChar(array []string) (csize C.int, carray **C.char, free func()) {

	size := len(array)
	csize = C.int(size)

	carray = C.make_argv(csize)
	pointers := make([]*C.char, size)
	for i, s := range array {
		pointers[i] = C.CString(s)
		C.set_arg(carray, C.int(i), pointers[i])
	}

	free = func() {
		for i := 0; i < size; i++ {
			C.free(unsafe.Pointer(pointers[i]))
		}
		C.free(unsafe.Pointer(carray))
	}
	return
}

func make3strings(str1, str2, str3 string) (cstr1, cstr2, cstr3 *C.char, free func()) {
	cstr1 = C.CString(str1)
	cstr2 = C.CString(str2)
	cstr3 = C.CString(str3)
	free = func() {
		C.free(unsafe.Pointer(cstr1))
		C.free(unsafe.Pointer(cstr2))
		C.free(unsafe.Pointer(cstr3))
	}
	return
}

func make2strings(str1, str2 string) (cstr1, cstr2 *C.char, free func()) {
	cstr1 = C.CString(str1)
	cstr2 = C.CString(str2)
	free = func() {
		C.free(unsafe.Pointer(cstr1))
		C.free(unsafe.Pointer(cstr2))
	}
	return
}

func make2stringsInt(str1, str2 string, val int) (cstr1, cstr2 *C.char, cval C.int, free func()) {
	cstr1 = C.CString(str1)
	cstr2 = C.CString(str2)
	cval = C.int(val)
	free = func() {
		C.free(unsafe.Pointer(cstr1))
		C.free(unsafe.Pointer(cstr2))
	}
	return
}

func make2stringsFloat(str1, str2 string, val float64) (cstr1, cstr2 *C.char, cval C.double, free func()) {
	cstr1 = C.CString(str1)
	cstr2 = C.CString(str2)
	cval = C.double(val)
	free = func() {
		C.free(unsafe.Pointer(cstr1))
		C.free(unsafe.Pointer(cstr2))
	}
	return
}

func check(status C.int, cerror *C.char) error {
	if int(status) != 0 {
		return errors.New(C.GoString(cerror))
	}
	return nil
}

func checkWithPanic(status C.int, cerror *C.char) {
	if int(status) != 0 {
		panic(C.GoString(cerror))
	}
}
