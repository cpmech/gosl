# Gosl &ndash; Go scientific library

Gosl is a computing library written in go language (golang) to help with the development of software
for scientific research. The library tries to be as general as possible. The use of concurrency for
multi-threaded applications and message passing for parallel computations are considered. Functions
and structures for geometry calculations, random numbers generation and probability distributions,
optimisation algorithms, and plotting and visualisation are implemented as well. This library helped
with the development of the results presented in [1-5].



## Content

1.  chk   &ndash; check and unit test
2.  io    &ndash; input/output
3.  utl   &ndash; utilities
4.  plt   &ndash; plotting
5.  mpi   &ndash; message passing interface
6.  lla   &ndash; linear algebra
7.  fdm   &ndash; finite differences method
8.  num   &ndash; numerical methods
9.  fun   &ndash; scalar functions of one scalar and one vector
10. gm    &ndash; geometry
11. graph &ndash; graph theory
12. ode   &ndash; ordinary differential equations
13. opt   &ndash; optimisation
14. rnd   &ndash; random numbers and probability distributions
15. tsr   &ndash; tensor algebra and definitions for continuum mechanics
16. vtk   &ndash; visualisation tool kit



## Examples

See examples here: https://github.com/cpmech/gosl/blob/master/examples/README.md



## Installation and documentation

```
mkdir -p $GOPATH/src/github.com/cpmech
cd $GOPATH/src/github.com/cpmech
git clone https://github.com/cpmech/gosl.git
cd gosl
./all.bash
```

The documentation for developers is available here: http://rawgit.com/cpmech/gosl/master/doc/index.html



## References

1. Pedroso DM (2015) A consistent u-p formulation for porous media with hysteresis. Int Journal for Numerical Methods in Engineering, 101(8) 606-634 http://dx.doi.org/10.1002/nme.4808
2. Pedroso DM (2015) A solution to transient seepage in unsaturated porous media. Computer Methods in Applied Mechanics and Engineering, 285 791-816 http://dx.doi.org/10.1016/j.cma.2014.12.009



## Acknowledgements

Funding from the Australian Research Council is gratefully acknowledged.



## License

Unless otherwise noted, the Gosl source files are distributed under the BSD-style license found in the LICENSE file.
