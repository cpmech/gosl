// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package al

import (
	"gosl/io"
)

// StringQueue implements a FIFO queue, a sequence where the first inserted will be the first removed.
// Think of arriving people at the Bank or DMV...
type StringQueue struct {
	bfSize   int       // guessed buffer size
	front    int       // index in ring of member at front
	back     int       // index in ring of member at back
	nMembers int       // current number of members
	ring     []*string // ring holds all data in a "ring fashion"
	Debug    bool      // debug flag
}

// NewStringQueue returns a new object
func NewStringQueue(guessedBufferSize int) (o *StringQueue) {
	o = new(StringQueue)
	o.bfSize = guessedBufferSize
	o.front = -1 // indicates first ring
	o.back = -1  // indicates first ring
	return
}

// Front returns the member @ front of queue (close to the DMV window...) or nil if empty
func (o *StringQueue) Front() *string {
	if o.nMembers == 0 {
		return nil
	}
	return o.ring[o.front]
}

// Back returns the member @ back (unlucky guy/girl...) or nil if empty.
// It is always the last item in the data array
func (o *StringQueue) Back() *string {
	if o.nMembers == 0 {
		return nil
	}
	return o.ring[o.back]
}

// Nmembers returns the length of queue; i.e. the number of members
func (o *StringQueue) Nmembers() int {
	return o.nMembers
}

// In receives a new member arrival
// TODO: implement use of different grow rates
func (o *StringQueue) In(member string) {

	// debug
	if o.Debug {
		io.Pfgrey("in  : before: F=%d B=%d N=%d ring=%q\n", o.front, o.back, o.nMembers, sliceStringToString(o.ring))
		defer func() {
			io.Pfyel("in  : after : F=%d B=%d N=%d ring=%q\n", o.front, o.back, o.nMembers, sliceStringToString(o.ring))
		}()
	}

	// first ring
	if o.front < 0 {
		o.ring = make([]*string, 1, o.bfSize+1)
		o.ring[0] = &member
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
	o.ring[o.back] = &member
	o.nMembers++
}

// Out removes the member @ front and returns a pointer to him/her
// TODO: implement memory recovery
func (o *StringQueue) Out() (member *string) {

	// debug
	if o.Debug {
		io.Pfpink("out : before: F=%d B=%d N=%d ring=%q\n", o.front, o.back, o.nMembers, sliceStringToString(o.ring))
		defer func() {
			io.Pfpink("out : after : F=%d B=%d N=%d ring=%q\n", o.front, o.back, o.nMembers, sliceStringToString(o.ring))
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
func (o *StringQueue) String() (l string) {
	if o.nMembers == 0 {
		return "[]"
	}
	if o.back < o.front {
		left := o.ring[o.front:]
		right := o.ring[:o.back+1]
		return "[" + sliceStringToString(left) + " " + sliceStringToString(right) + "]"
	}
	return "[" + sliceStringToString(o.ring[o.front:o.back+1]) + "]"
}

// auxiliary ////////////////////////////////////////////////////////////////////////////////////

// grow grows ring
func (o *StringQueue) grow() {

	// debug
	if o.Debug {
		io.Pfblue("grow: before: F=%d B=%d N=%d ring=%q\n", o.front, o.back, o.nMembers, sliceStringToString(o.ring))
		defer func() {
			io.Pfblue("grow: after : F=%d B=%d N=%d ring=%q\n", o.front, o.back, o.nMembers, sliceStringToString(o.ring))
		}()
	}

	// temporary array
	tmp := make([]*string, o.nMembers+1, o.bfSize+o.nMembers+1)

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

// sliceStringToString string converts slice of pointers to string
func sliceStringToString(input []*string) (l string) {
	for i := 0; i < len(input); i++ {
		if input[i] == nil {
			continue
		}
		if i > 0 {
			l += " "
		}
		l += io.Sf("%v", *input[i])
	}
	return l
}
