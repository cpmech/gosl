# Gosl. gm. Geometry algorithms and structures

[![GoDoc](https://godoc.org/github.com/cpmech/gosl/gm?status.svg)](https://godoc.org/github.com/cpmech/gosl/gm) 

More information is available in **[the documentation of this package](https://godoc.org/github.com/cpmech/gosl/gm).**

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



### Generate and draw a NURBS curve


**Quarter of Circle Curve**

```go
// set of available vertices => to define control points
verts := [][]float64{
    {0.0, 0.0, 0, 1},
    {1.0, 0.2, 0, 1},
    {0.5, 1.5, 0, 1},
    {2.5, 2.0, 0, 1},
    {2.0, 0.4, 0, 1},
    {3.0, 0.0, 0, 1},
}

// NURBS knots
knots := [][]float64{
    {0, 0, 0, 0, 0.3, 0.7, 1, 1, 1, 1},
}

// NURBS curve
gnd := 1           // geometry number of dimensions: 1=>curve, 2=>surface, 3=>volume
orders := []int{3} // 3rd order along the only dimension
curve := new(gm.Nurbs)
curve.Init(gnd, orders, knots)
curve.SetControl(verts, utl.IntRange(len(verts)))

// refine NURBS
refined := curve.Krefine([][]float64{
    {0.15, 0.5, 0.85}, // new knots along the only dimension
})

// configuration
argsCtrlA := &plt.A{C: "k", Ls: "--", L: "control"}
argsCtrlB := &plt.A{C: "green", L: "refined: control"}
argsElemsA := &plt.A{C: "b", L: "curve"}
argsElemsB := &plt.A{C: "orange", Ls: "none", M: "*", Me: 20, L: "refined: curve"}
argsIdsA := &plt.A{C: "k", Fsz: 7}
argsIdsB := &plt.A{C: "green", Fsz: 7}

// plot
ndim := 2
npts := 41
plt.Reset(true, &plt.A{WidthPt: 400})
curve.DrawCtrl(ndim, true, argsCtrlA, argsIdsA)
curve.DrawElems(ndim, npts, true, argsElemsA, nil)
refined.DrawCtrl(ndim, true, argsCtrlB, argsIdsB)
refined.DrawElems(ndim, npts, false, argsElemsB, nil)
plt.AxisOff()
plt.Equal()
plt.LegendX([]*plt.A{argsCtrlA, argsCtrlB, argsElemsA, argsElemsB}, &plt.A{LegOut: true, LegNcol: 2})
plt.Save("/tmp/gosl", "gm_nurbs01")
```

Source code: <a href="../examples/gm_nurbs01.go">../examples/gm_nurbs01.go</a>

<div id="container">
<p><img src="../examples/figs/gm_nurbs01.png" width="500"></p>
</div>



**Circle Curve using Factory**

```go
// curve
xc, yc, r := 0.5, 0.5, 1.5
curve := gm.FactoryNurbs.Curve2dCircle(xc, yc, r)

// configuration
argsIdsA := &plt.A{C: "k", Fsz: 7}
argsCtrlA := &plt.A{C: "k", M: ".", Ls: "--", L: "control"}
argsElemsA := &plt.A{C: "b", L: "curve"}

// plot
ndim := 2
npts := 41
plt.Reset(true, &plt.A{WidthPt: 400})
curve.DrawCtrl(ndim, true, argsCtrlA, argsIdsA)
curve.DrawElems(ndim, npts, true, argsElemsA, nil)
plt.HideAllBorders()
plt.Equal()
plt.AxisRange(-2.5, 2.5, -2.5, 2.5)
plt.Save("/tmp/gosl", "gm_nurbs02")
```

Source code: <a href="../examples/gm_nurbs02.go">../examples/gm_nurbs02.go</a>

<div id="container">
<p><img src="../examples/figs/gm_nurbs02.png" width="500"></p>
</div>
