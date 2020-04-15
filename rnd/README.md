# Gosl. rnd. Random numbers and probability distributions

[![GoDoc](https://pkg.go.dev/github.com/cpmech/gosl/rnd?status.svg)](https://pkg.go.dev/github.com/cpmech/gosl/rnd) 

More information is available in **[the documentation of this package](https://pkg.go.dev/github.com/cpmech/gosl/rnd).**

The `rnd` package assists on computations involving stochastic processes. The package has many
functions to generate pseudo-random numbers, probability distributions, and sampling techniques such
as the Latin hypercube algorithm.

## Pseudo random numbers

In this package, some standard Go functions from package `rand` are wrapped for convenience. Most
random generation functions have an equivalent function wrapping the Mersenne Twister code from
Mutsuo Saito and Makoto Matsumoto.

Some useful functions are:
1. `Init` initialise the system with a seed
2. `Int`, `Ints`, `Float64`, `Float64s` to generate integers and floats
3. Shuffle and GetUnique functions to shuffle slices and filter slices with unique values,
   respectively.

## Probability distributions

The probability distributions in the `rnd` package are initialised with the help of the `VarData`
structure that contains the following main fields:
```go
// input
D DistType // type of distribution
M float64  // mean
S float64  // standard deviation

// input: Frechet
L float64 // location
C float64 // scale
A float64 // shape

// input: uniform
Min float64 // min value
Max float64 // max value
```

The currently available distributions are:
1. `NormalKind`    Normal distribution
2. `LognormalKind` Lognormal distribution
3. `GumbelKind`    Type I Extreme Value distribution
4. `FrechetKind`   Type II Extreme Value distribution
5. `UniformKind`   Uniform distribution

## Sampling algorithms: Halton and Latin Hypercube methods

The `HaltonPoints` function is a simple way to generate combinations of point coordinates in a
hypercube.

The `LatinIHS` function implements the Latin improved distributed hypercube sampling method. The
results are the indices of points. The point coordinates can be computed with the `HypercubeCoords`
function.



## Examples

### Generate 100,000 integers and draw Histogram

Source code: <a href="../examples/rnd_ints01.go">../examples/rnd_ints01.go</a>

Output:

```
time elapsed = 3.506121ms
  [0,1) |  10085 ############################################################
  [1,2) |   9874 ###########################################################
  [2,3) |  10078 ############################################################
  [3,4) |   9998 ############################################################
  [4,5) |   9937 ###########################################################
  [5,6) |  10003 ############################################################
  [6,7) |  10119 #############################################################
  [7,8) |   9795 ###########################################################
  [8,9) |  10026 ############################################################
 [9,10) |  10085 ############################################################
  count = 100000
time elapsed = 3.259988ms
  [0,1) |  10077 ############################################################
  [1,2) |  10017 ############################################################
  [2,3) |   9910 ###########################################################
  [3,4) |  10092 ############################################################
  [4,5) |   9853 ###########################################################
  [5,6) |   9976 ############################################################
  [6,7) |  10096 #############################################################
  [7,8) |  10058 ############################################################
  [8,9) |   9905 ###########################################################
 [9,10) |  10016 ############################################################
  count = 100000
```


### Generate samples based on the Lognormal distribution

Source code: <a href="../examples/rnd_lognormalDistribution.go">../examples/rnd_lognormalDistribution.go</a>

<div id="container">
<p><img src="../examples/figs/rnd_lognormalDistribution.png" width="500"></p>
</div>


### Example: sampling algorithms

Source code: <a href="../examples/rnd_haltonAndLatin01.go">../examples/rnd_haltonAndLatin01.go</a>

<div id="container">
<p><img src="../examples/figs/rnd_haltonAndLatin01.png" width="400"></p>
</div>
