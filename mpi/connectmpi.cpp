// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include <cstdlib>
#include <mpi.h>
#include <complex.h>

extern "C" {

void abortmpi() {
    MPI::COMM_WORLD.Abort(666);
}

int ison() {
    if (MPI::Is_initialized()) {
        return 1;
    } else {
        return 0;
    }
}

void startmpi(int debug) {
    MPI::Init();
    if (debug > 0) {
        int myid = MPI::COMM_WORLD.Get_rank();
        printf("startmpi: hi from %d. debug = %d\n", myid, debug);
    }
}

void stopmpi(int debug) {
    if (debug > 0) {
        int myid = MPI::COMM_WORLD.Get_rank();
        printf("stopmpi: goodbye from %d. debug = %d\n", myid, debug);
    }
    MPI::Finalize();
}

int mpirank() {
    return MPI::COMM_WORLD.Get_rank();
}

int mpisize() {
    return MPI::COMM_WORLD.Get_size();
}

void barrier() {
    MPI::COMM_WORLD.Barrier();
}

void sumtoroot(double *dest, double *orig, int n) {
    MPI::COMM_WORLD.Reduce(orig, dest, n, MPI::DOUBLE, MPI::SUM, 0); // 0 => dest
}

void sumtorootC(double complex *dest, double complex *orig, int n) {
    MPI::COMM_WORLD.Reduce(orig, dest, n, MPI::DOUBLE_COMPLEX, MPI::SUM, 0); // 0 => dest
}

void bcastfromroot(double *x, int n) {
    MPI::COMM_WORLD.Bcast(x, n, MPI::DOUBLE, 0); // 0 => from
}

void bcastfromrootC(double complex *x, int n) {
    MPI::COMM_WORLD.Bcast(x, n, MPI::DOUBLE_COMPLEX, 0); // 0 => from
}

void allreducesum(double *dest, double *orig, int n) {
    MPI::COMM_WORLD.Allreduce(orig, dest, n, MPI::DOUBLE, MPI::SUM);
}

void allreducemin(double *dest, double *orig, int n) {
    MPI::COMM_WORLD.Allreduce(orig, dest, n, MPI::DOUBLE, MPI::MIN);
}

void allreducemax(double *dest, double *orig, int n) {
    MPI::COMM_WORLD.Allreduce(orig, dest, n, MPI::DOUBLE, MPI::MAX);
}

void intallreducemax(int *dest, int *orig, int n) {
    MPI::COMM_WORLD.Allreduce(orig, dest, n, MPI::LONG, MPI::MAX);
}

const int TAG_SINGLEINTSENDRECV = 1000;
const int TAG_INTSENDRECV       = 1001;
const int TAG_DBLSENDRECV       = 1002;

void intsend(int *vals, int n, int to_proc) {
    MPI::COMM_WORLD.Send(vals, n, MPI::LONG, to_proc, TAG_INTSENDRECV);
}

void intrecv(int *vals, int n, int from_proc) {
    MPI::COMM_WORLD.Recv(vals, n, MPI::LONG, from_proc, TAG_INTSENDRECV);
}

void dblsend(double *vals, int n, int to_proc) {
    MPI::COMM_WORLD.Send(vals, n, MPI::DOUBLE, to_proc, TAG_DBLSENDRECV);
}

void dblrecv(double *vals, int n, int from_proc) {
    MPI::COMM_WORLD.Recv(vals, n, MPI::DOUBLE, from_proc, TAG_DBLSENDRECV);
}

} // extern "C"
