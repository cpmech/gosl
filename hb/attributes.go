// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hb

import "gosl/chk"

// String //////////////////////////////////////////////////////////////////////////////////////

// SetStringAttribute sets a string attibute
func (o *File) SetStringAttribute(path, key, val string) {
	if o.gobReading {
		chk.Panic("cannot put %q because file is open for READONLY", path)
	}
	o.gobEnc.Encode("SetStringAttribute")
	o.gobEnc.Encode(path)
	o.gobEnc.Encode(key)
	o.gobEnc.Encode(val)
	return
}

// GetStringAttribute gets string attribute
func (o *File) GetStringAttribute(path, key string) (val string) {
	var cmd string
	o.gobDec.Decode(&cmd)
	if cmd != "SetStringAttribute" {
		chk.Panic("wrong command => %q\n(r/w commands need to be called in the same order)", cmd)
	}
	var rpath string
	o.gobDec.Decode(&rpath)
	if rpath != path {
		chk.Panic("cannot read path: %s != %s\n(r/w commands need to be called in the same order)", path, rpath)
	}
	var rkey string
	o.gobDec.Decode(&rkey)
	if rkey != key {
		chk.Panic("cannot read key: %s != %s\n(r/w commands need to be called in the same order)", key, rkey)
	}
	o.gobDec.Decode(&val)
	return
}

// Int /////////////////////////////////////////////////////////////////////////////////////////

// SetIntAttribute sets int attibute
func (o *File) SetIntAttribute(path, key string, val int) {
	if o.gobReading {
		chk.Panic("cannot put %q because file is open for READONLY", path)
	}
	o.gobEnc.Encode("SetIntAttribute")
	o.gobEnc.Encode(path)
	o.gobEnc.Encode(key)
	o.gobEnc.Encode(val)
	return
}

// GetIntAttribute gets int attribute
func (o *File) GetIntAttribute(path, key string) (val int) {
	var cmd string
	o.gobDec.Decode(&cmd)
	if cmd != "SetIntAttribute" {
		chk.Panic("wrong command => %q\n(r/w commands need to be called in the same order)", cmd)
	}
	var rpath string
	o.gobDec.Decode(&rpath)
	if rpath != path {
		chk.Panic("cannot read path: %s != %s\n(r/w commands need to be called in the same order)", path, rpath)
	}
	var rkey string
	o.gobDec.Decode(&rkey)
	if rkey != key {
		chk.Panic("cannot read key: %s != %s\n(r/w commands need to be called in the same order)", key, rkey)
	}
	o.gobDec.Decode(&val)
	return
}

// Ints /////////////////////////////////////////////////////////////////////////////////////////

// SetIntsAttribute sets slice-of-ints attibute
func (o *File) SetIntsAttribute(path, key string, vals []int) {
	if o.gobReading {
		chk.Panic("cannot put %q because file is open for READONLY", path)
	}
	o.gobEnc.Encode("SetIntsAttribute")
	o.gobEnc.Encode(path)
	o.gobEnc.Encode(key)
	o.gobEnc.Encode(vals)
	return
}

// GetIntsAttribute gets slice-of-ints attribute
func (o *File) GetIntsAttribute(path, key string) (vals []int) {
	var cmd string
	o.gobDec.Decode(&cmd)
	if cmd != "SetIntsAttribute" {
		chk.Panic("wrong command => %q\n(r/w commands need to be called in the same order)", cmd)
	}
	var rpath string
	o.gobDec.Decode(&rpath)
	if rpath != path {
		chk.Panic("cannot read path: %s != %s\n(r/w commands need to be called in the same order)", path, rpath)
	}
	var rkey string
	o.gobDec.Decode(&rkey)
	if rkey != key {
		chk.Panic("cannot read key: %s != %s\n(r/w commands need to be called in the same order)", key, rkey)
	}
	o.gobDec.Decode(&vals)
	return
}
