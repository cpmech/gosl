// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plt

// C returns a color from a default palette
func C(i, palette int) string {
	if palette < 0 || palette >= len(palettes) {
		return ""
	}
	p := palettes[palette]
	return p[i%len(p)]
}

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
