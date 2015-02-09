// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

import (
    "bytes"
    "time"
)

// StatHeader returns the header of a statistics table with:
//  _ _ _ _ _ _  min max ave sum nitems
func StatHeader(b *bytes.Buffer, lines bool) {
    if b == nil {
        return
    }
    if lines {
        Ff(b, " =========================================================================================\n")
    }
    Ff(b, " %10s %17s%17s%17s%17s%10s\n", "", "min", "max", "ave", "sum", "nitems")
    if lines {
        Ff(b, " -----------------------------------------------------------------------------------------\n")
    }
    return
}

// StatFooter returns the horizontal line closing the statistics table
func StatFooter(b *bytes.Buffer, lines bool) {
    if b == nil {
        return
    }
    if lines {
        Ff(b, " =========================================================================================\n")
    }
}

// DurStat generates the line for the statistics table corresponding to a duration
func DurStat(b *bytes.Buffer, durs []time.Duration, name string) (min, max, ave, sum time.Duration) {
    if len(durs) == 0 {
        return
    }
    min, max = durs[0], durs[0]
    for _, d := range durs {
        if d < min { min = d }
        if d > max { max = d }
        sum += d
    }
    ave = sum / time.Duration(int64(len(durs)))
    if name != "" && b != nil {
        Ff(b, " %10s:%17v%17v%17v%17v%10d\n", name, min, max, ave, sum, len(durs))
    }
    return
}

// IntStat generates the line for the statistics table corresponding to results collected in 'ints'
func IntStat(b *bytes.Buffer, ints []int, name string) (min, max, ave, sum int) {
    if len(ints) == 0 {
        return
    }
    min, max = ints[0], ints[0]
    for _, i := range ints {
        if i < min { min = i }
        if i > max { max = i }
        sum += i
    }
    ave = sum / len(ints)
    if name != "" && b != nil {
        Ff(b, " %10s:%17v%17v%17v%17v%10d\n", name, min, max, ave, sum, len(ints))
    }
    return
}

// DblStat generates the line for the statistics table corresponding to results collected in 'flts'
func DblStat(b *bytes.Buffer, flts []float64, name string) (min, max, ave, sum float64) {
    if len(flts) == 0 {
        return
    }
    min, max = flts[0], flts[0]
    for _, f := range flts {
        if f < min { min = f }
        if f > max { max = f }
        sum += f
    }
    ave = sum / float64(len(flts))
    if name != "" && b != nil {
        Ff(b, " %10s:%17.5e%17.5e%17.5e%17.5e%10d\n", name, min, max, ave, sum, len(flts))
    }
    return
}
