// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// OpenCV:
#include "opencv2/imgcodecs.hpp"
#include "opencv2/highgui.hpp"
using namespace cv;

// STL:
#include <map>
#include <cstddef>
#include <cstring>
#include <string>
#include <vector>
#include <exception>
#include <iostream>
using namespace std;

// Local:
#include "connectopencv.h"

// Cgo:
#include "_cgo_export.h"

ptMat new_mat() {
    try {
        Mat* pt = new Mat;
        return (void*)pt;
    } catch (const std::exception& e) {
        // returning NULL pointer
    }
    return NULL;
}

void free_mat(ptMat pt) {
    if (pt != NULL) {
        GET_MAT(mat, pt);
        delete mat;
    }
}

int new_zeros_mat(char* err, ptMat newMat, int rows, int cols, int type) {
    CV_RUN(
        GET_MAT(newmat, newMat);
        *newmat = Mat::zeros(rows, cols, type);
    )
}

int new_similar_mat(char* err, ptMat newMat, ptMat oldMat) {
    CV_RUN(
        GET_MAT(oldmat, oldMat);
        GET_MAT(newmat, newMat);
        *newmat = Mat::zeros(oldmat->size(), oldmat->type());
    )
}

int cv_imread_new(char* err, ptMat output, const char* filename) {
    CV_RUN(
        GET_MAT(mat, output);
        *mat = imread(filename);
    )
}

int cv_namedWindow(char* err, char* winname, int flags) { CV_RUN(namedWindow(winname,flags);) }
int cv_imshow(char* err, char* winname, ptMat mat) { CV_RUN(imshow(winname,REF_MAT(mat));) }
int cv_waitKey(char* err, int delay) { CV_RUN(waitKey(delay);) } 

std::vector<std::string> functionCodes;

void _trackbar_callbacks(int pos, void* userdata) {
    char* code = (char*)userdata;
    _callfromCwindowFunction(pos, code);
}

int cv_createTrackbar(char* err, const char* trackbarname, const char* winname, int* value, int count, const char* onChangeCode) {
    functionCodes.push_back(onChangeCode);
    void* code = (void*)functionCodes.back().c_str(); 
    CV_RUN(
        createTrackbar(trackbarname, winname, value, count, _trackbar_callbacks, code);
    )
}

void cv_initialise_enums(
    // WindowFlags
    int *window_normal         ,
    int *window_autosize       ,
    int *window_opengl         ,
    int *window_fullscreen     ,
    int *window_freeratio      ,
    int *window_keepratio      ,
    int *window_gui_expanded   ,
    int *window_gui_normal     ,

    // WindowPropertyFlags
    int *wnd_prop_fullscreen   ,
    int *wnd_prop_autosize     ,
    int *wnd_prop_aspect_ratio ,
    int *wnd_prop_opengl       ,
    int *wnd_prop_visible      ,

    // MouseEventTypes
    int *event_mousemove       ,
    int *event_lbuttondown     ,
    int *event_rbuttondown     ,
    int *event_mbuttondown     ,
    int *event_lbuttonup       ,
    int *event_rbuttonup       ,
    int *event_mbuttonup       ,
    int *event_lbuttondblclk   ,
    int *event_rbuttondblclk   ,
    int *event_mbuttondblclk   ,
    int *event_mousewheel      ,
    int *event_mousehwheel     ,

    // MouseEventFlags
    int *event_flag_lbutton    ,
    int *event_flag_rbutton    ,
    int *event_flag_mbutton    ,
    int *event_flag_ctrlkey    ,
    int *event_flag_shiftkey   ,
    int *event_flag_altkey     ,

    // QtFontWeights
    int *qt_font_light         ,
    int *qt_font_normal        ,
    int *qt_font_demibold      ,
    int *qt_font_bold          ,
    int *qt_font_black         ,

    // QtFontStyles
    int *qt_style_normal       ,
    int *qt_style_italic       ,
    int *qt_style_oblique      ,

    // QtButtonTypes
    int *qt_push_button        ,
    int *qt_checkbox           ,
    int *qt_radiobox           ,
    int *qt_new_buttonbar      ) {

    // WindowFlags
    *window_normal         = WINDOW_NORMAL        ;
    *window_autosize       = WINDOW_AUTOSIZE      ;
    *window_opengl         = WINDOW_OPENGL        ;
    *window_fullscreen     = WINDOW_FULLSCREEN    ;
    *window_freeratio      = WINDOW_FREERATIO     ;
    *window_keepratio      = WINDOW_KEEPRATIO     ;
    *window_gui_expanded   = WINDOW_GUI_EXPANDED  ;
    *window_gui_normal     = WINDOW_GUI_NORMAL    ;

    // WindowPropertyFlags
    *wnd_prop_fullscreen   = WND_PROP_FULLSCREEN  ;
    *wnd_prop_autosize     = WND_PROP_AUTOSIZE    ;
    *wnd_prop_aspect_ratio = WND_PROP_ASPECT_RATIO;
    *wnd_prop_opengl       = WND_PROP_OPENGL      ;
    *wnd_prop_visible      = WND_PROP_VISIBLE     ;

    // MouseEventTypes
    *event_mousemove       = EVENT_MOUSEMOVE      ;
    *event_lbuttondown     = EVENT_LBUTTONDOWN    ;
    *event_rbuttondown     = EVENT_RBUTTONDOWN    ;
    *event_mbuttondown     = EVENT_MBUTTONDOWN    ;
    *event_lbuttonup       = EVENT_LBUTTONUP      ;
    *event_rbuttonup       = EVENT_RBUTTONUP      ;
    *event_mbuttonup       = EVENT_MBUTTONUP      ;
    *event_lbuttondblclk   = EVENT_LBUTTONDBLCLK  ;
    *event_rbuttondblclk   = EVENT_RBUTTONDBLCLK  ;
    *event_mbuttondblclk   = EVENT_MBUTTONDBLCLK  ;
    *event_mousewheel      = EVENT_MOUSEWHEEL     ;
    *event_mousehwheel     = EVENT_MOUSEHWHEEL    ;

    // MouseEventFlags
    *event_flag_lbutton    = EVENT_FLAG_LBUTTON   ;
    *event_flag_rbutton    = EVENT_FLAG_RBUTTON   ;
    *event_flag_mbutton    = EVENT_FLAG_MBUTTON   ;
    *event_flag_ctrlkey    = EVENT_FLAG_CTRLKEY   ;
    *event_flag_shiftkey   = EVENT_FLAG_SHIFTKEY  ;
    *event_flag_altkey     = EVENT_FLAG_ALTKEY    ;

    // QtFontWeights
    *qt_font_light         = QT_FONT_LIGHT        ;
    *qt_font_normal        = QT_FONT_NORMAL       ;
    *qt_font_demibold      = QT_FONT_DEMIBOLD     ;
    *qt_font_bold          = QT_FONT_BOLD         ;
    *qt_font_black         = QT_FONT_BLACK        ;

    // QtFontStyles
    *qt_style_normal       = QT_STYLE_NORMAL      ;
    *qt_style_italic       = QT_STYLE_ITALIC      ;
    *qt_style_oblique      = QT_STYLE_OBLIQUE     ;

    // QtButtonTypes
    *qt_push_button        = QT_PUSH_BUTTON       ;
    *qt_checkbox           = QT_CHECKBOX          ;
    *qt_radiobox           = QT_RADIOBOX          ;
    *qt_new_buttonbar      = QT_NEW_BUTTONBAR     ;
}
