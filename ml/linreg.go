// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ml

// LinReg implements a linear regression model
type LinReg struct {
	data   *Data      // X-y data
	stat   *Stat      // statistics
	params *ParamsReg // θ and b
}

// NewLinReg returns a new LinReg object
//   Input:
//     data -- X,y data
//     params -- θ and b
//     name -- unique name of this object
func NewLinReg(data *Data, params *ParamsReg, name string) (o *LinReg) {
	o = new(LinReg)
	o.data = data
	o.params = params
	o.stat = NewStat(data, name)
	o.stat.Update()
	return
}
