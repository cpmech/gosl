// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// plt contains functions for plotting
package plt

import (
	"bytes"
	"os"
	"os/exec"

	"github.com/cpmech/gosl/utl"
)

var bb bytes.Buffer // the buffer
var ea bytes.Buffer // extra artists

func init() {
	Reset()
}

func Reset() {
	bb.Reset()
	ea.Reset()
	utl.Ff(&bb, "from gosl import *\n")
	utl.Ff(&ea, "ea = []\n")
}

func SetXlog() {
	utl.Ff(&bb, "SetXlog()\n")
}

func SetYlog() {
	utl.Ff(&bb, "SetYlog()\n")
}

func SetXnticks(num int) {
	utl.Ff(&bb, "SetXnticks(%d)\n", num)
}

func SetYnticks(num int) {
	utl.Ff(&bb, "SetYnticks(%d)\n", num)
}

func SetScientific(axis string, min_order, max_order int) {
	utl.Ff(&bb, "SetScientificFmt(axis='%s', min_order=%d, max_order=%d)\n", axis, min_order, max_order)
}

func Arrow(xi, yi, xf, yf float64, args string) {
	cmd := utl.Sf("Arrow(%g,%g, %g,%g", xi, yi, xf, yf)
	if len(args) > 0 {
		cmd += utl.Sf(",%s", args)
	}
	utl.Ff(&bb, "%s)\n", cmd)
}

func AxHline(y float64, args string) {
	cmd := utl.Sf("axhline(%g", y)
	if len(args) > 0 {
		cmd += utl.Sf(",%s", args)
	}
	utl.Ff(&bb, "%s)\n", cmd)
}

func AxVline(x float64, args string) {
	cmd := utl.Sf("axvline(%g", x)
	if len(args) > 0 {
		cmd += utl.Sf(",%s", args)
	}
	utl.Ff(&bb, "%s)\n", cmd)
}

func HideTRframe() {
	utl.Ff(&bb, "HideFrameLines()\n")
}

func Annotate(x, y float64, txt string, args string) {
	cmd := utl.Sf("annotate(%s, xy=(%g,%g)", txt, x, y)
	if len(args) > 0 {
		cmd += utl.Sf(",%s", args)
	}
	utl.Ff(&bb, "%s)\n", cmd)
}

func AnnotateXlabels(x float64, txt string, args string) {
	cmd := utl.Sf("AnnotateXlabels(%g, %s", x, txt)
	if len(args) > 0 {
		cmd += utl.Sf(",%s", args)
	}
	utl.Ff(&bb, "%s)\n", cmd)
}

func SupTitle(txt, args string) {
	n := bb.Len()
	if len(args) > 0 {
		utl.Ff(&bb, "st%d = suptitle('%s',%s)\n", n, txt, args)
	} else {
		utl.Ff(&bb, "st%d = suptitle('%s')\n", n, txt)
	}
	utl.Ff(&bb, "ea.append(st%d)\n", n)
}

func Title(txt, args string) {
	if len(args) > 0 {
		utl.Ff(&bb, "title('%s',%s)\n", txt, args)
	} else {
		utl.Ff(&bb, "title('%s')\n", txt)
	}
}

func Text(x, y float64, txt, args string) {
	if len(args) > 0 {
		utl.Ff(&bb, "text(%g,%g,'%s',%s)\n", x, y, txt, args)
	} else {
		utl.Ff(&bb, "text(%g,%g,'%s')\n", x, y, txt)
	}
}

func Cross() {
	utl.Ff(&bb, "Cross()\n")
}

func SplotGap(w, h float64) {
	utl.Ff(&bb, "SplotGap(%g, %g)\n", w, h)
}

func Subplot(i, j, k int) {
	utl.Ff(&bb, "subplot(%d,%d,%d)\n", i, j, k)
}

func SetHspace(hspace float64) {
	utl.Ff(&bb, "subplots_adjust(hspace=%g)\n", hspace)
}

func SetVspace(vspace float64) {
	utl.Ff(&bb, "subplots_adjust(vspace=%g)\n", vspace)
}

func AxisXmin(xmin float64) {
	utl.Ff(&bb, "axis([%g, axis()[1], axis()[2], axis()[3]])\n", xmin)
}

func AxisXmax(xmax float64) {
	utl.Ff(&bb, "axis([axis()[0], %g, axis()[2], axis()[3]])\n", xmax)
}

func AxisYmin(ymin float64) {
	utl.Ff(&bb, "axis([axis()[0], axis()[1], %g, axis()[3]])\n", ymin)
}

func AxisYmax(ymax float64) {
	utl.Ff(&bb, "axis([axis()[0], axis()[1], axis()[2], %g])\n", ymax)
}

func AxisXrange(xmin, xmax float64) {
	utl.Ff(&bb, "axis([%g, %g, axis()[2], axis()[3]])\n", xmin, xmax)
}

func AyisYrange(ymin, ymax float64) {
	utl.Ff(&bb, "axis([axis()[0], axis()[1], %g, %g])\n", ymin, ymax)
}

func AxisRange(xmin, xmax, ymin, ymax float64) {
	utl.Ff(&bb, "axis([%g, %g, %g, %g])\n", xmin, xmax, ymin, ymax)
}

func AxisLims(lims []float64) {
	utl.Ff(&bb, "axis([%g, %g, %g, %g])\n", lims[0], lims[1], lims[2], lims[3])
}

func Plot(x, y []float64, args string) {
	n := bb.Len()
	sx := utl.Sf("x%d", n)
	sy := utl.Sf("y%d", n)
	utl.Gen2Arrays(&bb, sx, sy, x, y)
	if len(args) > 0 {
		utl.Ff(&bb, "plot(%s,%s,%s)\n", sx, sy, args)
	} else {
		utl.Ff(&bb, "plot(%s,%s)\n", sx, sy)
	}
}

func PlotOne(x, y float64, args string) {
	if len(args) > 0 {
		utl.Ff(&bb, "plot(%23.15e,%23.15e,%s)\n", x, y, args)
	} else {
		utl.Ff(&bb, "plot(%23.15e,%23.15e)\n", x, y)
	}
}

func Contour(x, y, z [][]float64, args string) {
	n := bb.Len()
	sx := utl.Sf("x%d", n)
	sy := utl.Sf("y%d", n)
	sz := utl.Sf("z%d", n)
	utl.GenMat(&bb, sx, x)
	utl.GenMat(&bb, sy, y)
	utl.GenMat(&bb, sz, z)
	cmd := utl.Sf("Contour(%s,%s,%s", sx, sy, sz)
	if len(args) > 0 {
		cmd += utl.Sf(",%s", args)
	}
	utl.Ff(&bb, "%s)\n", cmd)
}

func ContourSimple(x, y, z [][]float64, args string) {
	n := bb.Len()
	sx := utl.Sf("x%d", n)
	sy := utl.Sf("y%d", n)
	sz := utl.Sf("z%d", n)
	utl.GenMat(&bb, sx, x)
	utl.GenMat(&bb, sy, y)
	utl.GenMat(&bb, sz, z)
	cmd := utl.Sf("contour(%s,%s,%s", sx, sy, sz)
	if len(args) > 0 {
		cmd += utl.Sf(",%s", args)
	}
	utl.Ff(&bb, "%s)\n", cmd)
}

func Quiver(x, y, gx, gy [][]float64, args string) {
	n := bb.Len()
	sx := utl.Sf("x%d", n)
	sy := utl.Sf("y%d", n)
	sgx := utl.Sf("gx%d", n)
	sgy := utl.Sf("gy%d", n)
	utl.GenMat(&bb, sx, x)
	utl.GenMat(&bb, sy, y)
	utl.GenMat(&bb, sgx, gx)
	utl.GenMat(&bb, sgy, gy)
	cmd := utl.Sf("quiver(%s,%s,%s,%s", sx, sy, sgx, sgy)
	if len(args) > 0 {
		cmd += utl.Sf(",%s", args)
	}
	utl.Ff(&bb, "%s)\n", cmd)
}

func Grid(args string) {
	utl.Ff(&bb, "grid(%s)\n", args)
}

func Gll(xl, yl string, args string) {
	n := bb.Len()
	cmd := utl.Sf("lg%d = Gll(r'%s',r'%s'", n, xl, yl)
	if len(args) > 0 {
		cmd += utl.Sf(",%s", args)
	}
	utl.Ff(&bb, "%s)\nea.append(lg%d)\n", cmd, n)
}

func Clf() {
	utl.Ff(&bb, "clf()\n")
}

func SetFontSize(args string) {
	utl.Ff(&bb, "SetFontSize(%s)\n", args)
}

func SetForEps(prop, widpt float64) {
	utl.Ff(&bb, "SetForEps(%g,%g)\n", prop, widpt)
}

func SetForPng(prop, widpt float64, dpi int) {
	utl.Ff(&bb, "SetForPng(%g,%g,%d)\n", prop, widpt, dpi)
}

func Save(fname string) {
	var buf bytes.Buffer
	utl.Ff(&buf, "Save('%s', ea=ea, verbose=1)\n", fname)
	run(&buf)
}

func SaveD(dirout, fname string) {
	os.MkdirAll(dirout, 0777)
	var buf bytes.Buffer
	utl.Ff(&buf, "Save('%s/%s', ea=ea, verbose=1)\n", dirout, fname)
	run(&buf)
}

func SetAxis(xmin, xmax, ymin, ymax float64) {
	utl.Ff(&bb, "axis([%g, %g, %g, %g])\n", xmin, xmax, ymin, ymax)
}

func AxisOff() {
	utl.Ff(&bb, "axis('off')\n")
}

func Equal() {
	utl.Ff(&bb, "axis('equal')\n")
}

func Show() {
	utl.Ff(&bb, "show()\n")
	run(nil)
}

func Circle(xc, yc, r float64, args string) {
	cmd := utl.Sf("Circle(%g,%g,%g", xc, yc, r)
	if len(args) > 0 {
		cmd += utl.Sf(",%s", args)
	}
	utl.Ff(&bb, "%s)\n", cmd)
}

func run(extra *bytes.Buffer) {
	fn := "/tmp/gosl_mplotlib_go.py"
	if extra != nil {
		utl.WriteFile(fn, &ea, &bb, extra)
	} else {
		utl.WriteFile(fn, &ea, &bb)
	}
	cmd := exec.Command("python", fn)
	var out, serr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &serr
	err := cmd.Run()
	if err != nil {
		utl.Panic(_mplotlib_err1, serr.String())
	}
	utl.Pf("%s", out.String())
}

var (
	_mplotlib_err1 = "mplotlib.go: python failed:\n%v"
)
