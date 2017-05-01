// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// package ocv connects with the OpenCV library for image analyses and machine learning, in
// particular for computer vision algorithms.
package ocv

/*
#cgo linux   CFLAGS: -O2 -I/usr/include
#cgo linux   CFLAGS: -O2 -I/usr/local/include
#cgo windows CFLAGS: -O2 -IC:/Gosl/include
#cgo linux   LDFLAGS: -lopencv_core -L/local/lib
#cgo darwin  LDFLAGS: -lopencv_core -L/usr/local/lib
#cgo windows LDFLAGS: -lopencv_core -LC:/Gosl/lib
#ifdef WIN32
#define LONG long long
#else
#define LONG long
#endif

#include "auxiliary.h"
#include "connectopencv.h"
*/
import "C"

import (
	"unsafe"

	"github.com/cpmech/gosl/io"
)

// Mat implements the OpenCV Mat structure
type Mat struct {
	pt  C.ptMat // pointer to OpenCV Mat
	err *C.char // buffer to hold error messages from C++
}

// NewMat allocates a new OpenCV Mat structure
func NewMat() (o *Mat) {
	o = new(Mat)
	o.pt = C.new_mat()
	o.err = C.make_buffer(C.ERROR_BUFFER_SIZE)
	return
}

// Free deallocates Mat object
func (o *Mat) Free() {
	if o.pt != nil {
		C.free_mat(o.pt)
		o.pt = nil
	}
	if o.err != nil {
		C.free_buffer(o.err)
		o.err = nil
	}
}

// NewMatFromFile creates a new Mat with image data from file
func NewMatFromFile(filename string) (o *Mat) {

	// new structure
	o = new(Mat)
	o.pt = C.new_mat()
	o.err = C.make_buffer(C.ERROR_BUFFER_SIZE)

	// c-strings
	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))

	// call OpenCV
	status := C.cv_imread_new(o.err, o.pt, cfilename)
	checkWithPanic(status, o.err)
	return
}

// NewZerosMat creates a new Mat filled with zeros
func NewZerosMat(rows, cols, Type int) (o *Mat) {

	// new structure
	o = new(Mat)
	o.pt = C.new_mat()
	o.err = C.make_buffer(C.ERROR_BUFFER_SIZE)

	// call OpenCV
	status := C.new_zeros_mat(o.err, o.pt, C.int(rows), C.int(cols), C.int(Type))
	checkWithPanic(status, o.err)
	return
}

// NewSimilarMat creates a new Mat with the same size and type as 'mat', but filled with zeros
func NewSimilarMat(mat *Mat) (o *Mat) {

	// new structure
	o = new(Mat)
	o.pt = C.new_mat()
	o.err = C.make_buffer(C.ERROR_BUFFER_SIZE)

	// call OpenCV
	status := C.new_similar_mat(o.err, o.pt, mat.pt)
	checkWithPanic(status, o.err)
	return
}

// window /////////////////////////////////////////////////////////////////////////////////////////

// NamedWindow creates a window
func NamedWindow(winName string, flags int) {

	// c-strings
	cwinName := C.CString(winName)
	defer C.free(unsafe.Pointer(cwinName))

	// error string
	cerr := C.make_buffer(C.ERROR_BUFFER_SIZE)
	defer C.free_buffer(cerr)

	// call OpenCV
	status := C.cv_namedWindow(cerr, cwinName, C.int(flags))
	checkWithPanic(status, cerr)
}

// Imshow Displays an image in the specified window.
func Imshow(winName string, mat *Mat) {

	// c-strings
	cwinName := C.CString(winName)
	defer C.free(unsafe.Pointer(cwinName))

	// call OpenCV
	status := C.cv_imshow(cwinName, mat.err, mat.pt)
	checkWithPanic(status, mat.err)
}

// WaitKey waits for a pressed key.
//   delay -- Delay in milliseconds. 0 is the special value that means "forever".
func WaitKey(delay int) {

	// error string
	cerr := C.make_buffer(C.ERROR_BUFFER_SIZE)
	defer C.free_buffer(cerr)

	// call OpenCV
	status := C.cv_waitKey(cerr, C.int(delay))
	checkWithPanic(status, cerr)
}

// CreateTrackbar creates a trackbar and attaches it to the specified window.
//   trackbarName -- Name of the created trackbar.
//   winName -- Name of the window that will be used as a parent of the created trackbar.
//   value -- Optional pointer to an integer variable whose value reflects the position of the slider. Upon creation, the slider position is defined by this variable.
//   count -- Maximum position of the slider. The minimal position is always 0.
//   onChange -- Pointer to the function to be called every time the slider changes position.
func CreateTrackbar(trackbarName, winName string, value *int, count int, onChange TrackbarCallback) {

	// c-strings
	ctrackbarName := C.CString(trackbarName)
	defer C.free(unsafe.Pointer(ctrackbarName))
	cwinName := C.CString(winName)
	defer C.free(unsafe.Pointer(cwinName))

	// error string
	cerr := C.make_buffer(C.ERROR_BUFFER_SIZE)
	defer C.free_buffer(cerr)

	// c-variables
	cvalue := (*C.int)(unsafe.Pointer(value))

	// callback
	code := winName + "." + trackbarName + "." + io.Sf("%d", len(trackbarCallbacks))
	trackbarCallbacks[code] = onChange
	ccode := C.CString(code)
	defer C.free(unsafe.Pointer(ccode))

	// call OpenCV
	status := C.cv_createTrackbar(cerr, ctrackbarName, cwinName, cvalue, C.int(count), ccode)
	checkWithPanic(status, cerr)
}

// functions to be called by C++ //////////////////////////////////////////////////////////////////

// trackbarCallbacks holds callback functions by name
var trackbarCallbacks = make(map[string]TrackbarCallback)

// _callfromCwindowFunction will be called from C code
//export _callfromCwindowFunction
func _callfromCwindowFunction(cpos C.int, ccode *C.char) {
	code := C.GoString(ccode)
	if fcn, ok := trackbarCallbacks[code]; ok {
		if fcn != nil {
			fcn(int(cpos))
		}
	}
}

// callbacks //////////////////////////////////////////////////////////////////////////////////////

// Callback function for mouse events. see cv::setMouseCallback
//   event one of the cv::MouseEventTypes constants.
//   x The x-coordinate of the mouse event.
//   y The y-coordinate of the mouse event.
//   flags one of the cv::MouseEventFlags constants.
//   userdata The optional parameter.
type MouseCallback func(event, x, y, flags int)

// Callback function for Trackbar see cv::createTrackbar
//   pos current position of the specified trackbar.
//   userdata The optional parameter.
type TrackbarCallback func(pos int)

// Callback function defined to be called every frame. See cv::setOpenGlDrawCallback
//   userdata The optional parameter.
type OpenGlDrawCallback func()

// Callback function for a button created by cv::createButton
//   state current state of the button. It could be -1 for a push button, 0 or 1 for a check/radio box button.
//   userdata The optional parameter.
type ButtonCallback func(state int)

// enums //////////////////////////////////////////////////////////////////////////////////////////

// constants
var (

	// from opencv2/highgui.cpp

	// WindowFlags
	WINDOW_NORMAL       int // the user can resize the window (no constraint) / also use to switch a fullscreen window to a normal size.
	WINDOW_AUTOSIZE     int // the user cannot resize the window, the size is constrained by the image displayed.
	WINDOW_OPENGL       int // window with opengl support.
	WINDOW_FULLSCREEN   int // change the window to fullscreen.
	WINDOW_FREERATIO    int // the image expends as much as it can (no ratio constraint).
	WINDOW_KEEPRATIO    int // the ratio of the image is respected.
	WINDOW_GUI_EXPANDED int // status bar and tool bar
	WINDOW_GUI_NORMAL   int // old fashion way

	// WindowPropertyFlags
	WND_PROP_FULLSCREEN   int // fullscreen property    (can be WINDOW_NORMAL or WINDOW_FULLSCREEN).
	WND_PROP_AUTOSIZE     int // autosize property      (can be WINDOW_NORMAL or WINDOW_AUTOSIZE).
	WND_PROP_ASPECT_RATIO int // window's aspect ration (can be set to WINDOW_FREERATIO or WINDOW_KEEPRATIO).
	WND_PROP_OPENGL       int // opengl support.
	WND_PROP_VISIBLE      int // checks whether the window exists and is visible

	// MouseEventTypes
	EVENT_MOUSEMOVE     int // indicates that the mouse pointer has moved over the window.
	EVENT_LBUTTONDOWN   int // indicates that the left mouse button is pressed.
	EVENT_RBUTTONDOWN   int // indicates that the right mouse button is pressed.
	EVENT_MBUTTONDOWN   int // indicates that the middle mouse button is pressed.
	EVENT_LBUTTONUP     int // indicates that left mouse button is released.
	EVENT_RBUTTONUP     int // indicates that right mouse button is released.
	EVENT_MBUTTONUP     int // indicates that middle mouse button is released.
	EVENT_LBUTTONDBLCLK int // indicates that left mouse button is double clicked.
	EVENT_RBUTTONDBLCLK int // indicates that right mouse button is double clicked.
	EVENT_MBUTTONDBLCLK int // indicates that middle mouse button is double clicked.
	EVENT_MOUSEWHEEL    int // positive and negative values mean forward and backward scrolling, respectively.
	EVENT_MOUSEHWHEEL   int // positive and negative values mean right and left scrolling, respectively.

	// MouseEventFlags
	EVENT_FLAG_LBUTTON  int // indicates that the left mouse button is down.
	EVENT_FLAG_RBUTTON  int // indicates that the right mouse button is down.
	EVENT_FLAG_MBUTTON  int // indicates that the middle mouse button is down.
	EVENT_FLAG_CTRLKEY  int // indicates that CTRL Key is pressed.
	EVENT_FLAG_SHIFTKEY int // indicates that SHIFT Key is pressed.
	EVENT_FLAG_ALTKEY   int // indicates that ALT Key is pressed.

	// QtFontWeights
	QT_FONT_LIGHT    int // Weight of 25
	QT_FONT_NORMAL   int // Weight of 50
	QT_FONT_DEMIBOLD int // Weight of 63
	QT_FONT_BOLD     int // Weight of 75
	QT_FONT_BLACK    int // Weight of 87

	// QtFontStyles
	QT_STYLE_NORMAL  int // Normal font.
	QT_STYLE_ITALIC  int // Italic font.
	QT_STYLE_OBLIQUE int // Oblique font.

	// QtButtonTypes
	QT_PUSH_BUTTON   int // Push button.
	QT_CHECKBOX      int // Checkbox button.
	QT_RADIOBOX      int // Radiobox button.
	QT_NEW_BUTTONBAR int // Button should create a new buttonbar
)

// initialise constants (enums)
func init() {
	C.cv_initialise_enums(
		// WindowFlags
		(*C.int)(unsafe.Pointer(&WINDOW_NORMAL)),
		(*C.int)(unsafe.Pointer(&WINDOW_AUTOSIZE)),
		(*C.int)(unsafe.Pointer(&WINDOW_OPENGL)),
		(*C.int)(unsafe.Pointer(&WINDOW_FULLSCREEN)),
		(*C.int)(unsafe.Pointer(&WINDOW_FREERATIO)),
		(*C.int)(unsafe.Pointer(&WINDOW_KEEPRATIO)),
		(*C.int)(unsafe.Pointer(&WINDOW_GUI_EXPANDED)),
		(*C.int)(unsafe.Pointer(&WINDOW_GUI_NORMAL)),

		// WindowPropertyFlags
		(*C.int)(unsafe.Pointer(&WND_PROP_FULLSCREEN)),
		(*C.int)(unsafe.Pointer(&WND_PROP_AUTOSIZE)),
		(*C.int)(unsafe.Pointer(&WND_PROP_ASPECT_RATIO)),
		(*C.int)(unsafe.Pointer(&WND_PROP_OPENGL)),
		(*C.int)(unsafe.Pointer(&WND_PROP_VISIBLE)),

		// MouseEventTypes
		(*C.int)(unsafe.Pointer(&EVENT_MOUSEMOVE)),
		(*C.int)(unsafe.Pointer(&EVENT_LBUTTONDOWN)),
		(*C.int)(unsafe.Pointer(&EVENT_RBUTTONDOWN)),
		(*C.int)(unsafe.Pointer(&EVENT_MBUTTONDOWN)),
		(*C.int)(unsafe.Pointer(&EVENT_LBUTTONUP)),
		(*C.int)(unsafe.Pointer(&EVENT_RBUTTONUP)),
		(*C.int)(unsafe.Pointer(&EVENT_MBUTTONUP)),
		(*C.int)(unsafe.Pointer(&EVENT_LBUTTONDBLCLK)),
		(*C.int)(unsafe.Pointer(&EVENT_RBUTTONDBLCLK)),
		(*C.int)(unsafe.Pointer(&EVENT_MBUTTONDBLCLK)),
		(*C.int)(unsafe.Pointer(&EVENT_MOUSEWHEEL)),
		(*C.int)(unsafe.Pointer(&EVENT_MOUSEHWHEEL)),

		// MouseEventFlags
		(*C.int)(unsafe.Pointer(&EVENT_FLAG_LBUTTON)),
		(*C.int)(unsafe.Pointer(&EVENT_FLAG_RBUTTON)),
		(*C.int)(unsafe.Pointer(&EVENT_FLAG_MBUTTON)),
		(*C.int)(unsafe.Pointer(&EVENT_FLAG_CTRLKEY)),
		(*C.int)(unsafe.Pointer(&EVENT_FLAG_SHIFTKEY)),
		(*C.int)(unsafe.Pointer(&EVENT_FLAG_ALTKEY)),

		// QtFontWeights
		(*C.int)(unsafe.Pointer(&QT_FONT_LIGHT)),
		(*C.int)(unsafe.Pointer(&QT_FONT_NORMAL)),
		(*C.int)(unsafe.Pointer(&QT_FONT_DEMIBOLD)),
		(*C.int)(unsafe.Pointer(&QT_FONT_BOLD)),
		(*C.int)(unsafe.Pointer(&QT_FONT_BLACK)),

		// QtFontStyles
		(*C.int)(unsafe.Pointer(&QT_STYLE_NORMAL)),
		(*C.int)(unsafe.Pointer(&QT_STYLE_ITALIC)),
		(*C.int)(unsafe.Pointer(&QT_STYLE_OBLIQUE)),

		// QtButtonTypes
		(*C.int)(unsafe.Pointer(&QT_PUSH_BUTTON)),
		(*C.int)(unsafe.Pointer(&QT_CHECKBOX)),
		(*C.int)(unsafe.Pointer(&QT_RADIOBOX)),
		(*C.int)(unsafe.Pointer(&QT_NEW_BUTTONBAR)),
	)
}
