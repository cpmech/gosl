# Gosl. rnd. Random numbers and probability distributions

More information is available in **[the documentation of this package](http://rawgit.com/cpmech/gosl/master/doc/xxrnd.html).**

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

### Examples

**Generate 100,000 integers and draw Histogram**

```go
const NSAMPLES = 100000

// initialise seed with fixed number; use 0 to use current time
rnd.Init(1234)

// allocate slice for integers
nints := 10
vals := make([]int, NSAMPLES)

// using the rnd.Int function
t0 := time.Now()
for i := 0; i < NSAMPLES; i++ {
    vals[i] = rnd.Int(0, nints-1)
}
io.Pf("time elapsed = %v\n", time.Now().Sub(t0))

// text histogram
hist := rnd.IntHistogram{Stations: utl.IntRange(nints + 1)}
hist.Count(vals, true)
io.Pf(rnd.TextHist(hist.GenLabels("%d"), hist.Counts, 60))

// using the rnd.Ints function
t0 = time.Now()
rnd.Ints(vals, 0, nints-1)
io.Pf("time elapsed = %v\n", time.Now().Sub(t0))

// text histogram
hist.Count(vals, true)
io.Pf(rnd.TextHist(hist.GenLabels("%d"), hist.Counts, 60))
```

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

Source code: <a href="../examples/rnd_ints01.go">../examples/rnd_ints01.go</a>



## Probability distributions

**Generate samples based on the Lognormal distribution**

```go
// initialise generator
rnd.Init(1234)

// parameters
μ := 1.0
σ := 0.25

// generate samples
nsamples := 1000
X := make([]float64, nsamples)
for i := 0; i < nsamples; i++ {
    X[i] = rnd.Lognormal(μ, σ)
}

// constants
nstations := 41        // number of bins + 1
xmin, xmax := 0.0, 3.0 // limits for histogram

// build histogram: count number of samples within each bin
var hist rnd.Histogram
hist.Stations = utl.LinSpace(xmin, xmax, nstations)
hist.Count(X, true)

// compute area of density diagram
area := hist.DensityArea(nsamples)
io.Pf("area = %v\n", area)

// plot lognormal distribution
var dist rnd.DistLogNormal
dist.Init(&rnd.VarData{M: μ, S: σ})

// compute lognormal points for plot
n := 101
x := utl.LinSpace(0, 3, n)
y := make([]float64, n)
Y := make([]float64, n)
for i := 0; i < n; i++ {
    y[i] = dist.Pdf(x[i])
    Y[i] = dist.Cdf(x[i])
}

// plot density
plt.Reset(false, nil)
plt.Subplot(2, 1, 1)
plt.Plot(x, y, nil)
hist.PlotDensity(nil)
plt.AxisYmax(1.8)
plt.Gll("$x$", "$f(x)$", nil)

// plot cumulative function
plt.Subplot(2, 1, 2)
plt.Plot(x, Y, nil)
plt.Gll("$x$", "$F(x)$", nil)

// save figure
plt.Save("/tmp/gosl", "rnd_lognormalDistribution")
```

Source code: <a href="../examples/rnd_lognormalDistribution.go">../examples/rnd_lognormalDistribution.go</a>

<div id="container">
<p><img src="../examples/figs/rnd_lognormalDistribution.png" width="400"></p>
</div>

## Sampling algorithm: Halton points

## Sampling algorithm: Latin hypercube
