// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/la/oblas"
)

/*  In the following code:
      start, endp1 := (icpu*len(v))/ncpu, ((icpu+1)*len(v))/ncpu
    multiplication has to be done first. Nonetheless, parentheses are
    acually optional since in Go:
     "Binary operators of the same precedence associate from left to right.
      For instance, x / y * z is the same as (x / y) * z" [http://golang.org/ref/spec]
*/

// variables for parallel code
var Pll bool = false // use parallel version of routines?
var NCPU int = 16    // number of CPUs to use

// --------------------------------------------------------------------------------------------------
// vector -------------------------------------------------------------------------------------------
// --------------------------------------------------------------------------------------------------

// VecFill fills a vector with a single number s:
//  v := s*ones(len(v))  =>  vi = s
func VecFill(v []float64, s float64) {
	if Pll {
		ncpu := imin(len(v), NCPU)
		ch := make(chan int, ncpu)
		for icpu := 0; icpu < ncpu; icpu++ {
			start, endp1 := (icpu*len(v))/ncpu, ((icpu+1)*len(v))/ncpu
			go func() {
				for i := start; i < endp1; i++ {
					v[i] = s
				}
				ch <- 1
			}()
		}
		for icpu := 0; icpu < ncpu; icpu++ {
			<-ch
		}
	} else {
		for i := 0; i < len(v); i++ {
			v[i] = s
		}
	}
}

// VecFillC fills a complex vector with a single number s:
//  v := s*ones(len(v))  =>  vi = s
func VecFillC(v []complex128, s complex128) {
	if Pll {
		ncpu := imin(len(v), NCPU)
		ch := make(chan int, ncpu)
		for icpu := 0; icpu < ncpu; icpu++ {
			start, endp1 := (icpu*len(v))/ncpu, ((icpu+1)*len(v))/ncpu
			go func() {
				for i := start; i < endp1; i++ {
					v[i] = s
				}
				ch <- 1
			}()
		}
		for icpu := 0; icpu < ncpu; icpu++ {
			<-ch
		}
	} else {
		for i := 0; i < len(v); i++ {
			v[i] = s
		}
	}
}

// VecApplyFunc runs a function over all components of a vector
//  vi = f(i,vi)
func VecApplyFunc(v []float64, f func(i int, x float64) float64) {
	if Pll {
		ncpu := imin(len(v), NCPU)
		ch := make(chan int, ncpu)
		for icpu := 0; icpu < ncpu; icpu++ {
			start, endp1 := (icpu*len(v))/ncpu, ((icpu+1)*len(v))/ncpu
			go func() {
				for i := start; i < endp1; i++ {
					v[i] = f(i, v[i])
				}
				ch <- 1
			}()
		}
		for icpu := 0; icpu < ncpu; icpu++ {
			<-ch
		}
	} else {
		for i := 0; i < len(v); i++ {
			v[i] = f(i, v[i])
		}
	}
}

// VecGetMapped returns a new vector after applying a function over all of its components
//  new: vi = f(i)
func VecGetMapped(dim int, f func(i int) float64) (v []float64) {
	v = make([]float64, dim)
	if Pll {
		ncpu := imin(len(v), NCPU)
		ch := make(chan int, ncpu)
		for icpu := 0; icpu < ncpu; icpu++ {
			start, endp1 := (icpu*len(v))/ncpu, ((icpu+1)*len(v))/ncpu
			go func() {
				for i := start; i < endp1; i++ {
					v[i] = f(i)
				}
				ch <- 1
			}()
		}
		for icpu := 0; icpu < ncpu; icpu++ {
			<-ch
		}
	} else {
		for i := 0; i < len(v); i++ {
			v[i] = f(i)
		}
	}
	return
}

// VecClone allocates a clone of a vector
//  b := a
func VecClone(a []float64) (b []float64) {
	b = make([]float64, len(a))
	for i := 0; i < len(a); i++ {
		b[i] = a[i]
	}
	return
}

// VecAccum sum/accumulates all components in a vector
//  sum := Σ_i v[i]
func VecAccum(v []float64) (sum float64) {
	if Pll {
		ncpu := imin(len(v), NCPU)
		chsum := make(chan float64, ncpu)
		for icpu := 0; icpu < ncpu; icpu++ {
			start, endp1 := (icpu*len(v))/ncpu, ((icpu+1)*len(v))/ncpu
			go func() {
				var mysum float64
				for i := start; i < endp1; i++ {
					mysum += v[i]
				}
				chsum <- mysum
			}()
		}
		sum = <-chsum
		for icpu := 1; icpu < ncpu; icpu++ {
			sum += <-chsum
		}
	} else {
		for i := 0; i < len(v); i++ {
			sum += v[i]
		}
	}
	return
}

// VecNorm returns the Euclidian norm of a vector:
//  nrm := ||v||
func VecNorm(v []float64) (nrm float64) {
	if Pll {
		ncpu := imin(len(v), NCPU)
		chnrm := make(chan float64, ncpu)
		for icpu := 0; icpu < ncpu; icpu++ {
			start, endp1 := (icpu*len(v))/ncpu, ((icpu+1)*len(v))/ncpu
			go func() {
				var mynrm float64
				for i := start; i < endp1; i++ {
					mynrm += v[i] * v[i]
				}
				chnrm <- mynrm
			}()
		}
		nrm = <-chnrm
		for icpu := 1; icpu < ncpu; icpu++ {
			nrm += <-chnrm
		}
	} else {
		for i := 0; i < len(v); i++ {
			nrm += v[i] * v[i]
		}
	}
	nrm = math.Sqrt(nrm)
	return
}

// VecNormDiff returns the Euclidian norm of the difference:
//  nrm := ||u - v||
func VecNormDiff(u, v []float64) (nrm float64) {
	if Pll {
		ncpu := imin(len(v), NCPU)
		chnrm := make(chan float64, ncpu)
		for icpu := 0; icpu < ncpu; icpu++ {
			start, endp1 := (icpu*len(v))/ncpu, ((icpu+1)*len(v))/ncpu
			go func() {
				var mynrm float64
				for i := start; i < endp1; i++ {
					mynrm += (u[i] - v[i]) * (u[i] - v[i])
				}
				chnrm <- mynrm
			}()
		}
		nrm = <-chnrm
		for icpu := 1; icpu < ncpu; icpu++ {
			nrm += <-chnrm
		}
	} else {
		for i := 0; i < len(v); i++ {
			nrm += (u[i] - v[i]) * (u[i] - v[i])
		}
	}
	nrm = math.Sqrt(nrm)
	return
}

// VecDot returns the dot product between two vectors:
//  s := u dot v
func VecDot(u, v []float64) (res float64) {
	if Pll {
		ncpu := imin(len(u), NCPU)
		ch := make(chan float64, ncpu)
		for icpu := 0; icpu < ncpu; icpu++ {
			start, endp1 := (icpu*len(u))/ncpu, ((icpu+1)*len(u))/ncpu
			go func() {
				myres := u[start] * v[start]
				for i := start + 1; i < endp1; i++ {
					myres += u[i] * v[i]
				}
				ch <- myres
			}()
		}
		res = <-ch
		for icpu := 1; icpu < ncpu; icpu++ {
			res += <-ch
		}
	} else {
		for i := 0; i < len(u); i++ {
			res += u[i] * v[i]
		}
	}
	return
}

// VecCopy copies a vector "b" into vector "a" (scaled):
//  a := α * b  =>  ai := α * bi
func VecCopy(a []float64, α float64, b []float64) {
	if Pll {
		ncpu := imin(len(a), NCPU)
		ch := make(chan int, ncpu)
		for icpu := 0; icpu < ncpu; icpu++ {
			start, endp1 := (icpu*len(a))/ncpu, ((icpu+1)*len(a))/ncpu
			go func() {
				for i := start; i < endp1; i++ {
					a[i] = α * b[i]
				}
				ch <- 1
			}()
		}
		for icpu := 0; icpu < ncpu; icpu++ {
			<-ch
		}
	} else {
		for i := 0; i < len(b); i++ {
			a[i] = α * b[i]
		}
	}
}

// VecAdd adds to vector "a", another vector "b" (scaled):
//  a += α * b  =>  ai += α * bi
func VecAdd(a []float64, α float64, b []float64) {
	if Pll {
		ncpu := imin(len(a), NCPU)
		ch := make(chan int, ncpu)
		for icpu := 0; icpu < ncpu; icpu++ {
			start, endp1 := (icpu*len(a))/ncpu, ((icpu+1)*len(a))/ncpu
			go func() {
				for i := start; i < endp1; i++ {
					a[i] += α * b[i]
				}
				ch <- 1
			}()
		}
		for icpu := 0; icpu < ncpu; icpu++ {
			<-ch
		}
	} else {
		for i := 0; i < len(b); i++ {
			a[i] += α * b[i]
		}
	}
}

// VecAdd2 adds two vectors (scaled):
//  u := α*a + β*b  =>  ui := α*ai + β*bi
func VecAdd2(u []float64, α float64, a []float64, β float64, b []float64) {
	if Pll {
		ncpu := imin(len(a), NCPU)
		ch := make(chan int, ncpu)
		for icpu := 0; icpu < ncpu; icpu++ {
			start, endp1 := (icpu*len(a))/ncpu, ((icpu+1)*len(a))/ncpu
			go func() {
				for i := start; i < endp1; i++ {
					u[i] = α*a[i] + β*b[i]
				}
				ch <- 1
			}()
		}
		for icpu := 0; icpu < ncpu; icpu++ {
			<-ch
		}
	} else {
		for i := 0; i < len(b); i++ {
			u[i] = α*a[i] + β*b[i]
		}
	}
}

// VecMin returns the minimum component of a vector
func VecMin(v []float64) (min float64) {
	if Pll {
		ncpu := imin(len(v), NCPU)
		chmin := make(chan float64, ncpu)
		for icpu := 0; icpu < ncpu; icpu++ {
			start, endp1 := (icpu*len(v))/ncpu, ((icpu+1)*len(v))/ncpu
			go func() {
				mymin := v[start]
				for i := start + 1; i < endp1; i++ {
					if v[i] < mymin {
						mymin = v[i]
					}
				}
				chmin <- mymin
			}()
		}
		min = <-chmin
		for icpu := 1; icpu < ncpu; icpu++ {
			othermin := <-chmin
			if othermin < min {
				min = othermin
			}
		}
	} else {
		min = v[0]
		for i := 1; i < len(v); i++ {
			if v[i] < min {
				min = v[i]
			}
		}
	}
	return
}

// VecMax returns the maximum component of a vector
func VecMax(v []float64) (max float64) {
	if Pll {
		ncpu := imin(len(v), NCPU)
		chmax := make(chan float64, ncpu)
		for icpu := 0; icpu < ncpu; icpu++ {
			start, endp1 := (icpu*len(v))/ncpu, ((icpu+1)*len(v))/ncpu
			go func() {
				mymax := v[start]
				for i := start + 1; i < endp1; i++ {
					if v[i] > mymax {
						mymax = v[i]
					}
				}
				chmax <- mymax
			}()
		}
		max = <-chmax
		for icpu := 1; icpu < ncpu; icpu++ {
			othermax := <-chmax
			if othermax > max {
				max = othermax
			}
		}
	} else {
		max = v[0]
		for i := 1; i < len(v); i++ {
			if v[i] > max {
				max = v[i]
			}
		}
	}
	return
}

// VecMinMax returns the min and max components of a vector
func VecMinMax(v []float64) (min, max float64) {
	if Pll {
		ncpu := imin(len(v), NCPU)
		chminmax := make(chan []float64, ncpu)
		for icpu := 0; icpu < ncpu; icpu++ {
			start, endp1 := (icpu*len(v))/ncpu, ((icpu+1)*len(v))/ncpu
			go func() {
				mymin, mymax := v[start], v[start]
				for i := start + 1; i < endp1; i++ {
					if v[i] < mymin {
						mymin = v[i]
					}
					if v[i] > mymax {
						mymax = v[i]
					}
				}
				chminmax <- []float64{mymin, mymax}
			}()
		}
		minmax := <-chminmax
		min, max = minmax[0], minmax[1]
		for icpu := 1; icpu < ncpu; icpu++ {
			otherminmax := <-chminmax
			if otherminmax[0] < min {
				min = otherminmax[0]
			}
			if otherminmax[1] > max {
				max = otherminmax[1]
			}
		}
	} else {
		min, max = v[0], v[0]
		for i := 1; i < len(v); i++ {
			if v[i] < min {
				min = v[i]
			}
			if v[i] > max {
				max = v[i]
			}
		}
	}
	return
}

// VecLargest returns the largest component (abs(u_i)) of a vector, normalised by den
func VecLargest(u []float64, den float64) (largest float64) {
	if Pll {
		ncpu := imin(len(u), NCPU)
		ch := make(chan float64, ncpu)
		for icpu := 0; icpu < ncpu; icpu++ {
			start, endp1 := (icpu*len(u))/ncpu, ((icpu+1)*len(u))/ncpu
			go func() {
				mylargest := math.Abs(u[start]) / den
				for i := start + 1; i < endp1; i++ {
					tmp := math.Abs(u[i]) / den
					if tmp > mylargest {
						mylargest = tmp
					}
				}
				ch <- mylargest
			}()
		}
		largest = <-ch
		for icpu := 1; icpu < ncpu; icpu++ {
			otherlargest := <-ch
			if otherlargest > largest {
				largest = otherlargest
			}
		}
	} else {
		largest = math.Abs(u[0]) / den
		for i := 1; i < len(u); i++ {
			tmp := math.Abs(u[i]) / den
			if tmp > largest {
				largest = tmp
			}
		}
	}
	return
}

// VecMaxDiff returns the maximum difference between components of two vectors
func VecMaxDiff(a, b []float64) (maxdiff float64) {
	if Pll {
		ncpu := imin(len(a), NCPU)
		ch := make(chan float64, ncpu)
		for icpu := 0; icpu < ncpu; icpu++ {
			start, endp1 := (icpu*len(a))/ncpu, ((icpu+1)*len(a))/ncpu
			go func() {
				mymaxdiff := math.Abs(a[start] - b[start])
				for i := start + 1; i < endp1; i++ {
					diff := math.Abs(a[i] - b[i])
					if diff > mymaxdiff {
						mymaxdiff = diff
					}
				}
				ch <- mymaxdiff
			}()
		}
		maxdiff = <-ch
		for icpu := 1; icpu < ncpu; icpu++ {
			othermaxdiff := <-ch
			if othermaxdiff > maxdiff {
				maxdiff = othermaxdiff
			}
		}
	} else {
		maxdiff = math.Abs(a[0] - b[0])
		for i := 1; i < len(a); i++ {
			diff := math.Abs(a[i] - b[i])
			if diff > maxdiff {
				maxdiff = diff
			}
		}
	}
	return
}

// VecMaxDiffC returns the maximum difference between components of two complex vectors
func VecMaxDiffC(a, b []complex128) (maxdiff float64) {
	if Pll {
		ncpu := imin(len(a), NCPU)
		ch := make(chan float64, ncpu)
		for icpu := 0; icpu < ncpu; icpu++ {
			start, endp1 := (icpu*len(a))/ncpu, ((icpu+1)*len(a))/ncpu
			go func() {
				mymaxdiff := math.Abs(real(a[start]) - real(b[start]))
				mymaxdiffz := math.Abs(imag(a[start]) - imag(b[start]))
				for i := start + 1; i < endp1; i++ {
					diff := math.Abs(real(a[i]) - real(b[i]))
					diffz := math.Abs(imag(a[i]) - imag(b[i]))
					if diff > mymaxdiff {
						mymaxdiff = diff
					}
					if diffz > mymaxdiffz {
						mymaxdiffz = diffz
					}
				}
				if mymaxdiff > mymaxdiffz {
					ch <- mymaxdiff
				} else {
					ch <- mymaxdiffz
				}
			}()
		}
		maxdiff = <-ch
		for icpu := 1; icpu < ncpu; icpu++ {
			othermaxdiff := <-ch
			if othermaxdiff > maxdiff {
				maxdiff = othermaxdiff
			}
		}
	} else {
		maxdiff = math.Abs(real(a[0]) - real(b[0]))
		maxdiffz := math.Abs(imag(a[0]) - imag(b[0]))
		for i := 1; i < len(a); i++ {
			diff := math.Abs(real(a[i]) - real(b[i]))
			diffz := math.Abs(imag(a[i]) - imag(b[i]))
			if diff > maxdiff {
				maxdiff = diff
			}
			if diffz > maxdiffz {
				maxdiffz = diffz
			}
		}
		if maxdiffz > maxdiff {
			maxdiff = maxdiffz
		}
	}
	return
}

// VecScale scales vector "v" using an absolute value (Atol) and a multiplier (Rtol)
//  res[i] := Atol + Rtol * v[i]
func VecScale(res []float64, Atol, Rtol float64, v []float64) {
	if Pll {
		ncpu := imin(len(v), NCPU)
		ch := make(chan int, ncpu)
		for icpu := 0; icpu < ncpu; icpu++ {
			start, endp1 := (icpu*len(v))/ncpu, ((icpu+1)*len(v))/ncpu
			go func() {
				for i := start; i < endp1; i++ {
					res[i] = Atol + Rtol*v[i]
				}
				ch <- 1
			}()
		}
		for i := 0; i < ncpu; i++ {
			<-ch
		}
	} else {
		for i := 0; i < len(v); i++ {
			res[i] = Atol + Rtol*v[i]
		}
	}
}

// VecScaleAbs scales vector abs(v) using an absolute value (Atol) and a multiplier (Rtol)
//  res[i] := Atol + Rtol * Abs(v[i])
func VecScaleAbs(res []float64, Atol, Rtol float64, v []float64) {
	if Pll {
		ncpu := imin(len(v), NCPU)
		ch := make(chan int, ncpu)
		for icpu := 0; icpu < ncpu; icpu++ {
			start, endp1 := (icpu*len(v))/ncpu, ((icpu+1)*len(v))/ncpu
			go func() {
				for i := start; i < endp1; i++ {
					res[i] = Atol + Rtol*math.Abs(v[i])
				}
				ch <- 1
			}()
		}
		for i := 0; i < ncpu; i++ {
			<-ch
		}
	} else {
		for i := 0; i < len(v); i++ {
			res[i] = Atol + Rtol*math.Abs(v[i])
		}
	}
}

// VecRms returns the root-mean-square of a vector u:
//  rms := sqrt(mean((u[:])^2))  ==  sqrt(sum_i((ui)^2)/n)
func VecRms(u []float64) (rms float64) {
	if Pll {
		ncpu := imin(len(u), NCPU)
		ch := make(chan float64, ncpu)
		for icpu := 0; icpu < ncpu; icpu++ {
			start, endp1 := (icpu*len(u))/ncpu, ((icpu+1)*len(u))/ncpu
			go func() {
				var mysum float64
				for i := start; i < endp1; i++ {
					mysum += u[i] * u[i]
				}
				ch <- mysum
			}()
		}
		rms = <-ch
		for icpu := 1; icpu < ncpu; icpu++ {
			rms += <-ch
		}
	} else {
		for i := 0; i < len(u); i++ {
			rms += u[i] * u[i]
		}
	}
	rms = math.Sqrt(rms / float64(len(u)))
	return
}

// VecRmsErr returns the root-mean-square of a vector u normalised by scal[i] = Atol + Rtol * |vi|
//  rms     := sqrt(sum_i((u[i]/scal[i])^2)/n)
//  scal[i] := Atol + Rtol * |v[i]|
func VecRmsErr(u []float64, Atol, Rtol float64, v []float64) (rms float64) {
	if Pll {
		ncpu := imin(len(v), NCPU)
		ch := make(chan float64, ncpu)
		for icpu := 0; icpu < ncpu; icpu++ {
			start, endp1 := (icpu*len(v))/ncpu, ((icpu+1)*len(v))/ncpu
			go func() {
				var scal, mysum float64
				for i := start; i < endp1; i++ {
					scal = Atol + Rtol*math.Abs(v[i])
					mysum += u[i] * u[i] / (scal * scal)
				}
				ch <- mysum
			}()
		}
		rms = <-ch
		for icpu := 1; icpu < ncpu; icpu++ {
			rms += <-ch
		}
	} else {
		var scal float64
		for i := 0; i < len(v); i++ {
			scal = Atol + Rtol*math.Abs(v[i])
			rms += u[i] * u[i] / (scal * scal)
		}
	}
	rms = math.Sqrt(rms / float64(len(u)))
	return
}

// VecRmsError returns the root-mean-square of a vector u normalised by scal[i] = Atol + Rtol * |vi|
//  rms     := sqrt(sum_i((|u[i]-w[i]|/scal[i])^2)/n)
//  scal[i] := Atol + Rtol * |v[i]|
func VecRmsError(u, w []float64, Atol, Rtol float64, v []float64) (rms float64) {
	if Pll {
		ncpu := imin(len(v), NCPU)
		ch := make(chan float64, ncpu)
		for icpu := 0; icpu < ncpu; icpu++ {
			start, endp1 := (icpu*len(v))/ncpu, ((icpu+1)*len(v))/ncpu
			go func() {
				var scal, mysum, e float64
				for i := start; i < endp1; i++ {
					scal = Atol + Rtol*math.Abs(v[i])
					e = math.Abs(u[i] - w[i])
					mysum += e * e / (scal * scal)
				}
				ch <- mysum
			}()
		}
		rms = <-ch
		for icpu := 1; icpu < ncpu; icpu++ {
			rms += <-ch
		}
	} else {
		var scal, e float64
		for i := 0; i < len(v); i++ {
			scal = Atol + Rtol*math.Abs(v[i])
			e = math.Abs(u[i] - w[i])
			rms += e * e / (scal * scal)
		}
	}
	rms = math.Sqrt(rms / float64(len(u)))
	return
}

// --------------------------------------------------------------------------------------------------
// matrix -------------------------------------------------------------------------------------------
// --------------------------------------------------------------------------------------------------

// MatAlloc allocates a matrix with size m,n:
//  a := 0  =>  aij = 0
func MatAlloc(m, n int) (mat [][]float64) {
	mat = make([][]float64, m)
	for k := 0; k < m; k++ {
		mat[k] = make([]float64, n)
	}
	return mat
}

// MatClone allocates and clones a matrix
//  b := a
func MatClone(a [][]float64) (b [][]float64) {
	b = make([][]float64, len(a))
	for i := 0; i < len(a); i++ {
		b[i] = make([]float64, len(a[i]))
		for j := 0; j < len(a[i]); j++ {
			b[i][j] = a[i][j]
		}
	}
	return
}

// MatFill fills a matrix with a single number s:
//  aij := s
func MatFill(a [][]float64, s float64) {
	if Pll {
		ncpu := imin(len(a), NCPU)
		ch := make(chan int, ncpu)
		for icpu := 0; icpu < ncpu; icpu++ {
			start, endp1 := (icpu*len(a))/ncpu, ((icpu+1)*len(a))/ncpu
			go func() {
				for i := start; i < endp1; i++ {
					for j := 0; j < len(a[i]); j++ {
						a[i][j] = s
					}
				}
				ch <- 1
			}()
		}
		for icpu := 0; icpu < ncpu; icpu++ {
			<-ch
		}
	} else {
		for i := 0; i < len(a); i++ {
			for j := 0; j < len(a[i]); j++ {
				a[i][j] = s
			}
		}
	}
}

// MatScale scales a matrix by a scalar s:
//  a := α * a  =>  aij := α * aij
func MatScale(a [][]float64, α float64) {
	if Pll {
		ncpu := imin(len(a), NCPU)
		ch := make(chan int, ncpu)
		for icpu := 0; icpu < ncpu; icpu++ {
			start, endp1 := (icpu*len(a))/ncpu, ((icpu+1)*len(a))/ncpu
			go func() {
				for i := start; i < endp1; i++ {
					for j := 0; j < len(a[i]); j++ {
						a[i][j] *= α
					}
				}
				ch <- 1
			}()
		}
		for icpu := 0; icpu < ncpu; icpu++ {
			<-ch
		}
	} else {
		for i := 0; i < len(a); i++ {
			for j := 0; j < len(a[i]); j++ {
				a[i][j] *= α
			}
		}
	}
}

// MatCopy copies to matrix "a", another matrix "b" (scaled):
//  a := α * b  =>  aij := α * bij
func MatCopy(a [][]float64, α float64, b [][]float64) {
	if Pll {
		ncpu := imin(len(a), NCPU)
		ch := make(chan int, ncpu)
		for icpu := 0; icpu < ncpu; icpu++ {
			start, endp1 := (icpu*len(a))/ncpu, ((icpu+1)*len(a))/ncpu
			go func() {
				for i := start; i < endp1; i++ {
					for j := 0; j < len(a[i]); j++ {
						a[i][j] = α * b[i][j]
					}
				}
				ch <- 1
			}()
		}
		for icpu := 0; icpu < ncpu; icpu++ {
			<-ch
		}
	} else {
		for i := 0; i < len(a); i++ {
			for j := 0; j < len(a[i]); j++ {
				a[i][j] = α * b[i][j]
			}
		}
	}
}

// MatSetDiag sets a diagonal matrix where the diagonal components are equal to s
func MatSetDiag(a [][]float64, s float64) {
	if Pll {
		ncpu := imin(len(a), NCPU)
		ch := make(chan int, ncpu)
		for icpu := 0; icpu < ncpu; icpu++ {
			start, endp1 := (icpu*len(a))/ncpu, ((icpu+1)*len(a))/ncpu
			go func() {
				for i := start; i < endp1; i++ {
					for j := 0; j < len(a[i]); j++ {
						a[i][j] = 0.0
						if i == j {
							a[i][j] = s
						}
					}
				}
				ch <- 1
			}()
		}
		for icpu := 0; icpu < ncpu; icpu++ {
			<-ch
		}
	} else {
		for i := 0; i < len(a); i++ {
			for j := 0; j < len(a[i]); j++ {
				a[i][j] = 0.0
				if i == j {
					a[i][j] = s
				}
			}
		}
	}
}

// MatMaxDiff returns the maximum difference between components of two matrices
func MatMaxDiff(a, b [][]float64) (maxdiff float64) {
	if Pll {
		ncpu := imin(len(a), NCPU)
		ch := make(chan float64, ncpu)
		for icpu := 0; icpu < ncpu; icpu++ {
			start, endp1 := (icpu*len(a))/ncpu, ((icpu+1)*len(a))/ncpu
			go func() {
				mymaxdiff := math.Abs(a[start][0] - b[start][0])
				for i := start; i < endp1; i++ {
					for j := 0; j < len(a[i]); j++ {
						diff := math.Abs(a[i][j] - b[i][j])
						if diff > mymaxdiff {
							mymaxdiff = diff
						}
					}
				}
				ch <- mymaxdiff
			}()
		}
		maxdiff = <-ch
		for icpu := 1; icpu < ncpu; icpu++ {
			othermaxdiff := <-ch
			if othermaxdiff > maxdiff {
				maxdiff = othermaxdiff
			}
		}
	} else {
		maxdiff = math.Abs(a[0][0] - b[0][0])
		for i := 0; i < len(a); i++ {
			for j := 0; j < len(a[i]); j++ {
				diff := math.Abs(a[i][j] - b[i][j])
				if diff > maxdiff {
					maxdiff = diff
				}
			}
		}
	}
	return
}

// MatLargest returns the largest component (abs(a_ij)) of a matrix, normalised by den
func MatLargest(a [][]float64, den float64) (largest float64) {
	if Pll {
		ncpu := imin(len(a), NCPU)
		ch := make(chan float64, ncpu)
		for icpu := 0; icpu < ncpu; icpu++ {
			start, endp1 := (icpu*len(a))/ncpu, ((icpu+1)*len(a))/ncpu
			go func() {
				mylargest := math.Abs(a[start][0])
				for i := start; i < endp1; i++ {
					for j := 0; j < len(a[i]); j++ {
						tmp := math.Abs(a[i][j])
						if tmp > mylargest {
							mylargest = tmp
						}
					}
				}
				ch <- mylargest
			}()
		}
		largest = <-ch
		for icpu := 1; icpu < ncpu; icpu++ {
			otherlargest := <-ch
			if otherlargest > largest {
				largest = otherlargest
			}
		}
	} else {
		largest = math.Abs(a[0][0])
		for i := 0; i < len(a); i++ {
			for j := 0; j < len(a[i]); j++ {
				tmp := math.Abs(a[i][j])
				if tmp > largest {
					largest = tmp
				}
			}
		}
	}
	return
}

// MatGetCol returns a column of matrix a into vector col:
//  col := a[:][j]
func MatGetCol(j int, a [][]float64) (col []float64) {
	col = make([]float64, len(a))
	if Pll {
		ncpu := imin(len(a), NCPU)
		ch := make(chan int, ncpu)
		for icpu := 0; icpu < ncpu; icpu++ {
			start, endp1 := (icpu*len(a))/ncpu, ((icpu+1)*len(a))/ncpu
			go func() {
				for i := start; i < endp1; i++ {
					col[i] = a[i][j]
				}
				ch <- 1
			}()
		}
		for icpu := 0; icpu < ncpu; icpu++ {
			<-ch
		}
	} else {
		for i := 0; i < len(a); i++ {
			col[i] = a[i][j]
		}
	}
	return
}

// MatNormF returns the Frobenious norm of a matrix:
//  ||A||_F := sqrt(sum_i sum_j aij*aij)
func MatNormF(a [][]float64) (res float64) {
	if Pll {
		ncpu := imin(len(a), NCPU)
		ch := make(chan float64, ncpu)
		for icpu := 0; icpu < ncpu; icpu++ {
			start, endp1 := (icpu*len(a))/ncpu, ((icpu+1)*len(a))/ncpu
			go func() {
				var mysum float64
				for i := start; i < endp1; i++ {
					for j := 0; j < len(a[i]); j++ {
						mysum += a[i][j] * a[i][j]
					}
				}
				ch <- mysum
			}()
		}
		for icpu := 0; icpu < ncpu; icpu++ {
			res += <-ch
		}
	} else {
		for i := 0; i < len(a); i++ {
			for j := 0; j < len(a[i]); j++ {
				res += a[i][j] * a[i][j]
			}
		}
	}
	return math.Sqrt(res)
}

// MatNormI returns the infinite norm of a matrix:
//  ||A||_∞ := max_i sum_j aij
func MatNormI(a [][]float64) (res float64) {
	if len(a) < 1 {
		return
	}
	if Pll {
		ncpu := imin(len(a), NCPU)
		ch := make(chan float64, ncpu)
		for icpu := 0; icpu < ncpu; icpu++ {
			start, endp1 := (icpu*len(a))/ncpu, ((icpu+1)*len(a))/ncpu
			go func() {
				var mymax float64
				for j := 0; j < len(a[start]); j++ {
					mymax += math.Abs(a[start][j])
				}
				var sumrow float64
				for i := start + 1; i < endp1; i++ {
					sumrow = 0.0
					for j := 0; j < len(a[i]); j++ {
						sumrow += math.Abs(a[i][j])
					}
					if sumrow > mymax {
						mymax = sumrow
					}
				}
				ch <- mymax
			}()
		}
		res = <-ch
		var othermax float64
		for icpu := 1; icpu < ncpu; icpu++ {
			othermax = <-ch
			if othermax > res {
				res = othermax
			}
		}
	} else {
		for j := 0; j < len(a[0]); j++ {
			res += math.Abs(a[0][j])
		}
		var sumrow float64
		for i := 0; i < len(a); i++ {
			sumrow = 0.0
			for j := 0; j < len(a[i]); j++ {
				sumrow += math.Abs(a[i][j])
				if sumrow > res {
					res = sumrow
				}
			}
		}
	}
	return
}

// --------------------------------------------------------------------------------------------------
// matrix-vector / matrix-matrix --------------------------------------------------------------------
// --------------------------------------------------------------------------------------------------

// MatVecMul returns the matrix-vector multiplication (scaled):
//  v := α * a * u  =>  vi = α * aij * uj
//  NOTE: not efficient implementation => use for small matrices
func MatVecMul(v []float64, α float64, a [][]float64, u []float64) {
	for i := 0; i < len(a); i++ {
		v[i] = 0.0
		for j := 0; j < len(u); j++ {
			v[i] += α * a[i][j] * u[j]
		}
	}
}

// MatVecMulAdd returns the matrix-vector multiplication with addition (scaled):
//  v += α * a * u  =>  vi += α * aij * uj
//  NOTE: not efficient implementation => use for small matrices
func MatVecMulAdd(v []float64, α float64, a [][]float64, u []float64) {
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(u); j++ {
			v[i] += α * a[i][j] * u[j]
		}
	}
}

// MatVecMulAddC (complex) returns the matrix-vector multiplication with addition (scaled):
//  v += α * a * u  =>  vi += α * aij * uj
//  NOTE: not efficient implementation => use for small matrices
func MatVecMulAddC(v []complex128, α complex128, a [][]complex128, u []complex128) {
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(u); j++ {
			v[i] += α * a[i][j] * u[j]
		}
	}
}

// MatTrVecMult returns the matrix-vector multiplication with "a" transposed (scaled):
//  v := α * transp(a) * u  =>  vi = α * aji * uj
//  NOTE: not efficient implementation => use for small matrices
func MatTrVecMul(v []float64, α float64, a [][]float64, u []float64) {
	for i := 0; i < len(a[0]); i++ {
		v[i] = 0.0
		for j := 0; j < len(u); j++ {
			v[i] += α * a[j][i] * u[j]
		}
	}
}

// MatTrVecMulAdd returns the matrix-vector multiplication with addition and "a" transposed (scaled):
//  v += α * transp(a) * u  =>  vi += α * aji * uj
//  NOTE: not efficient implementation => use for small matrices
func MatTrVecMulAdd(v []float64, α float64, a [][]float64, u []float64) {
	for i := 0; i < len(a[0]); i++ {
		for j := 0; j < len(u); j++ {
			v[i] += α * a[j][i] * u[j]
		}
	}
}

// MatVecMulCopyAdd returns the matrix-vector multiplication with copy and addition (scaled):
//  w = α*v + β*a*u  =>  wi = α*vi + β*aij*uj
//  NOTE: not efficient implementation => use for small matrices
func MatVecMulCopyAdd(w []float64, α float64, v []float64, β float64, a [][]float64, u []float64) {
	for i := 0; i < len(a); i++ {
		w[i] = α * v[i]
		for j := 0; j < len(u); j++ {
			w[i] += β * a[i][j] * u[j]
		}
	}
}

// MatMul returns the matrix multiplication (scaled):
//  c := α * a * b  =>  cij := α * aik * bkj
//  NOTE: not efficient implementation => use for small matrices
func MatMul(c [][]float64, α float64, a, b [][]float64) {
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(b[0]); j++ {
			c[i][j] = 0.0
			for k := 0; k < len(a[0]); k++ {
				c[i][j] += α * a[i][k] * b[k][j]
			}
		}
	}
}

// MatMulNew returns the matrix multiplication (scaled):
//  c := α * a * b  =>  cij := α * aik * bkj
//  NOTE: not efficient implementation => use for small matrices
func MatMulNew(c *oblas.Matrix, α float64, a, b *oblas.Matrix) {
	for i := 0; i < c.M*c.N; i++ {
		c.Data[i] = 0
	}
	err := oblas.Dgemm(false, false, a.M, b.N, a.N, α, a, a.M, b, b.M, 1, c, c.M)
	if err != nil {
		chk.Panic("%v\n", err)
	}
}

// MatMul3 returns the triple matrix multiplication:
//  d := α * a * b * c  =>  dij := α * aik * bkl * clj
//  NOTE: not efficient implementation => use for small matrices
func MatMul3(d [][]float64, α float64, a, b, c [][]float64) {
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(c[0]); j++ {
			d[i][j] = 0.0
			for k := 0; k < len(a[0]); k++ {
				for l := 0; l < len(b[0]); l++ {
					d[i][j] += α * a[i][k] * b[k][l] * c[l][j]
				}
			}
		}
	}
}

// MatTrMul3 returns the triple matrix multiplication with tranposed "a":
//  d := α * trans(a) * b * c  =>  dij := α * aki * bkl * clj
//  NOTE: not efficient implementation => use for small matrices
func MatTrMul3(d [][]float64, α float64, a, b, c [][]float64) {
	for i := 0; i < len(a[0]); i++ {
		for j := 0; j < len(c[0]); j++ {
			d[i][j] = 0.0
			for k := 0; k < len(a); k++ {
				for l := 0; l < len(b[0]); l++ {
					d[i][j] += α * a[k][i] * b[k][l] * c[l][j]
				}
			}
		}
	}
}

// MatTrMulAdd3 returns the triple matrix multiplication with addition and tranposed "a":
//  d += α * trans(a) * b * c  =>  dij += α * aki * bkl * clj
//  NOTE: not efficient implementation => use for small matrices
func MatTrMulAdd3(d [][]float64, α float64, a, b, c [][]float64) {
	for i := 0; i < len(a[0]); i++ {
		for j := 0; j < len(c[0]); j++ {
			for k := 0; k < len(a); k++ {
				for l := 0; l < len(b[0]); l++ {
					d[i][j] += α * a[k][i] * b[k][l] * c[l][j]
				}
			}
		}
	}
}

// VecOuterAdd returns the outer product between two vectors, with addition (scaled)
//  aij += α * u[i] * v[j]
//  NOTE: not efficient implementation => use for small matrices
func VecOuterAdd(a [][]float64, α float64, u, v []float64) {
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(a[0]); j++ {
			a[i][j] += α * u[i] * v[j]
		}
	}
}
