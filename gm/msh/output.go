// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package msh

import "github.com/cpmech/gosl/io"

// String returns a JSON representation of *Vert
func (o *Vertex) String() string {
	l := io.Sf("{\"i\":%d, \"t\":%d, \"x\":[", o.Id, o.Tag)
	for i, x := range o.X {
		if i > 0 {
			l += ", "
		}
		l += io.Sf("%g", x)
	}
	l += "] }"
	return l
}

// String returns a JSON representation of *Cell
func (o *Cell) String() string {
	l := io.Sf("{\"i\":%d, \"t\":%d, \"p\":%d, \"y\":%q, \"v\":[", o.Id, o.Tag, o.Part, o.TypeKey)
	for i, x := range o.V {
		if i > 0 {
			l += ", "
		}
		l += io.Sf("%d", x)
	}
	if len(o.EdgeTags) > 0 {
		l += "], \"et\":["
		for i, x := range o.EdgeTags {
			if i > 0 {
				l += ", "
			}
			l += io.Sf("%d", x)
		}
	}
	if len(o.FaceTags) > 0 {
		l += "], \"ft\":["
		for i, x := range o.FaceTags {
			if i > 0 {
				l += ", "
			}
			l += io.Sf("%d", x)
		}
	}
	l += "] }"
	return l
}
