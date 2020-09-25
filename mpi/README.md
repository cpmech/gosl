# Gosl. mpi. Message Passing Interface for parallel computing

The `mpi` package is a light wrapper to the [OpenMPI](https://www.open-mpi.org) C++ library designed
to develop algorithms for parallel computing.

This package allows parallel computations over the network.

## API

**go doc**

```
package mpi // import "gosl/mpi"

Package mpi wraps the Message Passing Interface for parallel computations

FUNCTIONS

func IsOn() bool
    IsOn tells whether MPI is on or not

        NOTE: this returns true even after Stop

func Start()
    Start initialises MPI

func StartThreadSafe() error
    StartThreadSafe initialises MPI thread safe

func Stop()
    Stop finalises MPI

func WorldRank() (rank int)
    WorldRank returns the processor rank/ID within the World communicator

func WorldSize() (size int)
    WorldSize returns the number of processors in the World communicator


TYPES

type Communicator struct {
	// Has unexported fields.
}
    Communicator holds the World communicator or a subset communicator

func NewCommunicator(ranks []int) (o *Communicator)
    NewCommunicator creates a new communicator or returns the World communicator

        ranks -- World indices of processors in this Communicator.
                 use nil or empty to get the World Communicator

func (o *Communicator) Abort()
    Abort aborts MPI

func (o *Communicator) AllReduceMax(dest, orig []float64)
    AllReduceMax combines all values from orig into dest picking minimum values

        NOTE (important): orig and dest must be different slices

func (o *Communicator) AllReduceMaxI(dest, orig []int)
    AllReduceMaxI combines all values from orig into dest picking minimum values
    (integer version)

        NOTE (important): orig and dest must be different slices

func (o *Communicator) AllReduceMin(dest, orig []float64)
    AllReduceMin combines all values from orig into dest picking minimum values

        NOTE (important): orig and dest must be different slices

func (o *Communicator) AllReduceMinI(dest, orig []int)
    AllReduceMinI combines all values from orig into dest picking minimum values
    (integer version)

        NOTE (important): orig and dest must be different slices

func (o *Communicator) AllReduceSum(dest, orig []float64)
    AllReduceSum combines all values from orig into dest summing values

        NOTE (important): orig and dest must be different slices

func (o *Communicator) AllReduceSumC(dest, orig []complex128)
    AllReduceSumC combines all values from orig into dest summing values
    (complex version)

        NOTE (important): orig and dest must be different slices

func (o *Communicator) Barrier()
    Barrier forces synchronisation

func (o *Communicator) BcastFromRoot(x []float64)
    BcastFromRoot broadcasts slice from root (Rank == 0) to all other processors

func (o *Communicator) BcastFromRootC(x []complex128)
    BcastFromRootC broadcasts slice from root (Rank == 0) to all other
    processors (complex version)

func (o *Communicator) Rank() (rank int)
    Rank returns the processor rank/ID

func (o *Communicator) Recv(vals []float64, fromID int)
    Recv receives values from processor fromId

func (o *Communicator) RecvC(vals []complex128, fromID int)
    RecvC receives values from processor fromId (complex version)

func (o *Communicator) RecvI(vals []int, fromID int)
    RecvI receives values from processor fromId (integer version)

func (o *Communicator) RecvOne(fromID int) (val float64)
    RecvOne receives one value from processor fromId

func (o *Communicator) RecvOneI(fromID int) (val int)
    RecvOneI receives one value from processor fromId (integer version)

func (o *Communicator) ReduceSum(dest, orig []float64)
    ReduceSum sums all values in 'orig' to 'dest' in root (Rank == 0) processor

        NOTE (important): orig and dest must be different slices

func (o *Communicator) ReduceSumC(dest, orig []complex128)
    ReduceSumC sums all values in 'orig' to 'dest' in root (Rank == 0) processor
    (complex version)

        NOTE (important): orig and dest must be different slices

func (o *Communicator) Send(vals []float64, toID int)
    Send sends values to processor toID

func (o *Communicator) SendC(vals []complex128, toID int)
    SendC sends values to processor toID (complex version)

func (o *Communicator) SendI(vals []int, toID int)
    SendI sends values to processor toID (integer version)

func (o *Communicator) SendOne(val float64, toID int)
    SendOne sends one value to processor toID

func (o *Communicator) SendOneI(val int, toID int)
    SendOneI sends one value to processor toID (integer version)

func (o *Communicator) Size() (size int)
    Size returns the number of processors

```
