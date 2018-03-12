// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plt

// C returns a color from a default palette
//  use palette < 0 for automatic color
func C(i, palette int) string {
	if palette < 0 || palette >= len(palettes) {
		return ""
	}
	p := palettes[palette]
	return p[i%len(p)]
}

// M returns a marker
//  use scheme < 0 for no marker
func M(i, scheme int) string {
	if scheme < 0 {
		return ""
	}
	if scheme >= len(markers) {
		scheme = 0
	}
	s := markers[scheme]
	return s[i%len(s)]
}

// palettes holds color palettes
var palettes = [][]string{
	{"#003fff", "#35b052", "#e8000b", "#8a2be2", "#ffc400", "#00d7ff"},
	{"blue", "green", "magenta", "orange", "red", "cyan", "black", "#de9700", "#89009d", "#7ad473", "#737ad4", "#d473ce", "#7e6322", "#462222", "#98ac9d", "#37a3e8", "yellow"},
	{"#4c72b0", "#55a868", "#c44e52", "#8172b2", "#ccb974", "#64b5cd"},
	{"#9b59b6", "#3498db", "#95a5a6", "#e74c3c", "#34495e", "#2ecc71"},
	{"#e41a1c", "#377eb8", "#4daf4a", "#984ea3", "#ff7f00", "#ffff33"},
	{"#7fc97f", "#beaed4", "#fdc086", "#ffff99", "#386cb0", "#f0027f", "#bf5b17"},
	{"#001c7f", "#017517", "#8c0900", "#7600a1", "#b8860b", "#006374"},
	{"#0072b2", "#009e73", "#d55e00", "#cc79a7", "#f0e442", "#56b4e9"},
	{"#4878cf", "#6acc65", "#d65f5f", "#b47cc7", "#c4ad66", "#77bedb"},
	{"#92c6ff", "#97f0aa", "#ff9f9a", "#d0bbff", "#fffea3", "#b0e0e6"},
}

// markers holds marker types
//    "   none
//    .   point
//    +   plus
//    x   x
//    *   star
//    d   thin_diamond
//    o   circle
//    s   square
//    ^   triangle_up
//    v   triangle_down
//    <   triangle_left
//    >   triangle_right
//    8   octagon
//    p   pentagon
//    h   hexagon1
//    H   hexagon2
//    D   diamond
//    |   vline
//    _   hline
//    1   tri_down
//    2   tri_up
//    3   tri_left
//    4   tri_right
//    ,   pixel
var markers = [][]string{
	{".", "+", "x", "*", "^", "s", "d", "p", "v"},
	{".", "+", "x", "*", "^", "s", "d", "p", "v", "|", "o", "_", "<", "1", ">", "2", "8", "3", "h", "4", "D"},
	{".", "s", "^", "+", "x", "o", "d", "p", "v", "|", "_", "<", "1", ">", "2", "8", "3", "h", "4", "D", "*"},
}
