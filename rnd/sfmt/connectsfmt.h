// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef GOSL_RANDOM_H
#define GOSL_RANDOM_H

void SfmtInit(long seed);
long SfmtRand(long lo, long hi);
void SfmtShuffle(long *values, long size);

#endif // GOSL_RANDOM_H
