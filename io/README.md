# Gosl. io. Input/output, read/write files, and print commands

This subpackage helps with reading and writing files, printing nice formatted messages (with
colours), and parsing strings.

## Examples

To write and read a file:
```
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

More information is available in [the documentation of this package](http://rawgit.com/cpmech/gosl/master/doc/xxio.html).
