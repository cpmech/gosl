# Gosl. hdf. Hierarchical Data Format (HDF5 Wrapper)

Hierarchical Data Format (HDF) is a set of file formats (HDF4, HDF5) designed to store and organize large amounts of data. 

[HDF](https://portal.hdfgroup.org/) enables the management of extremely large and complex data collections.

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/cpmech/gosl/hdf)

## Usage

HDF is quite simple to use. You basically give a string resembling a file path and the data to be stored or read.

Example of saving data:

```go
useGob := false // do not use standard Go binary file, i.e. use HDF5
f := Create("/tmp/gosl", "mydatafile", useGob)
f.PutArray("/displacements/u", []float64{4, 5, 6})
f.PutArray("/displacements/v", []float64{40, 50, 60})
f.PutInts("/someints", []int{100, 200, 300, 400})
f.PutInt("/data/oneint", 123)
f.PutFloat64("/data/onef64", 123.456)
f.Close()
```

Example of reading data:

```go
g := Open("/tmp/gosl", "mydatafile", useGob)
u := g.GetArray("/displacements/u")
v := g.GetArray("/displacements/v")
someints := g.GetInts("/someints")
oneint := g.GetInt("/data/oneint")
onef64 := g.GetFloat64("/data/onef64")
```
