// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rw

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
)

func Test_step01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("step01")

	dat := `
#94 = B_SPLINE_CURVE_WITH_KNOTS('',3,(#95,#96,#97,#98,#99,#100,#101,#102
    ,#103,#104,#105,#106,#107,#108,#109,#110,#111,#112,#113,#114,#115,
    #116,#117,#118),.UNSPECIFIED.,.F.,.F.,(4,2,2,2,2,2,2,2,2,2,2,4),(0.,
    6.25E-02,0.125,0.1875,0.25,0.375,0.5,0.625,0.75,0.875,0.9375,1.),
  .UNSPECIFIED.);
#95 = CARTESIAN_POINT('',(-101.6,22.574321695,10.));
#96 = CARTESIAN_POINT('',(-102.823664276,22.574321695,10.));
#97 = CARTESIAN_POINT('',(-103.997218852,22.153675049,10.));
#98 = CARTESIAN_POINT('',(-106.08086702,20.987537487,10.));
#99 = CARTESIAN_POINT('',(-107.039417042,20.231466666,10.));
#100 = CARTESIAN_POINT('',(-108.789255706,18.580147193,10.));
#101 = CARTESIAN_POINT('',(-109.592274608,17.667817128,10.));
#102 = CARTESIAN_POINT('',(-111.037380535,15.764228357,10.));
#103 = CARTESIAN_POINT('',(-111.683917439,14.768224001,10.));
#104 = CARTESIAN_POINT('',(-113.420358551,11.662050083,10.));
#105 = CARTESIAN_POINT('',(-114.315877571,9.399456447,10.));
#106 = CARTESIAN_POINT('',(-115.497833567,4.783654514,10.));
#107 = CARTESIAN_POINT('',(-115.797754021,2.416719426,10.));
#108 = CARTESIAN_POINT('',(-115.799425478,-2.39678873,10.));
#109 = CARTESIAN_POINT('',(-115.502828861,-4.759182837,10.));
#110 = CARTESIAN_POINT('',(-114.318091024,-9.395746276,10.));
#111 = CARTESIAN_POINT('',(-113.429172576,-11.644133428,10.));
#112 = CARTESIAN_POINT('',(-111.107704555,-15.801102362,10.));
#113 = CARTESIAN_POINT('',(-109.654523885,-17.763443825,10.));
#114 = CARTESIAN_POINT('',(-107.045203102,-20.226082032,10.));
#115 = CARTESIAN_POINT('',(-106.086648295,-20.984357267,10.));
#116 = CARTESIAN_POINT('',(-103.990686794,-22.157275434,10.));
#117 = CARTESIAN_POINT('',(-102.821037828,-22.574321695,10.));
#118 = CARTESIAN_POINT('',(-101.6,-22.574321695,10.));
`

	var stp STEP
	err := stp.ParseDATA(dat)
	if err != nil {
		tst.Errorf("Parse filed:\n%v", err)
		return
	}

	np := len(stp.Points)
	x := make([]float64, np)
	y := make([]float64, np)
	z := make([]float64, np)
	i := 0
	for _, p := range stp.Points {
		io.Pforan("p = %#v\n", p)
		x[i] = p.Coordinates[0]
		y[i] = p.Coordinates[1]
		z[i] = p.Coordinates[2]
		i++
	}

	for _, c := range stp.BScurves {
		io.Pf("c = %#v\n", c)
	}

	if false {
		plt.Plot3dPoints(x, y, z, "")
		plt.Show()
	}
}

func Test_step02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("step02")

	buf, err := io.ReadFile("data/beadpanel.step")
	if err != nil {
		tst.Errorf("cannot read file:\n%v", err)
		return
	}
	dat := string(buf)

	var stp STEP
	err = stp.ParseDATA(dat)
	if err != nil {
		tst.Errorf("Parse filed:\n%v", err)
		return
	}

	for _, c := range stp.BScurves {
		io.Pforan("c = %v\n", c)
	}
}
