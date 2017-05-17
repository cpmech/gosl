# Gosl &ndash; Go scientific library

Gosl is a computing library written in go language (golang) to help with the development of software
for scientific research. The library tries to be as general as possible. The use of concurrency for
multi-threaded applications and message passing for parallel computations are considered. Functions
and structures for geometry calculations, random numbers generation and probability distributions,
optimisation algorithms, and plotting and visualisation are implemented as well.

<div id="container">
<p><img src="examples/figs/gosl-collage1-sml.png"></p>
</div>


## Content

Gosl comprises several _subpackages_ as listed below:

1.  [chk](https://github.com/cpmech/gosl/tree/master/chk)           &ndash; Check code and unit test tools
2.  [io](https://github.com/cpmech/gosl/tree/master/io)             &ndash; Input/output, read/write files, and print commands
3.  [utl](https://github.com/cpmech/gosl/tree/master/utl)           &ndash; Utilities. Lists. Dictionaries. Simple Numerics
4.  [plt](https://github.com/cpmech/gosl/tree/master/plt)           &ndash; Plotting and drawing (png and eps)
5.  [mpi](https://github.com/cpmech/gosl/tree/master/mpi)           &ndash; Message Passing Interface for parallel computing
6.  [la](https://github.com/cpmech/gosl/tree/master/la)             &ndash; Linear Algebra and efficient sparse solvers
7.  [la/mkl](https://github.com/cpmech/gosl/tree/master/la/mkl)     &ndash; Lower level linear algebra using Intel MKL
8.  [la/oblas](https://github.com/cpmech/gosl/tree/master/la/oblas) &ndash; Lower level linear algebra using OpenBLAS
9.  [fdm](https://github.com/cpmech/gosl/tree/master/fdm)           &ndash; Simple finite differences method
10. [num](https://github.com/cpmech/gosl/tree/master/num)           &ndash; Fundamental numerical methods
11. [fun](https://github.com/cpmech/gosl/tree/master/fun)           &ndash; Scalar functions of one scalar and one vector
12. [gm](https://github.com/cpmech/gosl/tree/master/gm)             &ndash; Geometry algorithms and structures
13. [gm/msh](https://github.com/cpmech/gosl/tree/master/gm/msh)     &ndash; Mesh structures and interpolation functions for FEA
14. [gm/tri](https://github.com/cpmech/gosl/tree/master/gm/tri)     &ndash; Mesh generation: triangles and Delaunay triangulation
15. [gm/rw](https://github.com/cpmech/gosl/tree/master/gm/rw)       &ndash; Mesh generation: read/write routines
16. [graph](https://github.com/cpmech/gosl/tree/master/graph)       &ndash; Graph theory structures and algorithms
17. [ode](https://github.com/cpmech/gosl/tree/master/ode)           &ndash; Ordinary differential equations
18. [opt](https://github.com/cpmech/gosl/tree/master/opt)           &ndash; Solvers for optimisation problems
19. [rnd](https://github.com/cpmech/gosl/tree/master/rnd)           &ndash; Random numbers and probability distributions
20. [tsr](https://github.com/cpmech/gosl/tree/master/tsr)           &ndash; Tensor algebra and definitions for continuum mechanics
21. [vtk](https://github.com/cpmech/gosl/tree/master/vtk)           &ndash; 3D Visualisation with the VTK tool kit
22. [img](https://github.com/cpmech/gosl/tree/master/img)           &ndash; Image and machine learning algorithms for images
23. [img/ocv](https://github.com/cpmech/gosl/tree/master/img/ocv)   &ndash; Go wrapper to OpenCV

Check the **[developer's documentation](http://rawgit.com/cpmech/gosl/master/doc/index.html)** to
see what's available and how to call functions and methods.


## Examples

[Check out examples here](https://github.com/cpmech/gosl/blob/master/examples/README.md)

<div id="container">
<p><img src="examples/figs/gm_nurbs03.png" width="400"></p>
Construction of NURBS surfaces with the gm package.
</div>

.  

<div id="container">
<p><img src="examples/figs/rnd_normalDistribution.png" width="400"></p>
Normally distributed pseudo-random numbers. Using sub-package rnd
</div>



## 1 Installation on Windows

To install on Windows, [see instructions for Windows here](https://github.com/cpmech/gosl/blob/master/doc/InstallationOnWindows.md)

## 2 Installation on macOS

To install on macOS, [see instructions for macOS](https://github.com/cpmech/gosl/blob/master/doc/InstallationOnMacOS.md)

## 3 Installation on Linux (Debian/Ubuntu)

To install on Debian/Ubuntu/Linux, type the following commands:

### 3.1. Install dependencies:
```
sudo apt-get install libopenmpi-dev libhwloc-dev libsuitesparse-dev libmumps-dev 
sudo apt-get install gfortran libvtk6-dev python-scipy python-matplotlib dvipng
```



### 3.2. Clone and install Gosl
```
mkdir -p ${GOPATH%:*}/src/github.com/cpmech
cd ${GOPATH%:*}/src/github.com/cpmech
git clone https://github.com/cpmech/gosl.git
cd gosl
./all.bash
```



### 3.3 Set dynamic library flags if installing the next tools

To set LD\_LIBRARY\_PATH, add the following line to `.bashrc` or `.bash_aliases`:
```
export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/usr/local/lib
```
Alternatively, change `/etc/ld.so.conf` file as appropriate.



### 3.4 Optional: Install OpenBLAS
```
mkdir -p $HOME/xpkg && cd $HOME/xpkg
git clone https://github.com/xianyi/OpenBLAS.git
cd OpenBLAS
make -j4
sudo make PREFIX=/usr/local install
```

**Note**: Make sure to set the `/usr/local/lib` directory as a searchable LD\_LIBRARY\_PATH.
Otherwise, the following error may happen:
```
[...] libopenblas.so.0: cannot open shared object file: [...]
```
(see Section 3.3 above to fix this problem)

Now, install and test subpackage `la/oblas`:
```
cd ${GOPATH%:*}/src/github.com/cpmech/la/oblas
go install
go test
```



### 3.5 Optional: Install Intel MKL
Download MKL from [the intel MKL website](https://software.intel.com/en-us/intel-mkl), then:
```
mkdir -p $HOME/xpkg && cd $HOME/xpkg
tar xzvf l_mkl_2017.2.174.tgz
bash install_GUI.sh
```
and follow the instructions. These options have been tested:
1. Choose _Install as root using sudo_
2. Keep default install location: **/opt/intel**

Now, install and test subpackage `la/mkl`:
```
cd ${GOPATH%:*}/src/github.com/cpmech/la/mkl
go install
go test
```



### 3.6 Optional: Install OpenCV
```
sudo apt-get install libgtk2.0-dev pkg-config libavcodec-dev libavformat-dev libswscale-dev
mkdir -p $HOME/xpkg && cd $HOME/xpkg
git clone https://github.com/opencv/opencv.git
mkdir build_opencv
cd build_opencv
ccmake ../opencv
```
press `[c][c][g]`
```
make -j4
sudo make install
```

Now, install and test subpackge `img/ocv`:
```
cd ${GOPATH%:*}/src/github.com/cpmech/img/ocv
go install
go test
```



### Python dependencies

At the moment, the `plt` subpackage requires `python-scipy` and `python-matplotlib` (installed as
above in Section 3.1).



## Documentation

Here, we call _user-defined_ types as _structures_. These are simply Go `types` defined as `struct`.
Some may think of these _structures_ as _classes_. Gosl has several global functions as well and
tries to avoid complicated constructions.

An allocated structure is called an **object** and functions attached to this object are called
**methods**. The variable holding the pointer to an object is always named **o** in Gosl. This
variable is similar to `self` or `this` in other languages (Python, C++, respectively).

Some objects need to be initialised before use. In this case, functions named `Init` have to be
called (e.g. like `constructors`). Some structures require an explicit call to a function to release
allocated memory. This function is named `Free`. Functions that allocate a pointer to a structure
are prefixed with `New`; for instance: `NewIsoSurf`.

The directories corresponding to each package has a README.md file that should help with
understanding the library more in details. These include links to the definition of all functions
and structures (the developer's documentation, generated by `godoc`).

Again, Gosl has several functions and _structures_. Check the **[developer's
documentation](http://rawgit.com/cpmech/gosl/master/doc/index.html)** to see what's available and
how to call functions and methods.





## Bibliography

Gosl is partly included in the following works:

1. Pedroso DM (2017) FORM reliability analysis using a parallel evolutionary algorithm, Structural Safety 65:84-99 http://dx.doi.org/10.1016/j.strusafe.2017.01.001
2. Pedroso DM, Zhang YP, Ehlers W (2017) Solution of liquid-gas-solid coupled equations for porous media considering dynamics and hysteretic retention behaviour, Journal of Engineering Mechanics 04017021 http://dx.doi.org/10.1061/(ASCE)EM.1943-7889.0001208 
3. Pedroso DM (2015) A solution to transient seepage in unsaturated porous media. Computer Methods in Applied Mechanics and Engineering, 285:791-816 http://dx.doi.org/10.1016/j.cma.2014.12.009
4. Pedroso DM (2015) A consistent u-p formulation for porous media with hysteresis. Int. Journal for Numerical Methods in Engineering, 101(8):606-634 http://dx.doi.org/10.1002/nme.4808




## Authors and license

See the AUTHORS file.

Unless otherwise noted, the Gosl source files are distributed under the BSD-style license found in the LICENSE file.
