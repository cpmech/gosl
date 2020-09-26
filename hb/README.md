# Gosl. hb. Read/Write hb (go-binary gob) files

This package is inspired by the Hierarchical Data File (HDF) format but uses go-binary (gob) _stream of gobs_ instead.

Data is stored in binary streams and then to files by providing a **path** and the data. Here, the order of commands is important.

Example of writing data:

```go
f := Create("/tmp/gosl/h5", "basic01")
defer f.Close()
f.PutArray("/u", uSource)
f.PutArray("/displacements/u", []float64{4, 5, 6})
f.PutArray("/displacements/v", []float64{40, 50, 60})
f.PutArray("/time0/ip0/a0/u", []float64{7, 8, 9})
f.PutArray("/time1/ip0/b0/u", []float64{70, 80, 90})
f.PutInts("/someints", []int{100, 200, 300, 400})
f.PutInt("/data/oneint", 123)
f.PutFloat64("/data/onef64", 123.456)
f.Close()
```

Example of reading data (must be in the same order used during writing):

```go
g := Open("/tmp/gosl/h5", "basic01")
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

**go doc**

```
package hb // import "gosl/hb"

Package hb implements a pseudo hierarchical binary (hb) data file format

TYPES

type File struct {
	// Has unexported fields.
}
    File represents a HDF5 file

func Create(dirOut, fnameKey string) (o *File)
    Create creates a new file, deleting existent one

        Input:
          dirOut   -- directory name that will be created if non-existent
                      Note: dirOut may contain environment variables
          fnameKey -- filename key; e.g. without extension

        Output:
          returns a new File object where the filename will be:
            fnameKey + .gob

func Open(dirIn, fnameKey string) (o *File)
    Open opens an existent file for read only

        Input:
          dirIn    -- directory name where the file is located
                      Note: dirIn may contain environment variables
          fnameKey -- filename key; e.g. without extension

        Output:
          returns a new File object where the filename will be:
            fnameKey + .gob

func (o *File) Close()
    Close closes file

func (o File) Filename() string
    Filename returns the filename; i.e. fileNameKey + extension

func (o File) Filepath() string
    Filepath returns the full filepath, including directory name

func (o *File) GetArray(path string) (v []float64)
    GetArray gets an array from file. Memory will be allocated

func (o *File) GetDeep2(path string) (a [][]float64)
    GetDeep2 gets a Deep2 slice (that was serialized). Memory will be allocated

func (o *File) GetDeep2raw(path string) (m, n int, a []float64)
    GetDeep2raw returns the serialized data corresponding to a Deep2 slice

func (o *File) GetDeep3(path string) (a [][][]float64)
    GetDeep3 gets a deep slice with 3 levels from file. Memory will be allocated

func (o *File) GetFloat64(path string) float64
    GetFloat64 gets one float64 from file

        Note: this is a convenience function wrapping GetArray

func (o *File) GetInt(path string) int
    GetInt gets one integer from file

        Note: this is a convenience function wrapping GetInts

func (o *File) GetIntAttribute(path, key string) (val int)
    GetIntAttribute gets int attribute

func (o *File) GetInts(path string) (v []int)
    GetInts gets a slice of ints from file. Memory will be allocated

func (o *File) GetIntsAttribute(path, key string) (vals []int)
    GetIntsAttribute gets slice-of-ints attribute

func (o *File) GetStringAttribute(path, key string) (val string)
    GetStringAttribute gets string attribute

func (o *File) PutArray(path string, v []float64)
    PutArray puts an array with name described by path

        Input:
          path -- path such as "/myvec" or "/group/myvec"
          v    -- slice of float64

func (o *File) PutDeep2(path string, a [][]float64)
    PutDeep2 puts a Deep2 slice into file

        Input:
          path -- HDF5 path such as "/myvec" or "/group/myvec"
          a    -- slice of slices of float64
        Note: Slice will be serialized

func (o *File) PutDeep3(path string, a [][][]float64)
    PutDeep3 puts a deep slice with 3 levels and name described in path into
    HDF5 file

        Input:
          path -- HDF5 path such as "/myvec" or "/group/myvec"
          a    -- slice of slices of slices of float64
        Note: Slice will be serialized

func (o *File) PutFloat64(path string, val float64)
    PutFloat64 puts one float64 into file

        Input:
          path -- path such as "/myvec" or "/group/myvec"
          val  -- value
        Note: this is a convenience function wrapping PutArray

func (o *File) PutInt(path string, val int)
    PutInt puts one integer into file

        Input:
          path -- HDF5 path such as "/myvec" or "/group/myvec"
          val  -- value
        Note: this is a convenience function wrapping PutInts

func (o *File) PutInts(path string, v []int)
    PutInts puts a slice of integers into file

        Input:
          path -- HDF5 path such as "/myvec" or "/group/myvec"
          v    -- slice of integers

func (o *File) ReadArray(v []float64, path string) (dims []int)
    ReadArray reads an array from file into existent pre-allocated memory

        Input:
          path -- path such as "/myvec" or "/group/myvec"
        Output:
          array -- values in pre-allocated array => must know dimension
          dims  -- dimensions (for confirmation)

func (o *File) SetIntAttribute(path, key string, val int)
    SetIntAttribute sets int attibute

func (o *File) SetIntsAttribute(path, key string, vals []int)
    SetIntsAttribute sets slice-of-ints attibute

func (o *File) SetStringAttribute(path, key, val string)
    SetStringAttribute sets a string attibute

```
