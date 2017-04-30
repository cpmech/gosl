# Gosl. gm. Graph theory structures and algorithms

More information is available in **[the documentation of this package](http://rawgit.com/cpmech/gosl/master/doc/xxgraph.html).**

## Munkres (Hungarian algorithm): the assignment problem

The Munkres method, also known as the Hungarian algorithm, aims to solve the assignment problem;
i.e. problems such as the following example:

```
Minimise the cost of operation when assigning three employees (Fry, Leela, Bender) to three tasks,
where (the minimum cost is highlighted):

           $ | Clean  Sweep   Wash
      -------|--------------------
      Fry    |   [2]      3      3
      Leela  |     3    [2]      3
      Bender |     3      3    [2]
      minimum cost = 6
```

The code is based on [Bob Pilgrim' work](http://csclab.murraystate.edu/bob.pilgrim/445/munkres.html).

The method runs in O(n²), in the worst case; therefore is not efficient for large matrices.

The `Munkres` structure implements the solver.

### Examples

Solution of Euler's problem # 345 using the Hungarian algorithm.

```go
// problem matrix
C := [][]float64{
    {7, 53, 183, 439, 863, 497, 383, 563, 79, 973, 287, 63, 343, 169, 583},
    {627, 343, 773, 959, 943, 767, 473, 103, 699, 303, 957, 703, 583, 639, 913},
    {447, 283, 463, 29, 23, 487, 463, 993, 119, 883, 327, 493, 423, 159, 743},
    {217, 623, 3, 399, 853, 407, 103, 983, 89, 463, 290, 516, 212, 462, 350},
    {960, 376, 682, 962, 300, 780, 486, 502, 912, 800, 250, 346, 172, 812, 350},
    {870, 456, 192, 162, 593, 473, 915, 45, 989, 873, 823, 965, 425, 329, 803},
    {973, 965, 905, 919, 133, 673, 665, 235, 509, 613, 673, 815, 165, 992, 326},
    {322, 148, 972, 962, 286, 255, 941, 541, 265, 323, 925, 281, 601, 95, 973},
    {445, 721, 11, 525, 473, 65, 511, 164, 138, 672, 18, 428, 154, 448, 848},
    {414, 456, 310, 312, 798, 104, 566, 520, 302, 248, 694, 976, 430, 392, 198},
    {184, 829, 373, 181, 631, 101, 969, 613, 840, 740, 778, 458, 284, 760, 390},
    {821, 461, 843, 513, 17, 901, 711, 993, 293, 157, 274, 94, 192, 156, 574},
    {34, 124, 4, 878, 450, 476, 712, 914, 838, 669, 875, 299, 823, 329, 699},
    {815, 559, 813, 459, 522, 788, 168, 586, 966, 232, 308, 833, 251, 631, 107},
    {813, 883, 451, 509, 615, 77, 281, 613, 459, 205, 380, 274, 302, 35, 805},
}

// solver seeks to minimise cost, thus, multiply coefficients by -1
for i := 0; i < len(C); i++ {
    for j := 0; j < len(C[i]); j++ {
        C[i][j] *= -1
    }
}

// initialise solver
var mnk graph.Munkres
mnk.Init(len(C), len(C[0]))
mnk.SetCostMatrix(C)

// solve rpoblem
mnk.Run()

// results
io.Pf("links = %v\n", mnk.Links)
io.Pf("cost = %v  (optimal 13938)\n", -mnk.Cost)
```

Output:
```
links = [9 10 7 4 3 0 13 2 14 11 6 5 12 8 1]
cost = 13938  (optimal 13938)
```

Source code: <a href="../examples/graph_munkres01.go">../examples/graph_munkres01.go</a>
