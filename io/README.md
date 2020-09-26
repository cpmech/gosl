# Gosl. io. Input/output, read/write files, and print commands

This subpackage helps with reading and writing files, printing nice formatted messages (with
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

**go doc**

```
package io // import "gosl/io"

Package io (input/output) implements auxiliary functions for printing,
parsing, handling files, directories, etc.

VARIABLES

var (
	// Verbose activates display of messages on console
	Verbose = true

	// ColorsOn activates use of colors on console
	ColorsOn = true
)

FUNCTIONS

func AppendToFile(fn string, buffer ...*bytes.Buffer)
    AppendToFile appends data to an existent (or new) file

func ArgToBool(idxArg int, defaultValue bool) bool
    ArgToBool parses an argument as a boolean value

        Input:
         idxArg       -- index of argument; e.g. 0==first, 1==second, etc.
         defaultValue -- default value

func ArgToFilename(idxArg int, fnDefault, ext string, check bool) (filename, fnkey string)
    ArgToFilename parses an argument as a filename

        Input:
         idxArg    -- index of argument; e.g. 0==first, 1==second, etc.
         fnDefault -- default filename; can be ""
         ext       -- the file extension to be added; e.g. ".sim"
         check     -- check for null filename
        Output:
         filename -- the filename with extension added
         fnkey    -- filename key == filename without extension
        Notes:
         The first first argument may be a file with extension or not.
        Examples:
         If the first argument is "simulation.sim" or "simulation" (with ext=".sim")
         then the results are: filename="simulation.sim" and fnkey="simulation"

func ArgToFloat(idxArg int, defaultValue float64) float64
    ArgToFloat parses an argument as a float64 value

        Input:
         idxArg       -- index of argument; e.g. 0==first, 1==second, etc.
         defaultValue -- default value

func ArgToInt(idxArg int, defaultValue int) int
    ArgToInt parses an argument as an integer value

        Input:
         idxArg       -- index of argument; e.g. 0==first, 1==second, etc.
         defaultValue -- default value

func ArgToString(idxArg int, defaultValue string) string
    ArgToString parses an argument as a string

        Input:
         idxArg       -- index of argument; e.g. 0==first, 1==second, etc.
         defaultValue -- default value

func ArgsTable(title string, data ...interface{}) (table string)
    ArgsTable prints a nice table with input arguments

        Input:
         title -- title of table; e.g. INPUT ARGUMENTS
         data  -- sets of THREE items in the following order:
                       description, key, value, ...
                       description, key, value, ...
                            ...
                       description, key, value, ...

func Atob(val string) (bres bool)
    Atob converts string to bool

func Atof(val string) (res float64)
    Atof converts string to float64

func Atoi(val string) (res int)
    Atoi converts string to integer

func Btoa(flag bool) string
    Btoa converts flag to string

        Note: true  => "true"
              false => "false"

func Btoi(flag bool) int
    Btoi converts flag to interger

        Note: true  => 1
              false => 0

func CopyFileOver(destination, source string)
    CopyFileOver copies file (Linux only with cp), overwriting if it exists
    already

func DblSf(msg string, slice []float64) string
    DblSf is the Sprintf for a slice of float64 (without brackets)

func ExtractStrPair(pair, sep string) (key, val string)
    ExtractStrPair extracts the pair from, e.g., "key:val"

        Note: it returs empty strings if any is not found

func Ff(b *bytes.Buffer, msg string, prm ...interface{})
    Ff wraps Fprintf

func FnExt(fn string) string
    FnExt returns the extension of a file name. The extension is the suffix
    beginning at the final dot in the final element of path; it is empty if
    there is no dot.

func FnKey(fn string) string
    FnKey returns the file name key (without path and extension, if any)

func IntSf(msg string, slice []int) string
    IntSf is the Sprintf for a slice of integers (without brackets)

func Itob(val int) bool
    Itob converts from integer to bool

        Note: only zero returns false
              anything else returns true

func JoinKeys(keys []string) string
    JoinKeys join keys separeted by spaces

func JoinKeys3(k0, k1, k2 []string, sep string) (res string)
    JoinKeys3 joins keys from 3 slices into a string with sets separated by sep
    and keys separeted by spaces

func JoinKeys4(k0, k1, k2, k3 []string, sep string) (res string)
    JoinKeys4 joins keys from 4 slices into a string with sets separated by sep
    and keys separeted by spaces

func JoinKeysPre(prefix string, keys []string) (res string)
    JoinKeysPre join keys separeted by spaces with a prefix

func Keycode(String string, Type string) (keycode string, found bool)
    Keycode extracts a keycode from a string such as "!typeA:keycodeA
    !typeB:keycodeB!typeC:keycodeC"

        Note: String == "!keyA !typeB:valB" is also valid

func Keycodes(String string) (keycodes []string)
    Keycodes extracts keys from a keycode (extra) string

        Example: "!keyA !typeB:valB" => keycodes = [keyA, typeB]

func OpenFileR(fn string) (fil *os.File)
    OpenFileR opens a file for reading data

func PathKey(fn string) string
    PathKey returs the full path except the extension

func Pf(msg string, prm ...interface{})
    Pf prints formatted string

func PfBlue(msg string, prm ...interface{})
    PfBlue prints formatted string in high intensity blue

func PfCyan(msg string, prm ...interface{})
    PfCyan prints formatted string in high intensity cyan

func PfGreen(msg string, prm ...interface{})
    PfGreen prints formatted string in high intensity green

func PfMag(msg string, prm ...interface{})
    PfMag prints formatted string in high intensity magenta

func PfRed(msg string, prm ...interface{})
    PfRed prints formatted string in high intensity red

func PfWhite(msg string, prm ...interface{})
    PfWhite prints formatted string in high intensity white

func PfYel(msg string, prm ...interface{})
    PfYel prints formatted string in high intensity yello

func Pfblue(msg string, prm ...interface{})
    Pfblue prints formatted string in blue

func Pfblue2(msg string, prm ...interface{})
    Pfblue2 prints formatted string in another shade of blue

func Pfcyan(msg string, prm ...interface{})
    Pfcyan prints formatted string in cyan

func Pfcyan2(msg string, prm ...interface{})
    Pfcyan2 prints formatted string in another shade of cyan

func Pfdgreen(msg string, prm ...interface{})
    Pfdgreen prints formatted string in dark green

func Pfdyel(msg string, prm ...interface{})
    Pfdyel prints formatted string in dark yellow

func Pfdyel2(msg string, prm ...interface{})
    Pfdyel2 prints formatted string in another shade of dark yellow

func Pfgreen(msg string, prm ...interface{})
    Pfgreen prints formatted string in green

func Pfgreen2(msg string, prm ...interface{})
    Pfgreen2 prints formatted string in another shade of green

func Pfgrey(msg string, prm ...interface{})
    Pfgrey prints formatted string in grey

func Pfgrey2(msg string, prm ...interface{})
    Pfgrey2 prints formatted string in another shade of grey

func Pflmag(msg string, prm ...interface{})
    Pflmag prints formatted string in light magenta

func Pfmag(msg string, prm ...interface{})
    Pfmag prints formatted string in magenta

func Pforan(msg string, prm ...interface{})
    Pforan prints formatted string in orange

func Pfpink(msg string, prm ...interface{})
    Pfpink prints formatted string in pink

func Pfpurple(msg string, prm ...interface{})
    Pfpurple prints formatted string in purple

func Pfred(msg string, prm ...interface{})
    Pfred prints formatted string in red

func Pfyel(msg string, prm ...interface{})
    Pfyel prints formatted string in yellow

func Pipeline(cmds ...*exec.Cmd) (pipeLineOutput, collectedStandardError []byte)
    Pipeline strings together the given exec.Cmd commands in a similar fashion
    to the Unix pipeline. Each command's standard output is connected to the
    standard input of the next command, and the output of the final command in
    the pipeline is returned, along with the collected standard error of all
    commands and the first error found (if any).

    by Kyle Lemons

    To provide input to the pipeline, assign an io.Reader to the first's Stdin.

func Pl()
    Pl prints new line

func ReadFile(fn string) (b []byte)
    ReadFile reads bytes from a file

func ReadLines(fn string, cb ReadLinesCallback)
    ReadLines reads lines from a file and calls ReadLinesCallback to process
    each line being read

func ReadLinesFile(fil *os.File, cb ReadLinesCallback)
    ReadLinesFile reads lines from a file and calls ReadLinesCallback to process
    each line being read

func ReadMatrix(fn string) (M [][]float64)
    ReadMatrix reads a text file in which the float64 type of numeric values
    represent a matrix of data. The number of columns must be equal, including
    for the headers

func ReadTable(fn string) (keys []string, T map[string][]float64)
    ReadTable reads a text file in which the first line contains the headers and
    the next lines the float64 type of numeric values. The number of columns
    must be equal, including for the headers

func RemoveAll(key string)
    RemoveAll deletes all files matching filename specified by key (be careful)

func RunCmd(verbose bool, cmd string, args ...string) string
    RunCmd runs external command

func Sf(msg string, prm ...interface{}) string
    Sf wraps Sprintf

func SplitFloats(str string) (res []float64)
    SplitFloats splits space-separated float numbers

func SplitInts(str string) (res []int)
    SplitInts splits space-separated integers

func SplitKeys(keys string) []string
    SplitKeys splits keys separeted by spaces

func SplitKeys3(res string) (k0, k1, k2 []string)
    SplitKeys3 splits a string with three sets of keys separated by comma

func SplitKeys4(res string) (k0, k1, k2, k3 []string)
    SplitKeys4 splits a string with four sets of keys separated by comma

func SplitSpacesQuoted(str string) (res []string)
    SplitSpacesQuoted splits string with quoted substrings. e.g. " a,b, 'c',
    \"d\" "

func SplitWithinParentheses(s string) (res []string)
    SplitWithinParentheses extracts arguments (substrings) within brackets e.g.:
    "(arg1, (arg2.1, arg2.2), arg3, arg4, (arg5.1,arg5.2, arg5.3 ) )"

func StrSf(msg string, slice []string) string
    StrSf is the Sprintf for a slice of string (without brackets)

func StrSpaces(n int) (l string)
    StrSpaces returns a line with spaces

func StrThickLine(n int) (l string)
    StrThickLine returns a thick line (using '=')

func StrThinLine(n int) (l string)
    StrThinLine returns a thin line (using '-')

func UnColor(msg string) string
    UnColor removes console characters used to show colors

func WriteBytesToFile(fn string, b []byte)
    WriteBytesToFile writes slice of bytes to a new file

func WriteBytesToFileD(dirout, fn string, b []byte)
    WriteBytesToFileD writes slice of bytes to a new file after creating a
    directory

func WriteBytesToFileVD(dirout, fn string, b []byte)
    WriteBytesToFileVD writes slice of bytes to a new file, and print message,
    after creating a directory

func WriteFile(fn string, buffer ...*bytes.Buffer)
    WriteFile writes data to a new file with given bytes.Buffer(s)

func WriteFileD(dirout, fn string, buffer ...*bytes.Buffer)
    WriteFileD writes data to a new file after creating a directory

func WriteFileV(fn string, buffer ...*bytes.Buffer)
    WriteFileV writes data to a new file, and shows message

func WriteFileVD(dirout, fn string, buffer ...*bytes.Buffer)
    WriteFileVD writes data to a new file, and shows message after creating a
    directory

func WriteStringToFile(fn, data string)
    WriteStringToFile writes string to a new file

func WriteStringToFileD(dirout, fn, data string)
    WriteStringToFileD writes string to a new file after creating a directory

func WriteTableVD(dirout, fn string, headers []string, columns ...[]float64)
    WriteTableVD writes a text file in which the first line contains the
    headers, and the next lines contain the numeric values (float64). The number
    of columns must be equal to the number of headers.


TYPES

type ReadLinesCallback func(idx int, line string) (stop bool)
    ReadLinesCallback is used in ReadLines to process line by line during
    reading of a file

```
