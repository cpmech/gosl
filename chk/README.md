# Gosl. chk. Check code and unit test tools

Package `chk` provides tools to check numerical results and to perform unit tests.

## API

**go doc**

```
package chk // import "gosl/chk"

Package chk contains functions for checking and testing computations

VARIABLES

var (
	// AssertOn activates or deactivates asserts
	AssertOn = true

	// Verbose turn on verbose mode
	Verbose = false

	// ColorsOn turn on use of colours on console
	ColorsOn = true
)

FUNCTIONS

func AnaNum(tst *testing.T, msg string, tol, ana, num float64, verbose bool)
    AnaNum compares analytical versus numerical values

func AnaNumC(tst *testing.T, msg string, tol float64, ana, num complex128, verbose bool)
    AnaNumC compares analytical versus numerical values (complex version)

func Array(tst *testing.T, msg string, tol float64, a, b []float64)
    Array compares two array. The b slice may be nil indicating that all values
    are zero

func ArrayC(tst *testing.T, msg string, tol float64, a, b []complex128)
    ArrayC compares two slices of complex nummber. The b slice may be nil
    indicating that all values are zero

func Bool(tst *testing.T, msg string, a, b bool)
    Bool compares two bools

func Bools(tst *testing.T, msg string, a, b []bool)
    Bools compare two slices of bool. The b slice may be nil indicating that all
    values are false

func CallerInfo(idx int)
    CallerInfo returns the file and line positions where an error occurred

        idx -- use idx=2 to get the caller of Panic

func Complex128(tst *testing.T, msg string, tolNorm float64, a, b complex128)
    Complex128 compares two complex128 numbers

func Deep2(tst *testing.T, msg string, tol float64, a, b [][]float64)
    Deep2 compares two nested (depth=2) slice. The b slice may be nil indicating
    that all values are zero

func Deep2c(tst *testing.T, msg string, tol float64, a, b [][]complex128)
    Deep2c compares two nested (depth=2) slices. The b slice may be nil
    indicating that all values are zero

func Deep3(tst *testing.T, msg string, tol float64, a, b [][][]float64)
    Deep3 compares two deep3 slices. The b slice may be nil indicating that all
    values are zero

func Deep4(tst *testing.T, msg string, tol float64, a, b [][][][]float64)
    Deep4 compares two deep4 slices. The b slice may be nil indicating that all
    values are zero

func DerivScaSca(tst *testing.T, msg string, tol, gAna, xAt, h float64, verb bool, fcn func(x float64) float64)
    DerivScaSca checks the derivative of scalar w.r.t scalar by comparing with
    numerical solution obtained with central differences (5-point rule)

        Checks:
                  df │
              g = —— │      with   f:scalar,  x:scalar
                  dx │xAt          g:scalar
        Input:
          tst  -- testing.T structure
          msg  -- message about this test
          tol  -- tolerance to compare gAna with gNum
          gAna -- [scalar] analytical (or other kind) derivative dfdx
          xAt  -- [scalar] position to compute dfdx
          h    -- initial stepsize; e.g. 1e-1
          verb -- verbose: show messages
          fcn  -- [scalar] function f(x). x is scalar

func DerivScaVec(tst *testing.T, msg string, tol float64, gAna, xAt []float64, h float64,
	verb bool, fcn func(x []float64) float64)
    DerivScaVec checks the derivative of scalar w.r.t vector by comparing with
    numerical solution obtained with central differences (5-point rule)

        Check:
                    df  │               f:scalar   {x}:vector
             {g} = ———— │        with   {g}:vector
                   d{x} │{xAt}          len(g) == len(x) == len(xAt)
        Input:
          tst  -- testing.T structure
          msg  -- message about this test
          tol  -- tolerance to compare gAna with gNum
          gAna -- [vector] analytical (or other kind) derivative dfdx. size=len(x)=len(xAt)
          xAt  -- [vector] position to compute dfdx
          h    -- initial stepsize; e.g. 1e-1
          verb -- verbose: show messages
          fcn  -- [scalar] function f(x). x is vector

func DerivVecSca(tst *testing.T, msg string, tol float64, gAna []float64, xAt, h float64,
	verb bool, fcn func(f []float64, x float64))
    DerivVecSca checks the derivative of vector w.r.t scalar by comparing with
    numerical solution obtained with central differences (5-point rule)

        Check:
                   d{f} │             {f}:vector   x:scalar
             {g} = ———— │      with   {g}:vector
                    dx  │xAt          len(g) == len(f)
        Input:
          tst  -- testing.T structure
          msg  -- message about this test
          tol  -- tolerance to compare gAna with gNum
          gAna -- [vector] analytical (or other kind) derivative dfdx. size=len(f)
          xAt  -- [scalar] position to compute dfdx
          h    -- initial stepsize; e.g. 1e-1
          verb -- verbose: show messages
          fcn  -- [vector] function f(x). x is scalar

func DerivVecVec(tst *testing.T, msg string, tol float64, gAna [][]float64, xAt []float64, h float64,
	verb bool, fcn func(f, x []float64))
    DerivVecVec checks the derivative of vector w.r.t vector by comparing with
    numerical solution obtained with central differences (5-point rule)

        Checks:
                   d{f} │               {f}:vector   {x}:vector
             [g] = ———— │        with   [g]:matrix
                   d{x} │{xAt}          rows(g)==len(f)  cols(g)==len(x)==len(xAt)
        Input:
          tst  -- testing.T structure
          msg  -- message about this test
          tol  -- tolerance to compare gAna with gNum
          gAna -- [matrix] analytical (or other kind) derivative dfdx. size=(len(f),len(x))
          xAt  -- [vector] position to compute dfdx
          h    -- initial stepsize; e.g. 1e-1
          verb -- verbose: show messages
          fcn  -- [vector] function f(x). x is vector

func Err(msg string, prm ...interface{}) error
    Err returns a new error

func Float64(tst *testing.T, msg string, tol, a, b float64)
    Float64 compares two float64 numbers

func Float64assert(a, b float64)
    Float64assert asserts that a is equal to b (floats)

func Int(tst *testing.T, msg string, a, b int)
    Int compares two ints

func Int32(tst *testing.T, msg string, a, b int32)
    Int32 compares two int32

func Int32s(tst *testing.T, msg string, a, b []int32)
    Int32s compares two slices of 32 integer. The b slice may be nil indicating
    that all values are zero

func Int64(tst *testing.T, msg string, a, b int64)
    Int64 compares two int64

func Int64s(tst *testing.T, msg string, a, b []int64)
    Int64s compares two slices of 64 integer. The b slice may be nil indicating
    that all values are zero

func IntAssert(a, b int)
    IntAssert asserts that a is equal to b (ints)

func IntAssertLessThan(a, b int)
    IntAssertLessThan asserts that a < b (ints)

func IntAssertLessThanOrEqualTo(a, b int)
    IntAssertLessThanOrEqualTo asserts that a ≤ b (ints)

func IntDeep2(tst *testing.T, msg string, a, b [][]int)
    IntDeep2 compares nested slices of ints. The b slice may be nil indicating
    that all values are zero

func Ints(tst *testing.T, msg string, a, b []int)
    Ints compares two slices of integer. The b slice may be nil indicating that
    all values are zero

func Panic(msg string, prm ...interface{})
    Panic calls CallerInfo and panicks

func PanicSimple(msg string, prm ...interface{})
    PanicSimple panicks without calling CallerInfo

func PrintAnaNum(msg string, tol, ana, num float64, verbose bool) (e error)
    PrintAnaNum formats the output of analytical versus numerical comparisons

func PrintAnaNumC(msg string, tol float64, ana, num complex128, verbose bool) (e error)
    PrintAnaNumC formats the output of analytical versus numerical comparisons
    (complex version)

func PrintOk(msg string, prm ...interface{})
    PrintOk prints "OK" in green (if ColorsOn==true)

func PrintTitle(title string)
    PrintTitle returns the Test Title

func Recover()
    Recover catches panics and call os.Exit(1) on 'panic'

func RecoverTst(tst *testing.T)
    RecoverTst catches panics in tests. Test will fail on 'panic'

func RecoverTstPanicIsOK(tst *testing.T)
    RecoverTstPanicIsOK catches panics in tests. Test must 'panic' to be OK

func StrAssert(a, b string)
    StrAssert asserts that a is equal to b (strings)

func StrDeep2(tst *testing.T, msg string, a, b [][]string)
    StrDeep2 compares nested slices of strings. The b slice may be nil
    indicating that all values are zero

func String(tst *testing.T, a, b string)
    String compares two strings

func Strings(tst *testing.T, msg string, a, b []string)
    Strings compare two slices of string. The b slice may be nil indicating that
    all values are "" (empty)

func Symmetry(tst *testing.T, msg string, X []float64)
    Symmetry checks symmetry of SEGMENTS in an even or odd slice of float64

        NOTE: values in X must be sorted ascending

func TestDiffC(tst *testing.T, msg string, tol float64, a, b complex128, showOK bool) (failed bool)
    TestDiffC tests difference between complex128. It also prints "FAIL" or "OK"

func TstDiff(tst *testing.T, msg string, tol, a, b float64, showOK bool) (failed bool)
    TstDiff tests difference between float64

func TstFail(tst *testing.T, msg string, prm ...interface{})
    TstFail calls tst.Errorf() with msg and parameters. It also prints "FAIL" in
    red (if ColorsOn==true)

```
