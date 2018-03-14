// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package al

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func TestLinkedList01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("LinkedList01. basic functionality")

	io.Pf("create list and insert first node\n")
	list := NewLinkedList()
	root := list.Insert(FromString("root"))
	chk.String(tst, ToString(root.Data), "root")
	chk.String(tst, list.String(), "[root]")

	io.Pf("\ninsert next node \"A\" and traverse list\n")
	node := list.Insert(FromString("A"))
	chk.String(tst, ToString(node.Data), "A")
	chk.String(tst, list.String(), "[root A]")

	io.Pf("\ninsert next node \"B\" and traverse list\n")
	node = list.Insert(FromString("B"))
	chk.String(tst, ToString(node.Data), "B")
	chk.String(tst, list.String(), "[root B A]")

	io.Pf("\ninsert next node \"C\" and traverse list\n")
	node = list.Insert(FromString("C"))
	chk.String(tst, ToString(node.Data), "C")
	chk.String(tst, list.String(), "[root C B A]")

	io.Pf("\nfind node \"B\" in list\n")
	B := list.Find(func(node *LinkedNode) bool {
		if ToString(node.Data) == "B" {
			return true
		}
		return false
	})
	if B == nil {
		tst.Errorf("cannot find \"B\"\n")
		return
	}
	chk.String(tst, ToString(B.Data), "B")
	chk.String(tst, list.String(), "[root C B A]")

	io.Pf("\nfind node \"root\" in list\n")
	res := list.Find(func(node *LinkedNode) bool {
		if ToString(node.Data) == "root" {
			return true
		}
		return false
	})
	if res == nil {
		tst.Errorf("cannot find \"root\"\n")
		return
	}
	chk.String(tst, ToString(res.Data), "root")
	chk.String(tst, list.String(), "[root C B A]")

	io.Pf("\nremove node \"root\" from list. \"C\" becomes root\n")
	list.Remove(root)
	chk.String(tst, ToString(list.root.Data), "C")
	chk.String(tst, list.String(), "[C B A]")

	io.Pf("\nfind node \"C\" in list\n")
	C := list.Find(func(node *LinkedNode) bool {
		if ToString(node.Data) == "C" {
			return true
		}
		return false
	})
	if C == nil {
		tst.Errorf("cannot find \"C\"\n")
		return
	}
	chk.String(tst, ToString(C.Data), "C")
	chk.String(tst, list.String(), "[C B A]")
}

func TestLinkedList02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("LinkedList02. edge cases")

	list := NewLinkedList()
	list.Remove(nil)
	root := list.Insert(FromString("root"))
	chk.String(tst, list.String(), "[root]")
	list.Remove(root)
	chk.String(tst, list.String(), "[]")
	res := list.Find(func(node *LinkedNode) bool {
		return false // will find nothing
	})
	if res != nil {
		tst.Errorf("Find should have returned <nil>\n")
		return
	}
}
