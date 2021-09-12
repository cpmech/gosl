# Gosl. io. Input/output, read/write files, and print commands

[![Go Reference](https://pkg.go.dev/badge/github.com/cpmech/gosl/io.svg)](https://pkg.go.dev/github.com/cpmech/gosl/io)

This sub-package helps with reading and writing files, printing nice formatted messages (with
colours), and parsing strings.

## Examples

### Read and write files

To write and read a file:

```go
theline := "Hello World !!!"
io.WriteFileSD("/tmp/gosl", "filestring.txt", theline)

f, err := io.OpenFileR("/tmp/gosl/filestring.txt")
if err != nil {
    chk.Panic("%v", err)
}

io.ReadLinesFile(f, func(idx int, line string) (stop bool) {
    Pforan("line = %v\n", line)
    chk.String(tst, line, theline)
    return
})
```

## API

[Please see the documentation here](https://pkg.go.dev/github.com/cpmech/gosl/io)
