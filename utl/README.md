# Gosl. utl. Utilities. Lists. Dictionaries. Simple Numerics

More information is available in **[the documentation of this package](https://godoc.org/github.com/cpmech/gosl/utl).**

This package implements functions for simplifying numeric calculations such as finding the maximum
and minimum of lists (i.e. slices), allocation of _deep_ structures such as slices of slices, and
generation of _arrays_. It also contains functions for sorting quantities and updating dictionaries
(i.e. maps).

This package does not aim for high performance linear algebra computations. For that purpose, we
have the `la` and a future `lax` packages. Nonetheless, `utl` package is OK for _small computations_
such as for vectors in the 3D space. It also tries to use the best algorithms for sorting that were
implemented in the standard Go library.

Some _numerical_ functions are (a few were inspired by [NumPy](http://www.numpy.org)):
1. `IntRange` to generate slices of integers
2. `IntFill` `IntVals`, `IntsAlloc`, `IntsClone` for slices of integers
3. `ArgMinMax`, `Copy`, `GetSorted` for slices of _Float64_
4. `Imin`, `Imax`, `Min`, `Max` to find minima and maxima of integers and doubles
5. `ParetoFront` to generate the Pareto front
6. `Cross3d` to compute the cross product
7. `LinSpace` to generate a sequence of equally spaced numbers
8. `IntUnique` to filter a slice of integers by leaving unique items only (i.e. a _set_)
9. `MeshGrid2D` to generate a table of coordinates that define a 2D grid
10. many more!

Some functions to deal with slices of slices (of slices... of slices..), i.e. _deep_ slices are:
1. `Deep3alloc` to allocate the a slice of slice of slice, e.g. [][][]float64
2. `Deep3set`, `Deep3Serialize` and more to handle _deep3_ slices

Some functions to deal with maps (dictionaries) are:
1. `StrBoolMapSort` to sort the keys of a `string=>bool` map
2. `StrFltMapAppend` to append an item to a slice mapped to a key in a `string=>slice` map
3. `StrIntMapSortSplit` to extract the keys and values form a `string=>int` into sorted slices

Other functions help with CPU and memory profiling, such as:
1. `ProfCPU` and `ProfMEM`


## Examples

### Numerical functions

Generate lists of integers
```go
I := make([]int, 5)
utl.IntFill(I, 666)
J := utl.IntVals(5, 666)
M := utl.IntsAlloc(3, 4)
A := utl.IntRange(-1)
a := utl.IntRange2(0, 0)
b := utl.IntRange2(0, 1)
c := utl.IntRange2(0, 5)
C := utl.IntRange3(0, -5, -1)
d := utl.IntRange2(2, 5)
D := utl.IntRange2(-2, 5)
e := utl.IntAddScalar(D, 2)
g := []int{1, 2, 3, 4, 3, 4, 2, 1, 1, 2, 3, 4, 4, 2, 3, 7, 8, 3, 8, 3, 9, 0, 11, 23, 1, 2, 32, 12, 4, 32, 4, 11, 37}
h := utl.IntUnique(g)
G := []int{1, 2, 3, 38, 3, 5, 3, 1, 2, 15, 38, 1, 11}
H := utl.IntUnique(D, C, G, []int{16, 39})
```

Generate lists of doubles (float64)
```go
N := utl.Alloc(3, 4)
f := utl.Ones(5)
ff := utl.Vals(5, 666)
X, Y := utl.MeshGrid2D(3, 6, 10, 20, 4, 3)
b := utl.LinSpaceOpen(2.0, 3.0, n)
```

Cumulative sums
```go
p := []float64{1, 2, 3, 4, 5}
cs := make([]float64, len(p))
utl.CumSum(cs, p)
io.Pforan("cs = %v\n", cs)
chk.Vector(tst, "cumsum", 1e-17, cs, []float64{1, 3, 6, 10, 15})
```

Handling tables of doubles (float64)
```go
v := [][]float64{
    {1, 2, 3, 4},
    {-1, 2, 3, 0},
    {-1, 2, 1, 4},
    {1, -2, 3, 8},
    {1, 1, -3, 4},
    {0, 2, 9, -4},
}

x := utl.GetColumn(0, v)
chk.Vector(tst, "v[:,0]", 1e-17, x, []float64{1, -1, -1, 1, 1, 0})

x = utl.GetColumn(1, v)
chk.Vector(tst, "v[:,1]", 1e-17, x, []float64{2, 2, 2, -2, 1, 2})

x = utl.GetColumn(2, v)
chk.Vector(tst, "v[:,2]", 1e-17, x, []float64{3, 3, 1, 3, -3, 9})

x = utl.GetColumn(3, v)
chk.Vector(tst, "v[:,3]", 1e-17, x, []float64{4, 0, 4, 8, 4, -4})
```

Pareto fronts
```go
u := []float64{1, 2, 3, 4, 5, 6}
v := []float64{1, 2, 3, 4, 5, 6}
io.Pforan("u = %v\n", u)
io.Pfblue2("v = %v\n", v)

u_dominates, v_dominates := utl.ParetoMin(u, v)
io.Pfpink("u_dominates = %v\n", u_dominates)
io.Pfpink("v_dominates = %v\n", v_dominates)
```

### Slices and deep slices

Allocate deep slices
```go
a := utl.Deep3alloc(3, 2, 4)
utl.Deep3set(a, 666)
chk.Deep3(tst, "a", 1e-16, a, [][][]float64{
    {{666, 666, 666, 666}, {666, 666, 666, 666}},
    {{666, 666, 666, 666}, {666, 666, 666, 666}},
    {{666, 666, 666, 666}, {666, 666, 666, 666}},
})
io.Pf("a = %v\n", a)

b := utl.Deep4alloc(3, 2, 1, 2)
utl.Deep4set(b, 666)
chk.Deep4(tst, "b", 1e-16, b, [][][][]float64{
    {{{666, 666}}, {{666, 666}}},
    {{{666, 666}}, {{666, 666}}},
    {{{666, 666}}, {{666, 666}}},
})
io.Pf("b = %v\n", b)
```

Serialization
```go
A := [][][]float64{
    {{100, 101, 102}, {103}, {104, 105}},
    {{106}, {107}},
    {{108}, {109, 110}},
    {{111}},
    {{112, 113, 114, 115}, {116}, {117, 118}, {119, 120, 121}},
}

// serialize
utl.PrintDeep3("A", A)
I, P, S := utl.Deep3Serialize(A)
utl.Deep3GetInfo(I, P, S, true)

// check serialization
chk.Ints(tst, "I", I, []int{0, 0, 0, 1, 1, 2, 2, 3, 4, 4, 4, 4})
chk.Ints(tst, "P", P, []int{0, 3, 4, 6, 7, 8, 9, 11, 12, 16, 17, 19, 22})
```

Sorting
```go
i := []int{33, 0, 7, 8}
x := []float64{1000.33, 0, -77.7, 88.8}
y := []float64{1e-5, 1e-7, 1e-2, 1e-9}
z := []float64{-8000, -7000, 0, -1}

io.Pforan("by 'i'\n")
I, X, Y, Z, err := utl.SortQuadruples(i, x, y, z, "i")
if err != nil {
    tst.Errorf("%v\n", err)
}
chk.Ints(tst, "i", I, []int{0, 7, 8, 33})
chk.Vector(tst, "x", 1e-16, X, []float64{0, -77.7, 88.8, 1000.33})
chk.Vector(tst, "y", 1e-16, Y, []float64{1e-7, 1e-2, 1e-9, 1e-5})
chk.Vector(tst, "z", 1e-16, Z, []float64{-7000, 0, -1, -8000})
```

### Maps and dictionaries

Appending to maps of slices of float64
```go
m := map[string][]float64{
    "a": []float64{100, 101},
    "b": []float64{1000},
    "c": []float64{200, 300, 400},
}
io.Pforan("m (before) = %v\n", m)

utl.StrMapAppend(&m, "a", 102)
io.Pfpink("m (after) = %v\n", m)

chk.Vector(tst, "m[\"a\"]", 1e-16, m["a"], []float64{100, 101, 102})
chk.Vector(tst, "m[\"b\"]", 1e-16, m["b"], []float64{1000})
chk.Vector(tst, "m[\"c\"]", 1e-16, m["c"], []float64{200, 300, 400})

utl.StrFltsMapAppend(&m, "d", 666)
io.Pfcyan("m (after) = %v\n", m)

chk.Vector(tst, "m[\"a\"]", 1e-16, m["a"], []float64{100, 101, 102})
chk.Vector(tst, "m[\"b\"]", 1e-16, m["b"], []float64{1000})
chk.Vector(tst, "m[\"c\"]", 1e-16, m["c"], []float64{200, 300, 400})
chk.Vector(tst, "m[\"d\"]", 1e-16, m["d"], []float64{666})
chk.Vector(tst, "m[\"e\"]", 1e-16, m["e"], nil)
```

### Other


Finding the _best square_ for given `size = numberOfRows * numberOfColumns`
```go
for size := 1; size <= 12; size++ {
    nrow, ncol := utl.BestSquare(size)
    io.Pforan("nrow, ncol, nrow*ncol = %2d, %2d, %2d\n", nrow, ncol, nrow*ncol)
    if nrow*ncol != size {
        chk.Panic("BestSquare failed")
    }
}
```
