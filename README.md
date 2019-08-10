# Gosl &ndash; Go scientific library

[![GoDoc](https://godoc.org/github.com/cpmech/gosl?status.svg)](https://godoc.org/github.com/cpmech/gosl)
[![Build Status](https://travis-ci.org/cpmech/gosl.svg?branch=master)](https://travis-ci.org/cpmech/gosl)
[![Go Report Card](https://goreportcard.com/badge/github.com/cpmech/gosl)](https://goreportcard.com/report/github.com/cpmech/gosl)
[![Awesome](https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg)](https://github.com/avelino/awesome-go)

Gosl is a [Go](https://golang.org) library to develop Artificial Intelligence and High-Performance
Scientific Computations.

<div id="container">
<p><a href="examples/figs/gosl-collage1.jpg"><img src="examples/figs/gosl-collage1-sml.png"></a></p>
</div>

The library tries to be as general and as _easy_ as possible. Gosl considers the use of both Go
concurrency routines and parallel computing using the message passing interface (MPI). Gosl has
several modules (sub-packages) for a variety of tasks in scientific computing, image analysis, and
data post-processing.

Gosl includes high-performant linear algebra functions (wrapping MKL, OpenBLAS, LAPACK, SuiteSparse,
UMFPACK...), fast Fourier transform algorithms (wrapping FFTW), numerical integration (wrapping
QUADPACK), functions and structures for geometry calculations (e.g. 3D transfinite interpolation,
grid search, octree...), random numbers generation (SFMT and DSFMT) and probability distributions,
optimisation and graph algorithms, plotting and visualisation using the VTK, and much more.

Gosl has also solvers to (stiff or not) ordinary differential equations and several tools for 2D/3D
mesh generation to assist on the development of solvers for partial differential equations.

We now give focus to Machine Learning (see `ml` package) and Big Data (see `h5` package). We are
also planning wrappers to powerful tools such as CNTK, TensorFlow, and Hadoop. We have a wrapper to
OpenCV in the works as well.

## Resources

1. [Examples](examples/README.md) and [benchmarks](examples/benchmark/README.md)
2. [White papers](https://github.com/cpmech/gosl/tree/master/doc)
3. [Documentation](https://godoc.org/github.com/cpmech/gosl)
4. [Contributing and TODO](https://github.com/cpmech/gosl/blob/master/CONTRIBUTING.md)

## Installation

Gosl works on Windows, macOS, and Linux (Debian, Ubuntu).

For installation instructions, click on:
<a href="https://github.com/cpmech/gosl/blob/master/doc/InstallationOnWindows.md"><img src="doc/icon-windows.png" alt="Installation on Windows" align="middle"></a>
<a href="https://github.com/cpmech/gosl/blob/master/doc/InstallationOnMacOS.md"><img src="doc/icon-macos.png" alt="Installation on macOS" align="middle"></a>
<a href="https://github.com/cpmech/gosl/blob/master/doc/InstallationOnUbuntu.md"><img src="doc/icon-ubuntu.png" alt="Installation on Linux/Debian/Ubuntu" align="middle"></a>

Alternatively, to use Docker, click on:
<a href="https://hub.docker.com/r/gosl/gosl"><img src="doc/icon-docker.png" alt="Docker image" align="middle"></a>

The Docker image contains Ubuntu and all the dependencies, including MPI and VTK.

```bash
docker pull gosl/gosl
```

## Usage

Type:

```bash
cd /tmp && mkdir -p try-gosl && cd try-gosl
go mod init try-gosl
```

Create the file `try-gosl/main.go`:

```go
package main

import (
	"github.com/cpmech/gosl/fun"
	"github.com/cpmech/gosl/io"
)

func main() {
	i0 := fun.ModBesselI0(-10)
	io.Pf("Modified Bessel I0(-10) = %g\n", i0)
}
```

Run:

```bash
go run main.go
```

Output:

```
Modified Bessel I0(-10) = 2815.716628466253
```

## Sub-packages

Gosl includes the following _sub-packages_:

1.  [chk](https://github.com/cpmech/gosl/tree/master/chk) &ndash; Check code and unit test tools
2.  [io](https://github.com/cpmech/gosl/tree/master/io) &ndash; Input/output, read/write files, and print commands
3.  [io/h5](https://github.com/cpmech/gosl/tree/master/io/h5) &ndash; Read/write HDF5 (Big Data) files
4.  [utl](https://github.com/cpmech/gosl/tree/master/utl) &ndash; Utilities. Lists. Dictionaries. Simple Numerics
5.  [utl/al](https://github.com/cpmech/gosl/tree/master/utl/al) &ndash; Utilities. (Naive) Implementation of Classic algorithms and structures (e.g. Linked Lists)
6.  [plt](https://github.com/cpmech/gosl/tree/master/plt) &ndash; Plotting and drawing (png and eps)
7.  [mpi](https://github.com/cpmech/gosl/tree/master/mpi) &ndash; Message Passing Interface for parallel computing
8.  [la](https://github.com/cpmech/gosl/tree/master/la) &ndash; Linear Algebra: vector, matrix, efficient sparse solvers, eigenvalues, decompositions, etc.
9.  [la/mkl](https://github.com/cpmech/gosl/tree/master/la/mkl) &ndash; Lower level linear algebra using Intel MKL
10. [la/oblas](https://github.com/cpmech/gosl/tree/master/la/oblas) &ndash; Lower level linear algebra using OpenBLAS
11. [num/qpck](https://github.com/cpmech/gosl/tree/master/num/qpck) &ndash; Go wrapper to QUADPACK for numerical integration
12. [num](https://github.com/cpmech/gosl/tree/master/num) &ndash; Fundamental numerical methods such as root solvers, non-linear solvers, numerical derivatives and quadrature
13. [fun](https://github.com/cpmech/gosl/tree/master/fun) &ndash; Special functions, DFT, FFT, Bessel, elliptical integrals, orthogonal polynomials, interpolators
14. [fun/dbf](https://github.com/cpmech/gosl/tree/master/fun/dbf) &ndash; Database of functions of a scalar and a vector like f(t,{x}) (e.g. time-space)
15. [fun/fftw](https://github.com/cpmech/gosl/tree/master/fun/fftw) &ndash; Go wrapper to FFTW for fast Fourier Transforms
16. [gm](https://github.com/cpmech/gosl/tree/master/gm) &ndash; Geometry algorithms and structures
17. [gm/msh](https://github.com/cpmech/gosl/tree/master/gm/msh) &ndash; Mesh structures and interpolation functions for FEA, including quadrature over polyhedra
18. [gm/tri](https://github.com/cpmech/gosl/tree/master/gm/tri) &ndash; Mesh generation: triangles and Delaunay triangulation (wrapping Triangle)
19. [gm/rw](https://github.com/cpmech/gosl/tree/master/gm/rw) &ndash; Mesh generation: read/write routines
20. [graph](https://github.com/cpmech/gosl/tree/master/graph) &ndash; Graph theory structures and algorithms
21. [opt](https://github.com/cpmech/gosl/tree/master/opt) &ndash; Numerical optimization: Interior Point, Conjugate Gradients, Powell, Grad Descent, more
22. [rnd](https://github.com/cpmech/gosl/tree/master/rnd) &ndash; Random numbers and probability distributions
23. [rnd/dsfmt](https://github.com/cpmech/gosl/tree/master/rnd/dsfmt) &ndash; Go wrapper to dSIMD-oriented Fast Mersenne Twister
24. [rnd/sfmt](https://github.com/cpmech/gosl/tree/master/rnd/sfmt) &ndash; Go wrapper to SIMD-oriented Fast Mersenne Twister
25. [vtk](https://github.com/cpmech/gosl/tree/master/vtk) &ndash; 3D Visualisation with the VTK tool kit
26. [ode](https://github.com/cpmech/gosl/tree/master/ode) &ndash; Solvers for ordinary differential equations
27. [ml](https://github.com/cpmech/gosl/tree/master/ml) &ndash; Machine learning algorithms
28. [ml/imgd](https://github.com/cpmech/gosl/tree/master/ml/imgd) &ndash; Machine learning. Auxiliary functions for handling images
29. [pde](https://github.com/cpmech/gosl/tree/master/pde) &ndash; Solvers for partial differential equations (FDM, Spectral, FEM)
30. [tsr](https://github.com/cpmech/gosl/tree/master/tsr) &ndash; Tensors, continuum mechanics, and tensor algebra (e.g. eigendyads)

We are currently working on the following additional packages:

<ol start="31">
<li>img - Image and machine learning algorithms for images</li>
<li>img/ocv - Wrapper to OpenCV</li>
</ol>

## About the filenames

1. `t_something_test.go` is a **unit test**. We have several of them! Some usage
   information can be learned from these files.
2. `t_something_main.go` is a test with a main function to be run with `go run ...` or `mpirun -np ? go run ...` (replace ? with the number of cpus).
3. `t_b_something_test.go` is a benchmark test. Run benchmarks with `go test -run=XXX -bench=.`

## Design strategies

Here, we call _structure_ any _user-defined type_. These are simply Go `types` defined as `struct`.
One may think of these _structures_ as _classes_. Gosl has several global functions as well and
tries to avoid complicated constructions.

An allocated structure (instance) is called an **object** and functions attached to this object are
called **methods**. In Gosl, the variable holding the pointer to an object is always named **o**
(lower case "o"). This variable is similar to the `self` or `this` keywords in other languages
(Python, C++, respectively).

Functions that allocate a pointer to a structure are prefixed with `New`; for instance:
`NewIsoSurf`. Some structures require an explicit call to another function to release allocated
memory. Be aware of this requirement! In this case, the function is named `Free` and appears in a
few sub-packages that use CGO. Also, some objects may need to be initialized before use. In this
case, functions named `Init` have to be called.

The directories corresponding to each package have a README.md file that should help with
understanding the library. Also, there are links to `godoc.org` where all functions, structures, and
variables are well explained.

## Test coverage

We aim for a 100% test coverage! Despite trying our best to accomplish this goal, full coverage is
difficult, in particular with (sub)packages that involve `Panic` or `figure generation`.
Nonetheless, critical algorithms are completely tested.

We use the following `bash` macro frequently to check our test coverage:

```bash
gocov() {
    go test -coverprofile=/tmp/cv.out
    go tool cover -html=/tmp/cv.out
}
```

### Some results from cover.run

- [chk](https://cover.run/go/github.com/cpmech/gosl/chk) test coverage
- [io](https://cover.run/go/github.com/cpmech/gosl/io) test coverage
- [utl](https://cover.run/go/github.com/cpmech/gosl/utl) test coverage
- [utl/al](https://cover.run/go/github.com/cpmech/gosl/utl/al) test coverage
- [fun/dbf](https://cover.run/go/github.com/cpmech/gosl/fun/dbf) test coverage
- [gm/tri](https://cover.run/go/github.com/cpmech/gosl/gm/tri) test coverage
- [gm/rw](https://cover.run/go/github.com/cpmech/gosl/gm/rw) test coverage
- [rnd](https://cover.run/go/github.com/cpmech/gosl/rnd) test coverage
- [rnd/dsfmt](https://cover.run/go/github.com/cpmech/gosl/rnd/dsfmt) test coverage
- [rnd/sfmt](https://cover.run/go/github.com/cpmech/gosl/rnd/sfmt) test coverage

## Bibliography

The following works take advantage of Gosl:

1. Pedroso DM, Bonyadi MR, Gallagher M (2017) Parallel evolutionary algorithm for single and multi-objective optimisation: differential evolution and constraints handling, Applied Soft Computing http://dx.doi.org/10.1016/j.asoc.2017.09.006 [paper available here](https://github.com/cpmech/goga/blob/master/doc/goga.pdf)
2. Pedroso DM (2017) FORM reliability analysis using a parallel evolutionary algorithm, Structural Safety 65:84-99 http://dx.doi.org/10.1016/j.strusafe.2017.01.001
3. Pedroso DM, Zhang YP, Ehlers W (2017) Solution of liquid-gas-solid coupled equations for porous media considering dynamics and hysteretic retention behaviour, Journal of Engineering Mechanics 04017021 http://dx.doi.org/10.1061/(ASCE)EM.1943-7889.0001208
4. Pedroso DM (2015) A solution to transient seepage in unsaturated porous media. Computer Methods in Applied Mechanics and Engineering, 285:791-816 http://dx.doi.org/10.1016/j.cma.2014.12.009
5. Pedroso DM (2015) A consistent u-p formulation for porous media with hysteresis. Int. Journal for Numerical Methods in Engineering, 101(8):606-634 http://dx.doi.org/10.1002/nme.4808

## Authors and license

See the AUTHORS file.

Unless otherwise noted, the Gosl source files are distributed under the BSD-style license found in the LICENSE file.
