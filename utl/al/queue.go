// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package al

import (
	"strings"

	"github.com/cpmech/gosl/io"
)

// Queue implements a FIFO queue, a sequence where the first inserted will be the first removed.
// Think of arriving people at the Bank or DMV...
type Queue struct {
	bfSize   int     // guessed buffer size
	front    int     // index in ring of member at front
	back     int     // index in ring of member at back
	nMembers int     // current number of members
	ring     []Adata // ring holds all data in a "ring fashion"
	Debug    bool    // debug flag
}

// NewQueue returns a new object
func NewQueue(guessedBufferSize int) (o *Queue) {
	o = new(Queue)
	o.bfSize = guessedBufferSize
	o.front = -1 // indicates first ring
	o.back = -1  // indicates first ring
	return
}

// Front returns the member @ front of queue (close to the DMV window...) or nil if empty
func (o *Queue) Front() Adata {
	if o.nMembers == 0 {
		return nil
	}
	return o.ring[o.front]
}

// Back returns the member @ back (unlucky guy/girl...) or nil if empty.
// It is always the last item in the data array
func (o *Queue) Back() Adata {
	if o.nMembers == 0 {
		return nil
	}
	return o.ring[o.back]
}

// Nmembers returns the length of queue; i.e. the number of members
func (o *Queue) Nmembers() int {
	return o.nMembers
}

// In receives a new member arrival
// TODO: implement use of different grow rates
func (o *Queue) In(member Adata) {

	// debug
	if o.Debug {
		io.Pfgrey("in  : before: F=%d B=%d N=%d ring=%v\n", o.front, o.back, o.nMembers, o.ring)
		defer func() {
			io.Pfyel("in  : after : F=%d B=%d N=%d ring=%v\n", o.front, o.back, o.nMembers, o.ring)
		}()
	}

	// first ring
	if o.front < 0 {
		o.ring = make([]Adata, 1, o.bfSize+1)
		o.ring[0] = member
		o.front = 0
		o.back = o.front
		o.nMembers = 1
		return
	}

	// no space available ⇒ grow ring
	if o.nMembers+1 > len(o.ring) {
		o.grow()
	}

	// updates
	o.back = (o.back + 1) % len(o.ring) // cyclic increment
	o.ring[o.back] = member
	o.nMembers++
}

// Out removes the member @ front and returns a pointer to him/her
// TODO: implement memory recovery
func (o *Queue) Out() (member Adata) {

	// debug
	if o.Debug {
		io.Pfpink("out : before: F=%d B=%d N=%d ring=%v\n", o.front, o.back, o.nMembers, o.ring)
		defer func() {
			io.Pfpink("out : after : F=%d B=%d N=%d ring=%v\n", o.front, o.back, o.nMembers, o.ring)
		}()
	}

	// no members
	if o.nMembers == 0 {
		return nil
	}

	// simply move Front pointer
	member = o.Front()
	o.front = (o.front + 1) % len(o.ring) // cyclic increment
	o.nMembers--
	return
}

// String returns the string representation of this object
func (o *Queue) String() (l string) {
	if o.nMembers == 0 {
		return "[]"
	}
	if o.back < o.front {
		left := o.ring[o.front:]
		right := o.ring[:o.back+1]
		return strings.Replace(io.Sf("%v", left)+io.Sf("%v", right), "][", " ", 1)
	}
	return io.Sf("%v", o.ring[o.front:o.back+1])
}

// auxiliary ////////////////////////////////////////////////////////////////////////////////////

// grow grows ring
func (o *Queue) grow() {

	// debug
	if o.Debug {
		io.Pfblue("grow: before: F=%d B=%d N=%d ring=%v\n", o.front, o.back, o.nMembers, o.ring)
		defer func() {
			io.Pfblue("grow: after : F=%d B=%d N=%d ring=%v\n", o.front, o.back, o.nMembers, o.ring)
		}()
	}

	// temporary array
	tmp := make([]Adata, o.nMembers+1, o.bfSize+o.nMembers+1)

	// members are at different sides
	if o.back < o.front {
		left := o.ring[o.front:]
		right := o.ring[:o.back+1]
		copy(tmp, left)
		copy(tmp[len(left):], right)

		// members are a the same side
	} else {
		copy(tmp, o.ring[o.front:o.back+1])
	}

	// set indices and replace ring. Note: nMembers remains unchanged
	o.front = 0
	o.back = o.nMembers - 1
	o.ring = tmp
}
