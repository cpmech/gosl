// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fdm

import (
    "fmt"
    "sort"
)

/* Equations
   =========
   K11 K12 => unknowns
   K21 K22 => prescribed
    |   +---> prescribed
    +-------> unknowns

   "RF" means "reduced to full" 
   "FR" means "full to reduced"
   reduced 1: is a reduced system of equations where only unknown    equations are present
   reduced 2: is a reduced system of equations where only prescribed equations are present
   full     : corresponds to all equations, unknown and prescribed

   Example:

   N   = 9         => total number of equations
   peq = [0  3  6] => prescribed equations

     0  1  2  3  4  5  6  7  8                       0   1  2   3   4  5   6   7  8
   0 +--------+--------+------ 0 -- prescribed -- 0 [22] 21 21 [22] 21 21 [22] 21 21 0
   1 |        |        |       1                  1  12  .. ..  12  .. ..  12  .. .. 1
   2 |        |        |       2                  2  12  .. ..  12  .. ..  12  .. .. 2
   3 +--------+--------+------ 3 -- prescribed -- 3 [22] 21 21 [22] 21 21 [22] 21 21 3
   4 |        |        |       4                  4  12  .. ..  12  .. ..  12  .. .. 4
   5 |        |        |       5                  5  12  .. ..  12  .. ..  12  .. .. 5
   6 +--------+--------+------ 6 -- prescribed -- 6 [22] 21 21 [22] 21 21 [22] 21 21 6
   7 |        |        |       7                  7  12  .. ..  12  .. ..  12  .. .. 7
   8 |        |        |       8                  8  12  .. ..  12  .. ..  12  .. .. 8
     0  1  2  3  4  5  6  7  8                       0   1  2   3   4  5   6   7  8
     |        |        |
    pre      pre      pre                                 .. => 11 equations

   N1 = 6 => number of equations of type 1 (unknowns)
   N2 = 3 => number of equations of type 2 (prescribed)
   N1 + N2 == N

          0  1  2  3  4  5
   RF1 = [1  2  4  5  7  8]           => ex:  RF1[3] = full system equation # 5

          0  1  2  3  4  5  6  7  8
   FR1 = [   0  1     2  3     4  5]  => ex:  FR1[5] = reduced system equation # 3
         -1       -1       -1         => indicates 'value not set'

          0  1  2
   RF2 = [0  3  6]                    => ex:  RF2[1] = full system equation # 3

          0  1  2  3  4  5  6  7  8
   FR2 = [0        1        2      ]  => ex:  FR2[3] = reduced system equation # 1
            -1 -1    -1 -1    -1 -1   => indicates 'value not set'
*/
type Equations struct {
    N1, N2, N int   // unknowns, prescribed, total numbers
    RF1, FR1  []int // reduced=>full/full=>reduced maps for unknowns
    RF2, FR2  []int // reduced=>full/full=>reduced maps for prescribed
}

func (e *Equations) Init(n int, peq_notsorted []int) {
    peq := make([]int, len(peq_notsorted))
    copy(peq, peq_notsorted)
    sort.Ints(peq)
    e.N   = n
    e.N2  = len(peq)
    e.N1  = e.N - e.N2
    e.RF1 = make([]int, e.N1)
    e.FR1 = make([]int, e.N)
    e.RF2 = make([]int, e.N2)
    e.FR2 = make([]int, e.N)
    var i1, i2 int
    for eq := 0; eq < n; eq++ {
        e.FR1[eq] = -1 // indicates 'not set'
        e.FR2[eq] = -1
        idx := sort.SearchInts(peq, eq)
        if idx < len(peq) && peq[idx] == eq { // found => prescribed
            e.RF2[i2] = eq
            e.FR2[eq] = i2
            i2 += 1
        } else { // not found => unknowns
            e.RF1[i1] = eq
            e.FR1[eq] = i1
            i1 += 1
        }
    }
}

func (e *Equations) Print() {
    fmt.Printf("N1 = %v, N2 = %v, N = %v\n", e.N1, e.N2, e.N)
    fmt.Printf("RF1 (unknown) =\n %v\n",     e.RF1)
    fmt.Printf("FR1 = \n%v\n",               e.FR1)
    fmt.Printf("RF2 (prescribed) =\n %v\n",  e.RF2)
    fmt.Printf("FR2 = \n%v\n",               e.FR2)
}
