// Copyright 2017 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

double fcn(double* x, int* fid) {
    extern double gofcn(double, int);
    return gofcn(*x, *fid);
}
