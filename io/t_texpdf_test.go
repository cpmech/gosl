// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package io

import (
	"bytes"
	"testing"
	"time"

	"github.com/cpmech/gosl/chk"
)

func Test_texpdf01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("texpdf01")

	l1 := TexNum("", 123.456, true)
	l2 := TexNum("", 123.456e-8, true)
	l3 := TexNum("%.1e", 123.456e+8, true)
	l4 := TexNum("%.2e", 123.456e-8, false)
	Pforan("l1 = %v\n", l1)
	Pforan("l2 = %v\n", l2)
	Pforan("l3 = %v\n", l3)
	Pforan("l4 = %v\n", l4)
	chk.String(tst, l1, "123.456")
	chk.String(tst, l2, "1.23456\\cdot 10^{-6}")
	chk.String(tst, l3, "1.2\\cdot 10^{10}")
	chk.String(tst, l4, "1.23e-06")

	Pl()
	var d1 time.Duration = 123.081823e6
	Pf("d1 = %v\n", d1)
	d1round := RoundDuration(d1, -1)
	d2round := RoundDuration(d1, 1e6)
	d3round := RoundDuration(d1, 1e5)
	chk.String(tst, Sf("%v", d1round), "123ms")
	chk.String(tst, Sf("%v", d2round), "123ms")
	chk.String(tst, Sf("%v", d3round), "123.1ms")
	var d2 time.Duration = -123.081823e6
	d4round := RoundDuration(d2, 1e6)
	chk.String(tst, Sf("%v", d4round), "-123ms")
}

func Test_texpdf02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("texpdf02")

	keys, T := ReadTable("data/table01.dat")

	key2tex := map[string]string{
		"a": `$a = \int x \, \mathrm{d}x$`,
		"b": `$b$`,
		"c": `cval`,
		"d": `$d = \sum_{i=0}^{n} v$`,
	}

	key2convert := map[string]FcnConvertNum{
		"a": func(i int, x float64) string { return Sf("a:%g", x) },
		"b": func(i int, x float64) string { return Sf("b:%g", x) },
		"c": func(i int, x float64) string { return Sf("c:%g", x) },
		"d": func(i int, x float64) string { return Sf("d:%g", x) },
	}

	rpt := Report{
		Title:  "Gosl test",
		Author: "Gosl authors",
	}

	if !chk.Verbose {
		rpt.DoNotGeneratePDF = true
	}

	rpt.AddSection("Introduction", 0)
	rpt.AddTex("In this test, we add one table and one equation to the LaTeX document.")
	rpt.AddTex("Then we generate a PDF files in the temporary directory.")
	rpt.AddTex("The numbers in the rows of the table have a fancy format.")

	rpt.AddSection("MyTable", 1)
	rpt.AddTable("Results from simulation.", "results", "", keys, T, key2tex, key2convert)

	rpt.AddSection("Extra", 3)
	extra := new(bytes.Buffer)
	Ff(extra, `\begin{equation}`+"\n")
	Ff(extra, `\sigma = E \, \varepsilon`+"\n")
	Ff(extra, `\end{equation}`)

	rpt.WriteTexPdf("/tmp/gosl", "test_texpdf02", extra)
}

func Test_texpdf03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("texpdf03")

	keys, T := ReadTable("data/table01.dat")

	key2tex := map[string]string{
		"a": `$a^\star$`,
		"b": `$b$`,
		"c": `cval`,
		"d": `$d^\circ$`,
	}

	F := map[string]FcnRow{
		"a": func(i int) string { return Sf("a:%g", T["a"][i]) },
		"b": func(i int) string { return Sf("b:%g", T["b"][i]) },
		"c": func(i int) string { return Sf("c:%g", T["c"][i]) },
		"d": func(i int) string { return Sf("d:%g", T["d"][i]) },
	}

	rpt := Report{
		Title:  "Gosl test",
		Author: "Gosl authors",
	}

	if !chk.Verbose {
		rpt.DoNotGeneratePDF = true
	}

	notes := "$^\\star$First value\n"
	notes += "$^\\circ$Last value"

	nrows := len(T["a"])
	rpt.AddTableF("Results from simulation.", "results", notes, keys, nrows, F, key2tex)

	T["a"] = append(T["a"], -1)
	T["b"] = append(T["b"], -2)
	T["c"] = append(T["c"], -3)
	T["d"] = append(T["d"], -4)
	nrows = len(T["a"])

	rpt.RowGapPt = 20
	rpt.RowGapStep = 2
	rpt.AddTableF("Results from simulation (again).", "resultsAgain", notes, keys, nrows, F, key2tex)

	rpt.WriteTexPdf("/tmp/gosl", "test_texpdf03", nil)
	rpt.WriteTexTables("/tmp/gosl", map[string]string{"resultsAgain": "test_texpdf03_results"})
}
