// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

import (
    "os"
)

type RTEhandler func() (stop bool)

// CatchRTE catch run time error
func CatchRTE(msg string, handler RTEhandler) {
    if err := recover(); err != nil {
        Pfred("\n\n%v:\n%v\n", msg, err)
        if handler() {
            PfRed("rterrors.go: CatchRTE: STOP\n")
            os.Exit(1)
        }
    }
}
