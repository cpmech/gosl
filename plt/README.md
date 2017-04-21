# Gosl. plt. Plotting and drawing (png and eps)

More information is available in [the documentation of this package](http://rawgit.com/cpmech/gosl/master/doc/xxplt.html).

This package provides several functions to draw figures, plot results and annotate graphs. `plt` has
been largelly based on [matplotlib](https://matplotlib.org) and is, currently, a sort of wrapper to
Python/pyplot in the sense that it generates python scripts to be run by an external `os` call.

A future implementation will add an option to draw directly to the web browser, after generating
JavaScript codes. We also plan for a [QT](https://www.qt.io) version.

Many pyplot functions are _wrapped_. For instance:
1. `Axes`, `Clf`
2. `Hist`, `Plot`, `Text`
3. `Show`, `Grid`
4. `Contour`, `Quiver`

The `plt` package is in fact somewhat **more convenient** than the analogue pyplot because it comes
with a number of _higher level_ functions such as:
1. `Arrow`, `Circle`, `DrawPolyline`
2. `AutoScale`, `AxisOff`, `AxisRange`
3. `AxisXmin`, `AxisXmax`, `AxisYmin`, `AxisYmax`
4. `Camera`, `Cross` (indicating the origin)
5.  `SetScientific`, `SetTicksX`, `SetTicksY`, `SetXlog`, `SetYlog`
6. `Gll` (grid-labels-legend)

Functions to draw and handle 3D graphs are also available:
1. `Plot3dLine`, `Plot3dPoints`
2. `Wireframe`, `Surface`
3. `AxisRange3d` 
Nonetheless fancier graphs are better developed with the `vtk` subpackage.

To view and save figures, the following commands are available:
1. `SetForEps` prepares figure for generating an eps file, e.g. by giving the proportion and size
2. `SetForPng` prepares figure for generating a png file, e.g. by giving the proportion, size and resolution
3. `Save` and `SaveD` that saves the figure and save after creating a directory, respectively.

Most functions take an extra argument that is passed to python for further customisation.

In fact, you can do everything here (in Go) as you would do in python, because, if needed, you can
write python scripts of even load a python file to be run when showing or saving figures. The
following commands allow this:
1. `PyCmds` run python commands
2. `PyFile` load and run a python file

## Examples



### Drawing a polygon

```go
// point coordinates
P := [][]float64{
    {-2.5, 0.0},
    {-5.5, 4.0},
    {0.0, 10.0},
    {5.5, 4.0},
    {2.5, 0.0},
}

// formatting/styling data
// Fc: face color, Ec: edge color, Lw: linewidth
stPlines := &plt.A{Fc: "#c1d7cf", Ec: "#4db38e", Lw: 4.5, Closed: true, NoClip: true}
stCircles := &plt.A{Fc: "#b2cfa5", Ec: "#5dba35", Z: 1}
stArrows := &plt.A{Fc: "cyan", Ec: "blue", Z: 2, Scale: 50, Style: "fancy"}

// clear drawing area, with defaults
setDefault := true
plt.Reset(setDefault, nil)

// draw polyline
plt.Polyline(P, stPlines)

// draw circle
plt.Circle(0, 4, 2.0, stCircles)

// draw arrow
plt.Arrow(-4, 2, 4, 7, stArrows)

// draw arc
plt.Arc(0, 4, 3, 0, 90, nil)

// autoscale axes
plt.AutoScale(P)

// enforce same scales
plt.Equal()

// draw a _posteriori_ legend
plt.LegendX([]*plt.A{
    &plt.A{C: "red", M: "o", Ls: "-", Lw: 1, Ms: -1, L: "first", Me: -1},
    &plt.A{C: "green", M: "s", Ls: "-", Lw: 2, Ms: 0, L: "second", Me: -1},
    &plt.A{C: "blue", M: "+", Ls: "-", Lw: 3, Ms: 10, L: "third", Me: -1},
}, nil)

// save figure (default is PNG)
err := plt.Save("/tmp/gosl", "plt_polygon01")
if err != nil {
    io.Pf("error: %v\n", err)
}
```

<div id="container">
<p><img src="../examples/figs/plt_polygon01.png" width="400"></p>
Polygon
</div>



### Plotting a contour
```go
// grid size
xmin, xmax, N := -math.Pi/2.0+0.1, math.Pi/2.0-0.1, 21

// mesh grid, scalar and vector field
X, Y, F, U, V := utl.MeshGrid2dFG(xmin, xmax, xmin, xmax, N, N, func(x, y float64) (f, u, v float64) {

    // scalar field
    m := math.Pow(math.Cos(x), 2.0) + math.Pow(math.Cos(y), 2.0)
    f = -math.Pow(m, 2.0)

    // gradient. u=dfdx, v=dfdy
    u = 4.0 * math.Cos(x) * math.Sin(x) * m
    v = 4.0 * math.Cos(y) * math.Sin(y) * m
    return
})

// plot
plt.Reset(false, nil)
plt.ContourF(X, Y, F, &plt.A{CmapIdx: 4, Nlevels: 15})
plt.Quiver(X, Y, U, V, &plt.A{C: "r"})
plt.Gll("x", "y", nil)
plt.Equal()
plt.Save("/tmp/gosl", "plt_contour01")
```

<div id="container">
<p><img src="../examples/figs/plt_contour01.png" width="400"></p>
Contour and vector field
</div>
