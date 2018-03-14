// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package al

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func TestQueue01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Queue01")

	guessedMaxSize := 20
	qu := NewQueue(guessedMaxSize)
	member := qu.Out()
	if member != nil {
		tst.Errorf("returned member should be nil in an empty Queue\n")
		return
	}
	qu.Debug = true

	// add
	io.PfYel("In(l)\n")
	chk.Int(tst, "len(queue)", qu.Nmembers(), 0)
	qu.In(FromString("l"))
	chk.String(tst, "[l]", io.Sf("%v", qu.ring))
	chk.String(tst, "[l]", qu.String())
	chk.String(tst, ToString(qu.Front()), "l")
	chk.String(tst, ToString(qu.Back()), "l")
	chk.Int(tst, "len(queue)", qu.Nmembers(), 1)

	io.PfYel("\nIn(o)\n")
	qu.In(FromString("o"))
	chk.String(tst, "[l o]", io.Sf("%v", qu.ring))
	chk.String(tst, "[l o]", qu.String())
	chk.String(tst, ToString(qu.Front()), "l")
	chk.String(tst, ToString(qu.Back()), "o")
	chk.Int(tst, "len(queue)", qu.Nmembers(), 2)

	io.PfYel("\nIn(v)\n")
	qu.In(FromString("v"))
	chk.String(tst, "[l o v]", io.Sf("%v", qu.ring))
	chk.String(tst, "[l o v]", qu.String())
	chk.String(tst, ToString(qu.Front()), "l")
	chk.String(tst, ToString(qu.Back()), "v")
	chk.Int(tst, "len(queue)", qu.Nmembers(), 3)

	// remove
	io.PfYel("\nOut(l)\n")
	res := ToString(qu.Out())
	chk.String(tst, "[l o v]", io.Sf("%v", qu.ring))
	chk.String(tst, "[o v]", qu.String())
	chk.String(tst, res, "l")
	chk.Int(tst, "len(queue)", qu.Nmembers(), 2)

	io.PfYel("\nOut(o)\n")
	res = ToString(qu.Out())
	chk.String(tst, "[l o v]", io.Sf("%v", qu.ring))
	chk.String(tst, "[v]", qu.String())
	chk.String(tst, res, "o")
	chk.Int(tst, "len(queue)", qu.Nmembers(), 1)

	io.PfYel("\nOut(v)\n")
	res = ToString(qu.Out())
	chk.String(tst, "[l o v]", io.Sf("%v", qu.ring))
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
	qu.In(FromString("a"))
	chk.String(tst, "[a o v]", io.Sf("%v", qu.ring))
	chk.String(tst, "[a]", qu.String())
	chk.String(tst, ToString(qu.Front()), "a")
	chk.String(tst, ToString(qu.Back()), "a")
	chk.Int(tst, "len(queue)", qu.Nmembers(), 1)

	io.PfYel("\nIn(b)\n")
	qu.In(FromString("b"))
	chk.String(tst, "[a b v]", io.Sf("%v", qu.ring))
	chk.String(tst, "[a b]", qu.String())
	chk.String(tst, ToString(qu.Front()), "a")
	chk.String(tst, ToString(qu.Back()), "b")
	chk.Int(tst, "len(queue)", qu.Nmembers(), 2)

	io.PfYel("\nOut(a)\n")
	res = ToString(qu.Out())
	chk.String(tst, "[a b v]", io.Sf("%v", qu.ring))
	chk.String(tst, "[b]", qu.String())
	chk.String(tst, ToString(qu.Front()), "b")
	chk.String(tst, ToString(qu.Back()), "b")
	chk.Int(tst, "len(queue)", qu.Nmembers(), 1)

	io.PfYel("\nIn(a) again\n")
	qu.In(FromString("a"))
	chk.String(tst, "[a b a]", io.Sf("%v", qu.ring))
	chk.String(tst, "[b a]", qu.String())
	chk.String(tst, ToString(qu.Front()), "b")
	chk.String(tst, ToString(qu.Back()), "a")
	chk.Int(tst, "len(queue)", qu.Nmembers(), 2)

	io.PfYel("\nIn(c)\n")
	qu.In(FromString("c"))
	chk.String(tst, "[c b a]", io.Sf("%v", qu.ring))
	chk.String(tst, "[b a c]", qu.String())
	chk.String(tst, ToString(qu.Front()), "b")
	chk.String(tst, ToString(qu.Back()), "c")
	chk.Int(tst, "len(queue)", qu.Nmembers(), 3)

	io.PfYel("\nIn(x)\n")
	qu.In(FromString("x"))
	chk.String(tst, "[b a c x]", io.Sf("%v", qu.ring))
	chk.String(tst, "[b a c x]", qu.String())
	chk.String(tst, ToString(qu.Front()), "b")
	chk.String(tst, ToString(qu.Back()), "x")
	chk.Int(tst, "len(queue)", qu.Nmembers(), 4)

	io.PfYel("\nOut(b)\n")
	res = ToString(qu.Out())
	chk.String(tst, "[b a c x]", io.Sf("%v", qu.ring))
	chk.String(tst, "[a c x]", qu.String())
	chk.String(tst, ToString(qu.Front()), "a")
	chk.String(tst, ToString(qu.Back()), "x")
	chk.Int(tst, "len(queue)", qu.Nmembers(), 3)

	io.PfYel("\nOut(a)\n")
	res = ToString(qu.Out())
	chk.String(tst, "[b a c x]", io.Sf("%v", qu.ring))
	chk.String(tst, "[c x]", qu.String())
	chk.String(tst, ToString(qu.Front()), "c")
	chk.String(tst, ToString(qu.Back()), "x")
	chk.Int(tst, "len(queue)", qu.Nmembers(), 2)

	io.PfYel("\nOut(c)\n")
	res = ToString(qu.Out())
	chk.String(tst, "[b a c x]", io.Sf("%v", qu.ring))
	chk.String(tst, "[x]", qu.String())
	chk.String(tst, ToString(qu.Front()), "x")
	chk.String(tst, ToString(qu.Back()), "x")
	chk.Int(tst, "len(queue)", qu.Nmembers(), 1)

	io.PfYel("\nOut(x)\n")
	res = ToString(qu.Out())
	chk.String(tst, "[b a c x]", io.Sf("%v", qu.ring))
	chk.String(tst, "[]", qu.String())
	chk.Int(tst, "len(queue)", qu.Nmembers(), 0)

	io.PfYel("\nOut(nothing)\n")
	chk.String(tst, "[b a c x]", io.Sf("%v", qu.ring))
	chk.String(tst, "[]", qu.String())
	chk.Int(tst, "len(queue)", qu.Nmembers(), 0)

	io.PfYel("\nIn(i)\n")
	qu.In(FromString("i"))
	chk.String(tst, "[i a c x]", io.Sf("%v", qu.ring))
	chk.String(tst, "[i]", qu.String())
	chk.String(tst, ToString(qu.Front()), "i")
	chk.String(tst, ToString(qu.Back()), "i")
	chk.Int(tst, "len(queue)", qu.Nmembers(), 1)

	io.PfYel("\nIn(j)\n")
	qu.In(FromString("j"))
	chk.String(tst, "[i j c x]", io.Sf("%v", qu.ring))
	chk.String(tst, "[i j]", qu.String())
	chk.String(tst, ToString(qu.Front()), "i")
	chk.String(tst, ToString(qu.Back()), "j")
	chk.Int(tst, "len(queue)", qu.Nmembers(), 2)
}
