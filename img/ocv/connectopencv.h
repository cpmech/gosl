// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef CONNECTOPENCV_H
#define CONNECTOPENCV_H

#ifdef __cplusplus
extern "C" {
#endif

// max number of characters for error message text from C++ to Go
#define ERROR_BUFFER_SIZE 1024

// pointers
typedef void* ptMat;
typedef void* ptUserData;

// constructors
ptMat new_mat();

// free memory at pointers
void free_mat(ptMat);

// macros to convert pointers
#define PT_MAT(pt) (cv::Mat*)pt
#define REF_MAT(pt) *((cv::Mat*)pt)
#define GET_MAT(ptOut, ptIn) cv::Mat* ptOut = (cv::Mat*)ptIn;

// macro to help with calling OpenCV and handling exception errors
//   err -- error message buffer
#define CV_RUN(commands)                                          \
    try {                                                         \
        commands                                                  \
    } catch (const std::exception& e) {                           \
        std::string message = std::string("ERROR: ") + e.what();  \
        std::strncpy(err, message.c_str(), ERROR_BUFFER_SIZE);    \
        return 1;                                                 \
    }                                                             \
    return 0;

// All functions return:
//  0 = OK
//  1 = failed; e.g. exception happend => see 'err'

int new_zeros_mat(char* err, ptMat newMat, int rows, int cols, int type);
int new_similar_mat(char* err, ptMat newMat, ptMat oldMat);
int cv_imread_new(char* err, ptMat output, const char* filename);
int cv_namedWindow(char* err, char* winname, int flags);
int cv_imshow(char* err, char* winname, ptMat mat);
int cv_waitKey(char* err, int delay);
int cv_createTrackbar(char* err, const char* trackbarname, const char* winname, int* value, int count, const char* onChangeFunctionName);

// enums
void cv_initialise_enums(int*, int*, int*, int*, int*, int*, int*, int*, int*, int*, int*, int*, int*, int*, int*, int*, int*, int*, int*, int*, int*, int*, int*, int*, int*, int*, int*, int*, int*, int*, int*, int*, int*, int*, int*, int*, int*, int*, int*, int*, int*, int*, int*);

#ifdef __cplusplus
} // extern "C"
#endif

#endif // CONNECTOPENCV_H
