#!/bin/bash

# find Gosl scripts path
GP=${GOPATH//\\//}
SS="$GP/src/github.com/cpmech/gosl/scripts/windows/gen-installer"

# create SourceDir in TEMP directory
rm -rf $TEMP/GOSLWIX
mkdir -p $TEMP/GOSLWIX/SourceDir
cmd //c 'mklink /J %TEMP%\GOSLWIX\SourceDir\MyGo C:\MyGo'
cmd //c 'mklink /J %TEMP%\GOSLWIX\SourceDir\Gcc64 C:\TDM-GCC-64'

# generate frags file
cd $TEMP/GOSLWIX
heat.exe dir SourceDir -gg -sfrag -sreg -srd -dr TARGETDIR -t $SS/trn.xslt -out frags.wxs

# compile wxs files
candle.exe frags.wxs $SS/app.wxs $SS/varsdlg.wxs $SS/ui.wxs

# generate installer
light.exe -ext WixUIExtension -ext WixUtilExtension -out gosl-installer.msi frags.wixobj app.wixobj varsdlg.wixobj ui.wixobj

# remove object files
rm *.wixobj *.wixpdb