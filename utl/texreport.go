// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

import (
	"strings"

	"github.com/cpmech/gosl/io"
)

// TexNum returns a string representation in TeX format of a real number.
// scientificNotation:
//   peforms the conversion of numbers into scientific notation where
//   the exponent notation with e{+-}{##} is converted into \cdot 10^{{+-}##}
func TexNum(fmt string, num float64, scientificNotation bool) (l string) {
	if fmt == "" {
		fmt = "%g"
	}
	l = io.Sf(fmt, num)
	if scientificNotation {
		s := strings.Split(l, "e-")
		if len(s) == 2 {
			e := s[1]
			if e == "00" {
				l = s[0]
				return
			}
			if e[0] == '0' {
				e = string(e[1])
			}
			l = s[0] + "\\cdot 10^{-" + e + "}"
		}
		s = strings.Split(l, "e+")
		if len(s) == 2 {
			e := s[1]
			if e == "00" {
				l = s[0]
				return
			}
			if e[0] == '0' {
				e = string(e[1])
			}
			l = s[0] + "\\cdot 10^{+" + e + "}"
		}
	}
	return
}
