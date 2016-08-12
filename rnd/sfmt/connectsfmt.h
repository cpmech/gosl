// Copyright 2015 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef GOSL_RANDOM_H
#define GOSL_RANDOM_H

#ifdef WIN32
#define LONG long long
#else
#define LONG long
#endif

void SfmtInit(LONG seed);
LONG SfmtRand(LONG lo, LONG hi);
void SfmtShuffle(LONG *values, LONG size);
void SfmtPrintIdString();

#endif // GOSL_RANDOM_H
