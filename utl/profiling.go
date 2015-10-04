// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

import (
	"flag"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

const (
	KBSIZE = 1024.0
	MBSIZE = 1048576.0
	GBSIZE = 1073741824.0
)

var do_prof_cpu = flag.Bool("cpuprof", false, "write cpu profile data to file")
var do_prof_mem = flag.Bool("memprof", false, "write mem profile data to file")

// PrintMemStat prints memory statistics
func PrintMemStat(msg string) {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	io.PfYel("%s\n", msg)
	io.Pfyel("Alloc      = %v [KB]  %v [MB]  %v [GB]\n", mem.Alloc/KBSIZE, mem.Alloc/MBSIZE, mem.Alloc/GBSIZE)
	io.Pfyel("HeapAlloc  = %v [KB]  %v [MB]  %v [GB]\n", mem.HeapAlloc/KBSIZE, mem.HeapAlloc/MBSIZE, mem.HeapAlloc/GBSIZE)
	io.Pfyel("Sys        = %v [KB]  %v [MB]  %v [GB]\n", mem.Sys/KBSIZE, mem.Sys/MBSIZE, mem.Sys/GBSIZE)
	io.Pfyel("HeapSys    = %v [KB]  %v [MB]  %v [GB]\n", mem.HeapSys/KBSIZE, mem.HeapSys/MBSIZE, mem.HeapSys/GBSIZE)
	io.Pfyel("TotalAlloc = %v [KB]  %v [MB]  %v [GB]\n", mem.TotalAlloc/KBSIZE, mem.TotalAlloc/MBSIZE, mem.TotalAlloc/GBSIZE)
	io.Pfyel("Mallocs    = %v\n", mem.Mallocs)
	io.Pfyel("Frees      = %v\n", mem.Frees)
}

// ProfCPU activates CPU profiling
//  Note: returns a "stop()" function to be called before shutting down
func ProfCPU(dirout, filename string, silent bool) func() {
	os.MkdirAll(dirout, 0777)
	fn := filepath.Join(dirout, filename)
	f, err := os.Create(fn)
	if err != nil {
		chk.Panic(_profiling_err1, "ProfCPU", err.Error())
	}
	if !silent {
		io.Pfcyan("CPU profiling => %s\n", fn)
	}
	pprof.StartCPUProfile(f)
	return func() {
		pprof.StopCPUProfile()
		f.Close()
		if !silent {
			io.Pfcyan("CPU profiling finished\n")
		}
	}
}

// ProfMEM activates memory profiling
//  Note: returns a "stop()" function to be called before shutting down
func ProfMEM(dirout, filename string, silent bool) func() {
	os.MkdirAll(dirout, 0777)
	fn := filepath.Join(dirout, filename)
	f, err := os.Create(fn)
	if err != nil {
		chk.Panic(_profiling_err1, "ProfMEM", err.Error())
	}
	if !silent {
		io.Pfcyan("MEM profiling => %s\n", fn)
	}
	return func() {
		pprof.WriteHeapProfile(f)
		f.Close()
		if !silent {
			io.Pfcyan("MEM profiling finished\n")
		}
	}
}

// DoProf runs either CPU profiling or MEM profiling
//  Notes:
//    1) based on input flags: -cpuprof (preferred); or
//                             -memprof
//    2) returns a "stop()" function to be called before shutting down
//    3) output files are saved to "/tmp/gosl/cpu.pprof"; or
//                                 "/tmp/gosl/mem.pprof"
//  Run analysis with (e.g.):
//   go tool pprof binary /tmp/gosl/mem.pprof
func DoProf(silent bool) func() {
	if *do_prof_cpu {
		return ProfCPU("/tmp/gosl/", "cpu.pprof", silent)
	} else if *do_prof_mem {
		return ProfMEM("/tmp/gosl/", "mem.pprof", silent)
	}
	return func() {}
}

// error messages
var (
	_profiling_err1 = "profiling.go: %s: cannot create file:\n%v"
)
