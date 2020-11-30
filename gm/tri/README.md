# Gosl. gm/tri. Mesh generation: triangles

[![PkgGoDev](https://pkg.go.dev/badge/github.com/cpmech/gosl/gm/tri)](https://pkg.go.dev/github.com/cpmech/gosl/gm/tri)

The `tri` package has functions go generate meshes of triangular elements and Delaunay
triangulations. The package is a light wrapper to the very efficient `Triangle` code by Jonathan
Shewchuk and available at [the Triangle website](https://www.cs.cmu.edu/~quake/triangle.html).

Triangle's licence details are in the <a href="triangle_README.txt">triangle_README.txt</a> file.

In addition of being very fast, `Triangle` can generate meshes with great quality; i.e. it is a
**quality-mesh triangulator**.

Here, the Cartesian coordinates of points are stored in continuous 1D arrays (slices) such as:

```
X = { x0, x1, x2, ... Npoints }
Y = { y0, y1, y2, ... Npoints }
```

The topology is simply defined by two slices usually named `V` for vertices and `C` for cells
(triangles):

```
V = { { x0, y0 }, { x1, y1 }, { x2, y2 } ... Nvertices }
C = { { id0, id1, id2 }, { id0, id1, id2 } ... Ncellls }
```

where `V` is the list of vertices (there are Nvertices vertices) and `C` is the list of triangles
(there are Ncells triangles). The ids (e.g. id0, id1, id2) in `C` are the indices in `V`.

## Draw mesh

![](data/tri_draw01.png)

For example, the set of triangles in the above figure are defined (and drawn) with:

```go
// vertices (points)
V := [][]float64{
    {0, 0}, {1, 0},
    {1, 1}, {0, 1},
}

// cells (triangles)
C := [][]int{
    {0, 1, 2},
    {2, 3, 0},
}
```

Source code: <a href="../../examples/tri_draw01.go">../../examples/tri_draw01.go</a>

## Delaunay triangulation

The Delaunay triangulation of a cloud of points in the `tri` package is easily computed with the
`Delaunay` command that takes as input the Cartesian coordinates.

<div id="container">
<p><img src="../../examples/figs/tri_delaunay01.png"></p>
</div>

Source code: <a href="../../examples/tri_delaunay01.go">../../examples/tri_delaunay01.go</a>

## API

[Please see the documentation here](https://pkg.go.dev/github.com/cpmech/gosl/gm/tri)
