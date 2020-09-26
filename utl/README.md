# Gosl. utl. Utilities. Lists. Dictionaries. Simple Numerics

This package implements functions for simplifying numeric calculations such as finding the maximum
and minimum of lists (i.e. slices), allocation of _deep_ structures such as slices of slices, and
generation of _arrays_. It also contains functions for sorting quantities and updating dictionaries
(i.e. maps).

This package does not aim for high performance linear algebra computations. For that purpose, we
have the `la` package. Nonetheless, `utl` package is OK for _small computations_ such as for vectors
in the 3D space. It also tries to use the best algorithms for sorting that are implemented in the
standard Go library and others.

Example of what the functions here can do:

- Generate lists of integers
- Generate lists of float64s
- Cumulative sums
- Handling tables of float64s
- Pareto fronts
- Slices and deep (nested) slices up to the 4th depth
- Allocate deep slices
- Serialization of deep slices
- Sorting
- Maps and _dictionaries_
- Append to maps of slices of float64
- Find the _best square_ for given `size = numberOfRows * numberOfColumns`
- ...

## API

**go doc**

```
package utl // import "gosl/utl"

Package utl implements functions for simplifying calculations and allocation
of structures such as slices and slices of slices. It also contains
functions for sorting quantities.

CONSTANTS

const (
	// SQ2 = sqrt(2) [https://oeis.org/A002193]
	SQ2 = 1.41421356237309504880168872420969807856967187537694807317667974

	// SQ3 = sqrt(3) [https://oeis.org/A002194]
	SQ3 = 1.73205080756887729352744634150587236694280525381038062805580698

	// SQ5 = sqrt(5) [https://oeis.org/A002163]
	SQ5 = 2.23606797749978969640917366873127623544061835961152572427089724

	// SQ6 = sqrt(6) [https://oeis.org/A010464]
	SQ6 = 2.44948974278317809819728407470589139196594748065667012843269257

	// SQ7 = sqrt(7) [https://oeis.org/A010465]
	SQ7 = 2.64575131106459059050161575363926042571025918308245018036833446

	// SQ8 = sqrt(8) [https://oeis.org/A010466]
	SQ8 = 2.82842712474619009760337744841939615713934375075389614635335947
)

FUNCTIONS

func AllFalse(values []bool) bool
    AllFalse returns true if all values are false

func AllTrue(values []bool) bool
    AllTrue returns true if all values are true

func Alloc(m, n int) (mat [][]float64)
    Alloc allocates a slice of slices of float64

func ArgMinMax(v []float64) (imin, imax int)
    ArgMinMax finds the indices of min and max arguments

func BestSquare(size int) (nrow, ncol int)
    BestSquare finds the best square for given size=Nrows * Ncolumns

func BestSquareApprox(size int) (nrow, ncol int)
    BestSquareApprox finds the best square for given size=Nrows * Ncolumns.
    Approximate version; i.e. nrow*ncol may not be equal to size

func Clone(a [][]float64) (b [][]float64)
    Clone allocates and clones a matrix of float64

func Cross3d(w, u, v []float64)
    Cross3d computes the cross product of two 3D vectors u and w

        w = u cross v
        Note: w must be pre-allocated

func CumSum(cs, p []float64)
    CumSum returns the cumulative sum of the elements in p

        Input:
         p -- values
        Output:
         cs -- cumulated sum; pre-allocated with len(cs) == len(p)

func Deep2checkSize(n1, n2 int, a [][]float64) bool
    Deep2checkSize checks if dimensions of Deep2 slice are correct

func Deep2transpose(a [][]float64) (aT [][]float64)
    Deep2transpose returns the transpose of a deep2 slice

func Deep3GetInfo(I, P []int, S []float64, verbose bool) (nitems, nrows, ncolsTot int, ncols []int)
    Deep3GetInfo returns information of serialized array of array of array

func Deep3alloc(n1, n2, n3 int) (a [][][]float64)
    Deep3alloc allocates a slice of slice of slice

func Deep3checkSize(n1, n2, n3 int, a [][][]float64) bool
    Deep3checkSize checks if dimensions of Deep3 slice are correct

func Deep3set(a [][][]float64, v float64)
    Deep3set sets deep slice of slice of slice with v values

func Deep4alloc(n1, n2, n3, n4 int) (a [][][][]float64)
    Deep4alloc allocates a slice of slice of slice of slice

func Deep4set(a [][][][]float64, v float64)
    Deep4set sets deep slice of slice of slice of slice with v values

func DeserializeDeep2(v []float64, m, n int) (a [][]float64)
    DeserializeDeep2 converts a column-major array to a matrix

func DeserializeDeep3(I, P []int, S []float64, debug bool) (A [][][]float64)
    DeserializeDeep3 deserializes an array of array of array in
    column-compressed format

func Digits(maxint int) (ndigits int, format string)
    Digits returns the nubmer of digits

func Dot3d(u, v []float64) (s float64)
    Dot3d returns the dot product between two 3D vectors

func DurSum(v []time.Duration) (seconds float64)
    DurSum sums all seconds in v

        NOTE: this is not efficient and should be used for small slices only

func Expon(val float64) (ndigits int)
    Expon returns the exponent

func Fill(s []float64, val float64)
    Fill fills a slice of float64

func FlipCoin(p float64) bool
    FlipCoin generates a Bernoulli variable; throw a coin with probability p

func FromFloat64s(a []float64) (b []int)
    FromFloat64s returns a new slice of int from a slice of float64

func FromInts(a []int) (b []float64)
    FromInts returns a new slice of float64 from a slice of ints

func FromString(s string) (r []float64)
    FromString splits a string with numbers separeted by spaces into float64

func FromStrings(s []string) (v []float64)
    FromStrings converts a slice of strings to a slice of float64

func GetColumn(j int, v [][]float64) (x []float64)
    GetColumn returns the column of a matrix of float64

func GetCopy(in []float64) (out []float64)
    GetCopy gets a copy of slice of float64

func GetITout(allOutputTimes, timeStationsOut []float64, tol float64) (I []int, T []float64)
    GetITout returns indices and output times

        Input:
          all_output_times  -- array with all output times. ex: [0,0.1,0.2,0.22,0.3,0.4]
          time_stations_out -- time stations for output: ex: [0,0.2,0.4]
          tol               -- tolerance to compare times
        Output:
          iout -- indices of times in all_output_times
          tout -- times corresponding to iout
        Notes:
          use -1 in all_output_times to enforce output of last timestep

func GetMapped(X []float64, filter func(x float64) float64) (Y []float64)
    GetMapped returns a new slice such that: Y[i] = filter(X[i])

func GetMapped2(X [][]float64, filter func(x float64) float64) (Y [][]float64)
    GetMapped2 returns a new slice of slice such that: Y[i][j] = filter(X[i][j])

        NOTE: each row in X may have different number of columns; i.e. len(X[0]) may be != len(X[1])

func GetReversed(in []float64) (out []float64)
    GetReversed return a copy with reversed items

func GetSorted(A []float64) (sortedA []float64)
    GetSorted returns a sorted (increasing) copy of 'A'

func GetStrides(nTotal, nReq int) (I []int)
    GetStrides returns nReq indices from 0 (inclusive) to nTotal (inclusive)

        Input:
          nTotal -- total number of intices
          nReq -- required indices
        Example:
          GetStrides(2001, 5) => [0 400 800 1200 1600 2000 2001]

        NOTE: GetStrides will always include nTotal as the last item in I

func GtPenalty(x, b, penaltyM float64) float64
    GtPenalty implements a 'greater than' penalty function where x must be
    greater than b; otherwise the error is magnified

func GtePenalty(x, b, penaltyM float64) float64
    GtePenalty implements a 'greater than or equal' penalty function where x
    must be greater than b or equal to be; otherwise the error is magnified

func Iabs(a int) int
    Iabs performs the absolute operation with ints

func Imax(a, b int) int
    Imax returns the maximum between two integers

func Imin(a, b int) int
    Imin returns the minimum between two integers

func IntAddScalar(a []int, s int) (res []int)
    IntAddScalar adds a scalar to all values in a slice of integers

func IntAlloc(m, n int) (mat [][]int)
    IntAlloc allocates a matrix of integers

func IntBoolMapSort(m map[int]bool) (sortedKeys []int)
    IntBoolMapSort returns sorted keys of map[int]bool

func IntClone(a [][]int) (b [][]int)
    IntClone allocates and clones a matrix of integers

func IntCopy(in []int) (out []int)
    IntCopy returns a copy of slice of ints

func IntFill(s []int, val int)
    IntFill fills a slice of integers

func IntFilter(a []int, out func(idx int) bool) (res []int)
    IntFilter filters out components in slice

        NOTE: this function is not efficient and should be used with small slices only

func IntGetSorted(A []int) (sortedA []int)
    IntGetSorted returns a sorted (increasing) copy of 'A'

func IntIndexSmall(a []int, val int) int
    IntIndexSmall finds the index of an item in a slice of ints

        NOTE: this function is not efficient and should be used with small slices only; say smaller than 20

func IntIntsMapAppend(m map[int][]int, key int, item int)
    IntIntsMapAppend appends a new item to a map of slice.

        Note: this function creates a new slice in the map if key is not found.

func IntMinMax(v []int) (mi, ma int)
    IntMinMax returns the maximum and minimum elements in v

        NOTE: this is not efficient and should be used for small slices only

func IntNegOut(a []int) (res []int)
    IntNegOut filters out negative components in slice

        NOTE: this function is not efficient and should be used with small slices only

func IntPy(a []int) (res string)
    IntPy returns a Python string representing a slice of integers

func IntRange(n int) (res []int)
    IntRange generates a slice of integers from 0 to n-1

func IntRange2(start, stop int) []int
    IntRange2 generates slice of integers from start to stop (but not stop)

func IntRange3(start, stop, step int) (res []int)
    IntRange3 generates a slice of integers from start to stop (but not stop),
    afer each 'step'

func IntSort3(a, b, c *int)
    IntSort3 sorts 3 values in ascending order

func IntSort4(a, b, c, d *int)
    IntSort4 sort four values in ascending order

func IntUnique(slices ...[]int) (res []int)
    IntUnique returns a unique and sorted slice of integers

func IntVals(n int, val int) (s []int)
    IntVals allocates a slice of integers with size==n, filled with val

func IsPowerOfTwo(n int) bool
    IsPowerOfTwo checks if n is power of 2; i.e. 2⁰, 2¹, 2², 2³, 2⁴, ...

func L2norm(p, q []float64) (dist float64)
    L2norm returns the Euclidean distance between p and q

func LinSpace(start, stop float64, num int) (res []float64)
    LinSpace returns evenly spaced numbers over a specified closed interval.

func LinSpaceOpen(start, stop float64, num int) (res []float64)
    LinSpaceOpen returns evenly spaced numbers over a specified open interval.

func Max(a, b float64) float64
    Max returns the maximum between two float point numbers

func MeshGrid2d(xmin, xmax, ymin, ymax float64, nx, ny int) (X, Y [][]float64)
    MeshGrid2d creates a grid with x-y coordinates

        X, Y -- [ny][nx]

func MeshGrid2dF(xmin, xmax, ymin, ymax float64, nx, ny int, f func(x, y float64) float64) (X, Y, Z [][]float64)
    MeshGrid2dF creates a grid with x-y coordinates and evaluates z=f(x,y)

        X, Y, Z -- [ny][nx]

func MeshGrid2dFG(xmin, xmax, ymin, ymax float64, nx, ny int, fg func(x, y float64) (z, u, v float64)) (X, Y, Z, U, V [][]float64)
    MeshGrid2dFG creates a grid with x-y coordinates and evaluates
    (z,u,v)=fg(x,y)

        X, Y, Z, U, V -- [ny][nx]

func MeshGrid2dV(xVals, yVals []float64) (X, Y [][]float64)
    MeshGrid2dV creates a grid with x-y coordinates given x and y values

        X, Y -- [len(yVals)][len(xVals)]

func Min(a, b float64) float64
    Min returns the minimum between two float point numbers

func MinMax(v []float64) (mi, ma float64)
    MinMax returns the maximum and minimum elements in v

        NOTE: this is not efficient and should be used for small slices only

func NonlinSpace(xa, xb float64, N int, R float64, symmetric bool) (x []float64)
    NonlinSpace generates N points such that the ratio between the last segment
    (or the middle segment) to the first one is equal to a given constant R.

        The ratio between the last (or middle largest) segment to the first one is:

            ΔxL   = Δx0 ⋅ R

        The ratio between successive segments is

                     k-1
            Δx[k] = α    ⋅ Δx[0]

        Unsymmetricc case:

          |--|-------|------------|-----------------|
           Δx0                            ΔxL

        Symmetric case with odd number of spacings:

          |---|--------|---------------|--------|---|
           Δx0                ΔxL                Δx0

        Symmetric case with even number of spacings:

          |---|----------------|----------------|---|
           Δx0       ΔxL               ΔxL       Δx0

func Ones(n int) (res []float64)
    Ones generates a slice of float64 with ones

func ParetoFront(Ovs [][]float64) (front []int)
    ParetoFront computes the Pareto optimal front

        Input:
         Ovs -- [nsamples][ndim] objective values
        Output:
         front -- indices of pareto front
        Note: this function is slow for large sets

func ParetoMin(u, v []float64) (uDominates, vDominates bool)
    ParetoMin compares two vectors using Pareto's optimal criterion

        Note: minimum dominates (is better)

func ParetoMinProb(u, v []float64, φ float64) (uDominates bool)
    ParetoMinProb compares two vectors using Pareto's optimal criterion φ ∃
    [0,1] is a scaling factor that helps v win even if it's not smaller. If
    φ==0, deterministic analysis is carried out. If φ==1, probabilistic analysis
    is carried out. As φ → 1, v "gets more help".

        Note: (1) minimum dominates (is better)
              (2) v dominates if !uDominates

func PrintDeep3(name string, A [][][]float64)
    PrintDeep3 prints an array of array of array

func PrintDeep4(name string, A [][][][]float64, format string)
    PrintDeep4 prints an array of array of array

func PrintMemStat(msg string)
    PrintMemStat prints memory statistics

func ProbContestSmall(u, v, φ float64) float64
    ProbContestSmall computes the probability for a contest between u and v
    where u wins if it's the smaller value. φ ∃ [0,1] is a scaling factor that
    helps v win even if it's not smaller. If φ==0, deterministic analysis is
    carried out. If φ==1, probabilistic analysis is carried out. As φ → 1, v
    "gets more help".

func Prof(mem, silent bool) func()
    Prof runs either CPU profiling or MEM profiling

        INPUT:
          mem    -- memory profiling instead of CPU profiling
          silent -- hide messages

        OUTPUT:
          returns a "stop()" function to be called before exiting
          output files are saved to "/tmp/gosl/cpu.pprof"; or
                                    "/tmp/gosl/mem.pprof"
        Run analysis with:
            go tool pprof <binary-goes-here> /tmp/gosl/cpu.pprof
        or
            go tool pprof <binary-goes-here> /tmp/gosl/mem.pprof

        Example of use (notice the last parentheses):

             func main() {
                defer utl.Prof(false, false)()
                ...
             }

func ProfCPU(dirout, filename string, silent bool) func()
    ProfCPU activates CPU profiling

        OUTPUT: returns a "stop()" function to be called before exiting

        Example of use (notice the last parentheses):

             func main() {
                defer ProfCPU("/tmp", "cpu.pprof", true)()
                ...
             }

        Run analysis with:
            go tool pprof <binary-goes-here> /tmp/cpu.pprof

func ProfMEM(dirout, filename string, silent bool) func()
    ProfMEM activates memory profiling

        OUTPUT: returns a "stop()" function to be called before exiting

        Example of use (notice the last parentheses):

             func main() {
                defer ProfMEM("/tmp", "mem.pprof", true)()
                ...
             }

        Run analysis with:
            go tool pprof <binary-goes-here> /tmp/mem.pprof

func Qsort(arr []float64)
    Qsort sort an array arr[0..n-1] into ascending numerical order using the
    Quicksort algorithm. arr is replaced on output by its sorted rearrangement.
    Normally, the optional argument m should be omitted, but if it is set to a
    positive value, then only the first m elements of arr are sorted.
    Implementation from [1]

        Reference:
        [1] Press WH, Teukolsky SA, Vetterling WT, Fnannery BP (2007) Numerical Recipes: The Art of
            Scientific Computing. Third Edition. Cambridge University Press. 1235p.

func Qsort2(arr, brr []float64)
    Qsort2 sorts an array arr[0..n-1] into ascending order using Quicksort,
    while making the corresponding rearrangment of the array brr[0..n-1].
    Implementation from [1]

        Reference:
        [1] Press WH, Teukolsky SA, Vetterling WT, Fnannery BP (2007) Numerical Recipes: The Art of
            Scientific Computing. Third Edition. Cambridge University Press. 1235p.

func Scaling(s, x []float64, ds, tol float64, reverse, useinds bool) (xmin, xmax float64)
    Scaling computes a scaled version of the input slice with results in [0.0,
    1.0]

        Input:
         x       -- values
         ds      -- δs value to be added to all 's' values
         tol     -- tolerance for capturing xmax ≅ xmin
         reverse -- compute reverse series;
                    i.e. 's' decreases from 1 to 0 while x goes from xmin to xmax
         useinds -- if (xmax-xmin)<tol, use indices to generate the 's' slice;
                    otherwise, 's' will be filled with δs + zeros
        Ouptut:
         s          -- scaled series; pre--allocated with len(s) == len(x)
         xmin, xmax -- min(x) and max(x)

func SerializeDeep2(a [][]float64) (v []float64)
    SerializeDeep2 converts a matrix into a column-major array

func SerializeDeep3(A [][][]float64) (I, P []int, S []float64)
    SerializeDeep3 serializes an array of array of array in column-compressed
    format

func Sort3(a, b, c *float64)
    Sort3 sorts 3 values in ascending order

func Sort3Desc(a, b, c *float64)
    Sort3Desc sorts 3 values in descending order

func Sort4(a, b, c, d *float64)
    Sort4 sort four values in ascending order

func SortQuadruples(i []int, x, y, z []float64, by string) (I []int, X, Y, Z []float64)
    SortQuadruples sorts i, x, y, and z by "i", "x", "y", or "z"

        Note: either i, x, y, or z can be nil; i.e. at least one of them must be non nil

func StrAlloc(m, n int) (mat [][]string)
    StrAlloc allocates a matrix of strings

func StrBoolMapSort(m map[string]bool) (sortedKeys []string)
    StrBoolMapSort returns sorted keys of map[string]bool

func StrBoolMapSortSplit(m map[string]bool) (sortedKeys []string, sortedVals []bool)
    StrBoolMapSortSplit returns sorted keys of map[string]bool and sorted values

func StrFltMapSort(m map[string]float64) (sortedKeys []string)
    StrFltMapSort returns sorted keys of map[string]float64

func StrFltMapSortSplit(m map[string]float64) (sortedKeys []string, sortedVals []float64)
    StrFltMapSortSplit returns sorted keys of map[string]float64 and sorted
    values

func StrFltsMapAppend(m map[string][]float64, key string, item float64)
    StrFltsMapAppend appends a new item to a map of slice.

        Note: this function creates a new slice in the map if key is not found.

func StrIndexSmall(a []string, val string) int
    StrIndexSmall finds the index of an item in a slice of strings

        NOTE: this function is not efficient and should be used with small slices only; say smaller than 20

func StrIntMapSort(m map[string]int) (sortedKeys []string)
    StrIntMapSort returns sorted keys of map[string]int

func StrIntMapSortSplit(m map[string]int) (sortedKeys []string, sortedVals []int)
    StrIntMapSortSplit returns sorted keys of map[string]int and sorted values

func StrIntsMapAppend(m map[string][]int, key string, item int)
    StrIntsMapAppend appends a new item to a map of slice.

        Note: this function creates a new slice in the map if key is not found.

func StrVals(n int, val string) (s []string)
    StrVals allocates a slice of strings with size==n, filled with val

func Sum(v []float64) (sum float64)
    Sum sums all items in v

        NOTE: this is not efficient and should be used for small slices only

func Swap(a, b *float64)
    Swap swaps two float64 numbers

func ToStrings(v []float64, format string) (s []string)
    ToStrings converts a slice of float64 to a slice of strings

func Vals(n int, v float64) (res []float64)
    Vals generates a slice of float64 filled with v


TYPES

type List struct {
	Vals [][]float64 // values
}
    List implements a tabular list with variable number of columns

        Example:
          Vals = [][]float64{
                   {0.0},
                   {1.0, 1.1, 1.2, 1.3},
                   {2.0, 2.1},
                   {3.0, 3.1, 3.2},
                 }

func (o *List) Append(rowidx int, value float64)
    Append appends items to List

type Observable struct {
	// Has unexported fields.
}
    Observable indicates that an object is observable; i.e. it has a list of
    interested observers

func (o *Observable) AddObserver(obs Observer)
    AddObserver adds an object to the list of interested observers

func (o *Observable) NotifyUpdate()
    NotifyUpdate notifies observers of updates

type Observer interface {
	Update() // the data observed by this observer is being update
}
    Observer is an interface to objects that need to observe something

type OutFcnType func(u []float64, t float64)
    OutFcnType is a function that sets the "u" array with an output @ time "t"

type Outputter struct {
	Dt     float64     // time step
	DtOut  float64     // time step for output
	Tmax   float64     // final time (eventually increased to accommodate all time steps)
	Nsteps int         // number of time steps
	Every  int         // increment to output after every time step
	Tidx   int         // the index of the time step for the next output
	Nmax   int         // max number of outputs
	Idx    int         // index in the output arrays for the next output == num.outputs in the end
	T      []float64   // saved t values len = Nmax
	U      [][]float64 // saved u values @ t. [Nmax][nu]
	Fcn    OutFcnType  // function to process output. may be nil
}
    Outputter helps with the output of (numerical) results

        The time loop can be something like:
          t := 0.0
          for tidx := 0; tidx < outper.Nsteps; tidx++ {
              t += outper.Dt
              ... // do something
              outper.MaybeNow(tidx, t)
          }

func NewOutputter(dt, dtOut, tmax float64, nu int, outFcn OutFcnType) (o *Outputter)
    NewOutputter creates a new Outputter

        dt     -- time step
        dtOut  -- time step for output
        tmax   -- final time (of simulation)
        outFcn -- callback function to perform output (may be nil)
        NOTE: a first output will be processed if outFcn != nil

func (o *Outputter) MaybeNow(tidx int, t float64)
    MaybeNow process the output if tidx == NextIdxT

type P struct {

	// input
	N      string  `json:"n"`      // name of parameter
	V      float64 `json:"v"`      // value of parameter
	Min    float64 `json:"min"`    // min value
	Max    float64 `json:"max"`    // max value
	S      float64 `json:"s"`      // standard deviation
	D      string  `json:"d"`      // probability distribution type
	U      string  `json:"u"`      // unit (not verified)
	Adj    int     `json:"adj"`    // adjustable: unique ID (greater than zero)
	Dep    int     `json:"dep"`    // depends on "adj"
	Extra  string  `json:"extra"`  // extra data
	Inact  bool    `json:"inact"`  // parameter is inactive in optimisation
	SetDef bool    `json:"setdef"` // tells model to use a default value

	// connected parameter
	Other *P // dependency: connected parameter

	// function
	Func func(t float64, x []float64) // a function y=f(t,x)

	// Has unexported fields.
}
    P holds numeric parameters defined by a name N and a value V.

    P is convenient to store the range of allowed values in Min and Max, and
    other information such as standard deviation S, probability distribution
    type D, among others.

    Dependent variables may be connected to P using Connect so when Set is
    called, the dependendt variable is updated as well.

    Other parameters can be linked to this one via the Other data member and
    Func may be useful to compute y=f(t,x)

func (o *P) Connect(V *float64)
    Connect connects parameter to variable

func (o *P) Set(V float64)
    Set sets parameter, including connected variables

type Params []*P
    Params holds many parameters

        A set of Params can be initialized as follows:

          var params Params
          params = []*P{
              {N: "klx", V: 1.0},
              {N: "kly", V: 2.0},
              {N: "klz", V: 3.0},
          }

        Alternatively, see NewParams function

func NewParams(pp ...interface{}) (o Params)
    NewParams returns a set of parameters

        This is an alternative to initializing Params by setting slice items

        A set of Params can be initialized as follows:

          params := NewParams(
              &P{N: "P1", V: 1},
              &P{N: "P2", V: 2},
              &P{N: "P3", V: 3},
          )

        Alternatively, you may set slice components directly (see Params definition)

func (o *Params) CheckAndGetValues(names []string) (values []float64)
    CheckAndGetValues check min/max limits and return values. Will panic if
    values are outside corresponding min/max range. Will also panic if a
    parameter name is not found.

func (o *Params) CheckAndSetVariables(names []string, variables []*float64)
    CheckAndSetVariables get parameter values and check limits defined in Min
    and Max Will panic if values are outside corresponding Min/Max range. Will
    also panic if a parameter name is not found.

func (o *Params) CheckLimits()
    CheckLimits check limits of variables given in Min/Max Will panic if values
    are outside corresponding Min/Max range.

func (o *Params) Connect(V *float64, name, caller string) (errorMessage string)
    Connect connects parameter

func (o *Params) ConnectSet(V []*float64, names []string, caller string) (errorMessage string)
    ConnectSet connects set of parameters

func (o *Params) ConnectSetOpt(V []*float64, names []string, optional []bool, caller string) (errorMessage string)
    ConnectSetOpt connects set of parameters with some being optional

func (o *Params) Find(name string) *P
    Find finds a parameter by name

        Note: returns nil if not found

func (o *Params) GetBool(name string) bool
    GetBool reads Boolean parameter or Panic Returns true if P[name] > 0;
    otherwise returns false Will panic if name does not exist in parameters set

func (o *Params) GetBoolOrDefault(name string, defaultValue bool) bool
    GetBoolOrDefault reads Boolean parameter or returns default value Returns
    true if P[name] > 0; otherwise returns false Will return defaultValue if
    name does not exist in parameters set

func (o *Params) GetIntOrDefault(name string, defaultInt int) int
    GetIntOrDefault reads parameter or returns default value Will return
    defaultInt if name does not exist in parameters set

func (o *Params) GetValue(name string) float64
    GetValue reads parameter or Panic Will panic if name does not exist in
    parameters set

func (o *Params) GetValueOrDefault(name string, defaultValue float64) float64
    GetValueOrDefault reads parameter or returns default value Will return
    defaultValue if name does not exist in parameters set

func (o *Params) GetValues(names []string) (values []float64, found []bool)
    GetValues get parameter values

func (o *Params) SetBool(name string, value float64)
    SetBool sets Boolean parameter or Panic Sets +1==true if value > 0;
    otherwise sets -1==false Will panic if name does not exist in parameters set

func (o *Params) SetValue(name string, value float64)
    SetValue sets parameter or Panic Will panic if name does not exist in
    parameters set

func (o Params) String() (l string)
    String returns a summary of parameters

type Quadruple struct {
	I int
	X float64
	Y float64
	Z float64
}
    Quadruple helps to sort a quadruple of 1 int and 3 float64s

type Quadruples []*Quadruple
    Quadruples helps to sort quadruples

func BuildQuadruples(i []int, x, y, z []float64) (q Quadruples)
    BuildQuadruples initialise Quadruples with i, x, y, and z

        Note: i, x, y, or z can be nil; but at least one of them must be non nil

func (o Quadruples) I() (i []int)
    I returns the 'i' in quadruples

func (o Quadruples) Len() int
    Len returns the length of Quadruples

func (o Quadruples) String() string
    String returns the string representation of Quadruples

func (o Quadruples) Swap(i, j int)
    Swap swaps two quadruples

func (o Quadruples) X() (x []float64)
    X returns the 'x' in quadruples

func (o Quadruples) Y() (y []float64)
    Y returns the 'y' in quadruples

func (o Quadruples) Z() (z []float64)
    Z returns the 'z' in quadruples

type QuadruplesByI struct{ Quadruples }
    QuadruplesByI defines struct to sort Quadruples by I

func (o QuadruplesByI) Less(i, j int) bool
    Less compares two QuadruplesI

type QuadruplesByX struct{ Quadruples }
    QuadruplesByX defines struct to sort Quadruples by X

func (o QuadruplesByX) Less(i, j int) bool
    Less compares two QuadruplesX

type QuadruplesByY struct{ Quadruples }
    QuadruplesByY Sort Quadruples by Y

func (o QuadruplesByY) Less(i, j int) bool
    Less compares two QuadruplesY

type QuadruplesByZ struct{ Quadruples }
    QuadruplesByZ defines struct to Sort Quadruples by Z

func (o QuadruplesByZ) Less(i, j int) bool
    Less compares two QuadruplesZ

type SerialList struct {
	Vals []float64 // values
	Ptrs []int     // pointers
}
    SerialList implements a tabular list with variable number of columns using a
    serial representation

        Example:
            0.0
            1.0  1.1  1.2  1.3
            2.0  2.1
            3.0  3.1  3.2
        becomes:
                   (0)   (1)    2    3    4   (5)    6   (7)    8    9   (10)
            Vals = 0.0 | 1.0  1.1  1.2  1.3 | 2.0  2.1 | 3.0  3.1  3.2 |
            Ptrs = 0 1 5 7 10
        Notes:
            len(Ptrs) = nrows + 1
            Ptrs[len(Ptrs)-1] = len(Vals)

func (o *SerialList) Append(startRow bool, value float64)
    Append appends item to SerialList

func (o SerialList) Print(fmt string)
    Print prints the souble-serial-list

type Sorter struct {
	Index []int
}
    Sorter builds an index list to sort arrays of any type.

        Reference:
        [1] Press WH, Teukolsky SA, Vetterling WT, Fnannery BP (2007) Numerical Recipes: The Art of
            Scientific Computing. Third Edition. Cambridge University Press. 1235p.

func (o *Sorter) GetSorted(a []float64) (b []float64)
    GetSorted returns a copy of array 'a' sorted according to the previously
    built index. NOTE: the copy may be smaller if the index was built with a
    smaller set

func (o *Sorter) GetSortedI(a []int) (b []int)
    GetSortedI returns a copy of array 'a' sorted according to the previously
    built index. NOTE: the copy may be smaller if the index was built with a
    smaller set

func (o *Sorter) Init(n int, less func(i, j int) bool)
    Init builds an index indx[0..n-1] to sort an array a[0..n-1] such that
    a[indx[j]] is in ascending order for j=0,1,...,n-1.

        Input:
          n    -- number of items in the array to be sorted. n must be ≤ len(a)
          less -- a function that returns true if a[i] < a[j]
        Reference:
        [1] Press WH, Teukolsky SA, Vetterling WT, Fnannery BP (2007) Numerical Recipes: The Art of
            Scientific Computing. Third Edition. Cambridge University Press. 1235p.

```
