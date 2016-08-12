// Copyright 2015 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "stdio.h"
#include "SFMT.h"

#ifdef WIN32
#define LONG long long
#else
#define LONG long
#endif

sfmt_t GLOBAL_SFMT;

void SfmtInit(LONG seed) {
    sfmt_init_gen_rand(&GLOBAL_SFMT, seed);
}

LONG SfmtRand(LONG lo, LONG hi) {
    return (sfmt_genrand_uint64(&GLOBAL_SFMT) % (hi-lo+1) + lo);
}

void SfmtShuffle(LONG *values, LONG size) {
    uint64_t j;
    LONG tmp;
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
