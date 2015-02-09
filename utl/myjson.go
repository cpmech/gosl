// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

import (
    "bytes"
    "encoding/json"
)

// WriteIntSlice writes to buffer a slice of integers
func WriteIntSlice(b *bytes.Buffer, v []int) {
    Ff(b, "[")
    for i, _ := range v {
        Ff(b, "%d", v[i])
        if i != len(v)-1 {
            Ff(b, ",")
        }
    }
    Ff(b, "]")
}

// JsonMapIntSlice writes to buffer a JSON string representation of a map of slice of integers
func JsonMapIntSlice(b *bytes.Buffer, indent, name string, m *map[int][]int) {
    Ff(b, "%s\"%s\" : {\n", indent, name)
    count := 0
    for key, vals := range *m {
        Ff(b, "%s%s\"%d\" : ", indent, indent, key)
        WriteIntSlice(b, vals)
        if count < len(*m)-1 {
            Ff(b, ",")
        }
        Ff(b, "\n")
        count += 1
    }
    Ff(b, "%s}", indent)
}

// GetMapIntSlice gets a map of slice of integers from interface (read by json)
func GetMapIntSlice(dict interface{}) (res map[int][]int) {
    res = make(map[int][]int)
    d  := dict.(map[string]interface{})
    for k, v := range d {
        //Pfpink("k=%v  v=%+#v\n", k, v)
        var ints []int
        switch vv := v.(type) {
        case []interface{}:
            //Pfgrey("  vv = %+#v\n", vv)
            for _, x := range vv {
                switch xx := x.(type) {
                case int:
                    //Pfgrey2("    i=%v x=%v (int)\n", i, xx)
                    ints = append(ints, int(xx))
                case float64:
                    //Pfgrey2("    i=%v x=%v (float64)\n", i, xx)
                    ints = append(ints, int(xx))
                default:
                    Panic(_myjson_err1, x)
                }
            }
        default:
            Panic(_myjson_err2, v)
        }
        res[Atoi(k)] = ints
    }
    return
}

// OpenAndParseJson opens and parse a json text file
func OpenAndParseJson(results interface{}, filepath, ext string) (b []byte) {

    // open file
    fname  := PathKey(filepath) + ext
    b, err := ReadFile(fname)
    if err != nil {
        Panic(_myjson_err3, fname)
    }

    // parse file
    err = json.Unmarshal(b, results)
    if err != nil {
        Panic(_myjson_err4, fname)
    }
    return
}

// error messages
var (
    _myjson_err1 = "myjson.go: GetMapIntSlice: cannot find type (int/float64) of inner item: %+#v"
    _myjson_err2 = "myjson.go: GetMapIntSlice: cannot find type of slice ([]int/[]float64): %+#v"
    _myjson_err3 = "myjson.go: OpenAndParseJson: cannot open file <%s>"
    _myjson_err4 = "myjson.go: OpenAndParseJson: cannot parse file <%s>"
)
