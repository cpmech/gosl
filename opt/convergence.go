package opt

import (
	"math"

	"gosl/fun"
	"gosl/fun/dbf"
	"gosl/la"
	"gosl/utl"
)

// Convergence assists in checking the convergence of numerical optimizers
// Convergence can be accessed to set convergence parameters, max iteration number,
// or to enable and access history of iterations.
type Convergence struct {

	// input
	Ffcn fun.Sv // objective function; scalar function of vector: y = f({x})
	Gfcn fun.Vv // gradient function: vector function of vector: g = dy/d{x} = deriv(f({x}), {x}) [may be nil]

	// configuration
	MaxIt   int     // max iterations
	Ftol    float64 // tolerance on f({x})
	Gtol    float64 // convergence criterion for the zero gradient test
	EpsF    float64 // small number to rectify the special case of converging to exactly zero function value
	UseHist bool    // save history
	Verbose bool    // show messages

	// statistics and History (e.g. for debugging)
	NumFeval int      // number of calls to Ffcn (function evaluations)
	NumGeval int      // number of calls to Gfcn (Jacobian evaluations)
	NumIter  int      // number of iterations from last call to Solve
	Hist     *History // history of optimization data (for debugging)

	// internal
	uhist la.Vector // direction of descents to be saved in History
}

// InitConvergence initialize convergence parameters
func (o *Convergence) InitConvergence(Ffcn fun.Sv, Gfcn fun.Vv) {
	o.Ffcn = func(x la.Vector) float64 {
		o.NumFeval++
		return Ffcn(x)
	}
	o.Gfcn = func(g, x la.Vector) {
		o.NumGeval++
		Gfcn(g, x)
	}
	o.MaxIt = 200
	o.Ftol = 1e-8
	o.Gtol = 1e-6
	o.EpsF = 1e-18
}

// InitHist initializes history
func (o *Convergence) InitHist(x0 la.Vector) {
	fmin := o.Ffcn(x0)
	o.Hist = NewHistory(o.MaxIt, fmin, x0, o.Ffcn)
	o.uhist = la.NewVector(len(x0))
}

// SetParams sets parameters
//   Example:
//             o.SetParams(dbf.NewParams(
//                 &dbf.P{N: "maxit", V: 1000},
//                 &dbf.P{N: "ftol", V: 1e-2},
//                 &dbf.P{N: "gtol", V: 1e-2},
//                 &dbf.P{N: "hist", V: 1},
//                 &dbf.P{N: "verb", V: 1},
//             ))
func (o *Convergence) SetParams(params dbf.Params) {
	o.MaxIt = params.GetIntOrDefault("maxit", o.MaxIt)
	o.Ftol = params.GetValueOrDefault("ftol", o.Ftol)
	o.Gtol = params.GetValueOrDefault("gtol", o.Gtol)
	o.UseHist = params.GetBoolOrDefault("hist", o.UseHist)
	o.Verbose = params.GetBoolOrDefault("verb", o.Verbose)
}

// SetConvParams sets convergence parameters
func (o *Convergence) SetConvParams(maxIt int, ftol, gtol float64) {
	o.MaxIt = maxIt
	o.Ftol = ftol
	o.Gtol = gtol
}

// SetUseHistory sets use history parameter
func (o *Convergence) SetUseHistory(useHist bool) {
	o.UseHist = useHist
}

// SetVerbose sets verbose mode
func (o *Convergence) SetVerbose(verbose bool) {
	o.Verbose = verbose
}

// AccessHistory gets access to History
func (o *Convergence) AccessHistory() *History {
	return o.Hist
}

// Fconvergence performs the check for f({x}) values
//
//   Input:
//     fprev -- a previous f({x}) value
//     fmin -- current f({x}) value
//
//   Output:
//     returns true if f values are not changing any longer
//
func (o *Convergence) Fconvergence(fprev, fmin float64) bool {
	if 2.0*math.Abs(fmin-fprev) <= o.Ftol*(math.Abs(fmin)+math.Abs(fprev)+o.EpsF) {
		return true // f values are not changing
	}
	return false // not yet
}

// Gconvergence performs the check for df/dx|({x}) values
//
//   Input:
//     fprev -- a previous f({x}) value (for normalization purposes)
//     x -- current {x} value
//     u -- current direction; e.g. dfdx
//
//   Output:
//     returns true if dfdy values are not changing any longer
//
func (o *Convergence) Gconvergence(fprev float64, x, u la.Vector) bool {
	var temp float64
	size := len(x)
	test := 0.0
	coef := utl.Max(fprev, 1.0)
	for j := 0; j < size; j++ {
		temp = math.Abs(u[j]) * utl.Max(math.Abs(x[j]), 1.0) / coef
		if temp > test {
			test = temp
		}
	}
	if test < o.Gtol {
		return true // dfdy values are not changing
	}
	return false // not yet
}
