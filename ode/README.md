# Gosl. ode. Ordinary differential equations

Package `ode` implements solution techniques to ordinary differential equations, such as the
Runge-Kutta method. Methods that can handle stiff problems are also available.

## Examples

### Robertson's Equation

From Hairer-Wanner VII-p3 Eq.(1.4) [2].

<div id="container">
<p><img src="../examples/figs/rober.png" width="400"></p>
Solution of Robertson's equation
</div>

## Output of Tests

### Convergence of explicit Runge-Kutta methods

Source code: <a href="t_erk_test.go">t_erk_test.go</a>

<div id="container">
<p><img src="../examples/figs/t_erk04.png" width="500"></p>
</div>

<div id="container">
<p><img src="../examples/figs/t_erk05.png" width="500"></p>
</div>

## References

[1] Hairer E, Norset SP, Wanner G. Solving O. Solving Ordinary Differential Equations I. Nonstiff
Problems, Springer. 1987

[2] Hairer E, Wanner G. Solving Ordinary Differential Equations II. Stiff and Differential-Algebraic
Problems, Second Revision Edition. Springer. 1996

## API

**go doc**

```
package ode // import "gosl/ode"

Package ode implements solvers for ordinary differential equations,
including explicit and implicit Runge-Kutta methods; e.g. the fantastic
Radau5 method by Hairer, Norsett & Wanner [1, 2].

    References:
      [1] Hairer E, Nørsett SP, Wanner G (1993). Solving Ordinary Differential Equations I:
          Nonstiff Problems. Springer Series in Computational Mathematics, Vol. 8, Berlin,
          Germany, 523 p.
      [2] Hairer E, Wanner G (1996). Solving Ordinary Differential Equations II: Stiff and
          Differential-Algebraic Problems. Springer Series in Computational Mathematics,
          Vol. 14, Berlin, Germany, 614 p.

FUNCTIONS

func Dopri5simple(fcn Func, y la.Vector, xf, tol float64)
    Dopri5simple solves ODE using DoPri5 method without options for saving
    results and others

func Dopri8simple(fcn Func, y la.Vector, xf, tol float64)
    Dopri8simple solves ODE using DoPri8 method without options for saving
    results and others

func Radau5simple(fcn Func, jac JacF, y la.Vector, xf, tol float64)
    Radau5simple solves ODE using Radau5 method without options for saving
    results and others

func Solve(method string, fcn Func, jac JacF, y la.Vector, xf, dx, atol, rtol float64,
	numJac, fixedStep, saveStep, saveDense bool) (stat *Stat, out *Output)
    Solve solves ODE problem using standard parameters

        INPUT:
         method    -- the method
         fcn       -- function d{y}/dx := {f}(h=dx, x, {y})
         jac       -- Jacobian function d{f}/d{y} := [J](h=dx, x, {y}) [may be nil]
         y         -- current {y} @ x=0
         xf        -- final x
         dx        -- stepsize. [may be used for dense output]
         atol      -- absolute tolerance; use 0 for default [default = 1e-4] (for fixedStp=false)
         rtol      -- relative tolerance; use 0 for default [default = 1e-4] (for fixedStp=false)
         numJac    -- use numerical Jacobian if if jac is non nil
         fixedStep -- fixed steps
         saveStep  -- save steps
         saveDense -- save many steps (dense output) [using dx]

        OUTPUT:
         y    -- modified y vector with final {y}
         stat -- statistics
         out  -- output with all steps results with save==true


TYPES

type BwEuler struct {
	// Has unexported fields.
}
    BwEuler implements the (implicit) Backward Euler method

func (o *BwEuler) Accept(y0 la.Vector, x0 float64) (dxnew float64)
    Accept accepts update

func (o *BwEuler) DenseOut(yout la.Vector, h, x float64, y la.Vector, xout float64)
    DenseOut produces dense output (after Accept)

func (o *BwEuler) Free()
    Free releases memory

func (o *BwEuler) Info() (fixedOnly, implicit bool, nstages int)
    Info returns information about this method

func (o *BwEuler) Init(ndim int, conf *Config, work *rkwork, stat *Stat, fcn Func, jac JacF, M *la.Triplet)
    Init initialises structure

func (o *BwEuler) Reject() (dxnew float64)
    Reject processes step rejection

func (o *BwEuler) Step(x0 float64, y0 la.Vector)
    Step steps update

type Config struct {

	// parameters
	Hmin       float64 // minimum H allowed
	IniH       float64 // initial H
	NmaxIt     int     // max num iterations (allowed)
	NmaxSS     int     // max num substeps
	Mmin       float64 // min step multiplier
	Mmax       float64 // max step multiplier
	Mfac       float64 // step multiplier factor
	MfirstRej  float64 // coefficient to multiply stepsize if first step is rejected [0 ⇒ use dxnew]
	PredCtrl   bool    // use Gustafsson's predictive controller
	Eps        float64 // smallest number satisfying 1.0 + ϵ > 1.0
	ThetaMax   float64 // max theta to decide whether the Jacobian should be recomputed or not
	C1h        float64 // c1 of HW-VII p124 => min ratio to retain previous h
	C2h        float64 // c2 of HW-VII p124 => max ratio to retain previous h
	LerrStrat  int     // strategy to select local error computation method
	GoChan     bool    // allow use of go channels (threaded); e.g. to solve R and C systems concurrently
	CteTg      bool    // use constant tangent (Jacobian) in BwEuler
	UseRmsNorm bool    // use RMS norm instead of Euclidean in BwEuler
	Verbose    bool    // show messages, e.g. during iterations
	ZeroTrial  bool    // always start iterations with zero trial values (instead of collocation interpolation)
	StabBeta   float64 // Lund stabilisation coefficient β

	// stiffness detection
	StiffNstp  int     // number of steps to check stiff situation. 0 ⇒ no check. [default = 1]
	StiffRsMax float64 // maximum value of ρs [default = 0.5]
	StiffNyes  int     // number of "yes" stiff steps allowed [default = 15]
	StiffNnot  int     // number of "not" stiff steps to disregard stiffness [default = 6]

	// configurations for linear solver
	LinSolConfig *la.SparseConfig // configurations for sparse linear solver

	// Has unexported fields.
}
    Config holds configuration parameters for the ODE solver

func NewConfig(method string, lsKind string, comm *mpi.Communicator) (o *Config)
    NewConfig returns a new [default] set of configuration parameters

        method -- the ODE method: e.g. fweuler, bweuler, radau5, moeuler, dopri5
        lsKind -- kind of linear solver: "umfpack" or "mumps" [may be empty]
        comm   -- communicator for the linear solver [may be nil]
        NOTE: (1) if comm == nil, the linear solver will be "umfpack" by default
              (2) if comm != nil and comm.Size() == 1, you can use either "umfpack" or "mumps"
              (3) if comm != nil and comm.Size() > 1, the linear solver will be set to "mumps" automatically

func (o *Config) SetDenseOut(save bool, dxOut, xf float64, out DenseOutF)
    SetDenseOut activates dense output

        save -- save all values
        out  -- function to be during dense output [may be nil]

func (o *Config) SetFixedH(dxApprox, xf float64)
    SetFixedH calculates the number of steps, the exact stepsize h, and set to
    use fixed stepsize

func (o *Config) SetStepOut(save bool, out StepOutF)
    SetStepOut activates output of (variable) steps

        save -- save all values
        out  -- function to be during step output [may be nil]

func (o *Config) SetTol(atolAndRtol float64)
    SetTol sets both tolerances: Atol and Rtol

func (o *Config) SetTols(atol, rtol float64)
    SetTols sets tolerances according to Hairer and Wanner' suggestions

        atol   -- absolute tolerance; use 0 for default [default = 1e-4]
        rtol   -- relative tolerance; use 0 for default [default = 1e-4]

type DenseOutF func(istep int, h, x float64, y la.Vector, xout float64, yout la.Vector) (stop bool)
    DenseOutF defines a function to produce a dense output (i.e. many equally
    spaced points, regardless of the actual stepsize)

        INPUT:
          istep -- index of step (0 is the very first output whereas 1 is the first accepted step)
          h     -- best (current) h
          x     -- current (just updated) x
          y     -- current (just updated) y
          xout  -- selected x to produce an output
          yout  -- y values computed @ xout

        OUTPUT:
          stop -- stop simulation (nicely)

type ExplicitRK struct {

	// constants
	FSAL     bool        // can use previous ks to compute k0; i.e. k0 := ks{previous]. first same as last [1, page 167]
	Embedded bool        // has embedded error estimator
	A        [][]float64 // A coefficients
	B        []float64   // B coefficients
	Be       []float64   // B coefficients [may be nil, e.g. if FSAL = false]
	C        []float64   // C coefficients
	E        []float64   // error coefficients. difference between B and Be: e = b - be (if be is not nil)
	Nstg     int         // number of stages = len(A) = len(B) = len(C)
	P        int         // order of y1 (corresponding to b)
	Q        int         // order of error estimator (embedded only); e.g. DoPri5(4) ⇒ q = 4 (=min(order(y1),order(y1bar))
	Ad       [][]float64 // A coefficients for dense output
	Cd       []float64   // C coefficients for dense output
	D        [][]float64 // dense output coefficients. [may be nil]

	// Has unexported fields.
}
    ExplicitRK implements explicit Runge-Kutta methods

         The methods available are:
           moeuler    -- 2(1) Modified-Euler 2(1) ⇒ q = 1
           rk2        -- 2 Runge, order 2 (mid-point). page 135 of [1]
           rk3        -- 3 Runge, order 3. page 135 of [1]
           heun3      -- 3 Heun, order 3. page 135 of [1]
           rk4        -- 4 "The" Runge-Kutta method. page 138 of [1]
           rk4-3/8    -- 4 Runge-Kutta method: 3/8-Rule. page 138 of [1]
           merson4    -- 4 Merson 4("5") method. "5" means that the order 5 is for linear equations with constant coefficients; otherwise the method is of order3. page 167 of [1]
           zonneveld4 -- 4 Zonneveld 4(3). page 167 of [1]
           fehlberg4  -- 4(5) Fehlberg 4(5) ⇒ q = 4
           dopri5     -- 5(4) Dormand-Prince 5(4) ⇒ q = 4
           verner6    -- 6(5) Verner 6(5) ⇒ q = 5
           fehlberg7  -- 7(8) Fehlberg 7(8) ⇒ q = 7
           dopri8     -- 8(5,3) Dormand-Prince 8 order with 5,3 estimator
        where p(q) means method of p-order with embedded estimator of q-order

        References:
          [1] E. Hairer, S. P. Nørsett, G. Wanner (2008) Solving Ordinary Differential Equations I.
              Nonstiff Problems. Second Revised Edition. Corrected 3rd printing 2008. Springer Series
              in Computational Mathematics ISSN 0179-3632, 528p

        NOTE: (1) Fehlberg's methods give identically zero error estimates for quadrature problems
                  y'=f(x); see page 180 of [1]

func (o *ExplicitRK) Accept(y0 la.Vector, x0 float64) (dxnew float64)
    Accept accepts update and computes next stepsize

func (o *ExplicitRK) DenseOut(yout la.Vector, h, x float64, y la.Vector, xout float64)
    DenseOut produces dense output (after Accept)

func (o *ExplicitRK) Free()
    Free releases memory

func (o *ExplicitRK) Info() (fixedOnly, implicit bool, nstages int)
    Info returns information about this method

func (o *ExplicitRK) Init(ndim int, conf *Config, work *rkwork, stat *Stat, fcn Func, jac JacF, M *la.Triplet)
    Init initialises structure

func (o *ExplicitRK) Reject() (dxnew float64)
    Reject processes step rejection and computes next stepsize

func (o *ExplicitRK) Step(xa float64, ya la.Vector)
    Step steps update

type Func func(f la.Vector, h, x float64, y la.Vector)
    Func defines the main function d{y}/dx = {f}(x, {y})

        Here, the "main" function receives the stepsize h as well, i.e.

          d{y}/dx := {f}(h=dx, x, {y})

        INPUT:
          h -- current stepsize = dx
          x -- current x
          y -- current {y}

        OUTPUT:
          f -- {f}(h, x, {y})

type FwEuler struct {
	// Has unexported fields.
}
    FwEuler implements the (explicit) Forward Euler method

func (o *FwEuler) Accept(y0 la.Vector, x0 float64) (dxnew float64)
    Accept accepts update

func (o *FwEuler) DenseOut(yout la.Vector, h, x float64, y la.Vector, xout float64)
    DenseOut produces dense output (after Accept)

func (o *FwEuler) Free()
    Free releases memory

func (o *FwEuler) Info() (fixedOnly, implicit bool, nstages int)
    Info returns information about this method

func (o *FwEuler) Init(ndim int, conf *Config, work *rkwork, stat *Stat, fcn Func, jac JacF, M *la.Triplet)
    Init initialises structure

func (o *FwEuler) Reject() (dxnew float64)
    Reject processes step rejection

func (o *FwEuler) Step(x0 float64, y0 la.Vector)
    Step steps update

type JacF func(dfdy *la.Triplet, h, x float64, y la.Vector)
    JacF defines the Jacobian matrix of Func

        Here, the Jacobian function receives the stepsize h as well, i.e.

        d{f}/d{y} := [J](h=dx, x, {y})

        INPUT:
          h -- current stepsize = dx
          x -- current x
          y -- current {y}

        OUTPUT:
          dfdy -- Jacobian matrix d{f}/d{y} := [J](h=dx, x, {y})

type Output struct {

	// discrete output at accepted steps
	StepIdx int         // current index in Xvalues and Yvalues == last output
	StepRS  []float64   // ρ{stiffness} ratio [IdxSave]
	StepH   []float64   // h values [IdxSave]
	StepX   []float64   // X values [IdxSave]
	StepY   []la.Vector // Y values [IdxSave][ndim]

	// dense output
	DenseIdx int         // current index in DenseS, DenseX, and DenseY arrays
	DenseS   []int       // index of step
	DenseX   []float64   // X values during dense output [DenseIdx]
	DenseY   []la.Vector // Y values during dense output [DenseIdx][ndim]

	// Has unexported fields.
}
    Output holds output data (prepared by Solver)

func (o *Output) GetDenseS() (S []int)
    GetDenseS returns all s (step-index) values from the dense output data

func (o *Output) GetDenseX() (X []float64)
    GetDenseX returns all x values from the dense output data

func (o *Output) GetDenseY(i int) (Y []float64)
    GetDenseY extracts the y[i] values for all output times from the dense
    output data

        i -- index of y component
        use to plot time series; e.g.:
           plt.Plot(o.GetDenseX(), o.GetDenseY(0), &plt.A{L:"y0"})

func (o *Output) GetDenseYtable() (Y [][]float64)
    GetDenseYtable returns a table with all y values such that Y[idxOut][dim]
    from the dense output data

func (o *Output) GetDenseYtableT() (Y [][]float64)
    GetDenseYtableT returns a (transposed) table with all y values such that
    Y[dim][idxOut] from the dense output data

func (o *Output) GetStepH() (X []float64)
    GetStepH returns all h values from the (accepted) steps output

func (o *Output) GetStepRs() (X []float64)
    GetStepRs returns all ρs (stiffness ratio) values from the (accepted) steps
    output

func (o *Output) GetStepX() (X []float64)
    GetStepX returns all x values from the (accepted) steps output

func (o *Output) GetStepY(i int) (Y []float64)
    GetStepY extracts the y[i] values for all output times from the (accepted)
    steps output

        i -- index of y component
        use to plot time series; e.g.:
           plt.Plot(o.GetStepX(), o.GetStepY(0), &plt.A{L:"y0"})

func (o *Output) GetStepYtable() (Y [][]float64)
    GetStepYtable returns a table with all y values such that Y[idxOut][dim]
    from the (accepted) steps output

func (o *Output) GetStepYtableT() (Y [][]float64)
    GetStepYtableT returns a (transposed) table with all y values such that
    Y[dim][idxOut] from the (accepted) steps output

type Problem struct {
	Yana YanaF       // analytical solution
	Fcn  Func        // the function f(x,y) = dy/df
	Jac  JacF        // df/dy function
	Dx   float64     // timestep for fixedStep tests
	Xf   float64     // final x
	Y    la.Vector   // initial (current) y vector
	Ndim int         // dimension == len(Y)
	M    *la.Triplet // "mass" matrix
	Ytmp la.Vector   // to use with Yana
}
    Problem defines the data for an ODE problems (e.g. for testing)

func ProbArenstorf() (o *Problem)
    ProbArenstorf returns the Arenstorf orbit problem

func ProbHwAmplifier() (o *Problem)
    ProbHwAmplifier returns the Hairer-Wanner VII-p376 Transistor Amplifier
    NOTE: from E. Hairer's website, not the equation in the book

func ProbHwEq11() (o *Problem)
    ProbHwEq11 returns the Hairer-Wanner problem from VII-p2 Eq.(1.1)

func ProbRobertson() (o *Problem)
    ProbRobertson returns the Robertson's Equation as given in Hairer-Wanner
    VII-p3 Eq.(1.4)

func ProbSimpleNdim2() (o *Problem)
    ProbSimpleNdim2 returns a simple 2-dim problem

func ProbSimpleNdim4a() (o *Problem)
    ProbSimpleNdim4a returns a simple 4-dim problem (a)

func ProbSimpleNdim4b() (o *Problem)
    ProbSimpleNdim4b returns a simple 4-dim problem (b)

func ProbVanDerPol(eps float64, stationary bool) (o *Problem)
    ProbVanDerPol returns the Van der Pol' Equation as given in Hairer-Wanner
    VII-p5 Eq.(1.5)

        eps  -- ε coefficient; use 0 for default value [=1e-6]
        stationary -- use eps=1 and compute period and amplitude such that
                      y = [A, 0] is a stationary point

func (o *Problem) CalcYana(idxY int, x float64) float64
    CalcYana computes component idxY of analytical solution @ x, if available

func (o *Problem) ConvergenceTest(tst *testing.T, dxmin, dxmax float64, ndx int, yExact la.Vector,
	methods []string, orders, tols []float64)
    ConvergenceTest runs convergence test

        yExact -- is the exact (reference) y @ xf

func (o *Problem) Solve(method string, fixedStp, numJac bool) (y la.Vector, stat *Stat, out *Output)
    Solve solves ODE problem using standard parameters NOTE: this solver doesn't
    change o.Y

type Radau5 struct {

	// constants
	C    []float64   // c coefficients
	T    [][]float64 // T matrix
	Ti   [][]float64 // inv(T) matrix
	Alp  float64     // alpha-hat
	Bet  float64     // beta-hat
	Gam  float64     // gamma-hat
	Gam0 float64     // gamma0 coefficient
	E0   float64     // e0 coefficient
	E1   float64     // e1 coefficient
	E2   float64     // e2 coefficient
	Mu1  float64     // collocation: C1    = (4.D0-SQ6)/10.D0
	Mu2  float64     // collocation: C2    = (4.D0+SQ6)/10.D0
	Mu3  float64     // collocation: C1M1  = C1-1.D0
	Mu4  float64     // collocation: C2M1  = C2-1.D0
	Mu5  float64     // collocation: C1MC2 = C1-C2
	// Has unexported fields.
}
    Radau5 implements the Radau5 implicit Runge-Kutta method

func (o *Radau5) Accept(y0 la.Vector, x0 float64) (dxnew float64)
    Accept accepts update and computes next stepsize

func (o *Radau5) DenseOut(yout la.Vector, h, x float64, y la.Vector, xout float64)
    DenseOut produces dense output (after Accept)

func (o *Radau5) Free()
    Free releases memory

func (o *Radau5) Info() (fixedOnly, implicit bool, nstages int)
    Info returns information about this method

func (o *Radau5) Init(ndim int, conf *Config, work *rkwork, stat *Stat, fcn Func, jac JacF, M *la.Triplet)
    Init initialises structure

func (o *Radau5) Reject() (dxnew float64)
    Reject processes step rejection and computes next stepsize

func (o *Radau5) Step(x0 float64, y0 la.Vector)
    Step steps update

type Solver struct {
	Out  *Output // output handler
	Stat *Stat   // statistics

	FixedOnly bool // method can only be used with fixed steps
	Implicit  bool // method is implicit
	// Has unexported fields.
}
    Solver implements an ODE solver

func NewSolver(ndim int, conf *Config, fcn Func, jac JacF, M *la.Triplet) (o *Solver)
    NewSolver returns a new ODE structure with default values and allocated
    slices

        INPUT:
          ndim -- problem dimension
          conf -- configuration parameters
          out  -- output handler [may be nil]
          fcn  -- f(x,y) = dy/dx function
          jac  -- Jacobian: df/dy function [may be nil ⇒ use numerical Jacobian, if necessary]
          M    -- "mass" matrix, such that M ⋅ dy/dx = f(x,y) [may be nil]

        NOTE: remember to call Free() to release allocated resources (e.g. from the linear solvers)

func (o *Solver) Free()
    Free releases allocated memory (e.g. by the linear solvers)

func (o *Solver) Solve(y la.Vector, x, xf float64)
    Solve solves dy/dx = f(x,y) from x to xf with initial y given in y

type Stat struct {
	Nfeval    int     // number of calls to fcn
	Njeval    int     // number of Jacobian matrix evaluations
	Nsteps    int     // total number of substeps
	Naccepted int     // number of accepted substeps
	Nrejected int     // number of rejected substeps
	Ndecomp   int     // number of matrix decompositions
	Nlinsol   int     // number of calls to linsolver
	Nitmax    int     // number max of iterations
	Hopt      float64 // optimal step size at the end
	LsKind    string  // kind of linear solver used
	Implicit  bool    // method is implicit

	// benchmark
	NanosecondsStep       int64 // maximum time elapsed during steps [nanoseconds]
	NanosecondsJeval      int64 // maximum time elapsed during Jacobian evaluation [nanoseconds]
	NanosecondsIniSol     int64 // maximum time elapsed during initialization of the linear solver [nanoseconds]
	NanosecondsFact       int64 // maximum time elapsed during factorization (if any) [nanoseconds]
	NanosecondsLinSol     int64 // maximum time elapsed during solution of linear system (if any) [nanoseconds]
	NanosecondsErrorEstim int64 // maximum time elapsed during the error estimate [nanoseconds]
	NanosecondsTotal      int64 // total time elapsed during solution of ODE system [nanosecons]
}
    Stat holds statistics and output data

func NewStat(lskind string, implicit bool) (o *Stat)
    NewStat returns a new structure

func (o *Stat) Print(options ...bool)
    Print prints information about the solution process options -- noExtraInfo,
    noElapsedTime [all default to true]

func (o *Stat) Reset()
    Reset initialises Stat

type StepOutF func(istep int, h, x float64, y la.Vector) (stop bool)
    StepOutF defines a callback function to be called when a step is accepted

        INPUT:
          istep -- index of step (0 is the very first output whereas 1 is the first accepted step)
          h     -- stepsize = dx
          x     -- scalar variable
          y     -- vector variable

        OUTPUT:
          stop -- stop simulation (nicely)

type YanaF func(res []float64, x float64)
    YanaF defines a function to be used when computing analytical solutions

```
