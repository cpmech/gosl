// Copyright 2015 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "stdio.h"
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

void SfmtPrintIdString() {
    printf("%s\n", sfmt_get_idstring(&GLOBAL_SFMT));
}
