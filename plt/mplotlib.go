// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// plt contains functions for plotting, drawing in 2D or 3D, and generationg PNG and EPS files
package plt

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

// default directory and temporary file name for python commands
var TEMPORARY = "/tmp/pltgosl.py"

// buffer holding Python commands
var bufferPy bytes.Buffer

// genUid returns an unique id for python variables
func genUid() int { return bufferPy.Len() }

// buffer holding Python extra artists commands
var bufferEa bytes.Buffer

// flag indicating that Python axes3d ('AX3D') has been created
var axes3dCreated bool

// init resets the buffers => ready to go
func init() {
	Reset(false, nil)
}

// fileExt holds the file extension, in case the user called Reset. Otherwise, a default is selected.
var fileExt string

// Reset resets drawing buffer (i.e. Python temporary file data) and sets figure data.
//
//   NOTE: This function is optional; i.e. plt works without calling this function.
//         Nonetheless, if fontsizes or figure sizes need to be specified, Reset can be called.
//
//   Input:
//     setDefault -- sets default values
//     args -- optional data (may be nil)
//
//   NOTE: Default values are selected if setDefault == true.
//         Otherwise, Python (matplotlib) will choose defaults.
//         Also, if args != nil, some values are set based on data in args.
//
//   The following data is set:
//     fontsizes:
//        args.Fsz     float64 // font size
//        args.FszLbl  float64 // font size of labels
//        args.FszLeg  float64 // font size of legend
//        args.FszXtck float64 // font size of x-ticks
//        args.FszYtck float64 // font size of y-ticks
//     figure data:
//        args.Dpi     int     // dpi to be used when saving figure. default = 96
//        args.Png     bool    // save png file
//        args.Eps     bool    // save eps file
//        args.Prop    float64 // proportion: height = width * prop
//        args.WidthPt float64 // width in points. Get this from LaTeX using \showthe\columnwidth
func Reset(setDefault bool, args *A) {

	// clear buffer and start python code
	bufferPy.Reset()
	bufferEa.Reset()
	io.Ff(&bufferPy, pythonHeader)
	axes3dCreated = false

	// set figure data
	if setDefault {
		txt, lbl, leg, xtck, ytck := argsFsz(args)
		figType, dpi, width, height := argsFigData(args)
		io.Ff(&bufferPy, "plt.rcdefaults()\n")
		io.Ff(&bufferPy, "plt.rcParams.update({\n")
		io.Ff(&bufferPy, "    'font.size'       : %g,\n", txt)
		io.Ff(&bufferPy, "    'axes.labelsize'  : %g,\n", lbl)
		io.Ff(&bufferPy, "    'legend.fontsize' : %g,\n", leg)
		io.Ff(&bufferPy, "    'xtick.labelsize' : %g,\n", xtck)
		io.Ff(&bufferPy, "    'ytick.labelsize' : %g,\n", ytck)
		io.Ff(&bufferPy, "    'figure.figsize'  : [%d,%d],\n", width, height)
		switch figType {
		case "eps":
			io.Ff(&bufferPy, "    'backend'            : 'ps',\n")
			io.Ff(&bufferPy, "    'text.usetex'        : True,\n")  // very IMPORTANT to avoid Type 3 fonts
			io.Ff(&bufferPy, "    'ps.useafm'          : True,\n")  // very IMPORTANT to avoid Type 3 fonts
			io.Ff(&bufferPy, "    'pdf.use14corefonts' : True})\n") // very IMPORTANT to avoid Type 3 fonts
			fileExt = ".eps"
		default:
			io.Ff(&bufferPy, "    'savefig.dpi'     : %d})\n", dpi)
			fileExt = ".png"
		}
	}
}

// PyCmds adds Python commands to be called when plotting
func PyCmds(text string) {
	io.Ff(&bufferPy, text)
}

// PyFile loads Python file and copy its contents to temporary buffer
func PyFile(filename string) (err error) {
	b, err := io.ReadFile(filename)
	if err != nil {
		return
	}
	io.Ff(&bufferPy, string(b))
	return
}

// DoubleYscale duplicates y-scale
func DoubleYscale(ylabelOrEmpty string) {
	io.Ff(&bufferPy, "plt.gca().twinx()\n")
	if ylabelOrEmpty != "" {
		io.Ff(&bufferPy, "plt.gca().set_ylabel('%s')\n", ylabelOrEmpty)
	}
}

// SetXlog sets x-scale to be log
func SetXlog() {
	io.Ff(&bufferPy, "plt.gca().set_xscale('log')\n")
}

// SetYlog sets y-scale to be log
func SetYlog() {
	io.Ff(&bufferPy, "plt.gca().set_yscale('log')\n")
}

// SetTicksXlist sets x-axis ticks with given list
func SetTicksXlist(values []float64) {
	io.Ff(&bufferPy, "plt.xticks(%v)\n", floats2list(values))
}

// SetTicksYlist sets y-ayis ticks with given list
func SetTicksYlist(values []float64) {
	io.Ff(&bufferPy, "plt.yticks(%v)\n", floats2list(values))
}

// SetXnticks sets number of ticks along x
func SetXnticks(num int) {
	if num == 0 {
		io.Ff(&bufferPy, "plt.gca().get_xaxis().set_ticks([])\n")
	} else {
		io.Ff(&bufferPy, "plt.gca().get_xaxis().set_major_locator(tck.MaxNLocator(%d))\n", num)
	}
}

// SetYnticks sets number of ticks along y
func SetYnticks(num int) {
	if num == 0 {
		io.Ff(&bufferPy, "plt.gca().get_yaxis().set_ticks([])\n")
	} else {
		io.Ff(&bufferPy, "plt.gca().get_yaxis().set_major_locator(tck.MaxNLocator(%d))\n", num)
	}
}

// SetTicksX sets ticks along x
func SetTicksX(majorEvery, minorEvery float64, majorFmt string) {
	uid := genUid()
	if majorEvery > 0 {
		io.Ff(&bufferPy, "majorLocator%d = tck.MultipleLocator(%g)\n", uid, majorEvery)
		io.Ff(&bufferPy, "nticks%d = (plt.gca().axis()[1] - plt.gca().axis()[0]) / %g\n", uid, majorEvery)
		io.Ff(&bufferPy, "if nticks%d < majorLocator%d.MAXTICKS * 0.9:\n", uid, uid)
		io.Ff(&bufferPy, "    plt.gca().xaxis.set_major_locator(majorLocator%d)\n", uid)
	}
	if minorEvery > 0 {
		io.Ff(&bufferPy, "minorLocator%d = tck.MultipleLocator(%g)\n", uid, minorEvery)
		io.Ff(&bufferPy, "nticks%d = (plt.gca().axis()[1] - plt.gca().axis()[0]) / %g\n", uid, minorEvery)
		io.Ff(&bufferPy, "if nticks%d < minorLocator%d.MAXTICKS * 0.9:\n", uid, uid)
		io.Ff(&bufferPy, "    plt.gca().xaxis.set_minor_locator(minorLocator%d)\n", uid)
	}
	if majorFmt != "" {
		io.Ff(&bufferPy, "majorFormatter%d = tck.FormatStrFormatter(r'%s')\n", uid, majorFmt)
		io.Ff(&bufferPy, "plt.gca().xaxis.set_major_formatter(majorFormatter%d)\n", uid)
	}
}

// SetTicksY sets ticks along y
func SetTicksY(majorEvery, minorEvery float64, majorFmt string) {
	uid := genUid()
	if majorEvery > 0 {
		io.Ff(&bufferPy, "majorLocator%d = tck.MultipleLocator(%g)\n", uid, majorEvery)
		io.Ff(&bufferPy, "nticks%d = (plt.gca().axis()[1] - plt.gca().axis()[0]) / %g\n", uid, majorEvery)
		io.Ff(&bufferPy, "if nticks%d < majorLocator%d.MAXTICKS * 0.9:\n", uid, uid)
		io.Ff(&bufferPy, "    plt.gca().yaxis.set_major_locator(majorLocator%d)\n", uid)
	}
	if minorEvery > 0 {
		io.Ff(&bufferPy, "minorLocator%d = tck.MultipleLocator(%g)\n", uid, minorEvery)
		io.Ff(&bufferPy, "nticks%d = (plt.gca().axis()[1] - plt.gca().axis()[0]) / %g\n", uid, minorEvery)
		io.Ff(&bufferPy, "if nticks%d < minorLocator%d.MAXTICKS * 0.9:\n", uid, uid)
		io.Ff(&bufferPy, "    plt.gca().yaxis.set_minor_locator(minorLocator%d)\n", uid)
	}
	if majorFmt != "" {
		io.Ff(&bufferPy, "majorFormatter%d = tck.FormatStrFormatter(r'%s')\n", uid, majorFmt)
		io.Ff(&bufferPy, "plt.gca().yaxis.set_major_formatter(majorFormatter%d)\n", uid)
	}
}

// SetScientificX sets scientific notation for ticks along x-axis
func SetScientificX(minOrder, maxOrder int) {
	uid := genUid()
	io.Ff(&bufferPy, "fmt%d = plt.ScalarFormatter(useOffset=True)\n", uid)
	io.Ff(&bufferPy, "fmt%d.set_powerlimits((%d,%d))\n", uid, minOrder, maxOrder)
	io.Ff(&bufferPy, "plt.gca().xaxis.set_major_formatter(fmt%d)\n", uid)
}

// SetScientificY sets scientific notation for ticks along y-axis
func SetScientificY(minOrder, maxOrder int) {
	uid := genUid()
	io.Ff(&bufferPy, "fmt%d = plt.ScalarFormatter(useOffset=True)\n", uid)
	io.Ff(&bufferPy, "fmt%d.set_powerlimits((%d,%d))\n", uid, minOrder, maxOrder)
	io.Ff(&bufferPy, "plt.gca().yaxis.set_major_formatter(fmt%d)\n", uid)
}

// SetTicksNormal sets normal ticks
func SetTicksNormal() {
	io.Ff(&bufferPy, "plt.gca().ticklabel_format(useOffset=False)\n")
}

// ReplaceAxes substitutes axis frame (see Axes in gosl.py)
//   ex: xDel, yDel := 0.04, 0.04
func ReplaceAxes(xi, yi, xf, yf, xDel, yDel float64, xLab, yLab string, argsArrow, argsText *A) {
	io.Ff(&bufferPy, "plt.axis('off')\n")
	Arrow(xi, yi, xf, yi, argsArrow)
	Arrow(xi, yi, xi, yf, argsArrow)
	Text(xf, yi-xDel, xLab, argsText)
	Text(xi-yDel, yf, yLab, argsText)
}

// AxHline adds horizontal line to axis
func AxHline(y float64, args *A) {
	io.Ff(&bufferPy, "plt.axhline(%g", y)
	updateBufferAndClose(&bufferPy, args, false, false)
}

// AxVline adds vertical line to axis
func AxVline(x float64, args *A) {
	io.Ff(&bufferPy, "plt.axvline(%g", x)
	updateBufferAndClose(&bufferPy, args, false, false)
}

// HideBorders hides frame borders
func HideBorders(args *A) {
	hide := getHideList(args)
	if hide != "" {
		io.Ff(&bufferPy, "for spine in %s: plt.gca().spines[spine].set_visible(0)\n", hide)
	}
}

// HideAllBorders hides all frame borders
func HideAllBorders() {
	io.Ff(&bufferPy, "for spine in ['left','right','bottom','top']: plt.gca().spines[spine].set_visible(0)\n")
}

// Annotate adds annotation to plot
func Annotate(x, y float64, txt string, args *A) {
	io.Ff(&bufferPy, "plt.annotate(r'%s', xy=(%g,%g)", txt, x, y)
	updateBufferAndClose(&bufferPy, args, false, false)
}

// AnnotateXlabels sets text of xlabels
func AnnotateXlabels(x float64, txt string, args *A) {
	fsz := 7.0
	if args != nil {
		if args.Fsz > 0 {
			fsz = args.Fsz
		}
	}
	io.Ff(&bufferPy, "plt.annotate('%s', xy=(%g, -%g-3), xycoords=('data', 'axes points'), va='top', ha='center', size=%g", txt, x, fsz, fsz)
	updateBufferAndClose(&bufferPy, args, false, false)
}

// SupTitle sets subplot title
func SupTitle(txt string, args *A) {
	uid := genUid()
	io.Ff(&bufferPy, "st%d = plt.suptitle(r'%s'", uid, txt)
	updateBufferAndClose(&bufferPy, args, false, false)
	io.Ff(&bufferPy, "addToEA(st%d)\n", uid)
}

// Title sets title
func Title(txt string, args *A) {
	io.Ff(&bufferPy, "plt.title(r'%s'", txt)
	updateBufferAndClose(&bufferPy, args, false, false)
}

// Text adds text to plot
func Text(x, y float64, txt string, args *A) {
	io.Ff(&bufferPy, "plt.text(%g,%g,r'%s'", x, y, txt)
	updateBufferAndClose(&bufferPy, args, false, false)
}

// Cross adds a vertical and horizontal lines @ (x0,y0) to plot (i.e. large cross)
func Cross(x0, y0 float64, args *A) {
	cl, ls, lw, z := "black", "dashed", 1.2, 0
	if args != nil {
		if args.C != "" {
			cl = args.C
		}
		if args.Lw > 0 {
			lw = args.Lw
		}
		if args.Ls != "" {
			ls = args.Ls
		}
		if args.Z > 0 {
			z = args.Z
		}
	}
	io.Ff(&bufferPy, "plt.axvline(%g, color='%s', linestyle='%s', linewidth=%g, zorder=%d)\n", x0, cl, ls, lw, z)
	io.Ff(&bufferPy, "plt.axhline(%g, color='%s', linestyle='%s', linewidth=%g, zorder=%d)\n", y0, cl, ls, lw, z)
}

// SplotGap sets gap between subplots
func SplotGap(w, h float64) {
	io.Ff(&bufferPy, "plt.subplots_adjust(wspace=%g, hspace=%g)\n", w, h)
}

// Subplot adds/sets a subplot
func Subplot(i, j, k int) {
	io.Ff(&bufferPy, "plt.subplot(%d,%d,%d)\n", i, j, k)
}

// Subplot adds/sets a subplot with given indices in I
func SubplotI(I []int) {
	if len(I) != 3 {
		return
	}
	io.Ff(&bufferPy, "plt.subplot(%d,%d,%d)\n", I[0], I[1], I[2])
}

// SetHspace sets horizontal space between subplots
func SetHspace(hspace float64) {
	io.Ff(&bufferPy, "plt.subplots_adjust(hspace=%g)\n", hspace)
}

// SetVspace sets vertical space between subplots
func SetVspace(vspace float64) {
	io.Ff(&bufferPy, "plt.subplots_adjust(vspace=%g)\n", vspace)
}

// Equal sets same scale for both axes
func Equal() {
	io.Ff(&bufferPy, "plt.axis('equal')\n")
}

// AxisOff hides axes
func AxisOff() {
	io.Ff(&bufferPy, "plt.axis('off')\n")
}

// SetAxis sets axes limits
func SetAxis(xmin, xmax, ymin, ymax float64) {
	io.Ff(&bufferPy, "plt.axis([%g, %g, %g, %g])\n", xmin, xmax, ymin, ymax)
}

// AxisXmin sets minimum x
func AxisXmin(xmin float64) {
	io.Ff(&bufferPy, "plt.axis([%g, plt.axis()[1], plt.axis()[2], plt.axis()[3]])\n", xmin)
}

// AxisXmax sets maximum x
func AxisXmax(xmax float64) {
	io.Ff(&bufferPy, "plt.axis([plt.axis()[0], %g, plt.axis()[2], plt.axis()[3]])\n", xmax)
}

// AxisYmin sets minimum y
func AxisYmin(ymin float64) {
	io.Ff(&bufferPy, "plt.axis([plt.axis()[0], plt.axis()[1], %g, plt.axis()[3]])\n", ymin)
}

// AxisYmax sets maximum y
func AxisYmax(ymax float64) {
	io.Ff(&bufferPy, "plt.axis([plt.axis()[0], plt.axis()[1], plt.axis()[2], %g])\n", ymax)
}

// AxisXrange sets x-range (i.e. limits)
func AxisXrange(xmin, xmax float64) {
	io.Ff(&bufferPy, "plt.axis([%g, %g, plt.axis()[2], plt.axis()[3]])\n", xmin, xmax)
}

// AxisYrange sets y-range (i.e. limits)
func AxisYrange(ymin, ymax float64) {
	io.Ff(&bufferPy, "plt.axis([plt.axis()[0], plt.axis()[1], %g, %g])\n", ymin, ymax)
}

// AxisRange sets x and y ranges (i.e. limits)
func AxisRange(xmin, xmax, ymin, ymax float64) {
	io.Ff(&bufferPy, "plt.axis([%g, %g, %g, %g])\n", xmin, xmax, ymin, ymax)
}

// AxisLims sets x and y limits
func AxisLims(lims []float64) {
	io.Ff(&bufferPy, "plt.axis([%g, %g, %g, %g])\n", lims[0], lims[1], lims[2], lims[3])
}

// Plot plots x-y series
func Plot(x, y []float64, args *A) (sx, sy string) {
	uid := genUid()
	sx = io.Sf("x%d", uid)
	sy = io.Sf("y%d", uid)
	gen2Arrays(&bufferPy, sx, sy, x, y)
	io.Ff(&bufferPy, "plt.plot(%s,%s", sx, sy)
	updateBufferAndClose(&bufferPy, args, false, false)
	return
}

// PlotOne plots one point @ (x,y)
func PlotOne(x, y float64, args *A) {
	io.Ff(&bufferPy, "plt.plot(%23.15e,%23.15e", x, y)
	updateBufferAndClose(&bufferPy, args, false, false)
}

// Hist draws histogram
func Hist(x [][]float64, labels []string, args *A) {
	uid := genUid()
	sx := io.Sf("x%d", uid)
	sy := io.Sf("y%d", uid)
	genList(&bufferPy, sx, x)
	genStrArray(&bufferPy, sy, labels)
	io.Ff(&bufferPy, "plt.hist(%s,label=r'%s'", sx, sy)
	updateBufferAndClose(&bufferPy, args, true, false)
}

// Grid2d draws grid lines (and points) of 2D grid
func Grid2d(X, Y [][]float64, withPoints bool, argsLines, argsPoints *A) {
	if argsLines == nil {
		argsLines = &A{C: "k", NoClip: true}
	}
	if argsPoints == nil {
		argsPoints = &A{C: "k", M: ".", Ls: "none", NoClip: true}
	}
	for i := 0; i < len(X); i++ {
		Plot(X[i], Y[i], argsLines)
	}
	nrows := len(X)
	if nrows < 1 {
		return
	}
	ncols := len(X[0])
	x, y := make([]float64, nrows), make([]float64, nrows)
	for j := 0; j < ncols; j++ {
		for i := 0; i < nrows; i++ {
			x[i], y[i] = X[i][j], Y[i][j]
			if withPoints {
				PlotOne(X[i][j], Y[i][j], argsPoints)
			}
		}
		Plot(x, y, argsLines)
	}
}

// ContourF draws filled contour and possibly with a contour of lines (if args.UnoLines=false)
func ContourF(x, y, z [][]float64, args *A) {
	uid := genUid()
	sx := io.Sf("x%d", uid)
	sy := io.Sf("y%d", uid)
	sz := io.Sf("z%d", uid)
	genMat(&bufferPy, sx, x)
	genMat(&bufferPy, sy, y)
	genMat(&bufferPy, sz, z)
	a, colors, levels := argsContour(args, z)
	io.Ff(&bufferPy, "c%d = plt.contourf(%s,%s,%s%s%s)\n", uid, sx, sy, sz, colors, levels)
	if !a.NoLines {
		io.Ff(&bufferPy, "cc%d = plt.contour(%s,%s,%s,colors=['k']%s,linewidths=[%g])\n", uid, sx, sy, sz, levels, a.Lw)
		if !a.NoLabels {
			io.Ff(&bufferPy, "plt.clabel(cc%d,inline=%d,fontsize=%g)\n", uid, pyBool(!a.NoInline), a.Fsz)
		}
	}
	if !a.NoCbar {
		io.Ff(&bufferPy, "cb%d = plt.colorbar(c%d, format='%s')\n", uid, uid, a.NumFmt)
		if a.CbarLbl != "" {
			io.Ff(&bufferPy, "cb%d.ax.set_ylabel(r'%s')\n", uid, a.CbarLbl)
		}
	}
	if a.SelectC != "" {
		io.Ff(&bufferPy, "ccc%d = plt.contour(%s,%s,%s,colors=['%s'],levels=[%g],linewidths=[%g],linestyles=['-'])\n", uid, sx, sy, sz, a.SelectC, a.SelectV, a.SelectLw)
	}
}

// ContourL draws a contour with lines only
func ContourL(x, y, z [][]float64, args *A) {
	uid := genUid()
	sx := io.Sf("x%d", uid)
	sy := io.Sf("y%d", uid)
	sz := io.Sf("z%d", uid)
	genMat(&bufferPy, sx, x)
	genMat(&bufferPy, sy, y)
	genMat(&bufferPy, sz, z)
	a, colors, levels := argsContour(args, z)
	io.Ff(&bufferPy, "c%d = plt.contour(%s,%s,%s%s%s)\n", uid, sx, sy, sz, colors, levels)
	if !a.NoLabels {
		io.Ff(&bufferPy, "plt.clabel(c%d,inline=%d,fontsize=%g)\n", uid, pyBool(!a.NoInline), a.Fsz)
	}
	if a.SelectC != "" {
		io.Ff(&bufferPy, "cc%d = plt.contour(%s,%s,%s,colors=['%s'],levels=[%g],linewidths=[%g],linestyles=['-'])\n", uid, sx, sy, sz, a.SelectC, a.SelectV, a.SelectLw)
	}
}

// Quiver draws vector field
func Quiver(x, y, gx, gy [][]float64, args *A) {
	uid := genUid()
	sx := io.Sf("x%d", uid)
	sy := io.Sf("y%d", uid)
	sgx := io.Sf("gx%d", uid)
	sgy := io.Sf("gy%d", uid)
	genMat(&bufferPy, sx, x)
	genMat(&bufferPy, sy, y)
	genMat(&bufferPy, sgx, gx)
	genMat(&bufferPy, sgy, gy)
	io.Ff(&bufferPy, "plt.quiver(%s,%s,%s,%s", sx, sy, sgx, sgy)
	updateBufferAndClose(&bufferPy, args, false, false)
}

// Grid adds grid to plot
func Grid(args *A) {
	io.Ff(&bufferPy, "plt.grid(")
	updateBufferFirstArgsAndClose(&bufferPy, args, false, false)
}

// Legend adds legend to plot
func Legend(args *A) {
	loc, ncol, hlen, fsz, frame, out, outX := argsLeg(args)
	uid := genUid()
	io.Ff(&bufferPy, "h%d, l%d = plt.gca().get_legend_handles_labels()\n", uid, uid)
	io.Ff(&bufferPy, "if len(h%d) > 0 and len(l%d) > 0:\n", uid, uid)
	if out == 1 {
		io.Ff(&bufferPy, "    d%d = %s\n", uid, outX)
		io.Ff(&bufferPy, "    l%d = plt.legend(bbox_to_anchor=d%d, ncol=%d, handlelength=%g, prop={'size':%g}, loc=3, mode='expand', borderaxespad=0.0, columnspacing=1, handletextpad=0.05)\n", uid, uid, ncol, hlen, fsz)
		io.Ff(&bufferPy, "    addToEA(l%d)\n", uid)
	} else {
		io.Ff(&bufferPy, "    l%d = plt.legend(loc=%s, ncol=%d, handlelength=%g, prop={'size':%g})\n", uid, loc, ncol, hlen, fsz)
		io.Ff(&bufferPy, "    addToEA(l%d)\n", uid)
	}
	if frame == 0 {
		io.Ff(&bufferPy, "    l%d.get_frame().set_linewidth(0.0)\n", uid)
	}
}

// LegendX adds legend to plot with given data instead of relying on labels
func LegendX(dat []*A, args *A) {
	loc, ncol, hlen, fsz, frame, out, outX := argsLeg(args)
	uid := genUid()
	io.Ff(&bufferPy, "h%d = [", uid)
	for i, d := range dat {
		if i > 0 {
			io.Ff(&bufferPy, ",\n")
		}
		if d != nil {
			io.Ff(&bufferPy, "lns.Line2D([], [], %s)", d.String(false, false))
		}
	}
	io.Ff(&bufferPy, "]\n")
	io.Ff(&bufferPy, "if len(h%d) > 0:\n", uid)
	if out == 1 {
		io.Ff(&bufferPy, "    d%d = %s\n", uid, outX)
		io.Ff(&bufferPy, "    l%d = plt.legend(handles=h%d, bbox_to_anchor=d%d, ncol=%d, handlelength=%g, prop={'size':%g}, loc=3, mode='expand', borderaxespad=0.0, columnspacing=1, handletextpad=0.05)\n", uid, uid, uid, ncol, hlen, fsz)
		io.Ff(&bufferPy, "    addToEA(l%d)\n", uid)
	} else {
		io.Ff(&bufferPy, "    l%d = plt.legend(handles=h%d, loc=%s, ncol=%d, handlelength=%g, prop={'size':%g})\n", uid, uid, loc, ncol, hlen, fsz)
		io.Ff(&bufferPy, "    addToEA(l%d)\n", uid)
	}
	if frame == 0 {
		io.Ff(&bufferPy, "    l%d.get_frame().set_linewidth(0.0)\n", uid)
	}
}

// Gll adds grid, labels, and legend to plot
func Gll(xl, yl string, args *A) {
	hide := getHideList(args)
	if hide != "" {
		io.Ff(&bufferPy, "for spine in %s: plt.gca().spines[spine].set_visible(0)\n", hide)
	}
	io.Ff(&bufferPy, "plt.grid(color='grey', zorder=-1000)\n")
	io.Ff(&bufferPy, "plt.xlabel(r'%s')\n", xl)
	io.Ff(&bufferPy, "plt.ylabel(r'%s')\n", yl)
	Legend(args)
}

// SetLabels sets x-y axes labels
func SetLabels(x, y string, args *A) {
	a := ""
	if args != nil {
		a = "," + args.String(false, false)
	}
	io.Ff(&bufferPy, "plt.xlabel(r'%s'%s);plt.ylabel(r'%s'%s)\n", x, a, y, a)
}

// Clf clears current figure
func Clf() {
	io.Ff(&bufferPy, "plt.clf()\n")
}

// SetFontSizes sets font sizes
func SetFontSizes(args *A) {
	txt, lbl, leg, xtck, ytck := argsFsz(args)
	io.Ff(&bufferPy, "plt.rcParams.update({\n")
	io.Ff(&bufferPy, "    'font.size'       : %g,\n", txt)
	io.Ff(&bufferPy, "    'axes.labelsize'  : %g,\n", lbl)
	io.Ff(&bufferPy, "    'legend.fontsize' : %g,\n", leg)
	io.Ff(&bufferPy, "    'xtick.labelsize' : %g,\n", xtck)
	io.Ff(&bufferPy, "    'ytick.labelsize' : %g})\n", ytck)
}

// ZoomWindow adds another axes to plot a figure within the figure; e.g. a zoom window
//  lef, bot, wid, hei -- normalised figure coordinates: left,bottom,width,height
//  asOld -- handle to the previous axes
//  axNew -- handle to the new axes
func ZoomWindow(lef, bot, wid, hei float64, args *A) (axOld, axNew string) {
	uid := genUid()
	clr := "#dcdcdc"
	if args != nil {
		clr = args.C
	}
	axOld = io.Sf("axOld%d", uid)
	io.Ff(&bufferPy, "%s = plt.gca()\n", axOld)
	axNew = io.Sf("axNew%d", uid)
	io.Ff(&bufferPy, "%s = plt.axes([%g,%g,%g,%g], axisbg='%s')\n", axNew, lef, bot, wid, hei, clr)
	return
}

// Sca sets current axes
func Sca(axName string) {
	io.Ff(&bufferPy, "plt.sca(%s)\n", axName)
}

// functions to save figure ///////////////////////////////////////////////////////////////////////

// Save saves figure after creating a directory
//  NOTE: the file name will be fnkey + .png (default) or .eps depending on the Reset function
func Save(dirout, fnkey string) (err error) {
	empty := dirout == "" || fnkey == ""
	if empty {
		return chk.Err("directory and filename key must not be empty\n")
	}
	err = os.MkdirAll(dirout, 0777)
	if err != nil {
		return chk.Err("cannot create directory to save figure file:\n%v\n", err)
	}
	if fileExt == "" {
		fileExt = ".png"
	}
	fn := filepath.Join(dirout, fnkey+fileExt)
	io.Ff(&bufferPy, "plt.savefig(r'%s', bbox_inches='tight', bbox_extra_artists=EXTRA_ARTISTS)\n", fn)
	return run(fn)
}

// Show shows figure
func Show() error {
	io.Ff(&bufferPy, "plt.show()\n")
	return run("")
}

// ShowSave shows figure and/or save figure
func ShowSave(dirout, fnkey string) (err error) {
	empty := dirout == "" || fnkey == ""
	if empty {
		return chk.Err("directory and filename key must not be empty\n")
	}
	uid := genUid()
	io.Ff(&bufferPy, "fig%d = plt.gcf()\n", uid)
	io.Ff(&bufferPy, "plt.show()\n")
	err = os.MkdirAll(dirout, 0777)
	if err != nil {
		return chk.Err("cannot create directory to save figure file:\n%v\n", err)
	}
	if fileExt == "" {
		fileExt = ".png"
	}
	fn := filepath.Join(dirout, fnkey+fileExt)
	io.Ff(&bufferPy, "fig%d.savefig(r'%s', bbox_inches='tight', bbox_extra_artists=EXTRA_ARTISTS)\n", uid, fn)
	return run("")
}

// generate arrays and matrices ///////////////////////////////////////////////////////////////////

// genMat generates matrix
func genMat(buf *bytes.Buffer, name string, a [][]float64) {
	io.Ff(buf, "%s=np.array([", name)
	for i, _ := range a {
		io.Ff(buf, "[")
		for j, _ := range a[i] {
			io.Ff(buf, "%g,", a[i][j])
		}
		io.Ff(buf, "],")
	}
	io.Ff(buf, "],dtype=float)\n")
}

// genList generates list
func genList(buf *bytes.Buffer, name string, a [][]float64) {
	io.Ff(buf, "%s=[", name)
	for i, _ := range a {
		io.Ff(buf, "[")
		for j, _ := range a[i] {
			io.Ff(buf, "%g,", a[i][j])
		}
		io.Ff(buf, "],")
	}
	io.Ff(buf, "]\n")
}

// genArray generates the NumPy text corresponding to an array of float point numbers
func genArray(buf *bytes.Buffer, name string, u []float64) {
	io.Ff(buf, "%s=np.array([", name)
	for i, _ := range u {
		io.Ff(buf, "%g,", u[i])
	}
	io.Ff(buf, "],dtype=float)\n")
}

// gen2Arrays generates the NumPy text corresponding to 2 arrays of float point numbers
func gen2Arrays(buf *bytes.Buffer, nameA, nameB string, a, b []float64) {
	genArray(buf, nameA, a)
	genArray(buf, nameB, b)
}

// genStrArray generates the NumPy text corresponding to an array of strings
func genStrArray(buf *bytes.Buffer, name string, u []string) {
	io.Ff(buf, "%s=[", name)
	for i, _ := range u {
		io.Ff(buf, "r'%s',", u[i])
	}
	io.Ff(buf, "]\n")
}

// call Python ////////////////////////////////////////////////////////////////////////////////////

// run calls Python to generate plot
func run(fn string) (err error) {

	// write file
	io.WriteFile(TEMPORARY, &bufferEa, &bufferPy)

	// set command
	cmd := exec.Command("python", TEMPORARY)
	var out, serr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &serr

	// call Python
	err = cmd.Run()
	if err != nil {
		return chk.Err("call to Python failed:\n%v\n", serr.String())
	}

	// show filename
	if fn != "" {
		io.Pf("file <%s> written\n", fn)
	}

	// show output
	io.Pf("%s", out.String())
	return
}

const pythonHeader = `### file generated by Gosl #################################################
import numpy as np
import matplotlib.pyplot as plt
import matplotlib.ticker as tck
import matplotlib.patches as pat
import matplotlib.path as pth
import matplotlib.patheffects as pff
import matplotlib.lines as lns
import mpl_toolkits.mplot3d as m3d
NaN = np.NaN
EXTRA_ARTISTS = []
def addToEA(obj):
    if obj!=None: EXTRA_ARTISTS.append(obj)
COLORMAPS = [plt.cm.bwr, plt.cm.RdBu, plt.cm.hsv, plt.cm.jet, plt.cm.terrain, plt.cm.pink, plt.cm.Greys]
def getCmap(idx): return COLORMAPS[idx %% len(COLORMAPS)]
`
