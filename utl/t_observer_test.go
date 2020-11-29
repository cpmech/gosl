// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

type observabledata struct {
	Observable
	x, y, z int
}

func (o *observabledata) Set(x, y, z int) {
	o.x, o.y, o.z = x, y, z
	o.NotifyUpdate()
}

type iminterested struct {
	desc         string
	notification string
	data         *observabledata
}

func (o *iminterested) Update() {
	o.notification = io.Sf("got x=%d y=%d z=%d", o.data.x, o.data.y, o.data.z)
}

func TestObserver01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Observer01. Observer and Observable")

	var data observabledata

	obs01 := &iminterested{"I'm interested in updates", "", &data}
	obs02 := &iminterested{"I'm interested in updates as well", "", &data}

	data.AddObserver(obs01)
	data.AddObserver(obs02)

	chk.String(tst, obs01.notification, "")
	chk.String(tst, obs02.notification, "")
	data.Set(1, 2, 3)
	chk.String(tst, obs01.notification, "got x=1 y=2 z=3")
	chk.String(tst, obs02.notification, "got x=1 y=2 z=3")
}
