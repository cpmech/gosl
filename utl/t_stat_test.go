// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

import (
    "bytes"
    "testing"
    "time"
)

func TestStat01(tst *testing.T) {
    defer func() {
        if err := recover(); err != nil {
            tst.Error("[1;31mSome error has happened:[0m\n", err)
        }
    }()

    TTitle("Stat 01")

    var b bytes.Buffer
    durs := make([]time.Duration, 4)
    durs[0], _ = time.ParseDuration("1s")
    durs[1], _ = time.ParseDuration("2s")
    durs[2], _ = time.ParseDuration("3s")
    durs[3], _ = time.ParseDuration("1m")
    min, max, ave, sum := DurStat(nil, durs, "")
    Pf("min=%v, max=%v, ave=%v, sum=%v (%d items)\n", min, max, ave, sum, len(durs))

    StatHeader(&b, true)
    DurStat(&b, durs, "duration")

    ints := []int{1,5,6,7}
    IntStat(&b, ints, "counter")

    flts := []float64{1,5,6,7}
    DblStat(&b, flts, "values")

    StatFooter(&b, true)
    Pfblue2("%v\n", b.String())
}
