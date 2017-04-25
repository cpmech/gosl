#!/bin/bash

set -e

examples="\
chk_example01.go
fdm_grid2d.go
fdm_problem01.go
fdm_problem02.go
gm_bezier01.go
gm_bspline01.go
gm_bspline02.go
gm_nurbs01.go
gm_nurbs02.go
gm_nurbs03.go
la_HLsparseComplex01.go
la_HLsparseReal01.go
la_sparseComplex01.go
la_sparseReal01.go
num_brent01.go
num_deriv01.go
num_newton01.go
plt_contour01.go
plt_polygon01.go
plt_zoomwindow01.go
README-to-html.go
rnd_normalDistribution.go
vtk_cone01.go
vtk_isosurf01.go
"

for ex in $examples; do
    echo
    echo
    echo "[1;32m>>> running $ex <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<[0m"
    go run "$ex"
done

