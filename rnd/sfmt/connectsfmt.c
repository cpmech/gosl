// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "SFMT.h"

sfmt_t GLOBAL_SFMT;

void SfmtInit(long seed) {
    sfmt_init_gen_rand(&GLOBAL_SFMT, seed);
}

long SfmtRand(long lo, long hi) {
    return (sfmt_genrand_uint64(&GLOBAL_SFMT) % (hi-lo+1) + lo);
}

void SfmtShuffle(long *values, long size) {
    uint64_t j;
    long tmp;
    for (uint64_t i=size-1; i>0; i--) {
        j = sfmt_genrand_uint64(&GLOBAL_SFMT) % i;
        tmp = values[j];
        values[j] = values[i];
        values[i] = tmp;
    }
}
