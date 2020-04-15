# Gosl. plt. Plotting and drawing (png and eps)

[![GoDoc](https://pkg.go.dev/github.com/cpmech/gosl/plt?status.svg)](https://pkg.go.dev/github.com/cpmech/gosl/plt) 

More information is available in **[the documentation of this package](https://pkg.go.dev/github.com/cpmech/gosl/plt).**

This package provides several functions to draw figures, plot results and annotate graphs. `plt` has
been largelly based on [matplotlib](https://matplotlib.org) and is, currently, a wrapper to
Python/pyplot by generating scripts to be run by an external `os` call.

Some basic functions are (similar to matplotlib ones to some extent):
1. `ReplaceAxes`, `Clf`
2. `Hist`, `Plot`, `Text`
3. `Show`, `Grid`
4. `Contour`, `Quiver`

The `plt` package has also some _higher level_ functions such as:
1. `Arrow`, `Circle`, `Polyline`
2. `AutoScale`, `AxisOff`, `AxisRange`
4. `Camera`, `Cross` (indicating the origin)
5.  `SetScientificX`, `SetTicksX`, `SetTicksY`, `SetXlog`, `SetYlog`
6. `Gll` (grid-labels-legend)

Functions to draw and handle 3D graphs are also available:
1. `Plot3dLine`, `Plot3dPoint`, `Plot3dPoints`
2. `Wireframe`, `Surface`, `Hemisphere`, `Superquadric`
3. `AxisRange3d`, `CylinderZ`, `ConeZ`

Nonetheless, interactive 3D graphs can also be developed with the `vtk` subpackage.

To initialise the figure, view and save a file (PNG or EPS), the following commands are available:
1. `Reset` initialises drawing space (optional)
2. `Show` show figure
3. `Save` saves the figure, after creating a directory.

All functions take a pointer to a structure holding optional arguments, the `A` structure that
belongs to the `plt` package, i.e. `plt.A`.


## Examples


### Drawing a polygon

Source code: <a href="../examples/plt_polygon01.go">../examples/plt_polygon01.go</a>

<div id="container">
<p><img src="../examples/figs/plt_polygon01.png" width="400"></p>
Polygon
</div>



### Plotting a contour

Source code: <a href="../examples/plt_contour01.go">../examples/plt_contour01.go</a>

<div id="container">
<p><img src="../examples/figs/plt_contour01.png" width="400"></p>
Contour and vector field
</div>



### Plotting with zoom window

Source code: <a href="../examples/plt_zoomwindow01.go">../examples/plt_zoomwindow01.go</a>

<div id="container">
<p><img src="../examples/figs/plt_zoomwindow01.png" width="400"></p>
</div>



### Drawing a box and 3D points

Source code: <a href="../examples/plt_boxandpoints.go">../examples/plt_boxandpoints.go</a>

<div id="container">
<p><img src="../examples/figs/plt_boxandpoints.png" width="450"></p>
</div>



### Waterfall graph

Source code: <a href="t_extra_test">t_extra_test</a>

<div id="container">
<p><img src="../examples/figs/t_waterfall01.png" width="450"></p>
</div>



### Drawing slope indicators

Source code: <a href="t_extra_test">t_extra_test</a>

<div id="container">
<p><img src="../examples/figs/t_slopeind01.png" width="450"></p>
</div>

<div id="container">
<p><img src="../examples/figs/t_slopeind02.png" width="450"></p>
</div>


## Output of Tests

Below, you will find some figures produced by the tests in <a href="t_plot01_test.go">t_plot01_test.go</a>.

**Test\_plot01**

<div id="container">
<p><img src="../examples/figs/t_plot01.png" width="400"></p>
</div>

**Test\_plot03**

<div id="container">
<p><img src="../examples/figs/t_plot03.png" width="400"></p>
</div>

**Test\_plot04**

<div id="container">
<p><img src="../examples/figs/t_plot04.png" width="400"></p>
</div>

**Test\_plot05**

<div id="container">
<p><img src="../examples/figs/t_plot05.png" width="400"></p>
</div>

**Test\_plot06**

<div id="container">
<p><img src="../examples/figs/t_plot06.png" width="450"></p>
</div>

**Test\_plot07**

<div id="container">
<p><img src="../examples/figs/t_plot07.png" width="450"></p>
</div>

**Test\_plot08**

<div id="container">
<p><img src="../examples/figs/t_plot08.png" width="450"></p>
</div>

**Test\_plot09**
<div id="container">
<p><img src="../examples/figs/t_plot09.png" width="450"></p>
</div>
