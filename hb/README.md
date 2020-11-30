# Gosl. hb. Read/Write hb (go-binary gob) files

[![PkgGoDev](https://pkg.go.dev/badge/github.com/cpmech/gosl/hb)](https://pkg.go.dev/github.com/cpmech/gosl/hb)

This package is inspired by the Hierarchical Data File (HDF) format but uses go-binary (gob) _stream of gobs_ instead.

Data is stored in binary streams and then to files by providing a **path** and the data. Here, the order of commands is important.

Example of writing data:

```go
f := Create("/tmp/gosl/hb", "basic01")
defer f.Close()
f.PutArray("/u", uSource)
f.PutArray("/displacements/u", []float64{4, 5, 6})
f.PutArray("/displacements/v", []float64{40, 50, 60})
f.PutArray("/time0/ip0/a0/u", []float64{7, 8, 9})
f.PutArray("/time1/ip0/b0/u", []float64{70, 80, 90})
f.PutInts("/someints", []int{100, 200, 300, 400})
f.PutInt("/data/oneint", 123)
f.PutFloat64("/data/onef64", 123.456)
```

Example of reading data (must be in the same order used during writing):

```go
g := Open("/tmp/gosl/hb", "basic01")
defer g.Close()
u := g.GetArray("/u")
du := g.GetArray("/displacements/u")
dv := g.GetArray("/displacements/v")
t0i0a0u := g.GetArray("/time0/ip0/a0/u")
t1i0b0u := g.GetArray("/time1/ip0/b0/u")
someints := g.GetInts("/someints")
oneint := g.GetInt("/data/oneint")
onef64 := g.GetFloat64("/data/onef64")
```

## API

[Please see the documentation here](https://pkg.go.dev/github.com/cpmech/gosl/hb)
