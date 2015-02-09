// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestFileIO1(tst *testing.T) {

	prevTs := Tsilent
	defer func() {
		Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//Tsilent = false
	TTitle("TestFile IO 1")

	fn := "test/dorival/file.sim"
	CheckString(tst, "file.sim", filepath.Base(fn))
	CheckString(tst, ".sim", filepath.Ext(fn))
	CheckString(tst, "file", FnKey(fn))
	CheckString(tst, ".sim", FnExt(fn))
	CheckString(tst, "test/dorival/file", PathKey(fn))

	gn := "test/dorival/file.h5"
	CheckString(tst, "file.h5", filepath.Base(gn))
	CheckString(tst, ".h5", filepath.Ext(gn))
	CheckString(tst, "file", FnKey(gn))
	CheckString(tst, ".h5", FnExt(gn))
	CheckString(tst, "test/dorival/file", PathKey(gn))

	Pf("\n")
	Pf("fn   = %s\n", fn)
	Pf("base = %s\n", filepath.Base(fn))
	Pf("ext  = %s\n", filepath.Ext(fn))
	Pf("fnk  = %s\n", FnKey(fn))
	Pf("\n")

	fn = "test/dorival/file"
	CheckString(tst, "file", filepath.Base(fn))
	CheckString(tst, "", filepath.Ext(fn))
	CheckString(tst, "file", FnKey(fn))
	CheckString(tst, "test/dorival/file", PathKey(fn))

	Pf("\n")
	Pf("fn   = %s\n", fn)
	Pf("base = %s\n", filepath.Base(fn))
	Pf("ext  = %s\n", filepath.Ext(fn))
	Pf("fnk  = %s\n", FnKey(fn))
	Pf("\n")

	fn = "test/dorival/file."
	CheckString(tst, "file.", filepath.Base(fn))
	CheckString(tst, ".", filepath.Ext(fn))
	CheckString(tst, "file", FnKey(fn))
	CheckString(tst, "test/dorival/file", PathKey(fn))

	Pf("\n")
	Pf("fn   = %s\n", fn)
	Pf("base = %s\n", filepath.Base(fn))
	Pf("ext  = %s\n", filepath.Ext(fn))
	Pf("fnk  = %s\n", FnKey(fn))
	Pf("\n")

	fn = "test/dorival/f.extension"
	CheckString(tst, "f.extension", filepath.Base(fn))
	CheckString(tst, ".extension", filepath.Ext(fn))
	CheckString(tst, "f", FnKey(fn))
	CheckString(tst, "test/dorival/f", PathKey(fn))

	Pf("\n")
	Pf("fn   = %s\n", fn)
	Pf("base = %s\n", filepath.Base(fn))
	Pf("ext  = %s\n", filepath.Ext(fn))
	Pf("fnk  = %s\n", FnKey(fn))
	Pf("pathkey = %s\n", PathKey(fn))
	Pf("\n")
}

func TestFileIO2(tst *testing.T) {

	prevTs := Tsilent
	defer func() {
		Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//Tsilent = false
	TTitle("TestFile IO 2")

	os.MkdirAll("/tmp/gosl", 0777)

	fn := "/tmp/gosl/gosl_t_01_fileio.res"
	var bout bytes.Buffer
	Ff(&bout, "just testing %g\n", 666.0)
	AppendToFile(fn, &bout)

	ReadLines(fn, func(idx int, line string) (stop bool) {
		if line != "just testing 666" {
			Panic("read wrong line: '%v'", line)
		}
		return false
	})
}

func TestFileIO3(tst *testing.T) {

	prevTs := Tsilent
	defer func() {
		Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//Tsilent = false
	TTitle("TestFile IO 3")

	type Test struct {
		Id     int
		Cells  []int
		Types  []string
		Values []float64
	}
	t := Test{0, []int{7, 3, 5}, []string{"a", "x", "p", "y"}, []float64{666}}
	Pf("t = %v\n", t)

	b, err := json.Marshal(&t)
	if err != nil {
		Panic("marshal failed for %+v", t)
	}
	WriteBytesToFileD("/tmp/gosl/", "gosl_jsontest.res", b)
	PfBlue("file written /tmp/gosl/gosl_jsontest.res\n")
}

func TestFileIO4(tst *testing.T) {

	prevTs := Tsilent
	defer func() {
		Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//Tsilent = false
	TTitle("TestFile IO 4")

	var b0, b1 bytes.Buffer
	Ff(&b0, "from gosl_fig import *\n")
	Ff(&b1, "x = linspace(0., 4., 101)\n")
	Ff(&b1, "y = exp(-2.0*x)\n")
	Ff(&b1, "plot(x,y, label='decay')\n")
	Ff(&b1, "Gll('x','y')\n")
	Ff(&b1, "show()\n")
	WritePython("/tmp/gosl/", "testfileio04", false, &b0, &b1)
}

func TestFileIO5(tst *testing.T) {

	prevTs := Tsilent
	defer func() {
		Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//Tsilent = false
	TTitle("TestFile IO 5")

	theline := "Hello World !!!"
	WriteFileSD("/tmp/gosl", "filestring.txt", theline)

	f, err := OpenFileR("/tmp/gosl/filestring.txt")
	if err != nil {
		Panic("%v", err)
	}

	ReadLinesFile(f, func(idx int, line string) (stop bool) {
		Pforan("line = %v\n", line)
		CheckString(tst, line, theline)
		return
	})
}
