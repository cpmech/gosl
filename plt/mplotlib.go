// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// plt contains functions for plotting
package plt

import (
	"bytes"
	"os"
	"os/exec"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

var bb bytes.Buffer // the buffer
var ea bytes.Buffer // extra artists

func init() {
	Reset()
}

func Reset() {
	bb.Reset()
	ea.Reset()
	io.Ff(&bb, "from gosl import *\n")
	io.Ff(&ea, "ea = []\n")
}

func PyCmds(cmds string) {
	io.Ff(&bb, cmds)
}

func DoubleYscale(ylabelOrEmpty string) {
	io.Ff(&bb, "gca().twinx()\n")
	if ylabelOrEmpty != "" {
		io.Ff(&bb, "gca().set_ylabel('%s')\n", ylabelOrEmpty)
	}
}

func PyFile(filename string) {
	b, err := io.ReadFile(filename)
	if err != nil {
		chk.Panic("PyFile failed:\n%v", err)
	}
	io.Ff(&bb, string(b))
}

func SetXlog() {
	io.Ff(&bb, "SetXlog()\n")
}

func SetYlog() {
	io.Ff(&bb, "SetYlog()\n")
}

func SetXnticks(num int) {
	io.Ff(&bb, "SetXnticks(%d)\n", num)
}

func SetYnticks(num int) {
	io.Ff(&bb, "SetYnticks(%d)\n", num)
}

func SetScientific(axis string, min_order, max_order int) {
	io.Ff(&bb, "SetScientificFmt(axis='%s', min_order=%d, max_order=%d)\n", axis, min_order, max_order)
}

func SetTicksNormal() {
	io.Ff(&bb, "gca().ticklabel_format(useOffset=False)\n")
}

func Arrow(xi, yi, xf, yf float64, args string) {
	cmd := io.Sf("Arrow(%g,%g, %g,%g", xi, yi, xf, yf)
	if len(args) > 0 {
		cmd += io.Sf(",%s", args)
	}
	io.Ff(&bb, "%s)\n", cmd)
}

func AxHline(y float64, args string) {
	cmd := io.Sf("axhline(%g", y)
	if len(args) > 0 {
		cmd += io.Sf(",%s", args)
	}
	io.Ff(&bb, "%s)\n", cmd)
}

func AxVline(x float64, args string) {
	cmd := io.Sf("axvline(%g", x)
	if len(args) > 0 {
		cmd += io.Sf(",%s", args)
	}
	io.Ff(&bb, "%s)\n", cmd)
}

func HideTRframe() {
	io.Ff(&bb, "HideFrameLines()\n")
}

func Annotate(x, y float64, txt string, args string) {
	cmd := io.Sf("annotate(%s, xy=(%g,%g)", txt, x, y)
	if len(args) > 0 {
		cmd += io.Sf(",%s", args)
	}
	io.Ff(&bb, "%s)\n", cmd)
}

func AnnotateXlabels(x float64, txt string, args string) {
	cmd := io.Sf("AnnotateXlabels(%g, %s", x, txt)
	if len(args) > 0 {
		cmd += io.Sf(",%s", args)
	}
	io.Ff(&bb, "%s)\n", cmd)
}

func SupTitle(txt, args string) {
	n := bb.Len()
	if len(args) > 0 {
		io.Ff(&bb, "st%d = suptitle('%s',%s)\n", n, txt, args)
	} else {
		io.Ff(&bb, "st%d = suptitle('%s')\n", n, txt)
	}
	io.Ff(&bb, "ea.append(st%d)\n", n)
}

func Title(txt, args string) {
	if len(args) > 0 {
		io.Ff(&bb, "title('%s',%s)\n", txt, args)
	} else {
		io.Ff(&bb, "title('%s')\n", txt)
	}
}

func Text(x, y float64, txt, args string) {
	if len(args) > 0 {
		io.Ff(&bb, "text(%g,%g,'%s',%s)\n", x, y, txt, args)
	} else {
		io.Ff(&bb, "text(%g,%g,'%s')\n", x, y, txt)
	}
}

func Cross(args string) {
	if len(args) > 0 {
		io.Ff(&bb, "Cross(%s)\n", args)
	} else {
		io.Ff(&bb, "Cross()\n")
	}
}

func SplotGap(w, h float64) {
	io.Ff(&bb, "SplotGap(%g, %g)\n", w, h)
}

func Subplot(i, j, k int) {
	io.Ff(&bb, "subplot(%d,%d,%d)\n", i, j, k)
}

func SubplotI(I []int) {
	if len(I) != 3 {
		return
	}
	io.Ff(&bb, "subplot(%d,%d,%d)\n", I[0], I[1], I[2])
}

func SetHspace(hspace float64) {
	io.Ff(&bb, "subplots_adjust(hspace=%g)\n", hspace)
}

func SetVspace(vspace float64) {
	io.Ff(&bb, "subplots_adjust(vspace=%g)\n", vspace)
}

func AxisXmin(xmin float64) {
	io.Ff(&bb, "axis([%g, axis()[1], axis()[2], axis()[3]])\n", xmin)
}

func AxisXmax(xmax float64) {
	io.Ff(&bb, "axis([axis()[0], %g, axis()[2], axis()[3]])\n", xmax)
}

func AxisYmin(ymin float64) {
	io.Ff(&bb, "axis([axis()[0], axis()[1], %g, axis()[3]])\n", ymin)
}

func AxisYmax(ymax float64) {
	io.Ff(&bb, "axis([axis()[0], axis()[1], axis()[2], %g])\n", ymax)
}

func AxisXrange(xmin, xmax float64) {
	io.Ff(&bb, "axis([%g, %g, axis()[2], axis()[3]])\n", xmin, xmax)
}

func AxisYrange(ymin, ymax float64) {
	io.Ff(&bb, "axis([axis()[0], axis()[1], %g, %g])\n", ymin, ymax)
}

func AxisRange(xmin, xmax, ymin, ymax float64) {
	io.Ff(&bb, "axis([%g, %g, %g, %g])\n", xmin, xmax, ymin, ymax)
}

func AxisRange3d(xmin, xmax, ymin, ymax, zmin, zmax float64) {
	io.Ff(&bb, "gca().set_xlim3d(%g,%g)\ngca().set_ylim3d(%g,%g)\ngca().set_zlim3d(%g,%g)\n", xmin, xmax, ymin, ymax, zmin, zmax)
}

func AxisLims(lims []float64) {
	io.Ff(&bb, "axis([%g, %g, %g, %g])\n", lims[0], lims[1], lims[2], lims[3])
}

func Plot(x, y []float64, args string) {
	n := bb.Len()
	sx := io.Sf("x%d", n)
	sy := io.Sf("y%d", n)
	Gen2Arrays(&bb, sx, sy, x, y)
	if len(args) > 0 {
		io.Ff(&bb, "plot(%s,%s,%s)\n", sx, sy, args)
	} else {
		io.Ff(&bb, "plot(%s,%s)\n", sx, sy)
	}
}

func PlotOne(x, y float64, args string) {
	if len(args) > 0 {
		io.Ff(&bb, "plot(%23.15e,%23.15e,%s)\n", x, y, args)
	} else {
		io.Ff(&bb, "plot(%23.15e,%23.15e)\n", x, y)
	}
}

func Hist(x [][]float64, labels []string, args string) {
	n := bb.Len()
	sx := io.Sf("x%d", n)
	sy := io.Sf("y%d", n)
	GenList(&bb, sx, x)
	GenStrArray(&bb, sy, labels)
	if len(args) > 0 {
		io.Ff(&bb, "hist(%s,label=%s,%s)\n", sx, sy, args)
	} else {
		io.Ff(&bb, "hist(%s,label=%s)\n", sx, sy)
	}
}

func Plot3dLine(x, y, z []float64, first bool, args string) {
	n := bb.Len()
	sx := io.Sf("x%d", n)
	sy := io.Sf("y%d", n)
	sz := io.Sf("z%d", n)
	GenArray(&bb, sx, x)
	GenArray(&bb, sy, y)
	GenArray(&bb, sz, z)
	ifirst := 0
	if first {
		ifirst = 1
	}
	cmd := io.Sf("Plot3dLine(%s,%s,%s,%d", sx, sy, sz, ifirst)
	if len(args) > 0 {
		cmd += io.Sf(",%s", args)
	}
	io.Ff(&bb, "%s)\n", cmd)
}

func Plot3dPoints(x, y, z []float64, args string) {
	n := bb.Len()
	sx := io.Sf("x%d", n)
	sy := io.Sf("y%d", n)
	sz := io.Sf("z%d", n)
	GenArray(&bb, sx, x)
	GenArray(&bb, sy, y)
	GenArray(&bb, sz, z)
	cmd := io.Sf("ax%d = Plot3dPoints(%s,%s,%s", n, sx, sy, sz)
	if len(args) > 0 {
		cmd += io.Sf(",%s", args)
	}
	io.Ff(&bb, "%s)\n", cmd)
	io.Ff(&bb, "ea.append(ax%d)\n", n)
}

func Wireframe(x, y, z [][]float64, args string) {
	n := bb.Len()
	sx := io.Sf("x%d", n)
	sy := io.Sf("y%d", n)
	sz := io.Sf("z%d", n)
	GenMat(&bb, sx, x)
	GenMat(&bb, sy, y)
	GenMat(&bb, sz, z)
	cmd := io.Sf("ax%d = Wireframe(%s,%s,%s", n, sx, sy, sz)
	if len(args) > 0 {
		cmd += io.Sf(",%s", args)
	}
	io.Ff(&bb, "%s)\n", cmd)
	io.Ff(&bb, "ea.append(ax%d)\n", n)
}

func Surface(x, y, z [][]float64, args string) {
	n := bb.Len()
	sx := io.Sf("x%d", n)
	sy := io.Sf("y%d", n)
	sz := io.Sf("z%d", n)
	GenMat(&bb, sx, x)
	GenMat(&bb, sy, y)
	GenMat(&bb, sz, z)
	cmd := io.Sf("Surface(%s,%s,%s", sx, sy, sz)
	if len(args) > 0 {
		cmd += io.Sf(",%s", args)
	}
	io.Ff(&bb, "%s)\n", cmd)
}

func Contour(x, y, z [][]float64, args string) {
	n := bb.Len()
	sx := io.Sf("x%d", n)
	sy := io.Sf("y%d", n)
	sz := io.Sf("z%d", n)
	GenMat(&bb, sx, x)
	GenMat(&bb, sy, y)
	GenMat(&bb, sz, z)
	cmd := io.Sf("Contour(%s,%s,%s", sx, sy, sz)
	if len(args) > 0 {
		cmd += io.Sf(",%s", args)
	}
	io.Ff(&bb, "%s)\n", cmd)
}

func ContourSimple(x, y, z [][]float64, withClabel bool, clabelFsz float64, args string) {
	n := bb.Len()
	sx := io.Sf("x%d", n)
	sy := io.Sf("y%d", n)
	sz := io.Sf("z%d", n)
	GenMat(&bb, sx, x)
	GenMat(&bb, sy, y)
	GenMat(&bb, sz, z)
	cmd := io.Sf("ctour%d = contour(%s,%s,%s", n, sx, sy, sz)
	if len(args) > 0 {
		cmd += io.Sf(",%s", args)
	}
	io.Ff(&bb, "%s)\n", cmd)
	if withClabel {
		io.Ff(&bb, "clabel(ctour%d,inline=1,fsz=%g)\n", n, clabelFsz)
	}
}

func Camera(elev, azim float64, args string) {
	cmd := io.Sf("gca().view_init(elev=%g, azim=%g", elev, azim)
	if len(args) > 0 {
		cmd += io.Sf(",%s", args)
	}
	io.Ff(&bb, "%s)\n", cmd)
}

func AxDist(dist float64) {
	io.Ff(&bb, "gca().dist = %g\n", dist)
}

func Quiver(x, y, gx, gy [][]float64, args string) {
	n := bb.Len()
	sx := io.Sf("x%d", n)
	sy := io.Sf("y%d", n)
	sgx := io.Sf("gx%d", n)
	sgy := io.Sf("gy%d", n)
	GenMat(&bb, sx, x)
	GenMat(&bb, sy, y)
	GenMat(&bb, sgx, gx)
	GenMat(&bb, sgy, gy)
	cmd := io.Sf("quiver(%s,%s,%s,%s", sx, sy, sgx, sgy)
	if len(args) > 0 {
		cmd += io.Sf(",%s", args)
	}
	io.Ff(&bb, "%s)\n", cmd)
}

func Grid(args string) {
	io.Ff(&bb, "grid(%s)\n", args)
}

func Gll(xl, yl string, args string) {
	n := bb.Len()
	cmd := io.Sf("lg%d = Gll(r'%s',r'%s'", n, xl, yl)
	if len(args) > 0 {
		cmd += io.Sf(",%s", args)
	}
	io.Ff(&bb, "%s)\nea.append(lg%d)\n", cmd, n)
}

func Clf() {
	io.Ff(&bb, "clf()\n")
}

func SetFontSize(args string) {
	io.Ff(&bb, "SetFontSize(%s)\n", args)
}

func SetForEps(prop, widpt float64) {
	io.Ff(&bb, "SetForEps(%g,%g)\n", prop, widpt)
}

func SetForPng(prop, widpt float64, dpi int) {
	io.Ff(&bb, "SetForPng(%g,%g,%d)\n", prop, widpt, dpi)
}

func Save(fname string) {
	var buf bytes.Buffer
	io.Ff(&buf, "Save('%s', ea=ea, verbose=1)\n", fname)
	run(&buf)
}

func SaveD(dirout, fname string) {
	os.MkdirAll(dirout, 0777)
	var buf bytes.Buffer
	io.Ff(&buf, "Save('%s/%s', ea=ea, verbose=1)\n", dirout, fname)
	run(&buf)
}

func SetAxis(xmin, xmax, ymin, ymax float64) {
	io.Ff(&bb, "axis([%g, %g, %g, %g])\n", xmin, xmax, ymin, ymax)
}

func AxisOff() {
	io.Ff(&bb, "axis('off')\n")
}

func Equal() {
	io.Ff(&bb, "axis('equal')\n")
}

func Show() {
	io.Ff(&bb, "show()\n")
	run(nil)
}

func Circle(xc, yc, r float64, args string) {
	cmd := io.Sf("Circle(%g,%g,%g", xc, yc, r)
	if len(args) > 0 {
		cmd += io.Sf(",%s", args)
	}
	io.Ff(&bb, "%s)\n", cmd)
}

func run(extra *bytes.Buffer) {
	fn := "/tmp/gosl_mplotlib_go.py"
	if extra != nil {
		io.WriteFile(fn, &ea, &bb, extra)
	} else {
		io.WriteFile(fn, &ea, &bb)
	}
	cmd := exec.Command("python", fn)
	var out, serr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &serr
	err := cmd.Run()
	if err != nil {
		chk.Panic(_mplotlib_err1, serr.String())
	}
	io.Pf("%s", out.String())
}

// GenMat generates matrix
func GenMat(buf *bytes.Buffer, name string, a [][]float64) {
	io.Ff(buf, "%s=array([", name)
	for i, _ := range a {
		io.Ff(buf, "[")
		for j, _ := range a[i] {
			io.Ff(buf, "%g,", a[i][j])
		}
		io.Ff(buf, "],")
	}
	io.Ff(buf, "],dtype=float)\n")
}

// GenList generates list
func GenList(buf *bytes.Buffer, name string, a [][]float64) {
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

// GenArray generates the NumPy text in 'b' corresponding to an array of float point numbers
func GenArray(b *bytes.Buffer, name string, u []float64) {
	io.Ff(b, "%s=array([", name)
	for i, _ := range u {
		io.Ff(b, "%g,", u[i])
	}
	io.Ff(b, "],dtype=float)\n")
}

// Gen2Arrays generates the NumPy text in 'b' corresponding to 2 arrays of float point numbers
func Gen2Arrays(buf *bytes.Buffer, nameA, nameB string, a, b []float64) {
	GenArray(buf, nameA, a)
	GenArray(buf, nameB, b)
}

// GenStrArray generates the NumPy text in 'b' corresponding to an array of strings
func GenStrArray(b *bytes.Buffer, name string, u []string) {
	io.Ff(b, "%s=[", name)
	for i, _ := range u {
		io.Ff(b, "%q,", u[i])
	}
	io.Ff(b, "]\n")
}

var (
	_mplotlib_err1 = "mplotlib.go: python failed:\n%v"
)
