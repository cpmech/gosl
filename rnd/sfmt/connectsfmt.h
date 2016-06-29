// Copyright 2015 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef GOSL_RANDOM_H
#define GOSL_RANDOM_H

void SfmtInit(long seed);
long SfmtRand(long lo, long hi);
void SfmtShuffle(long *values, long size);
void SfmtPrintIdString();

#endif // GOSL_RANDOM_H
