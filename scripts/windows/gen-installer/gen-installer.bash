#!/bin/bash

# find Gosl scripts path
GP=${GOPATH//\\//}
SS="$GP/src/github.com/cpmech/gosl/scripts/windows/gen-installer"

# variables
TMPDIR=$TEMP/GOSLWIX
SRCDIR=$TMPDIR/SourceDir
CPMPTH=MyGo/src/github.com/cpmech
echo "TMPDIR=$TMPDIR"
echo "SRCDIR=$SRCDIR"
echo "CPMPTH=$CPMPTH"

# create SourceDir in TEMP directory
rm -rf $TMPDIR
mkdir -p $SRCDIR/$CPMPTH
cmd //c "mklink /J %TEMP%\GOSLWIX\SourceDir\Gcc64 C:\TDM-GCC-64"
cmd //c "mklink /J %TEMP%\GOSLWIX\SourceDir\MyGo\src\github.com\cpmech\gosl C:\MyGo\src\github.com\cpmech\gosl"

# generate frags file
cd $TMPDIR
heat.exe dir SourceDir -gg -sfrag -sreg -srd -dr TARGETDIR -t $SS/trn.xslt -out frags.wxs

# compile wxs files
candle.exe frags.wxs $SS/app.wxs $SS/varsdlg.wxs $SS/ui.wxs

# generate installer
VER=1.0.0
FNK=gosl-installer
FN=$FNK-v$VER.msk
light.exe -ext WixUIExtension -ext WixUtilExtension -out $FN frags.wixobj app.wixobj varsdlg.wixobj ui.wixobj

# remove object files
rm *.wixobj *.wixpdb

# message
echo "file <$TMPDIR/$FN> generated"