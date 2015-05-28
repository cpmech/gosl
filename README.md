# Gosl &ndash; Go scientific library

Gosl contains a mix of routines for research on numerical methods. These include linear algebra,
linear solvers, nonlinear solvers, ordinary differential equation solvers, plotting routines,
geometry functions, utilities for making better use of Go (golang), and more.

## License

Unless otherwise noted, the Gosl source files are distributed
under the BSD-style license found in the LICENSE file.

## Download

```
sudo apt-get install libopenmpi-dev libhwloc-dev libsuitesparse-dev libmumps-dev 
sudo apt-get install gfortran libvtk5-dev python-scipy python-matplotlib dvipng
go get github.com/cpmech/gosl/...
```

## Configurations

1. gosl/scripts must be in PYTHONPATH

## Website

http://cpmech.github.io/gosl

## Documentation for developers

http://rawgit.com/cpmech/gosl/master/doc/index.html

## References
Most of Gosl was developed to support the numerical research that resulted in the following papers:

1. Pedroso DM (2015) A consistent u-p formulation for porous media with hysteresis. Int Journal for Numerical Methods in Engineering, 101(8) 606-634 http://dx.doi.org/10.1002/nme.4808
2. Pedroso DM (2015) A solution to transient seepage in unsaturated porous media. Computer Methods in Applied Mechanics and Engineering, 285 791-816 http://dx.doi.org/10.1016/j.cma.2014.12.009
3. Pedroso DM, Sheng D and Zhao, J (2009) The concept of reference curves for constitutive modelling in soil mechanics, Computers and Geotechnics, 36(1-2), 149-165, http://dx.doi.org/10.1016/j.compgeo.2008.01.009
4. Pedroso DM and Williams DJ (2010) A novel approach for modelling soil-water characteristic curves with hysteresis, Computers and Geotechnics, 37(3), 374-380, http://dx.doi.org/10.1016/j.compgeo.2009.12.004
5. Pedroso DM and Williams DJ (2011) Automatic Calibration of soil-water characteristic curves using genetic algorithms. Computers and Geotechnics, 38(3), 330-340, http://dx.doi.org/10.1016/j.compgeo.2010.12.004

## Acknowledgements
Funding from the Australian Research Council is gratefully acknowledged.
