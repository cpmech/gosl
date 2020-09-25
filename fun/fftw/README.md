# Gosl. fun/fftw. Wrapper to FFTW3

This package wraps the [Fast Fourier Transform library (FFTW)](http://www.fftw.org)

## API

```
package fftw // import "gosl/fun/fftw"

Package fftw wraps the FFTW library to perform Fourier Transforms using the
"fast" method by Cooley and Tukey

TYPES

type Plan1d struct {
	// Has unexported fields.
}
    Plan1d implements the FFTW3 plan structure; i.e. a "plan" to compute direct
    or inverse 1D FTs

        Computes:
                           N-1         -i 2 π j k / N                 __
          forward:  X[k] =  Σ  x[j] ⋅ e                     with i = √-1
                           j=0

                           N-1         +i 2 π j k / N
          inverse:  Y[k] =  Σ  y[j] ⋅ e                     thus x[k] = Y[k] / N
                           j=0

        NOTE: FFTW says "so you should initialize your input data after creating the plan."
              Therefore, the plan can be created and reused several times.
              [http://www.fftw.org/fftw3_doc/Planner-Flags.html]
              Also: "The plan can be reused as many times as needed. In typical high-performance
              applications, many transforms of the same size are computed"
              [http://www.fftw.org/fftw3_doc/Introduction.html]

              Create a new Plan1d with NewPlan1d(...) AND deallocate memory with Free()

func NewPlan1d(data []complex128, inverse, measure bool) (o *Plan1d)
    NewPlan1d allocates a new "plan" to compute 1D Fourier Transforms

        data    -- [modified] data is a complex array of length N.
        inverse -- will perform inverse transform; otherwise will perform direct
                   Note: both transforms are non-normalised;
                   i.e. the user will have to multiply by (1/n) if computing inverse transforms
        measure -- use the FFTW_MEASURE flag for better optimisation analysis (slower initialisation times)
                   Note: using this flag with given "data" as input will cause the allocation
                   of a temporary array and the execution of two copy commands with size len(data)

        NOTE: (1) the user must remember to call Free to deallocate FFTW data
              (2) data will be overwritten

func (o *Plan1d) Execute()
    Execute performs the Fourier transform

func (o *Plan1d) Free()
    Free frees internal FFTW data

type Plan2d struct {
	// Has unexported fields.
}
    Plan2d implements the FFTW3 plan structure; i.e. a "plan" to compute direct
    or inverse 2D FTs

        Computes:
                           N1-1 N0-1             -i 2 π k1 l1 / N1    -i 2 π k0 l0 / N0
                X[l0,l1] =   Σ    Σ  x[k0,k1] ⋅ e                  ⋅ e
                           k1=0 k0=0

func NewPlan2d(N0, N1 int, data []complex128, inverse, measure bool) (o *Plan2d)
    NewPlan2d allocates a new "plan" to compute 2D Fourier Transforms

        N0, N1  -- dimensions
        data    -- [modified] data is a complex array of length N0*N1 (row-major matrix)
        inverse -- will perform inverse transform; otherwise will perform direct
                   Note: both transforms are non-normalised;
                   i.e. the user will have to multiply by (1/n) if computing inverse transforms
        measure -- use the FFTW_MEASURE flag for better optimisation analysis (slower initialisation times)
                   Note: using this flag with given "data" as input will cause the allocation
                   of a temporary array and the execution of two copy commands with size len(data)

        NOTE: (1) the user must remember to call Free to deallocate FFTW data
              (2) data will be overwritten

        A = data is a ROW-MAJOR matrix
             _                          _
        A = |   A00→a0  A01→a1  A02→a2   |  ⇒ A[i][j]
            |_  A10→a3  A11→a4  A12→a5  _|
                                         (n0,n1)=(2,3)

            l = 0      1      2        3      4      5
        a = [  A00    A01    A02      A10    A11    A12 ]  ⇒ a[l]
              3⋅0+0  3⋅0+1  3⋅0+2    3⋅1+0  3⋅1+1  3⋅1+2

        l = n1⋅i + j       i = l // n1      j = l % n1

func (o *Plan2d) Execute()
    Execute performs the Fourier transform

func (o *Plan2d) Free()
    Free frees internal FFTW data

func (o *Plan2d) Get(i, j int) (v complex128)
    Get gets data value located at "i,j". NOTE: this method does not check for
    out-of-range indices

func (o *Plan2d) GetSlice() (out [][]complex128)
    GetSlice gets the output array as a nested slice

func (o *Plan2d) Set(i, j int, v complex128)
    Set sets data value located at "i,j". NOTE: this method does not check for
    out-of-range indices

```
