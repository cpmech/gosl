// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

import (
	"strings"
)

// Keycode extracts a keycode from a string such as "!typeA:keycodeA !typeB:keycodeB!typeC:keycodeC"
//  Note: String == "!keyA !typeB:valB" is also valid
func Keycode(String string, Type string) (keycode string, found bool) {
	if len(String) < 1 || Type == "" {
		return "", false
	}
	if String[0] != '!' {
		Panic(_parsing_err3, "Keycode", String)
	}
	for _, s := range strings.Split(String[1:], "!") { // [1:] => skip first "!"
		ss := strings.TrimSpace(s)
		sss := strings.Split(ss, ":")
		//Pforan("ss=%#v, sss=%#v\n", ss, sss)
		for i, val := range sss {
			sss[i] = strings.TrimSpace(val)
		}
		//if len(sss) < 1 { // TODO: remove this since this will not happen
		//Panic(_parsing_err4, "Keycode", String)
		//}
		//Pfgreen("ss=%#v, sss=%#v\n", ss, sss)
		if sss[0] == Type {
			if len(sss) > 1 {
				return sss[1], true
			} else {
				return "", true
			}
		}
	}
	return
}

// Keycodes extracts keys from a keycode (extra) string
//  Example: "!keyA !typeB:valB" => keycodes = [keyA, typeB]
func Keycodes(String string) (keycodes []string) {
	if len(String) < 1 {
		return
	}
	if String[0] != '!' {
		Panic(_parsing_err3, "Keycodes", String)
	}
	for _, s := range strings.Split(String[1:], "!") { // [1:] => skip first "!"
		ss := strings.TrimSpace(s)
		sss := strings.Split(ss, ":")
		for i, val := range sss {
			sss[i] = strings.TrimSpace(val)
		}
		//if len(sss) < 1 { // TODO: remove this since this will not happen
		//Panic(_parsing_err4, "Keycodes", String)
		//}
		keycodes = append(keycodes, sss[0])
	}
	return
}

// JoinKeys3 joins keys from 3 slices into a string with sets separated by sep
// and keys separeted by spaces
func JoinKeys3(k0, k1, k2 []string, sep string) (res string) {
	res = strings.Join(k0, " ") + sep
	if len(k1) > 0 {
		res += " " + strings.Join(k1, " ")
	}
	res += sep
	if len(k2) > 0 {
		res += " " + strings.Join(k2, " ")
	}
	return
}

// SplitKeys3 splits a string with three sets of keys separated by comma
func SplitKeys3(res string) (k0, k1, k2 []string) {
	sets := strings.Split(res, ",")
	if len(sets) != 3 {
		Panic(_parsing_err1, res, sets)
	}
	s0 := strings.TrimSpace(sets[0])
	s1 := strings.TrimSpace(sets[1])
	s2 := strings.TrimSpace(sets[2])
	if len(s0) > 0 {
		k0 = strings.Split(s0, " ")
	}
	if len(s1) > 0 {
		k1 = strings.Split(s1, " ")
	}
	if len(s2) > 0 {
		k2 = strings.Split(s2, " ")
	}
	return
}

// JoinKeys4 joins keys from 4 slices into a string with sets separated by sep
// and keys separeted by spaces
func JoinKeys4(k0, k1, k2, k3 []string, sep string) (res string) {
	res = strings.Join(k0, " ") + sep
	if len(k1) > 0 {
		res += " " + strings.Join(k1, " ")
	}
	res += sep
	if len(k2) > 0 {
		res += " " + strings.Join(k2, " ")
	}
	res += sep
	if len(k3) > 0 {
		res += " " + strings.Join(k3, " ")
	}
	return
}

// SplitKeys4 splits a string with four sets of keys separated by comma
func SplitKeys4(res string) (k0, k1, k2, k3 []string) {
	sets := strings.Split(res, ",")
	if len(sets) != 4 {
		Panic(_parsing_err2, res, sets)
	}
	s0 := strings.TrimSpace(sets[0])
	s1 := strings.TrimSpace(sets[1])
	s2 := strings.TrimSpace(sets[2])
	s3 := strings.TrimSpace(sets[3])
	if len(s0) > 0 {
		k0 = strings.Split(s0, " ")
	}
	if len(s1) > 0 {
		k1 = strings.Split(s1, " ")
	}
	if len(s2) > 0 {
		k2 = strings.Split(s2, " ")
	}
	if len(s3) > 0 {
		k3 = strings.Split(s3, " ")
	}
	return
}

// JoinKeys join keys separeted by spaces
func JoinKeys(keys []string) string {
	return strings.Join(keys, " ")
}

// JoinKeysPre join keys separeted by spaces with a prefix
func JoinKeysPre(prefix string, keys []string) (res string) {
	for i, val := range keys {
		res += prefix + val
		if i < len(keys)-1 {
			res += " "
		}
	}
	return
}

// SplitKeys splits keys separeted by spaces
func SplitKeys(keys string) []string {
	return strings.Split(keys, " ")
}

// error messages
var (
	_parsing_err1 = "parsing.go: SplitKeys3: string '%s' does not contain 3 subsets separated by 2 commas. sets = %v"
	_parsing_err2 = "parsing.go: SplitKeys4: string '%s' does not contain 4 subsets separated by 3 commas. sets = %v"
	_parsing_err3 = "parsing.go: %s: first character in keycode string must be an exclamation mark !\nstring = \"%v\""
	//_parsing_err4 = "parsing.go: %s: keycode string is invalid:\nstring = \"%v\""
)
