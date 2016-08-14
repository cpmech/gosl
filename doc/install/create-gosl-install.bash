#!/bin/bash

# generate GUID with
# C:\Program Files (x86)\Microsoft Visual Studio 14.0\Common7\Tools\guidgen.exe

# add these two to the PATH
# C:\Program Files (x86)\WiX Toolset v3.10\bin

heat.exe dir SourceDir -gg -sfrag -sreg -srd -dr GOSLDIR -cg GoslComp -out goslfrags.wxs
candle.exe goslfrags.wxs goslapp.wxs varsdlg.wxs myui.wxs
light.exe -ext WixUIExtension -ext WixUtilExtension -out gosl-installer.msi goslfrags.wixobj goslapp.wixobj varsdlg.wixobj myui.wixobj
rm *.wixobj *.wixpdb