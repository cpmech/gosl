// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package al

import (
	"testing"

	"gosl/chk"
	"gosl/io"
)

// int ///////////////////////////////////////////////////////////////////////////

func TestIntQueue01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("IntQueue01")

	guessedMaxSize := 20
	qu := NewIntQueue(guessedMaxSize)
	member := qu.Out()
	if member != nil {
		tst.Errorf("returned member should be nil in an empty Queue\n")
		return
	}
	qu.Debug = true

	// add
	io.PfYel("In(1)\n")
	chk.Int(tst, "len(queue)", qu.Nmembers(), 0)
	qu.In(1)
	chk.String(tst, "1", io.Sf("%v", sliceIntToString(qu.ring)))
	chk.String(tst, "[1]", qu.String())
	chk.Int(tst, "1", *qu.Front(), 1)
	chk.Int(tst, "1", *qu.Back(), 1)
	chk.Int(tst, "len(queue)", qu.Nmembers(), 1)

	io.PfYel("\nIn(2)\n")
	qu.In(2)
	chk.String(tst, "1 2", io.Sf("%v", sliceIntToString(qu.ring)))
	chk.String(tst, "[1 2]", qu.String())
	chk.Int(tst, "1", *qu.Front(), 1)
	chk.Int(tst, "2", *qu.Back(), 2)
	chk.Int(tst, "len(queue)", qu.Nmembers(), 2)

	io.PfYel("\nIn(3)\n")
	qu.In(3)
	chk.String(tst, "1 2 3", io.Sf("%v", sliceIntToString(qu.ring)))
	chk.String(tst, "[1 2 3]", qu.String())
	chk.Int(tst, "1", *qu.Front(), 1)
	chk.Int(tst, "3", *qu.Back(), 3)
	chk.Int(tst, "len(queue)", qu.Nmembers(), 3)

	// remove
	io.PfYel("\nOut(1)\n")
	res := *qu.Out()
	chk.String(tst, "1 2 3", io.Sf("%v", sliceIntToString(qu.ring)))
	chk.String(tst, "[2 3]", qu.String())
	chk.Int(tst, "1", res, 1)
	chk.Int(tst, "len(queue)", qu.Nmembers(), 2)

	io.PfYel("\nOut(2)\n")
	res = *qu.Out()
	chk.String(tst, "1 2 3", io.Sf("%v", sliceIntToString(qu.ring)))
	chk.String(tst, "[3]", qu.String())
	chk.Int(tst, "2", res, 2)
	chk.Int(tst, "len(queue)", qu.Nmembers(), 1)

	io.PfYel("\nOut(3)\n")
	res = *qu.Out()
	chk.String(tst, "1 2 3", io.Sf("%v", sliceIntToString(qu.ring)))
	chk.String(tst, "[]", qu.String())
	chk.Int(tst, "3", res, 3)
	chk.Int(tst, "len(queue)", qu.Nmembers(), 0)

	// try to remove more in empty queue
	io.PfYel("\nOut(nothing)\n")
	member = qu.Out()
	if member != nil {
		tst.Errorf("returned member should be nil in an empty Queue\n")
		return
	}
	chk.String(tst, "[]", qu.String())
	chk.Int(tst, "len(queue)", qu.Nmembers(), 0)

	// add again
	io.PfYel("\nIn(4)\n")
	qu.In(4)
	chk.String(tst, "4 2 3", io.Sf("%v", sliceIntToString(qu.ring)))
	chk.String(tst, "[4]", qu.String())
	chk.Int(tst, "4", *qu.Front(), 4)
	chk.Int(tst, "4", *qu.Back(), 4)
	chk.Int(tst, "len(queue)", qu.Nmembers(), 1)

	io.PfYel("\nIn(5)\n")
	qu.In(5)
	chk.String(tst, "4 5 3", io.Sf("%v", sliceIntToString(qu.ring)))
	chk.String(tst, "[4 5]", qu.String())
	chk.Int(tst, "4", *qu.Front(), 4)
	chk.Int(tst, "5", *qu.Back(), 5)
	chk.Int(tst, "len(queue)", qu.Nmembers(), 2)

	io.PfYel("\nOut(4)\n")
	res = *qu.Out()
	chk.Int(tst, "4", res, 4)
	chk.String(tst, "4 5 3", io.Sf("%v", sliceIntToString(qu.ring)))
	chk.String(tst, "[5]", qu.String())
	chk.Int(tst, "5", *qu.Front(), 5)
	chk.Int(tst, "5", *qu.Back(), 5)
	chk.Int(tst, "len(queue)", qu.Nmembers(), 1)

	io.PfYel("\nIn(6)\n")
	qu.In(6)
	chk.String(tst, "4 5 6", io.Sf("%v", sliceIntToString(qu.ring)))
	chk.String(tst, "[5 6]", qu.String())
	chk.Int(tst, "5", *qu.Front(), 5)
	chk.Int(tst, "6", *qu.Back(), 6)
	chk.Int(tst, "len(queue)", qu.Nmembers(), 2)

	io.PfYel("\nIn(7)\n")
	qu.In(7)
	chk.String(tst, "7 5 6", io.Sf("%v", sliceIntToString(qu.ring)))
	chk.String(tst, "[5 6 7]", qu.String())
	chk.Int(tst, "5", *qu.Front(), 5)
	chk.Int(tst, "7", *qu.Back(), 7)
	chk.Int(tst, "len(queue)", qu.Nmembers(), 3)

	io.PfYel("\nIn(8)\n")
	qu.In(8)
	chk.String(tst, "5 6 7 8", io.Sf("%v", sliceIntToString(qu.ring)))
	chk.String(tst, "[5 6 7 8]", qu.String())
	chk.Int(tst, "5", *qu.Front(), 5)
	chk.Int(tst, "8", *qu.Back(), 8)
	chk.Int(tst, "len(queue)", qu.Nmembers(), 4)

	io.PfYel("\nOut(5)\n")
	res = *qu.Out()
	chk.Int(tst, "5", res, 5)
	chk.String(tst, "5 6 7 8", io.Sf("%v", sliceIntToString(qu.ring)))
	chk.String(tst, "[6 7 8]", qu.String())
	chk.Int(tst, "6", *qu.Front(), 6)
	chk.Int(tst, "8", *qu.Back(), 8)
	chk.Int(tst, "len(queue)", qu.Nmembers(), 3)

	io.PfYel("\nOut(6)\n")
	res = *qu.Out()
	chk.Int(tst, "6", res, 6)
	chk.String(tst, "5 6 7 8", io.Sf("%v", sliceIntToString(qu.ring)))
	chk.String(tst, "[7 8]", qu.String())
	chk.Int(tst, "7", *qu.Front(), 7)
	chk.Int(tst, "8", *qu.Back(), 8)
	chk.Int(tst, "len(queue)", qu.Nmembers(), 2)

	io.PfYel("\nOut(7)\n")
	res = *qu.Out()
	chk.Int(tst, "7", res, 7)
	chk.String(tst, "5 6 7 8", io.Sf("%v", sliceIntToString(qu.ring)))
	chk.String(tst, "[8]", qu.String())
	chk.Int(tst, "8", *qu.Front(), 8)
	chk.Int(tst, "8", *qu.Back(), 8)
	chk.Int(tst, "len(queue)", qu.Nmembers(), 1)

	io.PfYel("\nOut(8)\n")
	res = *qu.Out()
	chk.Int(tst, "8", res, 8)
	chk.String(tst, "5 6 7 8", io.Sf("%v", sliceIntToString(qu.ring)))
	chk.String(tst, "[]", qu.String())
	chk.Int(tst, "len(queue)", qu.Nmembers(), 0)

	io.PfYel("\nOut(nothing)\n")
	chk.String(tst, "5 6 7 8", io.Sf("%v", sliceIntToString(qu.ring)))
	chk.String(tst, "[]", qu.String())
	chk.Int(tst, "len(queue)", qu.Nmembers(), 0)

	io.PfYel("\nIn(-1)\n")
	qu.In(-1)
	chk.String(tst, "-1 6 7 8", io.Sf("%v", sliceIntToString(qu.ring)))
	chk.String(tst, "[-1]", qu.String())
	chk.Int(tst, "-1", *qu.Front(), -1)
	chk.Int(tst, "-1", *qu.Back(), -1)
	chk.Int(tst, "len(queue)", qu.Nmembers(), 1)

	io.PfYel("\nIn(-2)\n")
	qu.In(-2)
	chk.String(tst, "-1 -2 7 8", io.Sf("%v", sliceIntToString(qu.ring)))
	chk.String(tst, "[-1 -2]", qu.String())
	chk.Int(tst, "-1", *qu.Front(), -1)
	chk.Int(tst, "-2", *qu.Back(), -2)
	chk.Int(tst, "len(queue)", qu.Nmembers(), 2)
}

// float64 ///////////////////////////////////////////////////////////////////////////

func TestFloat64Queue01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Float64Queue01")

	guessedMaxSize := 20
	qu := NewFloat64Queue(guessedMaxSize)
	member := qu.Out()
	if member != nil {
		tst.Errorf("returned member should be nil in an empty Queue\n")
		return
	}
	qu.Debug = true

	// add
	io.PfYel("In(1)\n")
	chk.Int(tst, "len(queue)", qu.Nmembers(), 0)
	qu.In(1)
	chk.String(tst, "1", io.Sf("%v", sliceFloat64ToString(qu.ring)))
	chk.String(tst, "[1]", qu.String())
	chk.Float64(tst, "1", 1e-15, *qu.Front(), 1)
	chk.Float64(tst, "1", 1e-15, *qu.Back(), 1)
	chk.Int(tst, "len(queue)", qu.Nmembers(), 1)

	io.PfYel("\nIn(2)\n")
	qu.In(2)
	chk.String(tst, "1 2", io.Sf("%v", sliceFloat64ToString(qu.ring)))
	chk.String(tst, "[1 2]", qu.String())
	chk.Float64(tst, "1", 1e-15, *qu.Front(), 1)
	chk.Float64(tst, "2", 1e-15, *qu.Back(), 2)
	chk.Int(tst, "len(queue)", qu.Nmembers(), 2)

	io.PfYel("\nIn(3)\n")
	qu.In(3)
	chk.String(tst, "1 2 3", io.Sf("%v", sliceFloat64ToString(qu.ring)))
	chk.String(tst, "[1 2 3]", qu.String())
	chk.Float64(tst, "1", 1e-15, *qu.Front(), 1)
	chk.Float64(tst, "3", 1e-15, *qu.Back(), 3)
	chk.Int(tst, "len(queue)", qu.Nmembers(), 3)

	// remove
	io.PfYel("\nOut(1)\n")
	res := *qu.Out()
	chk.String(tst, "1 2 3", io.Sf("%v", sliceFloat64ToString(qu.ring)))
	chk.String(tst, "[2 3]", qu.String())
	chk.Float64(tst, "1", 1e-15, res, 1)
	chk.Int(tst, "len(queue)", qu.Nmembers(), 2)

	io.PfYel("\nOut(2)\n")
	res = *qu.Out()
	chk.String(tst, "1 2 3", io.Sf("%v", sliceFloat64ToString(qu.ring)))
	chk.String(tst, "[3]", qu.String())
	chk.Float64(tst, "2", 1e-15, res, 2)
	chk.Int(tst, "len(queue)", qu.Nmembers(), 1)

	io.PfYel("\nOut(3)\n")
	res = *qu.Out()
	chk.String(tst, "1 2 3", io.Sf("%v", sliceFloat64ToString(qu.ring)))
	chk.String(tst, "[]", qu.String())
	chk.Float64(tst, "3", 1e-15, res, 3)
	chk.Int(tst, "len(queue)", qu.Nmembers(), 0)

	// try to remove more in empty queue
	io.PfYel("\nOut(nothing)\n")
	member = qu.Out()
	if member != nil {
		tst.Errorf("returned member should be nil in an empty Queue\n")
		return
	}
	chk.String(tst, "[]", qu.String())
	chk.Int(tst, "len(queue)", qu.Nmembers(), 0)

	// add again
	io.PfYel("\nIn(4)\n")
	qu.In(4)
	chk.String(tst, "4 2 3", io.Sf("%v", sliceFloat64ToString(qu.ring)))
	chk.String(tst, "[4]", qu.String())
	chk.Float64(tst, "4", 1e-15, *qu.Front(), 4)
	chk.Float64(tst, "4", 1e-15, *qu.Back(), 4)
	chk.Int(tst, "len(queue)", qu.Nmembers(), 1)

	io.PfYel("\nIn(5)\n")
	qu.In(5)
	chk.String(tst, "4 5 3", io.Sf("%v", sliceFloat64ToString(qu.ring)))
	chk.String(tst, "[4 5]", qu.String())
	chk.Float64(tst, "4", 1e-15, *qu.Front(), 4)
	chk.Float64(tst, "5", 1e-15, *qu.Back(), 5)
	chk.Int(tst, "len(queue)", qu.Nmembers(), 2)

	io.PfYel("\nOut(4)\n")
	res = *qu.Out()
	chk.Float64(tst, "4", 1e-15, res, 4)
	chk.String(tst, "4 5 3", io.Sf("%v", sliceFloat64ToString(qu.ring)))
	chk.String(tst, "[5]", qu.String())
	chk.Float64(tst, "5", 1e-15, *qu.Front(), 5)
	chk.Float64(tst, "5", 1e-15, *qu.Back(), 5)
	chk.Int(tst, "len(queue)", qu.Nmembers(), 1)

	io.PfYel("\nIn(6)\n")
	qu.In(6)
	chk.String(tst, "4 5 6", io.Sf("%v", sliceFloat64ToString(qu.ring)))
	chk.String(tst, "[5 6]", qu.String())
	chk.Float64(tst, "5", 1e-15, *qu.Front(), 5)
	chk.Float64(tst, "6", 1e-15, *qu.Back(), 6)
	chk.Int(tst, "len(queue)", qu.Nmembers(), 2)

	io.PfYel("\nIn(7)\n")
	qu.In(7)
	chk.String(tst, "7 5 6", io.Sf("%v", sliceFloat64ToString(qu.ring)))
	chk.String(tst, "[5 6 7]", qu.String())
	chk.Float64(tst, "5", 1e-15, *qu.Front(), 5)
	chk.Float64(tst, "7", 1e-15, *qu.Back(), 7)
	chk.Int(tst, "len(queue)", qu.Nmembers(), 3)

	io.PfYel("\nIn(8)\n")
	qu.In(8)
	chk.String(tst, "5 6 7 8", io.Sf("%v", sliceFloat64ToString(qu.ring)))
	chk.String(tst, "[5 6 7 8]", qu.String())
	chk.Float64(tst, "5", 1e-15, *qu.Front(), 5)
	chk.Float64(tst, "8", 1e-15, *qu.Back(), 8)
	chk.Int(tst, "len(queue)", qu.Nmembers(), 4)

	io.PfYel("\nOut(5)\n")
	res = *qu.Out()
	chk.Float64(tst, "5", 1e-15, res, 5)
	chk.String(tst, "5 6 7 8", io.Sf("%v", sliceFloat64ToString(qu.ring)))
	chk.String(tst, "[6 7 8]", qu.String())
	chk.Float64(tst, "6", 1e-15, *qu.Front(), 6)
	chk.Float64(tst, "8", 1e-15, *qu.Back(), 8)
	chk.Int(tst, "len(queue)", qu.Nmembers(), 3)

	io.PfYel("\nOut(6)\n")
	res = *qu.Out()
	chk.Float64(tst, "6", 1e-15, res, 6)
	chk.String(tst, "5 6 7 8", io.Sf("%v", sliceFloat64ToString(qu.ring)))
	chk.String(tst, "[7 8]", qu.String())
	chk.Float64(tst, "7", 1e-15, *qu.Front(), 7)
	chk.Float64(tst, "8", 1e-15, *qu.Back(), 8)
	chk.Int(tst, "len(queue)", qu.Nmembers(), 2)

	io.PfYel("\nOut(7)\n")
	res = *qu.Out()
	chk.Float64(tst, "7", 1e-15, res, 7)
	chk.String(tst, "5 6 7 8", io.Sf("%v", sliceFloat64ToString(qu.ring)))
	chk.String(tst, "[8]", qu.String())
	chk.Float64(tst, "8", 1e-15, *qu.Front(), 8)
	chk.Float64(tst, "8", 1e-15, *qu.Back(), 8)
	chk.Int(tst, "len(queue)", qu.Nmembers(), 1)

	io.PfYel("\nOut(8)\n")
	res = *qu.Out()
	chk.Float64(tst, "8", 1e-15, res, 8)
	chk.String(tst, "5 6 7 8", io.Sf("%v", sliceFloat64ToString(qu.ring)))
	chk.String(tst, "[]", qu.String())
	chk.Int(tst, "len(queue)", qu.Nmembers(), 0)

	io.PfYel("\nOut(nothing)\n")
	chk.String(tst, "5 6 7 8", io.Sf("%v", sliceFloat64ToString(qu.ring)))
	chk.String(tst, "[]", qu.String())
	chk.Int(tst, "len(queue)", qu.Nmembers(), 0)

	io.PfYel("\nIn(-1)\n")
	qu.In(-1)
	chk.String(tst, "-1 6 7 8", io.Sf("%v", sliceFloat64ToString(qu.ring)))
	chk.String(tst, "[-1]", qu.String())
	chk.Float64(tst, "-1", 1e-15, *qu.Front(), -1)
	chk.Float64(tst, "-1", 1e-15, *qu.Back(), -1)
	chk.Int(tst, "len(queue)", qu.Nmembers(), 1)

	io.PfYel("\nIn(-2)\n")
	qu.In(-2)
	chk.String(tst, "-1 -2 7 8", io.Sf("%v", sliceFloat64ToString(qu.ring)))
	chk.String(tst, "[-1 -2]", qu.String())
	chk.Float64(tst, "-1", 1e-15, *qu.Front(), -1)
	chk.Float64(tst, "-2", 1e-15, *qu.Back(), -2)
	chk.Int(tst, "len(queue)", qu.Nmembers(), 2)
}

// string ///////////////////////////////////////////////////////////////////////////

func TestStringQueue01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("StringQueue01")

	guessedMaxSize := 20
	qu := NewStringQueue(guessedMaxSize)
	member := qu.Out()
	if member != nil {
		tst.Errorf("returned member should be nil in an empty Queue\n")
		return
	}
	qu.Debug = true

	// add
	io.PfYel("In(l)\n")
	chk.Int(tst, "len(queue)", qu.Nmembers(), 0)
	qu.In("l")
	chk.String(tst, "l", io.Sf("%v", sliceStringToString(qu.ring)))
	chk.String(tst, "[l]", qu.String())
	chk.String(tst, *qu.Front(), "l")
	chk.String(tst, *qu.Back(), "l")
	chk.Int(tst, "len(queue)", qu.Nmembers(), 1)

	io.PfYel("\nIn(o)\n")
	qu.In("o")
	chk.String(tst, "l o", io.Sf("%v", sliceStringToString(qu.ring)))
	chk.String(tst, "[l o]", qu.String())
	chk.String(tst, *qu.Front(), "l")
	chk.String(tst, *qu.Back(), "o")
	chk.Int(tst, "len(queue)", qu.Nmembers(), 2)

	io.PfYel("\nIn(v)\n")
	qu.In("v")
	chk.String(tst, "l o v", io.Sf("%v", sliceStringToString(qu.ring)))
	chk.String(tst, "[l o v]", qu.String())
	chk.String(tst, *qu.Front(), "l")
	chk.String(tst, *qu.Back(), "v")
	chk.Int(tst, "len(queue)", qu.Nmembers(), 3)

	// remove
	io.PfYel("\nOut(l)\n")
	res := *qu.Out()
	chk.String(tst, "l o v", io.Sf("%v", sliceStringToString(qu.ring)))
	chk.String(tst, "[o v]", qu.String())
	chk.String(tst, res, "l")
	chk.Int(tst, "len(queue)", qu.Nmembers(), 2)

	io.PfYel("\nOut(o)\n")
	res = *qu.Out()
	chk.String(tst, "l o v", io.Sf("%v", sliceStringToString(qu.ring)))
	chk.String(tst, "[v]", qu.String())
	chk.String(tst, res, "o")
	chk.Int(tst, "len(queue)", qu.Nmembers(), 1)

	io.PfYel("\nOut(v)\n")
	res = *qu.Out()
	chk.String(tst, "l o v", io.Sf("%v", sliceStringToString(qu.ring)))
	chk.String(tst, "[]", qu.String())
	chk.String(tst, res, "v")
	chk.Int(tst, "len(queue)", qu.Nmembers(), 0)

	// try to remove more in empty queue
	io.PfYel("\nOut(nothing)\n")
	member = qu.Out()
	if member != nil {
		tst.Errorf("returned member should be nil in an empty Queue\n")
		return
	}
	chk.String(tst, "[]", qu.String())
	chk.Int(tst, "len(queue)", qu.Nmembers(), 0)

	// add again
	io.PfYel("\nIn(a)\n")
	qu.In("a")
	chk.String(tst, "a o v", io.Sf("%v", sliceStringToString(qu.ring)))
	chk.String(tst, "[a]", qu.String())
	chk.String(tst, *qu.Front(), "a")
	chk.String(tst, *qu.Back(), "a")
	chk.Int(tst, "len(queue)", qu.Nmembers(), 1)

	io.PfYel("\nIn(b)\n")
	qu.In("b")
	chk.String(tst, "a b v", io.Sf("%v", sliceStringToString(qu.ring)))
	chk.String(tst, "[a b]", qu.String())
	chk.String(tst, *qu.Front(), "a")
	chk.String(tst, *qu.Back(), "b")
	chk.Int(tst, "len(queue)", qu.Nmembers(), 2)

	io.PfYel("\nOut(a)\n")
	res = *qu.Out()
	chk.String(tst, res, "a")
	chk.String(tst, "a b v", io.Sf("%v", sliceStringToString(qu.ring)))
	chk.String(tst, "[b]", qu.String())
	chk.String(tst, *qu.Front(), "b")
	chk.String(tst, *qu.Back(), "b")
	chk.Int(tst, "len(queue)", qu.Nmembers(), 1)

	io.PfYel("\nIn(a) again\n")
	qu.In("a")
	chk.String(tst, "a b a", io.Sf("%v", sliceStringToString(qu.ring)))
	chk.String(tst, "[b a]", qu.String())
	chk.String(tst, *qu.Front(), "b")
	chk.String(tst, *qu.Back(), "a")
	chk.Int(tst, "len(queue)", qu.Nmembers(), 2)

	io.PfYel("\nIn(c)\n")
	qu.In("c")
	chk.String(tst, "c b a", io.Sf("%v", sliceStringToString(qu.ring)))
	chk.String(tst, "[b a c]", qu.String())
	chk.String(tst, *qu.Front(), "b")
	chk.String(tst, *qu.Back(), "c")
	chk.Int(tst, "len(queue)", qu.Nmembers(), 3)

	io.PfYel("\nIn(x)\n")
	qu.In("x")
	chk.String(tst, "b a c x", io.Sf("%v", sliceStringToString(qu.ring)))
	chk.String(tst, "[b a c x]", qu.String())
	chk.String(tst, *qu.Front(), "b")
	chk.String(tst, *qu.Back(), "x")
	chk.Int(tst, "len(queue)", qu.Nmembers(), 4)

	io.PfYel("\nOut(b)\n")
	res = *qu.Out()
	chk.String(tst, res, "b")
	chk.String(tst, "b a c x", io.Sf("%v", sliceStringToString(qu.ring)))
	chk.String(tst, "[a c x]", qu.String())
	chk.String(tst, *qu.Front(), "a")
	chk.String(tst, *qu.Back(), "x")
	chk.Int(tst, "len(queue)", qu.Nmembers(), 3)

	io.PfYel("\nOut(a)\n")
	res = *qu.Out()
	chk.String(tst, res, "a")
	chk.String(tst, "b a c x", io.Sf("%v", sliceStringToString(qu.ring)))
	chk.String(tst, "[c x]", qu.String())
	chk.String(tst, *qu.Front(), "c")
	chk.String(tst, *qu.Back(), "x")
	chk.Int(tst, "len(queue)", qu.Nmembers(), 2)

	io.PfYel("\nOut(c)\n")
	res = *qu.Out()
	chk.String(tst, res, "c")
	chk.String(tst, "b a c x", io.Sf("%v", sliceStringToString(qu.ring)))
	chk.String(tst, "[x]", qu.String())
	chk.String(tst, *qu.Front(), "x")
	chk.String(tst, *qu.Back(), "x")
	chk.Int(tst, "len(queue)", qu.Nmembers(), 1)

	io.PfYel("\nOut(x)\n")
	res = *qu.Out()
	chk.String(tst, res, "x")
	chk.String(tst, "b a c x", io.Sf("%v", sliceStringToString(qu.ring)))
	chk.String(tst, "[]", qu.String())
	chk.Int(tst, "len(queue)", qu.Nmembers(), 0)

	io.PfYel("\nOut(nothing)\n")
	chk.String(tst, "b a c x", io.Sf("%v", sliceStringToString(qu.ring)))
	chk.String(tst, "[]", qu.String())
	chk.Int(tst, "len(queue)", qu.Nmembers(), 0)

	io.PfYel("\nIn(i)\n")
	qu.In("i")
	chk.String(tst, "i a c x", io.Sf("%v", sliceStringToString(qu.ring)))
	chk.String(tst, "[i]", qu.String())
	chk.String(tst, *qu.Front(), "i")
	chk.String(tst, *qu.Back(), "i")
	chk.Int(tst, "len(queue)", qu.Nmembers(), 1)

	io.PfYel("\nIn(j)\n")
	qu.In("j")
	chk.String(tst, "i j c x", io.Sf("%v", sliceStringToString(qu.ring)))
	chk.String(tst, "[i j]", qu.String())
	chk.String(tst, *qu.Front(), "i")
	chk.String(tst, *qu.Back(), "j")
	chk.Int(tst, "len(queue)", qu.Nmembers(), 2)
}
