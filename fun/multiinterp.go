package fun

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/utl"
	"math"
)

type InterpType int

const (
	biLinear InterpType = 2
	biCubic  InterpType = 3
)

// Axis implements a type to hold an arbitrarily spaced discrete data
type Axis struct {

	// configuration data
	DisableHunt bool // do not use hunt code at all

	// input data
	data []float64 // data array

	// derived data
	n       int  // length of data
	m       int  // number of points of  interpolating formula
	jHunt   int  // temporary j to decide on using hunt
	djHunt  int  // increent of j to decide on using hunt function or locate
	useHunt bool // use hunt code instead of locate
	ascnd   bool // ascending order of values

}

// Builds a new Axis type from a data slice for an InterpType
func NewAxis(data []float64, interpType InterpType) (o *Axis) {
	o = new(Axis)
	o.n = len(data)
	o.data = data
	switch interpType {
	case biLinear:
		o.m = 2
	case biCubic:
		o.m = 3
	}

	// check that axis is strictly monotonic
	if o.n > 2 {
		inc, dec := true, true
		for i := 1; i < o.n; i++ {
			if o.data[i] > o.data[i-1] {
				dec = false
			} else if o.data[i-1] > o.data[i] {
				inc = false
			}
		}
		if !inc && !dec {
			chk.Panic("Your Axis is not monotonic\n")
		}
		o.ascnd = inc
	} else {
		chk.Panic("length of an axis must be at least 2, %d is invalid\n", o.n)
	}

	o.djHunt = utl.Imin(1, int(math.Pow(float64(o.n), 0.25)))
	o.useHunt = false

	return
}

// returns the value at data[i]
func (o *Axis) Get(i int) float64 {
	if i >= o.n || i < 0 {
		chk.Panic("Axis out of bounds %d is not a valid query for axis of length %d", i, o.n)
	}
	return o.data[i]
}

// locate returns a value j such that x is (insofar as possible) centered in the subrange
// xx[j..j+mm-1], where xx is the stored pointer. The values in xx must be monotonic, either
// increasing or decreasing. The returned value is not less than 0, nor greater than n-1.
func (o *Axis) bisect(x float64) int {
	jl := 0
	ju := o.n - 1
	for ju-jl > 1 {
		jm := (ju + jl) >> 1
		if x >= o.data[jm] == o.ascnd {
			jl = jm
		} else {
			ju = jm
		}
	}

	if utl.Iabs(jl-o.jHunt) > o.djHunt {
		o.useHunt = false
	} else {
		o.useHunt = true
	}

	o.jHunt = jl

	return utl.Imax(0, utl.Imin(o.n-o.m, jl-((o.m-2)>>1)))

}

// hunt returns a value j such that x is (insofar as possible) centered in the subrange
// xx[j..j+mm-1], where xx is the stored pointer. The values in xx must be monotonic, either
// increasing or decreasing. The returned value is not less than 0, nor greater than n-1.
func (o *Axis) hunt(x float64) int {

	// hunting
	jl := o.jHunt
	inc := 1
	var ju, jm int
	if jl < 0 || jl > o.n-1 { // input guess not useful. skip hunting
		jl = 0
		ju = o.n - 1
	} else {
		if x >= o.Get(jl) == o.ascnd { // hunt up
			for {
				ju = jl + inc
				if ju >= o.n-1 {
					ju = o.n - 1
					break // off end of table.
				} else if x < o.Get(ju) == o.ascnd {
					break // found bracket.
				} else { // not done, so double the increment and try again.
					jl = ju
					inc += inc
				}
			}
		} else { // hunt down
			ju = jl
			for {
				jl = jl - inc
				if jl <= 0 {
					jl = 0
					break // off end of table.
				} else if x >= o.Get(jl) == o.ascnd {
					break // found bracket.
				} else { // not done, so double the increment and try again.
					ju = jl
					inc += inc
				}
			}
		}
	}

	// hunt is done, so begin the final bisection phase:
	for ju-jl > 1 {
		jm = (ju + jl) >> 1
		if x >= o.Get(jm) == o.ascnd {
			jl = jm
		} else {
			ju = jm
		}
	}

	// set hunt flag
	if utl.Iabs(jl-o.jHunt) > o.djHunt {
		o.useHunt = false
	} else {
		o.useHunt = true
	}
	o.jHunt = jl

	// results
	return utl.Imax(0, utl.Imin(o.n-o.m, jl-((o.m-2)>>1)))
}

func (o *Axis) locate(x float64) int {
	var jlo int
	if o.useHunt && !o.DisableHunt {
		jlo = o.hunt(x)
	} else {
		jlo = o.bisect(x)
	}
	return jlo
}

// BiLinear two dimensional interpolant
type BiLinear struct {
	// input data
	data *la.Matrix // column major data array
	yy   *Axis      // "y"
	xx   *Axis      // "x"
}

// This function builds a two dimensional bi-linear interpolant
// Input:
// 		xx -- function sample points abscissas
// 		yy -- function sample points ordinates
// 		f  -- function values
//			f(i,j) is stored at f[len(xx)*j + i]
//
// Ref:
//  f(x,y) = x^2 + 2y^2		 	  2 h  	i  	j
//  xx = [0.00,0.50,1.00]    		|
//  yy = [0.00,1.00,2.00] 	 	  1 d  	e  	f
//  f  = [a:0.00,b:0.25,c:1.00,		|
//		  d:2.00,e:2,25,f:3.00,		a___b___c
//		  h:8.00,i:8.25,j:9.00]	   0   0.5 1.0

func NewBiLinear(f, xx, yy []float64) (o *BiLinear) {
	o = new(BiLinear)

	o.Reset(f, xx, yy)

	return
}

// Disable the hunt function for both axis
func (o *BiLinear) SetDisableHunt(disable bool) {
	o.xx.DisableHunt = disable
	o.yy.DisableHunt = disable
}

// (Re)Set the axis and matrix for the interpolant
func (o *BiLinear) Reset(f, xx, yy []float64) {
	o.xx = NewAxis(xx, biLinear)
	o.yy = NewAxis(yy, biLinear)

	if len(f) != len(xx)*len(yy) {
		chk.Panic("Length of function matrix %d is not equal to axis' lengths %d*%d",
			len(f), len(xx), len(yy))
	}

	o.data = la.NewMatrixRaw(len(xx), len(yy), f)
}

func (o *BiLinear) locate(x, y float64) (i, j int) {
	i = o.xx.locate(x)
	j = o.yy.locate(y)
	return
}

func (o *BiLinear) P(x, y float64) float64 {
	i, j := o.locate(x, y)

	t := (x - o.xx.Get(i)) / (o.xx.Get(i+1) - o.xx.Get(i))
	u := (y - o.yy.Get(j)) / (o.yy.Get(j+1) - o.yy.Get(j))

	f11 := o.data.Get(i, j)
	f21 := o.data.Get(i+1, j)
	f12 := o.data.Get(i, j+1)
	f22 := o.data.Get(i+1, j+1)

	return (1-t)*(1-u)*f11 + t*(1-u)*f21 + (1-t)*u*f12 + t*u*f22
}
