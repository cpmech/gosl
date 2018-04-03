#!/bin/bash

set -e

examples="\
chk_example01.go
fun_fft01.go
fun_orthopoly01.go
gm_bezier01.go
gm_bspline01.go
gm_bspline02.go
gm_nurbs01.go
gm_nurbs02.go
gm_nurbs03.go
graph_munkres01.go
graph_shortestpaths01.go
graph_siouxfalls01.go
la_HLsparseComplex01.go
la_HLsparseReal01.go
la_sparseComplex01.go
la_sparseReal01.go
ml_ang01.go
ml_ang02.go
ml_kmeans01.go
ml_mclass01.go
ml_simple01.go
num_brent01.go
num_deriv01.go
num_newton01.go
opt_comparison01.go
opt_conjgrad01.go
opt_ipm01.go
opt_ipm02.go
plt_bigo.go
plt_contour01.go
plt_polygon01.go
plt_zoomwindow01.go
rnd_haltonAndLatin01.go
rnd_ints01.go
rnd_lognormalDistribution.go
rnd_normalDistribution.go
tri_delaunay01.go
tri_draw01.go
tri_generate01.go
"

#vtk_isosurf01.go
#vtk_spheres01.go

for ex in $examples; do
    echo
    echo
    echo "[1;32m>>> running $ex <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<[0m"
    go run "$ex"
done
