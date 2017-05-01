// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef AUXILIARY_H
#define AUXILIARY_H

#ifdef __cplusplus
extern "C" {
#endif

#include "stdlib.h"

// allocate memory for an array argv
static inline char** make_argv(int argc) {
	return (char**)malloc(sizeof(char*) * argc);
}

// set argument in array argv
static inline void set_arg(char** argv, int i, char* str) {
	argv[i] = str;
}

// allocate memory for a char buffer
static inline char* make_buffer(int size) {
	return (char*)malloc(size * sizeof(char));
}

// free char buffer
static inline void free_buffer(char* buf) {
	free(buf);
}

// allocate memory for an array of chars
static inline char** make_array_char(int nitems, int nchars) {
	char** A = (char**)malloc(nitems * sizeof(char*));
	int i = 0;
	for (i = 0; i < nitems; i++) {
		A[i] = make_buffer(nchars);
	}
	return A;
}

// free array of chars memory
static inline void free_array_char(char** A, int nitems) {
	int i = 0;
	for (i = 0; i < nitems; i++) {
		free(A[i]);
	}
	free(A);
}

// return an item of array of chars
static inline char* array_item(char** A, int i) {
	return A[i];
}

#ifdef __cplusplus
} /* extern "C" */
#endif

#endif // AUXILIARY_H
