# Gosl - Go scientific library

[![Go Reference](https://pkg.go.dev/badge/github.com/cpmech/gosl.svg)](https://pkg.go.dev/github.com/cpmech/gosl)
[![Go Report Card](https://goreportcard.com/badge/github.com/cpmech/gosl)](https://goreportcard.com/report/github.com/cpmech/gosl)
[![Awesome](https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg)](https://github.com/avelino/awesome-go)

Gosl is a set of tools for developing scientific simulations using the Go language. We mainly consider the development of numerical methods and solvers for differential equations but also present some functions for fast Fourier transforms, the generation of random numbers, probability distributions, and computational geometry.

This library contains essential functions for linear algebra computations (operations between all combinations of vectors and matrices, eigenvalues and eigenvectors, linear solvers) and the development of numerical methods (e.g. numerical quadrature).

We link Gosl with existent libraries written in C and Fortran, such as OpenBLAS, LAPACK, UMFPACK, MUMPS, QUADPACK and FFTW3. These existing libraries have been fundamental for the development of high-performant simulations over many years. We believe that it is nearly impossible to rewrite these libraries in native Go and at the same time achieve the same speed delivered by them. Just for reference, a naive implementation of matrix-matrix multiplication in Go is more than 100 times slower than OpenBLAS.

## Installation

Because of the other libraries, the easiest way to work with Gosl is via Docker. Having Docker and VS Code installed, you can start developing powerful numerical simulations using Gosl in a matter of minutes. Furthermore, it works on Windows, Linux, and macOS out of the box.

### Containerized

1. Install Docker
2. Install Visual Studio Code
3. Install the Remote Development extension for VS Code
4. Clone https://github.com/cpmech/hello-gosl
5. Create your application within a container (see gif below)

Done. And your system will "remain clean."

![](zdocs/vscode-open-in-container.gif)

### Debian/Ubuntu GNU Linux

First, install Go as explained in https://golang.org/doc/install

Second, install some libraries:

```
sudo apt-get install \
  gcc \
  gfortran \
  libfftw3-dev \
  liblapacke-dev \
  libmetis-dev \
  libmumps-seq-dev \
  libopenblas-dev \
  libsuitesparse-dev
```

Finally, download and compile Gosl:

```
git clone https://github.com/cpmech/gosl.git
cd gosl
./all.bash
```

Done. Installation completed.

## Documentation

Gosl includes the following _essential_ packages:

- [chk](https://github.com/cpmech/gosl/tree/master/chk). To check numerical results and for unit testing
- [io](https://github.com/cpmech/gosl/tree/master/io). Input/output including printing to the terminal and handling files
- [utl](https://github.com/cpmech/gosl/tree/master/utl). To generate series (e.g. linspace) and other functions as in pylab/matlab/octave
- [la](https://github.com/cpmech/gosl/tree/master/la). Linear Algebra: vector, matrix, efficient sparse solvers, eigenvalues, decompositions

Gosl includes the following _main_ packages:

- [fun](https://github.com/cpmech/gosl/tree/master/fun). Special functions, DFT, FFT, Bessel, elliptical integrals, orthogonal polynomials, interpolators
- [gm](https://github.com/cpmech/gosl/tree/master/gm). Geometry algorithms and structures
- [hb](https://github.com/cpmech/gosl/tree/master/hb). Pseudo hierarchical binary (hb) data file format
- [num](https://github.com/cpmech/gosl/tree/master/num). Fundamental numerical methods such as root solvers, non-linear solvers, numerical derivatives and quadrature
- [ode](https://github.com/cpmech/gosl/tree/master/ode). Solvers for ordinary differential equations
- [opt](https://github.com/cpmech/gosl/tree/master/opt). Numerical optimization: Interior Point, Conjugate Gradients, Powell, Grad Descent
- [pde](https://github.com/cpmech/gosl/tree/master/pde). Solvers for partial differential equations (FDM, Spectral, FEM)
- [rnd](https://github.com/cpmech/gosl/tree/master/rnd). Random numbers and probability distributions

(see each subdirectory for more information)

For the sake of maintenance (see next section), we have removed the previous `mpi` sub-package. However, we recommend the external library [gompi](https://github.com/sbromberger/gompi) if you plan to use MPI.

## Previous version

The previous version, including more packages, is [available here ](https://github.com/cpmech/gosl/tree/stable-1.1.3) and can be used with the Docker image 1.1.3 as in this [hello gosl example](https://github.com/cpmech/hello-gosl-old-1.1.3).

These other packages, such as machine learning, plotting, etc., have been removed because they do not depend on CGO and may be developed independently. We can now maintain the core of Gosl more efficiently, which has a focus on the foundation for other scientific code.
