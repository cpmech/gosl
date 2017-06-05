// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package msh

// FuncTet4 calculates the shape functions (S) and derivatives of shape functions (dSdR) of tet4
// elements at {r,s,t} natural coordinates. The derivatives are calculated only if derivs==true.
//
//                  t
//                  |
//                  3
//                 /|`.
//                 ||  `,
//                / |    ',
//                | |      \
//               /  |       `.
//               |  |         `,
//              /   |           `,
//              |   |             \
//             /    |              `.
//             |    |                ',
//            /     |                  \
//            |     0.,,_               `.
//           |     /     ``'-.,,__        `.
//           |    /              ``''-.,,_  ',
//          |    /                        `` 2 ,,s
//          |  ,'                       ,.-``
//         |  ,                    _,-'`
//         ' /                 ,.'`
//        | /             _.-``
//        '/          ,-'`
//       |/      ,.-``
//       /  _,-``
//      1 '`
//     /
//    r
//
func FuncTet4(S []float64, dSdR [][]float64, R []float64, derivs bool) {

	r, s, t := R[0], R[1], R[2]
	S[0] = 1.0 - r - s - t
	S[1] = r
	S[2] = s
	S[3] = t

	if !derivs {
		return
	}

	dSdR[0][0] = -1.0
	dSdR[1][0] = 1.0
	dSdR[2][0] = 0.0
	dSdR[3][0] = 0.0

	dSdR[0][1] = -1.0
	dSdR[1][1] = 0.0
	dSdR[2][1] = 1.0
	dSdR[3][1] = 0.0

	dSdR[0][2] = -1.0
	dSdR[1][2] = 0.0
	dSdR[2][2] = 0.0
	dSdR[3][2] = 1.0
}

// FuncTet10 calculates the shape functions (S) and derivatives of shape functions (dSdR) of tet10
// elements at {r,s,t} natural coordinates. The derivatives are calculated only if derivs==true.
//
//                  t
//                  |
//                  3
//                 /|`.
//                 ||  `,
//                / |    ',
//                | |      \
//               /  |       `.
//               |  |         `,
//              /   7            9
//              |   |             \
//             /    |              `.
//             |    |                ',
//            8     |                  \
//            |     0 ,,_               `.
//           |     /     ``'-., 6         `.
//           |    /               `''-.,,_  ',
//          |    /                        ``'2 ,,s
//          |   '                       ,.-``
//         |   4                   _,-'`
//         ' /                 ,.'`
//        | /             _ 5 `
//        '/          ,-'`
//       |/      ,.-``
//       /  _,-``
//      1 '`
//     /
//    r
//
func FuncTet10(S []float64, dSdR [][]float64, R []float64, derivs bool) {

	r, s, t := R[0], R[1], R[2]
	u := 1.0 - r - s - t
	S[0] = u * (2.0*u - 1.0)
	S[1] = r * (2.0*r - 1.0)
	S[2] = s * (2.0*s - 1.0)
	S[3] = t * (2.0*t - 1.0)
	S[4] = 4.0 * u * r
	S[5] = 4.0 * r * s
	S[6] = 4.0 * s * u
	S[7] = 4.0 * u * t
	S[8] = 4.0 * r * t
	S[9] = 4.0 * s * t

	if !derivs {
		return
	}

	dSdR[0][0] = 4.0*(r+s+t) - 3.0
	dSdR[1][0] = 4.0*r - 1.0
	dSdR[2][0] = 0.0
	dSdR[3][0] = 0.0
	dSdR[4][0] = 4.0 - 8.0*r - 4.0*s - 4.0*t
	dSdR[5][0] = 4.0 * s
	dSdR[6][0] = -4.0 * s
	dSdR[7][0] = -4.0 * t
	dSdR[8][0] = 4.0 * t
	dSdR[9][0] = 0.0

	dSdR[0][1] = 4.0*(r+s+t) - 3.0
	dSdR[1][1] = 0.0
	dSdR[2][1] = 4.0*s - 1.0
	dSdR[3][1] = 0.0
	dSdR[4][1] = -4.0 * r
	dSdR[5][1] = 4.0 * r
	dSdR[6][1] = 4.0 - 4.0*r - 8.0*s - 4.0*t
	dSdR[7][1] = -4.0 * t
	dSdR[8][1] = 0.0
	dSdR[9][1] = 4.0 * t

	dSdR[0][2] = 4.0*(r+s+t) - 3.0
	dSdR[1][2] = 0.0
	dSdR[2][2] = 0.0
	dSdR[3][2] = 4.0*t - 1.0
	dSdR[4][2] = -4.0 * r
	dSdR[5][2] = 0.0
	dSdR[6][2] = -4.0 * s
	dSdR[7][2] = 4.0 - 4.0*r - 4.0*s - 8.0*t
	dSdR[8][2] = 4.0 * r
	dSdR[9][2] = 4.0 * s
}
