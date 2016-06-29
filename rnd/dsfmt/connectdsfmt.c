// Copyright 2015 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "dSFMT.h"

dsfmt_t GLOBAL_DSFMT;

void DsfmtInit(long seed) {
    dsfmt_init_gen_rand(&GLOBAL_DSFMT, seed);
}

double DsfmtRand(double lo, double hi) {
    return lo + (hi - lo) * dsfmt_genrand_close_open(&GLOBAL_DSFMT);
}

double DsfmtRand01() {
    return dsfmt_genrand_close_open(&GLOBAL_DSFMT);
}
