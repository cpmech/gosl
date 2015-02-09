// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
    "flag"
    "code.google.com/p/gosl/gm"
    "code.google.com/p/gosl/utl"
    "code.google.com/p/gosl/plt"
)

func main() {

    // catch errors
    utl.Tsilent = false
    defer func() {
        if err := recover(); err != nil {
            utl.PfRed("Some error has happened: %v\n", err)
        }
    }()

    // default input data
    fn   := "nurbs01.msh"
    ctrl := true
    ids  := true
    npts := 41

    // parse flags
	flag.Parse()
	if len(flag.Args()) > 0 { fn   =          flag.Arg(0)  }
    if len(flag.Args()) > 1 { ctrl = utl.Atob(flag.Arg(1)) }
    if len(flag.Args()) > 2 { ids  = utl.Atob(flag.Arg(2)) }
    if len(flag.Args()) > 3 { npts = utl.Atoi(flag.Arg(3)) }

    // print input data
    utl.Pforan("Input data\n")
    utl.Pforan("==========\n")
    utl.Pfblue2("  fn   = %v\n", fn)
    utl.Pfblue2("  ctrl = %v\n", ctrl)
    utl.Pfblue2("  ids  = %v\n", ids)
    utl.Pfblue2("  npts = %v\n", npts)

    // load nurbss
    fnk := utl.FnKey(fn)
    B   := gm.ReadMsh(fnk)

    // plot
    plt.SetForEps(0.75, 500)
    for _, b := range B {
        if ctrl {
            b.DrawCtrl2D(ids)
        }
        b.DrawElems2D(npts, ids, "", "")
    }
    plt.Equal()
    plt.Save(fnk + ".eps")
}
