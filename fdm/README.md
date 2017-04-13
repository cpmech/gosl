# Gosl. fdm. Simple finite differences method

This package implements few routines to help with the implementation of the finite differences
method (FDM).

Basically, it has a `Grid2D` structure and functions to _assemble_ the linear system, considering
the existence of _known_ equations due to boundary conditions.

Another key structure of this package is the `Equations` data type that organises the numbers (ids)
assigned to each equation originating from the FDM discretisation.

A typical FDM problem will require the solution of the following linear system:

**K** * **U** = **F**

where **K** is a matrix (sometimes called the _stiffness matrix_), **U** is the solution vector
(e.g. displacements), and **F** is the _right-hand-side_ vector (e.g. forces).

Now, some components of the **U** vector are known (given or prescribed) due to the boundary
conditions. To automatize the solution, the **K** matrix is _split_ as shown below:

```
K11 K12 => unknowns
K21 K22 => prescribed
 |   +---> prescribed
 +-------> unknowns
```

where the numbers **1** indicate unknown quantities and the numbers **2** indicated known
(prescribed) quantities.

The _ids_ (numbers, indices) of each row in **U**, i.e. the **equation numbers**, can be organised
in the `Equations` structure defined below:

```go
type Equations struct {
    N1, N2, N int   // unknowns, prescribed, total numbers
    RF1, FR1  []int // reduced=>full and full=>reduced maps for unknowns
    RF2, FR2  []int // reduced=>full and full=>reduced maps for prescribed
}
```

where

```
"RF" means "reduced to full"
"FR" means "full to reduced"
reduced 1: is a reduced system of equations where only unknown equations are present
reduced 2: is a reduced system of equations where only prescribed equations are present
full     : corresponds to all equations, unknown and prescribed
```

For instance:

```
N   = 9         => total number of equations
peq = [0  3  6] => prescribed equations

  0  1  2  3  4  5  6  7  8                       0   1  2   3   4  5   6   7  8
0 +--------+--------+------ 0 -- prescribed -- 0 [22] 21 21 [22] 21 21 [22] 21 21 0
1 |        |        |       1                  1  12  .. ..  12  .. ..  12  .. .. 1
2 |        |        |       2                  2  12  .. ..  12  .. ..  12  .. .. 2
3 +--------+--------+------ 3 -- prescribed -- 3 [22] 21 21 [22] 21 21 [22] 21 21 3
4 |        |        |       4                  4  12  .. ..  12  .. ..  12  .. .. 4
5 |        |        |       5                  5  12  .. ..  12  .. ..  12  .. .. 5
6 +--------+--------+------ 6 -- prescribed -- 6 [22] 21 21 [22] 21 21 [22] 21 21 6
7 |        |        |       7                  7  12  .. ..  12  .. ..  12  .. .. 7
8 |        |        |       8                  8  12  .. ..  12  .. ..  12  .. .. 8
  0  1  2  3  4  5  6  7  8                       0   1  2   3   4  5   6   7  8
  |        |        |
 pre      pre      pre                                 .. => 11 equations

N1 = 6 => number of equations of type 1 (unknowns)
N2 = 3 => number of equations of type 2 (prescribed)
N1 + N2 == N
```

Therefore:
```
       0  1  2  3  4  5
RF1 = [1  2  4  5  7  8]           => ex:  RF1[3] = full system equation # 5

       0  1  2  3  4  5  6  7  8
FR1 = [   0  1     2  3     4  5]  => ex:  FR1[5] = reduced system equation # 3
      -1       -1       -1         => indicates 'value not set'

       0  1  2
RF2 = [0  3  6]                    => ex:  RF2[1] = full system equation # 3

       0  1  2  3  4  5  6  7  8
FR2 = [0        1        2      ]  => ex:  FR2[3] = reduced system equation # 1
         -1 -1    -1 -1    -1 -1   => indicates 'value not set'
```

## Examples

```
// grid
var g fdm.Grid2D
g.Init(1.0, 1.0, 3, 3)

// equations numbering
var e fdm.Equations
e.Init(g.N, []int{0, 3, 6})

// K11 and K12
var K11, K12 la.Triplet
fdm.InitK11andK12(&K11, &K12, &e)

// assembly
F1 := make([]float64, e.N1)
fdm.Assemble(&K11, &K12, F1, nil, &g, &e)
```

will output:


// check
K11d := K11.ToMatrix(nil).ToDense()
K12d := K12.ToMatrix(nil).ToDense()
K11c := [][]float64{
    {16.0, -4.0, -8.0, 0.0, 0.0, 0.0},
    {-8.0, 16.0, 0.0, -8.0, 0.0, 0.0},
    {-4.0, 0.0, 16.0, -4.0, -4.0, 0.0},
    {0.0, -4.0, -8.0, 16.0, 0.0, -4.0},
    {0.0, 0.0, -8.0, 0.0, 16.0, -4.0},
    {0.0, 0.0, 0.0, -8.0, -8.0, 16.0},
}
K12c := [][]float64{
    {-4.0, 0.0, 0.0},
    {0.0, 0.0, 0.0},
    {0.0, -4.0, 0.0},
    {0.0, 0.0, 0.0},
    {0.0, 0.0, -4.0},
    {0.0, 0.0, 0.0},
}
chk.Matrix(tst, "K11: ", 1e-16, K11d, K11c)
chk.Matrix(tst, "K12: ", 1e-16, K12d, K12c)
```


More information is available in [the documentation of this package](http://rawgit.com/cpmech/gosl/master/doc/xxfdm.html).
