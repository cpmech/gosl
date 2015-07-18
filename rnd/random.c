// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "SFMT.h"

sfmt_t GLOBAL_SFMT;

void Init(long seed) {
    sfmt_init_gen_rand(&GLOBAL_SFMT, seed);
}

long IntRand(long lo, long hi) {
    return (sfmt_genrand_uint64(&GLOBAL_SFMT) % (hi-lo+1) + lo);
}
