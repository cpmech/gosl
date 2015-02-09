// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fdm

import (
	"bytes"
	"fmt"
	"os/exec"

	"code.google.com/p/gosl/utl"
)

type Grid2D struct {
	Lx, Ly     float64
	Nx, Ny, N  int
	Dx, Dy     float64
	Dxx, Dyy   float64
	L, R, B, T []int
}

func (g *Grid2D) Init(lx, ly float64, nx, ny int) {
	g.Lx, g.Ly = lx, ly
	g.Nx, g.Ny, g.N = nx, ny, nx*ny
	g.Dx, g.Dy = g.Lx/float64(nx-1), g.Ly/float64(ny-1)
	g.Dxx, g.Dyy = g.Dx*g.Dx, g.Dy*g.Dy

	g.L = utl.IntRange3(0, g.N, g.Nx)
	g.R = utl.IntAddScalar(g.L, g.Nx-1)
	g.B = utl.IntRange(g.Nx)
	g.T = utl.IntAddScalar(g.B, (g.Ny-1)*g.Nx)
}

func (g *Grid2D) Draw(dirout, fnkey string, show bool) {
	// write buffer
	var b bytes.Buffer
	fmt.Fprintf(&b, "from gosl import *\n")
	fmt.Fprintf(&b, "XY = array([")
	for j := 0; j < g.Ny; j++ {
		for i := 0; i < g.Nx; i++ {
			x := float64(i) * g.Dx
			y := float64(j) * g.Dy
			fmt.Fprintf(&b, "(%g, %g),", x, y)
		}
	}
	fmt.Fprintf(&b, "],dtype=float)\n")
	fmt.Fprintf(&b, "L = %v\n", utl.IntPy(g.L))
	fmt.Fprintf(&b, "R = %v\n", utl.IntPy(g.R))
	fmt.Fprintf(&b, "B = %v\n", utl.IntPy(g.B))
	fmt.Fprintf(&b, "T = %v\n", utl.IntPy(g.T))
	fmt.Fprintf(&b, "plot(XY[:,0], XY[:,1], 'ko', clip_on=False)\n")
	fmt.Fprintf(&b, "plot(XY[L,0], XY[L,1], 'rs', ms=15, clip_on=False)\n")
	fmt.Fprintf(&b, "plot(XY[R,0], XY[R,1], 'bs', ms=15, clip_on=False)\n")
	fmt.Fprintf(&b, "plot(XY[B,0], XY[B,1], 'yo', ms=12, clip_on=False)\n")
	fmt.Fprintf(&b, "plot(XY[T,0], XY[T,1], 'go', ms=12, clip_on=False)\n")
	fmt.Fprintf(&b, "axis('equal')\n")
	fmt.Fprintf(&b, "grid()\n")
	fmt.Fprintf(&b, "show()\n")
	// save file
	utl.WriteFileD(dirout, fnkey+".py", &b)
	if show {
		_, err := exec.Command("python", dirout+"/"+fnkey+".py").Output()
		if err != nil {
			utl.Panic("Grid2D:Draw failed when calling python\n%v", err)
		}
	}
}

// fxy or z must be nil
func (g *Grid2D) Contour(dirout, fnkey string, fxy Cb_fxy, z []float64, nlevels int, show bool) {
	// write buffer
	var b bytes.Buffer
	fmt.Fprintf(&b, "from gosl import *\n")
	fmt.Fprintf(&b, "XYZ = array([")
	for j := 0; j < g.Ny; j++ {
		for i := 0; i < g.Nx; i++ {
			x := float64(i) * g.Dx
			y := float64(j) * g.Dy
			if fxy == nil {
				fmt.Fprintf(&b, "(%g, %g, %g),", x, y, z[i+j*g.Nx])
			} else {
				fmt.Fprintf(&b, "(%g, %g, %g),", x, y, fxy(x, y))
			}
		}
	}
	fmt.Fprintf(&b, "],dtype=float)\n")
	fmt.Fprintf(&b, "X = XYZ[:,0].reshape(%d,%d)\n", g.Ny, g.Nx)
	fmt.Fprintf(&b, "Y = XYZ[:,1].reshape(%d,%d)\n", g.Ny, g.Nx)
	fmt.Fprintf(&b, "Z = XYZ[:,2].reshape(%d,%d)\n", g.Ny, g.Nx)
	fmt.Fprintf(&b, "Contour(X,Y,Z, nlevels=%d)\n", nlevels)
	fmt.Fprintf(&b, "axis('equal')\n")
	fmt.Fprintf(&b, "show()\n")
	// save file
	utl.WriteFileD(dirout, fnkey+".py", &b)
	if show {
		_, err := exec.Command("python", dirout+"/"+fnkey+".py").Output()
		if err != nil {
			utl.Panic("Grid2D:Draw failed when calling python\n%v", err)
		}
	}
}

// callbacks
type Cb_fxy func(x, y float64) (z float64) // z = f(x,y)

// auxiliary function
func prin(msg string, prm ...interface{}) { fmt.Printf(msg, prm...) }
