// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tsr

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func TestTensor4set01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Tensor4set01")

	A := NewTensor4(false, false)

	// 00..
	A.Set(0, 0, 0, 0, 100.0)
	A.Set(0, 0, 1, 1, 101.0)
	A.Set(0, 0, 2, 2, 102.0)

	A.Set(0, 0, 0, 1, 103.0)
	A.Set(0, 0, 1, 2, 104.0)
	A.Set(0, 0, 0, 2, 105.0)

	A.Set(0, 0, 1, 0, 106.0)
	A.Set(0, 0, 2, 1, 107.0)
	A.Set(0, 0, 2, 0, 108.0)

	// 11..
	A.Set(1, 1, 0, 0, 109.0)
	A.Set(1, 1, 1, 1, 110.0)
	A.Set(1, 1, 2, 2, 111.0)

	A.Set(1, 1, 0, 1, 112.0)
	A.Set(1, 1, 1, 2, 113.0)
	A.Set(1, 1, 0, 2, 114.0)

	A.Set(1, 1, 1, 0, 115.0)
	A.Set(1, 1, 2, 1, 116.0)
	A.Set(1, 1, 2, 0, 117.0)

	// 22..
	A.Set(2, 2, 0, 0, 118.0)
	A.Set(2, 2, 1, 1, 119.0)
	A.Set(2, 2, 2, 2, 120.0)

	A.Set(2, 2, 0, 1, 121.0)
	A.Set(2, 2, 1, 2, 122.0)
	A.Set(2, 2, 0, 2, 123.0)

	A.Set(2, 2, 1, 0, 124.0)
	A.Set(2, 2, 2, 1, 125.0)
	A.Set(2, 2, 2, 0, 126.0)

	// 01..
	A.Set(0, 1, 0, 0, 127.0)
	A.Set(0, 1, 1, 1, 128.0)
	A.Set(0, 1, 2, 2, 129.0)

	A.Set(0, 1, 0, 1, 130.0)
	A.Set(0, 1, 1, 2, 131.0)
	A.Set(0, 1, 0, 2, 132.0)

	A.Set(0, 1, 1, 0, 133.0)
	A.Set(0, 1, 2, 1, 134.0)
	A.Set(0, 1, 2, 0, 135.0)

	// 12..
	A.Set(1, 2, 0, 0, 136.0)
	A.Set(1, 2, 1, 1, 137.0)
	A.Set(1, 2, 2, 2, 138.0)

	A.Set(1, 2, 0, 1, 139.0)
	A.Set(1, 2, 1, 2, 140.0)
	A.Set(1, 2, 0, 2, 141.0)

	A.Set(1, 2, 1, 0, 142.0)
	A.Set(1, 2, 2, 1, 143.0)
	A.Set(1, 2, 2, 0, 144.0)

	// 02..
	A.Set(0, 2, 0, 0, 145.0)
	A.Set(0, 2, 1, 1, 146.0)
	A.Set(0, 2, 2, 2, 147.0)

	A.Set(0, 2, 0, 1, 148.0)
	A.Set(0, 2, 1, 2, 149.0)
	A.Set(0, 2, 0, 2, 150.0)

	A.Set(0, 2, 1, 0, 151.0)
	A.Set(0, 2, 2, 1, 152.0)
	A.Set(0, 2, 2, 0, 153.0)

	// 10..
	A.Set(1, 0, 0, 0, 154.0)
	A.Set(1, 0, 1, 1, 155.0)
	A.Set(1, 0, 2, 2, 156.0)

	A.Set(1, 0, 0, 1, 157.0)
	A.Set(1, 0, 1, 2, 158.0)
	A.Set(1, 0, 0, 2, 159.0)

	A.Set(1, 0, 1, 0, 160.0)
	A.Set(1, 0, 2, 1, 161.0)
	A.Set(1, 0, 2, 0, 162.0)

	// 21..
	A.Set(2, 1, 0, 0, 163.0)
	A.Set(2, 1, 1, 1, 164.0)
	A.Set(2, 1, 2, 2, 165.0)

	A.Set(2, 1, 0, 1, 166.0)
	A.Set(2, 1, 1, 2, 167.0)
	A.Set(2, 1, 0, 2, 168.0)

	A.Set(2, 1, 1, 0, 169.0)
	A.Set(2, 1, 2, 1, 170.0)
	A.Set(2, 1, 2, 0, 171.0)

	// 20..
	A.Set(2, 0, 0, 0, 172.0)
	A.Set(2, 0, 1, 1, 173.0)
	A.Set(2, 0, 2, 2, 174.0)

	A.Set(2, 0, 0, 1, 175.0)
	A.Set(2, 0, 1, 2, 176.0)
	A.Set(2, 0, 0, 2, 177.0)

	A.Set(2, 0, 1, 0, 178.0)
	A.Set(2, 0, 2, 1, 179.0)
	A.Set(2, 0, 2, 0, 180.0)

	io.Pforan("A =\n%v\n", A.data.Print("%5g"))

	chk.Deep2(tst, "A", 1e-15, A.data.GetDeep2(), [][]float64{
		{100, 101, 102, 103, 104, 105, 106, 107, 108},
		{109, 110, 111, 112, 113, 114, 115, 116, 117},
		{118, 119, 120, 121, 122, 123, 124, 125, 126},
		{127, 128, 129, 130, 131, 132, 133, 134, 135},
		{136, 137, 138, 139, 140, 141, 142, 143, 144},
		{145, 146, 147, 148, 149, 150, 151, 152, 153},
		{154, 155, 156, 157, 158, 159, 160, 161, 162},
		{163, 164, 165, 166, 167, 168, 169, 170, 171},
		{172, 173, 174, 175, 176, 177, 178, 179, 180},
	})

	// 00..
	chk.Float64(tst, "A00..", 1e-15, A.Get(0, 0, 0, 0), 100.0)
	chk.Float64(tst, "A00..", 1e-15, A.Get(0, 0, 1, 1), 101.0)
	chk.Float64(tst, "A00..", 1e-15, A.Get(0, 0, 2, 2), 102.0)

	chk.Float64(tst, "A00..", 1e-15, A.Get(0, 0, 0, 1), 103.0)
	chk.Float64(tst, "A00..", 1e-15, A.Get(0, 0, 1, 2), 104.0)
	chk.Float64(tst, "A00..", 1e-15, A.Get(0, 0, 0, 2), 105.0)

	chk.Float64(tst, "A00..", 1e-15, A.Get(0, 0, 1, 0), 106.0)
	chk.Float64(tst, "A00..", 1e-15, A.Get(0, 0, 2, 1), 107.0)
	chk.Float64(tst, "A00..", 1e-15, A.Get(0, 0, 2, 0), 108.0)

	// 11..
	chk.Float64(tst, "A11..", 1e-15, A.Get(1, 1, 0, 0), 109.0)
	chk.Float64(tst, "A11..", 1e-15, A.Get(1, 1, 1, 1), 110.0)
	chk.Float64(tst, "A11..", 1e-15, A.Get(1, 1, 2, 2), 111.0)

	chk.Float64(tst, "A11..", 1e-15, A.Get(1, 1, 0, 1), 112.0)
	chk.Float64(tst, "A11..", 1e-15, A.Get(1, 1, 1, 2), 113.0)
	chk.Float64(tst, "A11..", 1e-15, A.Get(1, 1, 0, 2), 114.0)

	chk.Float64(tst, "A11..", 1e-15, A.Get(1, 1, 1, 0), 115.0)
	chk.Float64(tst, "A11..", 1e-15, A.Get(1, 1, 2, 1), 116.0)
	chk.Float64(tst, "A11..", 1e-15, A.Get(1, 1, 2, 0), 117.0)

	// 22..
	chk.Float64(tst, "A22..", 1e-15, A.Get(2, 2, 0, 0), 118.0)
	chk.Float64(tst, "A22..", 1e-15, A.Get(2, 2, 1, 1), 119.0)
	chk.Float64(tst, "A22..", 1e-15, A.Get(2, 2, 2, 2), 120.0)

	chk.Float64(tst, "A22..", 1e-15, A.Get(2, 2, 0, 1), 121.0)
	chk.Float64(tst, "A22..", 1e-15, A.Get(2, 2, 1, 2), 122.0)
	chk.Float64(tst, "A22..", 1e-15, A.Get(2, 2, 0, 2), 123.0)

	chk.Float64(tst, "A22..", 1e-15, A.Get(2, 2, 1, 0), 124.0)
	chk.Float64(tst, "A22..", 1e-15, A.Get(2, 2, 2, 1), 125.0)
	chk.Float64(tst, "A22..", 1e-15, A.Get(2, 2, 2, 0), 126.0)

	// 01..
	chk.Float64(tst, "A01..", 1e-15, A.Get(0, 1, 0, 0), 127.0)
	chk.Float64(tst, "A01..", 1e-15, A.Get(0, 1, 1, 1), 128.0)
	chk.Float64(tst, "A01..", 1e-15, A.Get(0, 1, 2, 2), 129.0)

	chk.Float64(tst, "A01..", 1e-15, A.Get(0, 1, 0, 1), 130.0)
	chk.Float64(tst, "A01..", 1e-15, A.Get(0, 1, 1, 2), 131.0)
	chk.Float64(tst, "A01..", 1e-15, A.Get(0, 1, 0, 2), 132.0)

	chk.Float64(tst, "A01..", 1e-15, A.Get(0, 1, 1, 0), 133.0)
	chk.Float64(tst, "A01..", 1e-15, A.Get(0, 1, 2, 1), 134.0)
	chk.Float64(tst, "A01..", 1e-15, A.Get(0, 1, 2, 0), 135.0)

	// 12..
	chk.Float64(tst, "A12..", 1e-15, A.Get(1, 2, 0, 0), 136.0)
	chk.Float64(tst, "A12..", 1e-15, A.Get(1, 2, 1, 1), 137.0)
	chk.Float64(tst, "A12..", 1e-15, A.Get(1, 2, 2, 2), 138.0)

	chk.Float64(tst, "A12..", 1e-15, A.Get(1, 2, 0, 1), 139.0)
	chk.Float64(tst, "A12..", 1e-15, A.Get(1, 2, 1, 2), 140.0)
	chk.Float64(tst, "A12..", 1e-15, A.Get(1, 2, 0, 2), 141.0)

	chk.Float64(tst, "A12..", 1e-15, A.Get(1, 2, 1, 0), 142.0)
	chk.Float64(tst, "A12..", 1e-15, A.Get(1, 2, 2, 1), 143.0)
	chk.Float64(tst, "A12..", 1e-15, A.Get(1, 2, 2, 0), 144.0)

	// 02..
	chk.Float64(tst, "A02..", 1e-15, A.Get(0, 2, 0, 0), 145.0)
	chk.Float64(tst, "A02..", 1e-15, A.Get(0, 2, 1, 1), 146.0)
	chk.Float64(tst, "A02..", 1e-15, A.Get(0, 2, 2, 2), 147.0)

	chk.Float64(tst, "A02..", 1e-15, A.Get(0, 2, 0, 1), 148.0)
	chk.Float64(tst, "A02..", 1e-15, A.Get(0, 2, 1, 2), 149.0)
	chk.Float64(tst, "A02..", 1e-15, A.Get(0, 2, 0, 2), 150.0)

	chk.Float64(tst, "A02..", 1e-15, A.Get(0, 2, 1, 0), 151.0)
	chk.Float64(tst, "A02..", 1e-15, A.Get(0, 2, 2, 1), 152.0)
	chk.Float64(tst, "A02..", 1e-15, A.Get(0, 2, 2, 0), 153.0)

	// 10..
	chk.Float64(tst, "A10..", 1e-15, A.Get(1, 0, 0, 0), 154.0)
	chk.Float64(tst, "A10..", 1e-15, A.Get(1, 0, 1, 1), 155.0)
	chk.Float64(tst, "A10..", 1e-15, A.Get(1, 0, 2, 2), 156.0)

	chk.Float64(tst, "A10..", 1e-15, A.Get(1, 0, 0, 1), 157.0)
	chk.Float64(tst, "A10..", 1e-15, A.Get(1, 0, 1, 2), 158.0)
	chk.Float64(tst, "A10..", 1e-15, A.Get(1, 0, 0, 2), 159.0)

	chk.Float64(tst, "A10..", 1e-15, A.Get(1, 0, 1, 0), 160.0)
	chk.Float64(tst, "A10..", 1e-15, A.Get(1, 0, 2, 1), 161.0)
	chk.Float64(tst, "A10..", 1e-15, A.Get(1, 0, 2, 0), 162.0)

	// 21..
	chk.Float64(tst, "A21..", 1e-15, A.Get(2, 1, 0, 0), 163.0)
	chk.Float64(tst, "A21..", 1e-15, A.Get(2, 1, 1, 1), 164.0)
	chk.Float64(tst, "A21..", 1e-15, A.Get(2, 1, 2, 2), 165.0)

	chk.Float64(tst, "A21..", 1e-15, A.Get(2, 1, 0, 1), 166.0)
	chk.Float64(tst, "A21..", 1e-15, A.Get(2, 1, 1, 2), 167.0)
	chk.Float64(tst, "A21..", 1e-15, A.Get(2, 1, 0, 2), 168.0)

	chk.Float64(tst, "A21..", 1e-15, A.Get(2, 1, 1, 0), 169.0)
	chk.Float64(tst, "A21..", 1e-15, A.Get(2, 1, 2, 1), 170.0)
	chk.Float64(tst, "A21..", 1e-15, A.Get(2, 1, 2, 0), 171.0)

	// 20..
	chk.Float64(tst, "A20..", 1e-15, A.Get(2, 0, 0, 0), 172.0)
	chk.Float64(tst, "A20..", 1e-15, A.Get(2, 0, 1, 1), 173.0)
	chk.Float64(tst, "A20..", 1e-15, A.Get(2, 0, 2, 2), 174.0)

	chk.Float64(tst, "A20..", 1e-15, A.Get(2, 0, 0, 1), 175.0)
	chk.Float64(tst, "A20..", 1e-15, A.Get(2, 0, 1, 2), 176.0)
	chk.Float64(tst, "A20..", 1e-15, A.Get(2, 0, 0, 2), 177.0)

	chk.Float64(tst, "A20..", 1e-15, A.Get(2, 0, 1, 0), 178.0)
	chk.Float64(tst, "A20..", 1e-15, A.Get(2, 0, 2, 1), 179.0)
	chk.Float64(tst, "A20..", 1e-15, A.Get(2, 0, 2, 0), 180.0)
}
