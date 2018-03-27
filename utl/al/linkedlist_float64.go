// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package al implements classical algorithms such as Queue, Stack and others
package al

import "github.com/cpmech/gosl/io"

// Float64LinkedNode defines a node in a Doubly Linked List
type Float64LinkedNode struct {
	prev *Float64LinkedNode
	next *Float64LinkedNode
	Data *float64
}

// Float64LinkedList defines Doubly Linked List
type Float64LinkedList struct {
	root *Float64LinkedNode // root has prev always <nil>
}

// NewFloat64LinkedList returns a new Doubly Linked List object
func NewFloat64LinkedList() (o *Float64LinkedList) {
	o = new(Float64LinkedList)
	return
}

// Insert inserts data just after root and returns the inserted node
func (o *Float64LinkedList) Insert(data float64) (newNode *Float64LinkedNode) {
	if o.root == nil { // first node
		o.root = &Float64LinkedNode{nil, nil, &data}
		return o.root
	}
	newNode = &Float64LinkedNode{
		prev: o.root,
		next: o.root.next,
		Data: &data,
	}
	if o.root.next != nil { // make sure to coonnect root's next node
		o.root.next.prev = newNode
	}
	o.root.next = newNode
	return
}

// Remove removes node from Doubly Linked List
func (o *Float64LinkedList) Remove(node *Float64LinkedNode) {
	if node == nil { // nothing to remove
		return
	}
	if node.prev != nil { // fix prev
		node.prev.next = node.next
	}
	if node.next != nil { // fix next
		node.next.prev = node.prev
	}
	if node.prev == nil { // new root
		o.root = o.root.next
	}
}

// Traverse traverses the Doubly Linked List and executes action(node)
// Note action(node) may never be called if there aren't any nodes in the list
func (o *Float64LinkedList) Traverse(action func(node *Float64LinkedNode) (stop bool)) {
	if o.root == nil { // list is empty
		return
	}
	node := o.root
	for {
		stop := action(node)
		if stop {
			return
		}
		if node.next == nil {
			break
		}
		node = node.next
	}
}

// Find finds a node by traversing the list and comparing a to b
func (o *Float64LinkedList) Find(condition func(node *Float64LinkedNode) bool) (found *Float64LinkedNode) {
	if o.root == nil { // list is empty
		return
	}
	o.Traverse(func(node *Float64LinkedNode) (stop bool) {
		if condition(node) {
			found = node
			return true // stop
		}
		return false // continue
	})
	return
}

// String returns a string representation of this list, after traversing all nodes
func (o *Float64LinkedList) String() (l string) {
	first := true
	l = "["
	o.Traverse(func(node *Float64LinkedNode) (stop bool) {
		if !first {
			l += " "
		}
		first = false
		l += io.Sf("%v", *node.Data)
		return
	})
	l += "]"
	return
}
