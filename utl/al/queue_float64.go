// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package al

import (
	"gosl/io"
)

// Float64Queue implements a FIFO queue, a sequence where the first inserted will be the first removed.
// Think of arriving people at the Bank or DMV...
type Float64Queue struct {
	bfSize   int        // guessed buffer size
	front    int        // index in ring of member at front
	back     int        // index in ring of member at back
	nMembers int        // current number of members
	ring     []*float64 // ring holds all data in a "ring fashion"
	Debug    bool       // debug flag
}

// NewFloat64Queue returns a new object
func NewFloat64Queue(guessedBufferSize int) (o *Float64Queue) {
	o = new(Float64Queue)
	o.bfSize = guessedBufferSize
	o.front = -1 // indicates first ring
	o.back = -1  // indicates first ring
	return
}

// Front returns the member @ front of queue (close to the DMV window...) or nil if empty
func (o *Float64Queue) Front() *float64 {
	if o.nMembers == 0 {
		return nil
	}
	return o.ring[o.front]
}

// Back returns the member @ back (unlucky guy/girl...) or nil if empty.
// It is always the last item in the data array
func (o *Float64Queue) Back() *float64 {
	if o.nMembers == 0 {
		return nil
	}
	return o.ring[o.back]
}

// Nmembers returns the length of queue; i.e. the number of members
func (o *Float64Queue) Nmembers() int {
	return o.nMembers
}

// In receives a new member arrival
// TODO: implement use of different grow rates
func (o *Float64Queue) In(member float64) {

	// debug
	if o.Debug {
		io.Pfgrey("in  : before: F=%d B=%d N=%d ring=%q\n", o.front, o.back, o.nMembers, sliceFloat64ToString(o.ring))
		defer func() {
			io.Pfyel("in  : after : F=%d B=%d N=%d ring=%q\n", o.front, o.back, o.nMembers, sliceFloat64ToString(o.ring))
		}()
	}

	// first ring
	if o.front < 0 {
		o.ring = make([]*float64, 1, o.bfSize+1)
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
func (o *Float64Queue) Out() (member *float64) {

	// debug
	if o.Debug {
		io.Pfpink("out : before: F=%d B=%d N=%d ring=%q\n", o.front, o.back, o.nMembers, sliceFloat64ToString(o.ring))
		defer func() {
			io.Pfpink("out : after : F=%d B=%d N=%d ring=%q\n", o.front, o.back, o.nMembers, sliceFloat64ToString(o.ring))
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
func (o *Float64Queue) String() (l string) {
	if o.nMembers == 0 {
		return "[]"
	}
	if o.back < o.front {
		left := o.ring[o.front:]
		right := o.ring[:o.back+1]
		return "[" + sliceFloat64ToString(left) + " " + sliceFloat64ToString(right) + "]"
	}
	return "[" + sliceFloat64ToString(o.ring[o.front:o.back+1]) + "]"
}

// auxiliary ////////////////////////////////////////////////////////////////////////////////////

// grow grows ring
func (o *Float64Queue) grow() {

	// debug
	if o.Debug {
		io.Pfblue("grow: before: F=%d B=%d N=%d ring=%q\n", o.front, o.back, o.nMembers, sliceFloat64ToString(o.ring))
		defer func() {
			io.Pfblue("grow: after : F=%d B=%d N=%d ring=%q\n", o.front, o.back, o.nMembers, sliceFloat64ToString(o.ring))
		}()
	}

	// temporary array
	tmp := make([]*float64, o.nMembers+1, o.bfSize+o.nMembers+1)

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

// sliceFloat64ToString string converts slice of pointers to string
func sliceFloat64ToString(input []*float64) (l string) {
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
