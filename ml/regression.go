// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ml

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/rnd"
	"github.com/cpmech/gosl/utl"
)

// Regression defines the functions required to perform regression computations
type Regression interface {
	Model(x, θ la.Vector) float64        // model equation
	Cost(data *RegData) float64          // computes cost
	Deriv(dCdθ la.Vector, data *RegData) // computes dCdθ for given data
}

// RegData holds data for regression computations
type RegData struct {

	// input
	m int        // ndata: number of data points
	n int        // nparams: number of parameters = len(θ) = number of features + 1
	x *la.Matrix // [m][n] x-data
	y la.Vector  // [m] y-data
	l la.Vector  // [m] l = X ⋅ θ
	θ la.Vector  // [n] parameters θ

	// control
	statOk bool // indicates that Stat is OK; otherwise Stat() must be called because x,y changed

	// stat
	MinX, MaxX  []float64 // [nFeatures] x-limits
	SumX, MeanX []float64 // [nFeatures] sum and mean of x
	SigX        []float64 // [nFeatures] standard deviations of x
	DelX        []float64 // [nFeatures] MaxX - MinX
	MinY, MaxY  float64   // y-limits
	SumY, MeanY float64   // sum and mean of y
	SigY        float64   // standard deviation of y
	DelY        float64   // MaxY - MinY

	// plotting
	xx, yy       []float64   // [npts] curves
	thi, thj, cc [][]float64 // [npts][npts] meshgrid over θ
	xxi, xxj, zz [][]float64 // [npts][npts] meshgrid over x
}

// NewRegData returns a new structure to hold Regression Data
func NewRegData(nData, nFeatures int) (o *RegData) {
	o = new(RegData)
	o.m = nData
	o.n = nFeatures + 1
	o.x = la.NewMatrix(o.m, o.n)
	o.y = la.NewVector(o.m)
	o.l = la.NewVector(o.m)
	o.θ = la.NewVector(o.n)
	o.x.SetCol(0, 1.0)
	return
}

// NewRegDataTable sets X and Y values given table
//  table -- [ndata][nfeatures+1] data such that the last column correspond to y-values
func NewRegDataTable(table [][]float64) (o *RegData) {
	if len(table) < 1 {
		chk.Panic("at leat one row of data in table must be provided\n")
	}
	nData := len(table)
	nFeatures := len(table[0]) - 1
	o = NewRegData(nData, nFeatures)
	for i := 0; i < nData; i++ {
		for j := 0; j < nFeatures; j++ {
			o.SetX(i, j, table[i][j])
			o.SetY(i, table[i][j+1])
		}
	}
	o.Stat()
	return
}

// Ndata returns the number of data points
func (o *RegData) Ndata() int {
	return o.m
}

// Nparams returns the number of parameters = len(θ)
func (o *RegData) Nparams() int {
	return o.n
}

// Nfeatures returns the number of features = number of parameters - 1
func (o *RegData) Nfeatures() int {
	return o.n - 1
}

// SetX sets x-value
func (o *RegData) SetX(iData, jFeature int, value float64) {
	o.x.Set(iData, 1+jFeature, value)
	o.statOk = false
}

// SetY sets y-value
func (o *RegData) SetY(iData int, value float64) {
	o.y[iData] = value
	o.statOk = false
}

// GetFeature get all x-values of feature iFeature
func (o *RegData) GetFeature(iFeature int) (xValues []float64) {
	return o.x.GetCol(1 + iFeature)
}

// GetY returns all y-values
func (o *RegData) GetY() (yValues []float64) {
	yValues = make([]float64, o.m)
	copy(yValues, o.y)
	return
}

// Normalize normalizes x values
func (o *RegData) Normalize(useMinMax bool) {
	o.Stat()
	den := o.SigX
	if useMinMax {
		den = o.DelX
	}
	for i := 0; i < o.m; i++ {
		for J := 1; J < o.n; J++ {
			j := J - 1
			o.x.Set(i, J, (o.x.Get(i, J)-o.MeanX[j])/den[j])
		}
	}
	o.statOk = false
}

// Stat returns stat information about data
func (o *RegData) Stat() {
	if o.statOk {
		return
	}
	nf := o.Nfeatures()
	if len(o.MinX) != nf {
		o.MinX = make([]float64, nf)
		o.MaxX = make([]float64, nf)
		o.SumX = make([]float64, nf)
		o.MeanX = make([]float64, nf)
		o.SigX = make([]float64, nf)
		o.DelX = make([]float64, nf)
	}
	for J := 1; J < o.n; J++ {
		j := J - 1
		o.MinX[j] = o.x.Get(0, J)
		o.MaxX[j] = o.MinX[j]
		o.SumX[j] = 0.0
		for i := 0; i < o.m; i++ {
			x := o.x.Get(i, J)
			o.MinX[j] = utl.Min(o.MinX[j], x)
			o.MaxX[j] = utl.Max(o.MaxX[j], x)
			o.SumX[j] += x
		}
		o.MeanX[j] = o.SumX[j] / float64(o.m)
		o.SigX[j] = rnd.StatDevFirst(o.x.Col(J), o.MeanX[j], true)
		o.DelX[j] = o.MaxX[j] - o.MinX[j]
	}
	o.MinY = o.y[0]
	o.MaxY = o.MinY
	o.SumY = 0.0
	for i := 0; i < o.m; i++ {
		o.MinY = utl.Min(o.MinY, o.y[i])
		o.MaxY = utl.Max(o.MaxY, o.y[i])
		o.SumY += o.y[i]
	}
	o.MeanY = o.SumY / float64(o.m)
	o.SigY = rnd.StatDevFirst(o.y, o.MeanY, true)
	o.DelY = o.MaxY - o.MinY
	o.statOk = true
}

// PlotX plots X values
//   args -- maps int(y) values to plot arguments
func (o *RegData) PlotX(iFeature, jFeature int, args map[int]*plt.A) {
	a := &plt.A{C: "orange", M: ".", Ls: "none", NoClip: true}
	if args == nil {
		xx := make([]float64, o.m)
		yy := make([]float64, o.m)
		for k := 0; k < o.m; k++ {
			xx[k] = o.x.Get(k, 1+iFeature)
			yy[k] = o.x.Get(k, 1+jFeature)
		}
		plt.Plot(xx, yy, a)
	} else {
		var aa *plt.A
		for k := 0; k < o.m; k++ {
			aa = args[int(o.y[k])]
			if aa == nil {
				aa = a
			}
			aa.NoClip = true
			plt.PlotOne(o.x.Get(k, 1+iFeature), o.x.Get(k, 1+jFeature), aa)
		}
	}
	plt.Gll(io.Sf("$x_{%d}$", iFeature), io.Sf("$x_{%d}$", jFeature), nil)
}

// PlotModel plots regression model @ iFeature and all other values set to MinX
func (o *RegData) PlotModel(reg Regression, iFeature, npts int, args *plt.A) {
	if npts < 3 {
		npts = 11
	}
	if args == nil {
		args = &plt.A{C: plt.C(0, 0), NoClip: true}
	}
	o.Stat()
	nf := o.Nfeatures()
	xmdl := la.NewVector(nf)
	for j := 0; j < nf; j++ {
		xmdl[j] = o.MinX[j]
	}
	dxi := (o.MaxX[iFeature] - o.MinX[iFeature]) / float64(npts-1)
	if len(o.xx) != npts {
		o.xx = make([]float64, npts)
		o.yy = make([]float64, npts)
	}
	for k := 0; k < npts; k++ {
		o.xx[k] = o.MinX[iFeature] + float64(k)*dxi
		xmdl[iFeature] = o.xx[k]
		o.yy[k] = reg.Model(xmdl, o.θ)
	}
	plt.Plot(o.xx, o.yy, args)
}

// PlotContModel plots contour of model
//  minX and maxX may be nil
func (o *RegData) PlotContModel(reg Regression, iFeature, jFeature, npi, npj int, minX, maxX []float64, args *plt.A) {
	if len(o.xxi) != npj {
		o.xxi = utl.Alloc(npj, npi)
		o.xxj = utl.Alloc(npj, npi)
		o.zz = utl.Alloc(npj, npi)
	}
	o.Stat()
	nf := o.Nfeatures()
	xmdl := la.NewVector(nf)
	for j := 0; j < nf; j++ {
		xmdl[j] = o.MinX[j]
	}
	if minX == nil {
		minX = o.MinX
	}
	if maxX == nil {
		maxX = o.MaxX
	}
	dxi := (maxX[iFeature] - minX[iFeature]) / float64(npi-1)
	dxj := (maxX[jFeature] - minX[jFeature]) / float64(npj-1)
	for r := 0; r < npj; r++ {
		for c := 0; c < npi; c++ {
			o.xxi[r][c] = minX[iFeature] + float64(c)*dxi
			o.xxj[r][c] = minX[jFeature] + float64(r)*dxj
			xmdl[iFeature] = o.xxi[r][c]
			xmdl[jFeature] = o.xxj[r][c]
			o.zz[r][c] = reg.Model(xmdl, o.θ)
		}
	}
	if args == nil {
		args = &plt.A{Colors: []string{"k"}, Levels: []float64{0}}
	}
	plt.ContourL(o.xxi, o.xxj, o.zz, args)
	plt.Gll(io.Sf("$x_{%d}$", iFeature), io.Sf("$x_{%d}$", jFeature), nil)
}

// PlotContCost plots contour of cost
func (o *RegData) PlotContCost(reg Regression, iPrm, jPrm, npi, npj int, minθ, maxθ []float64, args *plt.A) {
	if len(o.thi) != npj {
		o.thi = utl.Alloc(npj, npi)
		o.thj = utl.Alloc(npj, npi)
		o.cc = utl.Alloc(npj, npi)
	}
	θcpy := o.θ.GetCopy()
	dthi := (maxθ[iPrm] - minθ[iPrm]) / float64(npi-1)
	dthj := (maxθ[jPrm] - minθ[jPrm]) / float64(npj-1)
	for r := 0; r < npj; r++ {
		for c := 0; c < npi; c++ {
			o.thi[r][c] = minθ[iPrm] + float64(c)*dthi
			o.thj[r][c] = minθ[jPrm] + float64(r)*dthj
			o.θ[iPrm] = o.thi[r][c]
			o.θ[jPrm] = o.thj[r][c]
			o.cc[r][c] = reg.Cost(o)
		}
	}
	o.θ.Apply(1, θcpy)
	plt.ContourF(o.thi, o.thj, o.cc, args)
	plt.Gll(io.Sf("$\\theta_{%d}$", iPrm), io.Sf("$\\theta_{%d}$", jPrm), nil)
}
