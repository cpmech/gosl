# Gosl. vtk. 3D Visualisation with the VTK tool kit

More information is available in **[the documentation of this package](http://rawgit.com/cpmech/gosl/master/doc/xxvtk.html).**

The `vtk` package presents routines to use the [Visualisation Toolkit](http://www.vtk.org) directly
from Go. Some few structures are defined, including:
1. `Scene` to hold all scene elements
2. `Arrow`, `Sphere` `Spheres` to draw arrows and spheres
3. `IsoSurf` to draw iso-surfaces

## Examples

### Drawing an isosurface

```go
// create a new VTK Scene
scn := vtk.NewScene()
scn.Reverse = true // start viewing the negative side of the x-y-z Cartesian system

// parameters
M := 1.0  // slope of line in p-q graph
pt := 0.0 // tensile p
a0 := 0.8 // size of surface

// limits and divisions for grid generation
pqth := []float64{pt, a0, 0, M * a0, 0, 360}
ndiv := []int{21, 21, 41}

// cone symbolising the Drucker-Prager criterion
cone := vtk.NewIsoSurf(func(x []float64) (f, vx, vy, vz float64) {
    p, q := calc_p(x), calc_q(x)
    f = q - M*p
    return
})
cone.Limits = pqth
cone.Ndiv = ndiv
cone.OctRotate = true
cone.GridShowPts = false
cone.Color = []float64{0, 1, 1, 1}
cone.CmapNclrs = 0 // use this to use specified color
cone.AddTo(scn)    // remember to add to Scene

// ellipsoid symbolising the Cam-clay yield surface
ellipsoid := vtk.NewIsoSurf(func(x []float64) (f, vx, vy, vz float64) {
    p, q := calc_p(x), calc_q(x)
    f = q*q + M*M*(p-pt)*(p-a0)
    return
})
ellipsoid.Limits = pqth
cone.Ndiv = ndiv
ellipsoid.OctRotate = true
ellipsoid.Color = []float64{1, 1, 0, 0.5}
ellipsoid.CmapNclrs = 0 // use this to use specified color
ellipsoid.AddTo(scn)    // remember to add to Scene

// illustrate use of Arrow
arr := vtk.NewArrow() // X0 is equal to origin
arr.V = []float64{-1, -1, -1}
arr.AddTo(scn)

// illustrate use of Sphere
sph := vtk.NewSphere()
sph.Cen = []float64{-a0, -a0, -a0}
sph.R = 0.05
sph.AddTo(scn)

// start interactive mode
scn.SavePng = true
scn.Fnk = "/tmp/vtk_isosurf01"
scn.Run()
```

Source code: <a href="../examples/vtk_isosurf01.go">../examples/vtk_isosurf01.go</a>

<div id="container">
<p><img src="../examples/figs/vtk_isosurf01.png" width="400"></p>
</div>


### Drawing a cone and spheres

```go
// create a new VTK Scene
scn := vtk.NewScene()
scn.HydroLine = true
scn.FullAxes = true
scn.AxesLen = 1.5

// parameters
α = 15.0
kα := math.Tan(α * math.Pi / 180.0)

// cone
cone := vtk.NewIsoSurf(func(x []float64) (f, vx, vy, vz float64) {
    f = cone_angle(x) - kα
    return
})
cone.Limits = []float64{0, -1, 0, 1, 0, 360}
cone.Ndiv = []int{21, 21, 41}
cone.OctRotate = true
cone.GridShowPts = false
cone.Color = []float64{0, 1, 0, 1}
cone.CmapNclrs = 0 // use this to use specified color
cone.AddTo(scn)    // remember to add to Scene

// spheres
sset := vtk.NewSpheresFromFile("data/points.dat")
if true {
    sset.AddTo(scn)
}

// start interactive mode
scn.SavePng = false
scn.Fnk = "/tmp/vtk_cone01"
scn.Run()
```

Source code: <a href="../examples/vtk_cone01.go">../examples/vtk_cone01.go</a>

<div id="container">
<p><img src="../examples/figs/vtk_cone01.png" width="400"></p>
</div>
