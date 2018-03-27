// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package al

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

// int //////////////////////////////////////////////////////////////////////////////////////////

func TestIntLinkedList01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("IntLinkedList01. basic functionality")

	io.Pf("create list and insert first node\n")
	list := NewIntLinkedList()
	root := list.Insert(0)
	chk.Int(tst, "root == 0", *root.Data, 0)
	chk.String(tst, list.String(), "[0]")

	io.Pf("\ninsert next node \"1\" and traverse list\n")
	node := list.Insert(1)
	chk.Int(tst, "1", *node.Data, 1)
	chk.String(tst, list.String(), "[0 1]")

	io.Pf("\ninsert next node \"2\" and traverse list\n")
	node = list.Insert(2)
	chk.Int(tst, "2", *node.Data, 2)
	chk.String(tst, list.String(), "[0 2 1]")

	io.Pf("\ninsert next node \"3\" and traverse list\n")
	node = list.Insert(3)
	chk.Int(tst, "3", *node.Data, 3)
	chk.String(tst, list.String(), "[0 3 2 1]")

	io.Pf("\nfind node \"2\" in list\n")
	node2 := list.Find(func(node *IntLinkedNode) bool {
		if *node.Data == 2 {
			return true
		}
		return false
	})
	if node2 == nil {
		tst.Errorf("cannot find \"2\"\n")
		return
	}
	chk.Int(tst, "2", *node2.Data, 2)
	chk.String(tst, list.String(), "[0 3 2 1]")

	io.Pf("\nfind node \"0\" in list\n")
	res := list.Find(func(node *IntLinkedNode) bool {
		if *node.Data == 0 {
			return true
		}
		return false
	})
	if res == nil {
		tst.Errorf("cannot find \"0\"\n")
		return
	}
	chk.Int(tst, "0", *res.Data, 0)
	chk.String(tst, list.String(), "[0 3 2 1]")

	io.Pf("\nremove node \"0\" from list. \"3\" becomes root\n")
	list.Remove(root)
	chk.Int(tst, "3", *list.root.Data, 3)
	chk.String(tst, list.String(), "[3 2 1]")

	io.Pf("\nfind node \"3\" in list\n")
	node3 := list.Find(func(node *IntLinkedNode) bool {
		if *node.Data == 3 {
			return true
		}
		return false
	})
	if node3 == nil {
		tst.Errorf("cannot find \"3\"\n")
		return
	}
	chk.Int(tst, "3", *node3.Data, 3)
	chk.String(tst, list.String(), "[3 2 1]")

	io.Pf("\nremove node \"3\" from list. \"2\" becomes root\n")
	list.Remove(node3)
	chk.Int(tst, "2", *list.root.Data, 2)
	chk.String(tst, list.String(), "[2 1]")

	io.Pf("\nfind node \"1\" in list\n")
	node1 := list.Find(func(node *IntLinkedNode) bool {
		if *node.Data == 1 {
			return true
		}
		return false
	})
	if node1 == nil {
		tst.Errorf("cannot find \"1\"\n")
		return
	}
	chk.Int(tst, "1", *node1.Data, 1)
	chk.String(tst, list.String(), "[2 1]")

	io.Pf("\nremove node \"1\" from list. \"2\" remains as root\n")
	list.Remove(node1)
	chk.Int(tst, "2", *list.root.Data, 2)
	chk.String(tst, list.String(), "[2]")
}

func TestIntLinkedList02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("IntLinkedList02. edge cases")

	list := NewIntLinkedList()
	list.Remove(nil)
	root := list.Insert(0)
	chk.String(tst, list.String(), "[0]")
	list.Remove(root)
	chk.String(tst, list.String(), "[]")
	res := list.Find(func(node *IntLinkedNode) bool {
		return false // will find nothing
	})
	if res != nil {
		tst.Errorf("Find should have returned <nil>\n")
		return
	}
}

// float64 //////////////////////////////////////////////////////////////////////////////////////////

func TestFloat64LinkedList01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Float64LinkedList01. basic functionality")

	io.Pf("create list and insert first node\n")
	list := NewFloat64LinkedList()
	root := list.Insert(0)
	chk.Float64(tst, "root == 0", 1e-15, *root.Data, 0)
	chk.String(tst, list.String(), "[0]")

	io.Pf("\ninsert next node \"1\" and traverse list\n")
	node := list.Insert(1)
	chk.Float64(tst, "1", 1e-15, *node.Data, 1)
	chk.String(tst, list.String(), "[0 1]")

	io.Pf("\ninsert next node \"2\" and traverse list\n")
	node = list.Insert(2)
	chk.Float64(tst, "2", 1e-15, *node.Data, 2)
	chk.String(tst, list.String(), "[0 2 1]")

	io.Pf("\ninsert next node \"3\" and traverse list\n")
	node = list.Insert(3)
	chk.Float64(tst, "3", 1e-15, *node.Data, 3)
	chk.String(tst, list.String(), "[0 3 2 1]")

	io.Pf("\nfind node \"2\" in list\n")
	node2 := list.Find(func(node *Float64LinkedNode) bool {
		if *node.Data == 2 {
			return true
		}
		return false
	})
	if node2 == nil {
		tst.Errorf("cannot find \"2\"\n")
		return
	}
	chk.Float64(tst, "2", 1e-15, *node2.Data, 2)
	chk.String(tst, list.String(), "[0 3 2 1]")

	io.Pf("\nfind node \"0\" in list\n")
	res := list.Find(func(node *Float64LinkedNode) bool {
		if *node.Data == 0 {
			return true
		}
		return false
	})
	if res == nil {
		tst.Errorf("cannot find \"0\"\n")
		return
	}
	chk.Float64(tst, "0", 1e-15, *res.Data, 0)
	chk.String(tst, list.String(), "[0 3 2 1]")

	io.Pf("\nremove node \"0\" from list. \"3\" becomes root\n")
	list.Remove(root)
	chk.Float64(tst, "3", 1e-15, *list.root.Data, 3)
	chk.String(tst, list.String(), "[3 2 1]")

	io.Pf("\nfind node \"3\" in list\n")
	node3 := list.Find(func(node *Float64LinkedNode) bool {
		if *node.Data == 3 {
			return true
		}
		return false
	})
	if node3 == nil {
		tst.Errorf("cannot find \"3\"\n")
		return
	}
	chk.Float64(tst, "3", 1e-15, *node3.Data, 3)
	chk.String(tst, list.String(), "[3 2 1]")

	io.Pf("\nremove node \"3\" from list. \"2\" becomes root\n")
	list.Remove(node3)
	chk.Float64(tst, "2", 1e-15, *list.root.Data, 2)
	chk.String(tst, list.String(), "[2 1]")

	io.Pf("\nfind node \"1\" in list\n")
	node1 := list.Find(func(node *Float64LinkedNode) bool {
		if *node.Data == 1 {
			return true
		}
		return false
	})
	if node1 == nil {
		tst.Errorf("cannot find \"1\"\n")
		return
	}
	chk.Float64(tst, "1", 1e-15, *node1.Data, 1)
	chk.String(tst, list.String(), "[2 1]")

	io.Pf("\nremove node \"1\" from list. \"2\" remains as root\n")
	list.Remove(node1)
	chk.Float64(tst, "2", 1e-15, *list.root.Data, 2)
	chk.String(tst, list.String(), "[2]")
}

func TestFloat64LinkedList02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Float64LinkedList02. edge cases")

	list := NewFloat64LinkedList()
	list.Remove(nil)
	root := list.Insert(0)
	chk.String(tst, list.String(), "[0]")
	list.Remove(root)
	chk.String(tst, list.String(), "[]")
	res := list.Find(func(node *Float64LinkedNode) bool {
		return false // will find nothing
	})
	if res != nil {
		tst.Errorf("Find should have returned <nil>\n")
		return
	}
}

// string //////////////////////////////////////////////////////////////////////////////////////////

func TestStringLinkedList01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("StringLinkedList01. basic functionality")

	io.Pf("create list and insert first node\n")
	list := NewStringLinkedList()
	root := list.Insert("root")
	chk.String(tst, *root.Data, "root")
	chk.String(tst, list.String(), "[root]")

	io.Pf("\ninsert next node \"A\" and traverse list\n")
	node := list.Insert("A")
	chk.String(tst, *node.Data, "A")
	chk.String(tst, list.String(), "[root A]")

	io.Pf("\ninsert next node \"B\" and traverse list\n")
	node = list.Insert("B")
	chk.String(tst, *node.Data, "B")
	chk.String(tst, list.String(), "[root B A]")

	io.Pf("\ninsert next node \"C\" and traverse list\n")
	node = list.Insert("C")
	chk.String(tst, *node.Data, "C")
	chk.String(tst, list.String(), "[root C B A]")

	io.Pf("\nfind node \"B\" in list\n")
	nodeB := list.Find(func(node *StringLinkedNode) bool {
		if *node.Data == "B" {
			return true
		}
		return false
	})
	if nodeB == nil {
		tst.Errorf("cannot find \"B\"\n")
		return
	}
	chk.String(tst, *nodeB.Data, "B")
	chk.String(tst, list.String(), "[root C B A]")

	io.Pf("\nfind node \"root\" in list\n")
	res := list.Find(func(node *StringLinkedNode) bool {
		if *node.Data == "root" {
			return true
		}
		return false
	})
	if res == nil {
		tst.Errorf("cannot find \"root\"\n")
		return
	}
	chk.String(tst, *res.Data, "root")
	chk.String(tst, list.String(), "[root C B A]")

	io.Pf("\nremove node \"root\" from list. \"C\" becomes root\n")
	list.Remove(root)
	chk.String(tst, *list.root.Data, "C")
	chk.String(tst, list.String(), "[C B A]")

	io.Pf("\nfind node \"C\" in list\n")
	nodeC := list.Find(func(node *StringLinkedNode) bool {
		if *node.Data == "C" {
			return true
		}
		return false
	})
	if nodeC == nil {
		tst.Errorf("cannot find \"C\"\n")
		return
	}
	chk.String(tst, *nodeC.Data, "C")
	chk.String(tst, list.String(), "[C B A]")

	io.Pf("\nremove node \"C\" from list. \"B\" becomes root\n")
	list.Remove(nodeC)
	chk.String(tst, *list.root.Data, "B")
	chk.String(tst, list.String(), "[B A]")

	io.Pf("\nfind node \"A\" in list\n")
	nodeA := list.Find(func(node *StringLinkedNode) bool {
		if *node.Data == "A" {
			return true
		}
		return false
	})
	if nodeA == nil {
		tst.Errorf("cannot find \"A\"\n")
		return
	}
	chk.String(tst, *nodeA.Data, "A")
	chk.String(tst, list.String(), "[B A]")

	io.Pf("\nremove node \"A\" from list. \"B\" remains as root\n")
	list.Remove(nodeA)
	chk.String(tst, *list.root.Data, "B")
	chk.String(tst, list.String(), "[B]")
}

func TestStringLinkedList02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("StringLinkedList02. edge cases")

	list := NewStringLinkedList()
	list.Remove(nil)
	root := list.Insert("root")
	chk.String(tst, list.String(), "[root]")
	list.Remove(root)
	chk.String(tst, list.String(), "[]")
	res := list.Find(func(node *StringLinkedNode) bool {
		return false // will find nothing
	})
	if res != nil {
		tst.Errorf("Find should have returned <nil>\n")
		return
	}
}
