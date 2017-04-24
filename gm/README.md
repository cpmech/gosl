# Gosl. gm. Geometry algorithms and structures

More information is available in **[the documentation of this package](http://rawgit.com/cpmech/gosl/master/doc/xxgm.html).**

The package `gm` provides some functions to help with the solution of geometry problems or simply to
perform computations with structures at least loosely related with geometry.

Some basic geometry structures and functions are:
1. `Point`, `Segment`
2. `DistPointPoint`, `DistPointLine`
3. `VecDot`, `VecNorm`

Some structures and functions to perform fast searches in the 3D space are:
1. `Bin`, `Bins`, `HashPoint`

The more _advanced_ structures in this package are:
1. `BezierQuad`: quadratic Bezier curve and algorithms
2. `Bspline`: B-spline and algorithms, including fast computation of derivatives
3. `Nurbs`: NURBS and algorithms, of derivatives. NURBS solids are also implemented.

## Examples

### Draw Bezier curve


```go
// quadratic Bezier
var bez gm.BezierQuad

// control points
bez.Q = [][]float64{
    {-1.0, +1.0},
    {+0.5, -2.0},
    {+2.0, +4.0},
}

// points on Bezier curve
np := 11
Xb, Yb, _ := bez.GetPoints(utl.LinSpace(0, 1, np))

// quadratic curve
Xc := utl.LinSpace(-1, 2, np)
Yc := utl.GetMapped(Xc, func(x float64) float64 { return x * x })

// control points
Xq, Yq, _ := bez.GetControlCoords()

// plot
plt.Reset(true, &plt.A{WidthPt: 300})
plt.Plot(Xq, Yq, &plt.A{C: "k", M: "*", NoClip: true, L: "control"})
plt.Plot(Xc, Yc, &plt.A{C: "r", M: "o", Void: true, Ms: 10, Ls: "none", L: "y=x*x", NoClip: true})
plt.Plot(Xb, Yb, &plt.A{C: "b", Ls: "-", M: ".", L: "Bezier", NoClip: true})
plt.HideAllBorders()
plt.Gll("x", "y", &plt.A{LegLoc: "upper left"})
plt.Equal()
plt.Save("/tmp/gosl", "gm_bezier01")
```

Source code: <a href="../examples/gm_bezier01.go">../examples/gm_bezier01.go</a>

<div id="container">
<p><img src="../examples/figs/gm_bezier01.png" width="400"></p>
</div>



### Draw B-spline curve and control net

```go
// order of B-spline
p := 3

// knots for clamped curve
startT := make([]float64, p+1) // p+1 zeros
endT := utl.Ones(p + 1)        // p+1 ones

// knots
T1 := append(append(startT, utl.LinSpace(0.1, 0.9, 9)...), endT...)

// B-spline
var b1 gm.Bspline
b1.Init(T1, p)

// set control points
b1.SetControl([][]float64{
    {0.5, 0.5},
    {1.0, 0.5},
    {1.0, 0.5}, // repeated
    {1.0, 0.5}, // repeated => 3x => discontinuity
    {1.0, 0.2},
    {0.7, 0.0},
    {0.3, 0.0},
    {0.0, 0.3},
    {0.0, 0.7},
    {0.3, 1.0},
    {0.7, 1.0},
    {0.9, 0.9},
    {1.0, 0.8},
})

// configuration
withCtrl := true
argsCurve := &plt.A{C: "r", Lw: 10, L: "curve", NoClip: true}
argsCtrl := &plt.A{C: "k", M: ".", L: "control", NoClip: true}

// plot
np := 101
plt.Reset(false, nil)
b1.Draw2d(np, 0, withCtrl, argsCurve, argsCtrl)
plt.Equal()
plt.HideAllBorders()
plt.AxisXmax(1.0)
plt.Save("/tmp/gosl", "gm_bspline01")
```

Source code: <a href="../examples/gm_bspline01.go">../examples/gm_bspline01.go</a>

<div id="container">
<p><img src="../examples/figs/gm_bspline01.png" width="400"></p>
</div>
