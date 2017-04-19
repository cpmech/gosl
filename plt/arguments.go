// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plt

import (
	"bytes"

	"github.com/cpmech/gosl/io"
)

// 'A' holds "arguments" to configure plots, including "style" data for shapes (e.g. polygons)
type A struct {

	// plot and basic options
	C      string  // color
	M      string  // marker
	Ls     string  // linestyle
	Lw     float64 // linewidth; -1 => default
	Ms     int     // marker size; -1 => default
	L      string  // label
	Me     int     // mark-every; -1 => default
	Z      int     // z-order
	Mec    string  // marker edge color
	Mew    float64 // marker edge width
	Void   bool    // void marker => markeredgecolor='C', markerfacecolor='none'
	NoClip bool    // turn clipping off

	// shapes
	Fc     string  // shapes: face color
	Ec     string  // shapes: edge color
	Scale  float64 // shapes: scale information
	Style  string  // shapes: style information
	Closed bool    // shapes: closed shape

	// text and extra arguments
	Ha      string  // horizontal alignment; e.g. 'center'
	Va      string  // vertical alignment; e.g. 'center'
	Rot     float64 // rotation
	Fsz     float64 // font size
	FszLbl  float64 // font size of labels
	FszLeg  float64 // font size of legend
	FszXtck float64 // font size of x-ticks
	FszYtck float64 // font size of y-ticks
	HideL   bool    // hide left frame border
	HideR   bool    // hide right frame border
	HideB   bool    // hide bottom frame border
	HideT   bool    // hide top frame border

	// other options
	FigFraction bool // the given x-y coordinates correspond to figure coords  "xycoords='figure fraction'") }

	// legend
	LegLoc   string    // legend: location
	LegNcol  int       // legend: number of columns
	LegHlen  float64   // legend: handle length
	LegFrame bool      // legend: frame on
	LegOut   bool      // legend: outside
	LegOutX  []float64 // legend: normalised coordinates to put legend outside frame

	// colors for contours or histograms
	Colors []string // contour or histogram: colors

	// contours
	Nlevels  int       // contour: number of levels (overridden by Levels when it's not nil)
	Levels   []float64 // contour: levels (may be nil)
	CmapIdx  int       // contour: colormap index
	NumFmt   string    // contour: number format; e.g. "%g" or "%.2f"
	NoLines  bool      // contour: do not add lines on top of filled contour
	NoLabels bool      // contour: do not add labels
	NoInline bool      // contour: do not draw labels 'inline'
	NoCbar   bool      // contour: do not add colorbar
	CbarLbl  string    // contour: colorbar label
	SelectV  float64   // contour: selected value
	SelectC  string    // contour: color to mark selected level. empty means no selected line
	SelectLw float64   // contour: zero level linewidth

	// histograms
	Type    string // histogram: type; e.g. "bar"
	Stacked bool   // histogram: stacked
	NoFill  bool   // histogram: do not fill bars
	Nbins   int    // histogram: number of bins
	Normed  bool   // histogram: normed

	// figures
	Dpi     int     // figure: dpi to be used when saving figure. default = 96
	Png     bool    // figure: save png file
	Eps     bool    // figure: save eps file
	Prop    float64 // figure: proportion: height = width * prop
	WidthPt float64 // figure: width in points. Get this from LaTeX using \showthe\columnwidth
}

// String returns a string representation of arguments
func (o A) String(forHistogram bool) (l string) {

	// plot and basic options
	addToCmd(&l, o.C != "", io.Sf("color='%s'", o.C))
	addToCmd(&l, o.M != "", io.Sf("marker='%s'", o.M))
	addToCmd(&l, o.Ls != "", io.Sf("ls='%s'", o.Ls))
	addToCmd(&l, o.Lw > 0, io.Sf("lw=%g", o.Lw))
	addToCmd(&l, o.Ms > 0, io.Sf("ms=%d", o.Ms))
	addToCmd(&l, o.L != "", io.Sf("label='%s'", o.L))
	addToCmd(&l, o.Me > 0, io.Sf("markevery=%d", o.Me))
	addToCmd(&l, o.Z > 0, io.Sf("zorder=%d", o.Z))
	addToCmd(&l, o.Mec != "", io.Sf("markeredgecolor='%s'", o.Mec))
	addToCmd(&l, o.Mew > 0, io.Sf("mew=%g", o.Mew))
	addToCmd(&l, o.Void, "markerfacecolor='none'")
	addToCmd(&l, o.Void && o.Mec == "", io.Sf("markeredgecolor='%s'", o.C))
	addToCmd(&l, o.NoClip, "clip_on=0")

	// shapes
	addToCmd(&l, o.Fc != "", io.Sf("facecolor='%s'", o.Fc))
	addToCmd(&l, o.Ec != "", io.Sf("edgecolor='%s'", o.Ec))

	// text and extra arguments
	addToCmd(&l, o.Ha != "", io.Sf("ha='%s'", o.Ha))
	addToCmd(&l, o.Va != "", io.Sf("va='%s'", o.Va))
	addToCmd(&l, o.Fsz > 0, io.Sf("fontsize=%g", o.Fsz))

	// other options
	addToCmd(&l, o.FigFraction, "xycoords='figure fraction'")

	// histograms
	if forHistogram {
		addToCmd(&l, len(o.Colors) > 0, io.Sf("color=%s", strings2list(o.Colors)))
		addToCmd(&l, len(o.Type) > 0, io.Sf("histtype='%s'", o.Type))
		addToCmd(&l, o.Stacked, "stacked=1")
		addToCmd(&l, o.NoFill, "fill=0")
		addToCmd(&l, o.Nbins > 0, io.Sf("bins=%d", o.Nbins))
		addToCmd(&l, o.Normed, "normed=1")
	}
	return
}

// addToCmd adds new option to list of commands separated with commas
func addToCmd(line *string, condition bool, delta string) {
	if condition {
		if len(*line) > 0 {
			*line += ","
		}
		*line += delta
	}
}

// updateBufferWithArgsAndClose updates buffer with arguments and close with ")\n". See updateBufferWithArgs too.
func updateBufferAndClose(buf *bytes.Buffer, args *A, forHistogram bool) {
	if buf == nil {
		return
	}
	if args == nil {
		io.Ff(buf, ")\n")
		return
	}
	txt := args.String(forHistogram)
	if txt == "" {
		io.Ff(buf, ")\n")
		return
	}
	io.Ff(buf, ", "+txt+")\n")
}

// floats2list converts slice of floats to string representing a Python list
func floats2list(vals []float64) (l string) {
	l = "["
	for i, v := range vals {
		if i > 0 {
			l += ","
		}
		l += io.Sf("%g", v)
	}
	l += "]"
	return
}

// strings2list converts slice of strings to string representing a Python list
func strings2list(vals []string) (l string) {
	l = "["
	for i, v := range vals {
		if i > 0 {
			l += ","
		}
		l += io.Sf("'%s'", v)
	}
	l += "]"
	return
}

// getHideList returns a string representing the "spines-to-remove" list in Python
func getHideList(args *A) (l string) {
	if args == nil {
		return
	}
	if args.HideL || args.HideR || args.HideB || args.HideT {
		c := ""
		addToCmd(&c, args.HideL, "'left'")
		addToCmd(&c, args.HideR, "'right'")
		addToCmd(&c, args.HideB, "'bottom'")
		addToCmd(&c, args.HideT, "'top'")
		l = "[" + c + "]"
	}
	return
}

// argsLeg returns legend arguments
func argsLeg(args *A) (loc string, ncol int, hlen, fsz float64, frame int, out int, outX string) {
	loc = "'best'"
	ncol = 1
	hlen = 3.0
	fsz = 8.0
	frame = 0
	out = 0
	outX = "[0.0, 1.02, 1.0, 0.102]"
	if args == nil {
		return
	}
	if args.LegLoc != "" {
		loc = io.Sf("'%s'", args.LegLoc)
	}
	if args.LegNcol > 0 {
		ncol = args.LegNcol
	}
	if args.LegHlen > 0 {
		hlen = args.LegHlen
	}
	if args.FszLeg > 0 {
		fsz = args.FszLeg
	}
	if args.LegFrame {
		frame = 1
	}
	if args.LegOut {
		out = 1
	}
	if len(args.LegOutX) == 4 {
		outX = io.Sf("[%g, %g, %g, %g]", args.LegOutX[0], args.LegOutX[1], args.LegOutX[2], args.LegOutX[3])
	}
	return
}

// argsFsz allocates args if nil, and sets default fontsizes
func argsFsz(args *A) (txt, lbl, leg, xtck, ytck float64) {
	txt, lbl, leg, xtck, ytck = 11, 10, 9, 8, 8
	if args == nil {
		return
	}
	if args.Fsz > 0 {
		txt = args.Fsz
	}
	if args.FszLbl > 0 {
		lbl = args.FszLbl
	}
	if args.FszLeg > 0 {
		leg = args.FszLeg
	}
	if args.FszXtck > 0 {
		xtck = args.FszXtck
	}
	if args.FszYtck > 0 {
		ytck = args.FszYtck
	}
	return
}

// argsFigsize returns figure size data. Defaults are selected if args == nil
func argsFigData(args *A) (figType string, dpi, width, height int) {
	figType, dpi = "png", 150
	prop := 0.75
	widthPt := 400.0
	if args != nil {
		if args.Dpi > 0 {
			dpi = args.Dpi
		}
		if args.Eps {
			figType = "eps"
		}
		if args.Prop > 0 {
			prop = args.Prop
		}
		if args.WidthPt > 0 {
			widthPt = args.WidthPt
		}
	}
	w := widthPt / 72.27 // width in inches
	h := w * prop
	width, height = int(w), int(h)
	return
}

// argsContour allocates args if nil, sets default parameters, and return formatted arguments
func argsContour(in *A, Z [][]float64) (out *A, colors, levels string) {
	out = in
	if out == nil {
		out = new(A)
	}
	if out.NumFmt == "" {
		out.NumFmt = "%g"
	}
	if out.SelectLw < 0.01 {
		out.SelectLw = 3.0
	}
	if out.Lw < 0.01 {
		out.Lw = 1.0
	}
	if out.Fsz < 0.01 {
		out.Fsz = 8.0
	}
	if len(out.Colors) > 0 {
		colors = io.Sf(",colors=%s", strings2list(out.Colors))
	} else {
		colors = io.Sf(",cmap=getCmap(%d)", out.CmapIdx)
	}
	if len(out.Levels) > 0 {
		levels = io.Sf(",levels=%s", floats2list(out.Levels))
	} else {
		if out.Nlevels > 1 {
			levels = ",levels=" + getContourLevels(out.Nlevels, Z)
		}
	}
	return
}

// getContourLevels computes the list of levels based on min and max values in Z
//  Note: the search for min and max is not very efficient for very large matrix
func getContourLevels(nlevels int, Z [][]float64) (l string) {
	if nlevels < 2 {
		return
	}
	if len(Z) < 1 {
		return
	}
	if len(Z[0]) < 1 {
		return
	}
	minZ, maxZ := Z[0][0], Z[0][0]
	for i := 0; i < len(Z); i++ {
		for j := 0; j < len(Z[i]); j++ {
			if Z[i][j] < minZ {
				minZ = Z[i][j]
			}
			if Z[i][j] > maxZ {
				maxZ = Z[i][j]
			}
		}
	}
	delZ := (maxZ - minZ) / float64(nlevels-1)
	l = "["
	for i := 0; i < nlevels; i++ {
		if i > 0 {
			l += ","
		}
		l += io.Sf("%g", minZ+float64(i)*delZ)
	}
	l += "]"
	return
}

// pyBool converts Go bool to Python bool
func pyBool(flag bool) int {
	if flag {
		return 1
	}
	return 0
}
