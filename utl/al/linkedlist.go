// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package al implements classical algorithms such as Queue, Stack and others
package al

import "github.com/cpmech/gosl/io"

// LinkedNode defines a node in a Doubly Linked List
type LinkedNode struct {
	prev *LinkedNode
	next *LinkedNode
	Data Adata
}

// LinkedList defines Doubly Linked List
type LinkedList struct {
	root *LinkedNode // root has prev always <nil>
}

// NewLinkedList returns a new Doubly Linked List object
func NewLinkedList() (o *LinkedList) {
	o = new(LinkedList)
	return
}

// Insert inserts data just after root and returns the inserted node
func (o *LinkedList) Insert(data Adata) (newNode *LinkedNode) {
	if o.root == nil { // first node
		o.root = &LinkedNode{nil, nil, data}
		return o.root
	}
	newNode = &LinkedNode{
		prev: o.root,
		next: o.root.next,
		Data: data,
	}
	if o.root.next != nil { // make sure to coonnect root's next node
		o.root.next.prev = newNode
	}
	o.root.next = newNode
	return
}

// Remove removes node from Doubly Linked List
func (o *LinkedList) Remove(node *LinkedNode) {
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
func (o *LinkedList) Traverse(action func(node *LinkedNode) (stop bool)) {
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
func (o *LinkedList) Find(condition func(node *LinkedNode) bool) (found *LinkedNode) {
	if o.root == nil { // list is empty
		return
	}
	o.Traverse(func(node *LinkedNode) (stop bool) {
		if condition(node) {
			found = node
			return true // stop
		}
		return false // continue
	})
	return
}

// String returns a string representation of this list, after traversing all nodes
func (o *LinkedList) String() (l string) {
	first := true
	l = "["
	o.Traverse(func(node *LinkedNode) (stop bool) {
		if !first {
			l += " "
		}
		first = false
		l += io.Sf("%v", node.Data)
		return
	})
	l += "]"
	return
}
