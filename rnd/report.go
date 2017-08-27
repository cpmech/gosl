// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnd

import (
	"bytes"
	"strings"

	"github.com/cpmech/gosl/io"
)

// SetOfVars defines a set of random variables
type SetOfVars struct {
	Name string
	Vars []*Variable
}

// SetsOfVars defines a set of sets of random variables
type SetsOfVars []*SetOfVars

// ReportVariables generates TeX report of sets of variables
func ReportVariables(dirout, fnkey string, sets SetsOfVars, genPDF bool) {

	// table header
	buf := new(bytes.Buffer)
	io.Ff(buf, `
\begin{table} \centering
\caption{Random variables.}

\scriptsize

\begin{tabular}[c]{ccccccc} \toprule
name & var & $\mu$ & $\sigma$ & D$^{\star}$ & min & max \\ \hline
`)

	// generate table
	for i, set := range sets {
		for j, v := range set.Vars {
			key := ""
			if j == 0 {
				key = strings.Replace(set.Name, "/", "-", -1)
				key = strings.Replace(key, "_", "-", -1)
			}
			txtM, txtS := "-", "-"
			if v.D != "U" {
				txtM = "$" + io.TexNum("", v.M, true) + "$"
				txtS = "$" + io.TexNum("", v.S, true) + "$"
			}
			io.Ff(buf, `%s & $x_{%d}$ & %s & %s & %s & $%s$ & $%s$ \\`, key, j, txtM, txtS, v.D, io.TexNum("", v.Min, true), io.TexNum("", v.Max, true))
			io.Ff(buf, "\n")
		}
		if i < len(sets)-1 {
			io.Ff(buf, `\multicolumn{7}{l}{} \\`)
			io.Ff(buf, "\n")
		} else {
			io.Ff(buf, " \\hline\n\n")
		}
	}

	// table footer
	io.Ff(buf, `
\multicolumn{7}{p{7cm}}{
	\scriptsize
	$^{\star}$N:Normal, L:Lognormal, G:Gumbel, F:Frechet, U:Uniform
} \\

\bottomrule
\end{tabular}
\label{tab:prms%s}
\end{table}
`, fnkey)

	// write table
	tex := fnkey + ".tex"
	io.WriteFileVD(dirout, tex, buf)

	// generate PDF
	if genPDF {

		// header
		header := new(bytes.Buffer)
		io.Ff(header, `\documentclass[a4paper,twocolumn]{article}

\usepackage{amsmath}
\usepackage{amssymb}
\usepackage{booktabs}

\usepackage[margin=1.5cm,footskip=0.5cm]{geometry}

\title{Gosl-rnd Report: Random Variables}
\author{The Author}

\begin{document}`)

		// footer
		footer := new(bytes.Buffer)
		io.Ff(footer, `
\end{document}`)

		// write temporary TeX file
		tex = "tmp_" + tex
		io.WriteFileD(dirout, tex, header, buf, footer)

		// run pdflatex
		_, err := io.RunCmd(false, "pdflatex", "-interaction=batchmode", "-halt-on-error", "-output-directory="+dirout, tex)
		if err != nil {
			io.PfRed("pdflatex failed: %v\n", err)
			return
		}
		io.PfBlue("file <%s/tmp_%s.pdf> generated\n", dirout, fnkey)
	}
}
