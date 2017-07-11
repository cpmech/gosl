// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package io

import "github.com/cpmech/gosl/chk"

// Enum defines enumeration constants
type Enum int

// enumData holds data about an enumeration constant
type enumData struct {
	Name   string // the name; e.g. "Normal"
	Prefix string // prefix to classify enumerations; e.g. "rnd" [optional]
	Key    string // a short key; e.g. "N" [optional]
	Desc   string // a description; e.g. "Normal Distribution" [optional]
}

// enums holds all enumerations
var enums []enumData

// enumsMap maps enum "Prefix.Name" to index
var enumsMap map[string]Enum

// String prints the name of an enumeration constant
func (o Enum) String() string {
	return enums[int(o)].Name
}

// Prefix prints the prefix string of an enumeration constant
func (o Enum) Prefix() string {
	return enums[int(o)].Prefix
}

// Key prints the short key corresponding to an enumeration constant
func (o Enum) Key() string {
	return enums[int(o)].Key
}

// Desc prints the description about an enumeration constant
func (o Enum) Desc() string {
	return enums[int(o)].Desc
}

// EnumsFind finds enum by "Prefix.Name"
// returns -1 if not found
func EnumsFind(prefix, name string) (enum Enum) {
	enum, found := enumsMap[prefix+"."+name]
	if !found {
		return -1
	}
	return
}

// NewEnum sets a new enumeration
//   arguments -- Name, Prefix [optional], Key [optional], Desc [optional]
//   Prefix is used to generate a unique map key such as:
//      "Prefix.Name"
func NewEnum(arguments ...string) Enum {
	narg := len(arguments)
	if narg < 1 {
		chk.Panic("at least the name of enum must be provided")
	}
	var name, prefix, key, desc string
	if narg > 0 {
		name = arguments[0]
	}
	if narg > 1 {
		prefix = arguments[1]
	}
	if narg > 2 {
		key = arguments[2]
	}
	if narg > 3 {
		desc = arguments[3]
	}
	if e := EnumsFind(prefix, name); e >= 0 {
		chk.Panic("enum \"%s.%s\" exists already", prefix, name)
	}
	enums = append(enums, enumData{name, prefix, key, desc})
	var enum = Enum(len(enums) - 1)
	if enumsMap == nil {
		enumsMap = make(map[string]Enum)
	}
	enumsMap[prefix+"."+name] = enum
	return enum
}
