#!/bin/bash

mkdir -p /tmp/gosl
gcc test_mat.c test_mat_prb.c -lm -o /tmp/gosl/test_mat && /tmp/gosl/test_mat
