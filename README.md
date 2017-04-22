# Gosl &ndash; Go scientific library

Gosl is a computing library written in go language (golang) to help with the development of software
for scientific research. The library tries to be as general as possible. The use of concurrency for
multi-threaded applications and message passing for parallel computations are considered. Functions
and structures for geometry calculations, random numbers generation and probability distributions,
optimisation algorithms, and plotting and visualisation are implemented as well.



## Content

Gosl comprises several _subpackages_ as listed below:

1.  [chk](https://github.com/cpmech/gosl/tree/master/chk)       &ndash; Check code and unit test tools
2.  [io](https://github.com/cpmech/gosl/tree/master/io)         &ndash; Input/output, read/write files, and print commands
3.  [utl](https://github.com/cpmech/gosl/tree/master/utl)       &ndash; Utilities. Lists. Dictionaries. Simple Numerics
4.  [plt](https://github.com/cpmech/gosl/tree/master/plt)       &ndash; Plotting and drawing (png and eps)
5.  [mpi](https://github.com/cpmech/gosl/tree/master/mpi)       &ndash; Message Passing Interface for parallel computing
6.  [la](https://github.com/cpmech/gosl/tree/master/la)         &ndash; Linear Algebra and efficient sparse solvers
7.  [fdm](https://github.com/cpmech/gosl/tree/master/fdm)       &ndash; Simple finite differences method
8.  [num](https://github.com/cpmech/gosl/tree/master/num)       &ndash; Fundamental numerical methods
9.  [fun](https://github.com/cpmech/gosl/tree/master/fun)       &ndash; Scalar functions of one scalar and one vector
10. [gm](https://github.com/cpmech/gosl/tree/master/gm)         &ndash; Geometry algorithms and structures
11. [gm/msh](https://github.com/cpmech/gosl/tree/master/gm/msh) &ndash; Mesh generation and Delaunay triangulation
12. [gm/tri](https://github.com/cpmech/gosl/tree/master/gm/tri) &ndash; Mesh generation: triangles
13. [gm/rw](https://github.com/cpmech/gosl/tree/master/gm/rw)   &ndash; Mesh generation: read/write
14. [graph](https://github.com/cpmech/gosl/tree/master/graph)   &ndash; Graph theory structures and algorithms
15. [ode](https://github.com/cpmech/gosl/tree/master/ode)       &ndash; Ordinary differential equations
16. [opt](https://github.com/cpmech/gosl/tree/master/opt)       &ndash; Optimisation problem solvers
17. [rnd](https://github.com/cpmech/gosl/tree/master/rnd)       &ndash; Random numbers and probability distributions
18. [tsr](https://github.com/cpmech/gosl/tree/master/tsr)       &ndash; Tensor algebra and definitions for continuum mechanics
19. [vtk](https://github.com/cpmech/gosl/tree/master/vtk)       &ndash; 3D Visualisation with the VTK tool kit

Check the **[developer's documentation](http://rawgit.com/cpmech/gosl/master/doc/index.html)** to
see what's available and how to call functions and methods.


## Examples

[Check out examples here](https://github.com/cpmech/gosl/blob/master/examples/README.md)


<div id="container">
<p><img src="examples/figs/rnd_normalDistribution.png" width="400"></p>
Normally distributed pseudo-random numbers. Using sub-package `rnd`
</div>



## Installation

1 To install on Windows, [see instructions for Windows here](https://github.com/cpmech/gosl/blob/master/doc/InstallationOnWindows.md)

2 To install on macOS, [see instructions for macOS](https://github.com/cpmech/gosl/blob/master/doc/InstallationOnMacOS.md)

3 To install on Debian/Ubuntu/Linux, type the following commands:

3.1. Install dependencies:
```
sudo apt-get install libopenmpi-dev libhwloc-dev libsuitesparse-dev libmumps-dev 
sudo apt-get install gfortran libvtk6-dev python-scipy python-matplotlib dvipng
```

3.2. Clone and install Gosl
```
mkdir -p ${GOPATH%:*}/src/github.com/cpmech
cd ${GOPATH%:*}/src/github.com/cpmech
git clone https://github.com/cpmech/gosl.git
cd gosl
./all.bash
```

3.2 Optional: OpenBLAS
```
mkdir -p $HOME/xpkg && cd $HOME/xpkg
git clone https://github.com/xianyi/OpenBLAS.git
cd OpenBLAS
make
sudo make PREFIX=/usr/local NO_SHARED=true install
```

**Note**: Make sure to use the `NO_SHARED` flag (as above) to avoid installing the shared libraries.
Alternatively, set the `/usr/local/lib` directory as a searchable LD\_LIBRARY\_path. Otherwise, the
following error may happen:
```
[...] libopenblas.so.0: cannot open shared object file: [...]
```
**Optional**: Add the following line to .bashrc or .bash\_aliases:
`export LD_LIBRARY_PATH=/lib:/usr/lib:/usr/local/lib`
or set `/etc/ld.so.conf` file as appropriate.



### Python dependencies

At the moment, the `plt` subpackage requires `python-scipy` and `python-matplotlib` (installed as
above).



## Documentation

Here, we call _user-defined_ types as _structures_. These are simply Go `types` defined as `struct`.
Some may think of these _structures_ as _classes_. Gosl has several global functions as well and
tries to avoid complicated constructions.

An allocated structure is called here an **object** and functions attached to this object are called
**methods**. The variable holding the pointer to an object is always named **o** in Gosl (e.g.
like `self` or `this`).

Some objects need to be initialised before usage. In this case, functions named `Init` have to be
called (e.g. like `constructors`). Some structures require an explicit call to a function to release
allocated memory. This function is named `Free`. Functions that allocate a pointer to a structure
are prefixed with `New`; for instance: `NewIsoSurf`.

The directories corresponding to each package has a README.md file that should help with
understanding the library, including functions and structures.

Again, Gosl has several functions and _structures_. Check the **[developer's
documentation](http://rawgit.com/cpmech/gosl/master/doc/index.html)** to see what's available and
how to call functions and methods.






## Authors and license

See the AUTHORS file.

Unless otherwise noted, the Gosl source files are distributed under the BSD-style license found in the LICENSE file.
