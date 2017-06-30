# Gosl Benchmark

## Summary

* Matrix-Matrix multiplication: OpenBLAS, MKL and Naïve approach
* Matrix-Matrix multiplication: OpenBLAS versus Naïve approach

These tests were run on:
* Ubuntu 16.04.2 LTS (from lsb_release -a)
* Intel(R) Core(TM) i7-4770 CPU @ 3.40GHz (from cat /proc/cpuinfo) and "standard flags"
with default OpenBLAS compilation.


# Matrix-Matrix multiplication: OpenBLAS, MKL and Naïve approach

Test ran at: 2017 Jun 23

Source code 1: <a href="oblas_dgemm01.go">oblas_dgemm01.go</a>
Source code 2: <a href="mkl_dgemm01.go">mkl_dgemm01.go</a>

NOTE: must use different files to avoid using the same library (e.g. due to linking)

Code for plotting: <a href="plot_mkl_oblas_naive.go">plot_mkl_oblas_naive.go</a>

## Larger matrices

Single-threaded

<div id="container">
<p><img src="figs/mkl-oblas-comparison-large.png" width="600"></p>
</div>

4 threads

<div id="container">
<p><img src="figs/mkl-oblas-comparison-large-4t.png" width="600"></p>
</div>

## Smaller matrices

Single-threaded

<div id="container">
<p><img src="figs/mkl-oblas-comparison-small.png" width="600"></p>
</div>

4 threads

<div id="container">
<p><img src="figs/mkl-oblas-comparison-small-4t.png" width="600"></p>
</div>



# Matrix-Matrix multiplication: OpenBLAS versus Naïve approach

Test ran at: 2017 Jun 23

Source code: <a href="oblas_dgemm01.go">oblas_dgemm01.go</a>

## Larger matrices

<div id="container">
<p><img src="figs/oblas-dgemm01b.png" width="600"></p>
</div>

Output (larger matrices):
```
   size   ┃     OpenBLAS dgemm     (Dt) ┃          naïve           (naiveDt) ┃ naiveDt/Dt
━━━━━━━━━━╋━━━━━━━━━━━━━━━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━
  16×  16 ┃  8.71 GFlops (       940ns) ┃ naive:  0.89 GFlops (     9.159µs) ┃ 9.744 
  80×  80 ┃ 27.10 GFlops (    37.787µs) ┃ naive:  0.85 GFlops (  1.211082ms) ┃ 32.050 
 144× 144 ┃ 35.45 GFlops (   168.442µs) ┃ naive:  0.79 GFlops (  7.536701ms) ┃ 44.744 
 208× 208 ┃ 38.53 GFlops (   467.142µs) ┃ naive:  0.73 GFlops ( 24.490851ms) ┃ 52.427 
 272× 272 ┃ 41.79 GFlops (   963.097µs) ┃ naive:  0.78 GFlops ( 51.918322ms) ┃ 53.908 
 336× 336 ┃ 43.48 GFlops (  1.744716ms) ┃ naive:  0.72 GFlops (105.320761ms) ┃ 60.366 
 400× 400 ┃ 44.65 GFlops (  2.866527ms) ┃ naive:  0.72 GFlops (178.326095ms) ┃ 62.210 
 464× 464 ┃ 45.44 GFlops (  4.396533ms) ┃ naive:  0.65 GFlops (309.134597ms) ┃ 70.313 
 528× 528 ┃ 44.20 GFlops (    6.6598ms) ┃ naive:  0.71 GFlops (412.563052ms) ┃ 61.948 
 592× 592 ┃ 44.25 GFlops (  9.378083ms) ┃ naive:  0.71 GFlops (585.060299ms) ┃ 62.386 
 656× 656 ┃ 44.99 GFlops ( 12.549349ms) ┃ naive:  0.32 GFlops (1.739039722s) ┃ 138.576 
 720× 720 ┃ 45.78 GFlops ( 16.306108ms) ┃ naive:  0.67 GFlops (1.109812566s) ┃ 68.061 
 784× 784 ┃ 45.55 GFlops ( 21.159115ms) ┃ naive: N/A                         ┃ N/A
 848× 848 ┃ 45.88 GFlops ( 26.582728ms) ┃ naive: N/A                         ┃ N/A
 912× 912 ┃ 46.49 GFlops ( 32.634519ms) ┃ naive: N/A                         ┃ N/A
 976× 976 ┃ 46.63 GFlops (  39.87703ms) ┃ naive: N/A                         ┃ N/A
1040×1040 ┃ 45.89 GFlops ( 49.024206ms) ┃ naive: N/A                         ┃ N/A
1104×1104 ┃ 46.25 GFlops ( 58.192775ms) ┃ naive: N/A                         ┃ N/A
1168×1168 ┃ 46.56 GFlops ( 68.441843ms) ┃ naive: N/A                         ┃ N/A
1232×1232 ┃ 46.73 GFlops ( 80.037818ms) ┃ naive: N/A                         ┃ N/A
1296×1296 ┃ 46.74 GFlops (  93.15229ms) ┃ naive: N/A                         ┃ N/A
1360×1360 ┃ 46.27 GFlops (108.736216ms) ┃ naive: N/A                         ┃ N/A
```

## Small matrices

<div id="container">
<p><img src="figs/oblas-dgemm01a.png" width="600"></p>
</div>

Output (small matrices):
```
   size   ┃     OpenBLAS dgemm     (Dt) ┃          naïve           (naiveDt) ┃ naiveDt/Dt
━━━━━━━━━━╋━━━━━━━━━━━━━━━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━
   2×   2 ┃  0.04 GFlops (       389ns) ┃ naive:  0.35 GFlops (        46ns) ┃ 0.118 
   4×   4 ┃  0.29 GFlops (       442ns) ┃ naive:  0.43 GFlops (       300ns) ┃ 0.679 
   6×   6 ┃  0.78 GFlops (       552ns) ┃ naive:  0.52 GFlops (       835ns) ┃ 1.513 
   8×   8 ┃  1.67 GFlops (       615ns) ┃ naive:  0.55 GFlops (     1.871µs) ┃ 3.042 
  10×  10 ┃  2.31 GFlops (       866ns) ┃ naive:  0.55 GFlops (     3.607µs) ┃ 4.165 
  12×  12 ┃  3.47 GFlops (       995ns) ┃ naive:  0.57 GFlops (     6.067µs) ┃ 6.097 
  14×  14 ┃  3.96 GFlops (     1.385µs) ┃ naive:  0.42 GFlops (    13.198µs) ┃ 9.529 
  16×  16 ┃  5.74 GFlops (     1.428µs) ┃ naive:  0.59 GFlops (    13.912µs) ┃ 9.742 
  18×  18 ┃  5.95 GFlops (      1.96µs) ┃ naive:  0.42 GFlops (    27.639µs) ┃ 14.102 
  20×  20 ┃  8.98 GFlops (     1.781µs) ┃ naive:  0.96 GFlops (    16.593µs) ┃ 9.317 
  22×  22 ┃ 12.65 GFlops (     1.684µs) ┃ naive:  0.84 GFlops (    25.442µs) ┃ 15.108 
  24×  24 ┃ 15.19 GFlops (      1.82µs) ┃ naive:  0.87 GFlops (    31.848µs) ┃ 17.499 
  26×  26 ┃ 14.28 GFlops (     2.461µs) ┃ naive:  0.87 GFlops (    40.324µs) ┃ 16.385 
  28×  28 ┃ 16.06 GFlops (     2.734µs) ┃ naive:  0.89 GFlops (    49.472µs) ┃ 18.095 
  30×  30 ┃ 15.59 GFlops (     3.464µs) ┃ naive:  0.87 GFlops (    61.833µs) ┃ 17.850 
  32×  32 ┃ 17.44 GFlops (     3.757µs) ┃ naive:  0.86 GFlops (     76.38µs) ┃ 20.330 
  34×  34 ┃ 16.88 GFlops (     4.656µs) ┃ naive:  0.87 GFlops (    90.112µs) ┃ 19.354 
  36×  36 ┃ 18.46 GFlops (     5.055µs) ┃ naive:  0.83 GFlops (   112.196µs) ┃ 22.195 
  38×  38 ┃ 17.83 GFlops (     6.154µs) ┃ naive:  0.84 GFlops (   131.115µs) ┃ 21.306 
  40×  40 ┃ 18.30 GFlops (     6.995µs) ┃ naive:  0.84 GFlops (   152.936µs) ┃ 21.864 
  42×  42 ┃ 19.41 GFlops (     7.635µs) ┃ naive:  0.84 GFlops (   176.249µs) ┃ 23.084 
  44×  44 ┃ 21.06 GFlops (     8.088µs) ┃ naive:  0.84 GFlops (   201.707µs) ┃ 24.939 
  46×  46 ┃ 20.07 GFlops (     9.702µs) ┃ naive:  0.85 GFlops (   229.871µs) ┃ 23.693 
  48×  48 ┃ 23.89 GFlops (      9.26µs) ┃ naive:  0.85 GFlops (   261.616µs) ┃ 28.252 
  50×  50 ┃ 21.37 GFlops (      11.7µs) ┃ naive:  0.85 GFlops (    295.49µs) ┃ 25.256 
  52×  52 ┃ 22.86 GFlops (    12.304µs) ┃ naive:  0.85 GFlops (   331.815µs) ┃ 26.968 
  54×  54 ┃ 21.47 GFlops (     14.67µs) ┃ naive:  0.85 GFlops (   372.515µs) ┃ 25.393 
  56×  56 ┃ 25.66 GFlops (    13.686µs) ┃ naive:  0.85 GFlops (   415.259µs) ┃ 30.342 
  58×  58 ┃ 24.15 GFlops (    16.161µs) ┃ naive:  0.85 GFlops (   460.571µs) ┃ 28.499 
  60×  60 ┃ 25.59 GFlops (    16.879µs) ┃ naive:  0.85 GFlops (   509.182µs) ┃ 30.167 
  62×  62 ┃ 24.02 GFlops (    19.842µs) ┃ naive:  0.84 GFlops (   566.528µs) ┃ 28.552 
  64×  64 ┃ 26.61 GFlops (    19.703µs) ┃ naive:  0.84 GFlops (   623.595µs) ┃ 31.650 
```
