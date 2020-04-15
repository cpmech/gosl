# Gosl. vtk. 3D Visualisation with the VTK tool kit

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/cpmech/gosl/vtk)

More information is available in **[the documentation of this package](https://pkg.go.dev/github.com/cpmech/gosl/vtk).**

This package uses the [Visualisation Toolkit](http://www.vtk.org) to generate interactive 3D
visualisations.

## Examples

### Drawing spheres

The `Sphere` structure draws one sphere whereas the `Spheres` structure draws many spheres at once.
A data file with the spheres coordinates and radii is considered.

Source code: <a href="../examples/vtk_spheres01.go">../examples/vtk_spheres01.go</a>

<div id="container">
<p><img src="../examples/figs/vtk_spheres01.png" width="400"></p>
</div>

### Drawing an isosurface

Source code: <a href="../examples/vtk_isosurf01.go">../examples/vtk_isosurf01.go</a>

<div id="container">
<p><img src="../examples/figs/vtk_isosurf01.png" width="400"></p>
</div>
