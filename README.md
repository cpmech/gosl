# Gosl — a Go library for scientific computing

[![Go Reference](https://pkg.go.dev/badge/github.com/cpmech/gosl.svg)](https://pkg.go.dev/github.com/cpmech/gosl)
[![Go Report Card](https://goreportcard.com/badge/github.com/cpmech/gosl)](https://goreportcard.com/report/github.com/cpmech/gosl)
[![Awesome](https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg)](https://github.com/avelino/awesome-go)

Gosl is a set of tools for developing scientific simulations in Go. Its focus is numerical methods and solvers for differential equations. It also provides fast Fourier transforms, random-number generation, probability distributions, and computational geometry.

The library covers the linear algebra that numerical work requires — operations between all combinations of vectors and matrices, eigenvalues and eigenvectors, and linear solvers — together with the building blocks of numerical methods, such as numerical quadrature.

Gosl links against C and Fortran libraries: OpenBLAS, LAPACK, UMFPACK, MUMPS, QUADPACK, and FFTW3. These libraries have underpinned high-performance simulation for decades, and a rewrite in native Go is unlikely to match their speed: in our benchmark, a naive Go matrix-matrix multiplication runs more than 100 times slower than OpenBLAS.

## Status

Gosl is mature and **in maintenance mode**. It is stable and is used in published work, but new
development has moved to [Russell](https://github.com/cpmech/russell), its successor — if you are
starting a new project, look there first. Reports about Gosl's existing functionality are welcome;
new features are unlikely.

## Installation

Because Gosl links against these libraries, Docker is the easiest way to work with it: the libraries
are installed in the image, nothing is added to your system, and no Go toolchain is needed on the
host. Docker is the only requirement, and it works the same on Linux, macOS and Windows.

### Containerized (recommended)

#### Build and test Gosl itself

```bash
git clone https://github.com/cpmech/gosl.git
cd gosl
docker build -f docker/Dockerfile -t gosl:ci .    # Ubuntu 24.04, a pinned Go, the libraries
./gosl all                                        # build, install and test every package
```

`./gosl` wraps `docker run`, so the Go toolchain never reaches your host:

| Command | What it does |
|---|---|
| `./gosl all` | build, install and test every package |
| `./gosl test`, `./gosl race` | the tests, with and without the race detector |
| `./gosl vet`, `./gosl fmt` | `go vet ./...`; `gofmt -s -l .` |
| `./gosl cover`, `./gosl bench` | coverage profile into `.docker-cache/`; benchmarks |
| `./gosl shell`, `./gosl go <args>` | a shell in the container; any other `go` command |

It runs as your own user, so everything it produces belongs to you, and it caches into the
gitignored `.docker-cache/` — `./gosl all` takes 4.6 s once warm and 14.2 s from an empty cache.
Build and test run with no network access, which costs nothing because Gosl has no Go dependencies.
`docker run --rm gosl:ci` prints the library versions the image ships.

#### Use Gosl in your own project

Any Ubuntu 24.04 image with the same libraries works; the image above is one, and it can be your
build environment:

```bash
docker run --rm -it --network=bridge --user "$(id -u):$(id -g)" \
  -e HOME=/tmp -e GOCACHE=/tmp/gocache -e GOMODCACHE=/tmp/gomod -e GOPATH=/tmp/gopath \
  -v "$PWD":/w -w /w gosl:ci bash
```

Inside the container, start a module and pull Gosl:

```bash
go mod init example.com/demo
go get github.com/cpmech/gosl
```

Then build your program as usual. As a complete example, solving a small symmetric linear system:

```go
package main

import (
	"fmt"

	"github.com/cpmech/gosl/la"
)

func main() {
	a := la.NewMatrix(2, 2)
	a.Set(0, 0, 2)
	a.Set(1, 1, 3)

	b := la.Vector{4, 9}
	x := la.Vector{0, 0}
	la.SolveRealLinSysSPD(x, a, b)

	fmt.Printf("x = [%.6g %.6g]\n", x[0], x[1])
}
```

```bash
$ go run .
x = [2 3]
```

Two caveats worth knowing:

- The build and module caches live in `/tmp` inside the container, so they are thrown away with it.
  To keep them between runs, add
  `-v "$HOME/.cache/gosl-docker":/cache -e GOCACHE=/cache/build -e GOMODCACHE=/cache/mod`.
- A binary built this way links against the container's shared libraries. Run it in a container, as
  above, or install the same libraries on the host.

### Native install (Debian/Ubuntu Linux)

This is the path the CI workflow uses, on Ubuntu with Go 1.20 or newer.

First, install Go as explained in https://go.dev/doc/install.

Second, install the C and Fortran libraries that Gosl links against:

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

Finally, clone, build, and test Gosl:

```
git clone https://github.com/cpmech/gosl.git
cd gosl
./all.bash
```

`./all.bash` builds every package and runs its tests; it prints `SUCCESS!` when it is done.

## Documentation

Gosl includes the following _essential_ packages:

- [chk](https://github.com/cpmech/gosl/tree/main/chk) — checks on numerical results, and helpers for unit testing.
- [io](https://github.com/cpmech/gosl/tree/main/io) — input/output, including printing to the terminal and handling files.
- [utl](https://github.com/cpmech/gosl/tree/main/utl) — series generation (e.g. linspace) and other functions as in pylab, MATLAB, and Octave.
- [la](https://github.com/cpmech/gosl/tree/main/la) — linear algebra: vectors, matrices, efficient sparse solvers, eigenvalues, and decompositions.

Gosl includes the following _main_ packages:

- [fun](https://github.com/cpmech/gosl/tree/main/fun) — special functions, DFT, FFT, Bessel functions, elliptic integrals, orthogonal polynomials, and interpolators.
- [gm](https://github.com/cpmech/gosl/tree/main/gm) — geometry algorithms and structures.
- [hb](https://github.com/cpmech/gosl/tree/main/hb) — the pseudo-hierarchical binary (hb) data file format.
- [num](https://github.com/cpmech/gosl/tree/main/num) — fundamental numerical methods, such as root solvers, nonlinear solvers, numerical derivatives, and quadrature.
- [ode](https://github.com/cpmech/gosl/tree/main/ode) — solvers for ordinary differential equations.
- [opt](https://github.com/cpmech/gosl/tree/main/opt) — numerical optimization: interior point, conjugate gradients, Powell, and gradient descent.
- [pde](https://github.com/cpmech/gosl/tree/main/pde) — solvers for partial differential equations (FDM, spectral, FEM).
- [rnd](https://github.com/cpmech/gosl/tree/main/rnd) — random numbers and probability distributions.

See each subdirectory for more information.

The previous `mpi` sub-package has been removed for maintenance reasons. If you need MPI, we recommend the external library [gompi](https://github.com/sbromberger/gompi).

## Licence

Gosl is released under the BSD-3-Clause licence — see [LICENSE](LICENSE).

One component is not covered by it: `gm/tri` is a wrapper around
[Triangle](https://www.cs.cmu.edu/~quake/triangle.html) by Jonathan Richard Shewchuk. Triangle is
free for private, research and institutional use, and may be redistributed provided its copyright
notices are kept and no compensation is received; distributing it as part of a commercial system
requires a direct arrangement with the author. The full terms are in
[gm/tri/triangle_README.txt](gm/tri/triangle_README.txt).

The vendored source is not upstream Triangle 1.6 verbatim: it was adapted for Gosl when it was
imported, and its file keeps Shewchuk's copyright and licence notice unchanged.

If you publish results obtained with Triangle, please acknowledge it and cite: Jonathan Richard
Shewchuk, "Triangle: Engineering a 2D Quality Mesh Generator and Delaunay Triangulator", in Applied
Computational Geometry: Towards Geometric Engineering, LNCS 1148, pages 203--222, Springer-Verlag,
May 1996.
