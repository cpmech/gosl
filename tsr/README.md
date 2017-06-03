# Gosl. tsr. Tensor algebra and definitions for continuum mechanics

[![GoDoc](https://godoc.org/github.com/cpmech/gosl/tsr?status.svg)](https://godoc.org/github.com/cpmech/gosl/tsr) 

More information is available in **[the documentation of this package](https://godoc.org/github.com/cpmech/gosl/tsr).**

Package `tsr` implements several functions to assist in the implementation of tensor algebra and
calculus with some focus on continuum mechanics. Basic and advanced features are available. For
example, the complexity of functions range from the computation of characteristic invariants to the
derivatives of eigenprojectors.

An important feature of this library is the use of Mandel representation for the components of
symmetric tensors. In this representation, second order tensors are simply a vector in a 6D space
and fourth order tensors are matrices in a 6D space. We use the term **Mandel tensor** here to refer
to components of the tensor modified according to Mandel's representation.

## Continuum Mechanics

### Deformation Tensors

1. `RightCauchyGreenDef`
2. `LeftCauchyGreenDef`

### Strain tensors

1. `GreenStrain`
2. `AlmansiStrain`
3. `LinStrain`

### Stress tensors

1. `PK1ToCauchy`
2. `CauchyToPK2`
3. `PK2ToCauchy`

### Transformations

1. `PushForward`
2. `PullBack`
3. `PushForwardB`
4. `PullBackB`

## Scalar Invariants

## Isotropic functions

## Eigenvalues, Eigenvectors, and Eigenprojectors

## Plotting
