#!/bin/bash

set -e

platform='unknown'
unamestr=`uname`
if [[ "$unamestr" == 'Linux' ]]; then
   platform='linux'
elif [[ "$unamestr" == 'MINGW32_NT-6.2' ]]; then
   platform='windows'
elif [[ "$unamestr" == 'MINGW64_NT-10.0' ]]; then
   platform='windows'
elif [[ "$unamestr" == 'Darwin' ]]; then
   platform='darwin'
fi

echo "   platform = $platform"

if [[ $platform != 'linux' ]]; then
    echo "govtk works on Linux only at the time, sorry"
    exit 1
fi

# this line will fail if the path does not exist (OK)
VTK_INCLUDE_PATH=`ls -d /usr/include/vtk*`

echo "   VTK_INCLUDE_PATH = $VTK_INCLUDE_PATH"

if [[ -z $VTK_INCLUDE_PATH ]]; then
    echo "cannot find /usr/include/vtk*"
    exit 1
fi

VTK_LIB=`basename $VTK_INCLUDE_PATH`

#echo "   VTK_LIB = $VTK_LIB"

VTK_VERSION="${VTK_LIB#*-}"

echo "   VTK_VERSION = $VTK_VERSION"

VTK_CXXFLAGS="-I/usr/include/vtk-$VTK_VERSION"
VTK_LDFLAGS="-lvtkCommonCore-$VTK_VERSION -lvtkFiltersHybrid-$VTK_VERSION \
-lvtkFiltersModeling-$VTK_VERSION -lvtkFiltersSources-$VTK_VERSION -lvtkFiltersGeneral-$VTK_VERSION \
-lvtkCommonMisc-$VTK_VERSION -lvtkCommonSystem-$VTK_VERSION \
-lvtkCommonComputationalGeometry-$VTK_VERSION -lvtkCommonMath-$VTK_VERSION \
-lvtkCommonTransforms-$VTK_VERSION -lvtkIOCore-$VTK_VERSION -lvtkIOLegacy-$VTK_VERSION \
-lvtkIOGeometry-$VTK_VERSION -lvtkIOExport-$VTK_VERSION -lvtkIOImage-$VTK_VERSION \
-lvtkInteractionStyle-$VTK_VERSION -lvtkRenderingAnnotation-$VTK_VERSION \
-lvtkRenderingFreeType-$VTK_VERSION \
-lvtkRenderingOpenGL-$VTK_VERSION -lvtkRenderingCore-$VTK_VERSION -lvtkFiltersCore-$VTK_VERSION \
-lvtkCommonExecutionModel-$VTK_VERSION -lvtkCommonDataModel-$VTK_VERSION \
-lvtkCommonTransforms-$VTK_VERSION -lvtkRenderingGL2PS-$VTK_VERSION -lvtkRenderingImage-$VTK_VERSION \
-lvtkalglib-$VTK_VERSION -lvtksys-$VTK_VERSION -lvtkRenderingLOD-$VTK_VERSION \
-lvtkViewsCore-$VTK_VERSION -lvtkInteractionStyle-$VTK_VERSION"

if [[ $VTK_VERSION == "6.2" ]]; then
    VTK_LDFLAGS="$VTK_LDFLAGS -lvtkRenderingFreeTypeOpenGL-$VTK_VERSION"
fi

#echo "VTK_CXXFLAGS = $VTK_CXXFLAGS"
#echo "VTK_LDFLAGS = $VTK_LDFLAGS"

sed -i -e "s@CHANGE_HERE_WITH_VTK_CXXFLAGS@$VTK_CXXFLAGS@g" \
       -e "s@CHANGE_HERE_WITH_VTK_LDFLAGS@$VTK_LDFLAGS@g" \
          govtk.go
