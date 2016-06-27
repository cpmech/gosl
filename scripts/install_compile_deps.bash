#!/bin/bash

set -e

echo
echo "*************************************************************************"
echo "*                                                                       *"
echo "* You can call this script with options to force recompiling all        *"
echo "* dependencies and/or download their code again:                        *"
echo "*                                                                       *"
echo "*                                             recompile     download    *"
echo "* Example:                                            |     |           *"
echo "*                                                     V     V           *"
echo "*   bash ~/mechsys/scripts/install_compile_deps.sh {0,1} {0,1} {NPROCS} *"
echo "*                                                                       *"
echo "* By default, a dependency will not be re-compiled if the corresponding *"
echo "* directory exists under the PKG_HOME/pkg directory.                    *"
echo "*                                                                       *"
echo "* By default, -j4 (4 processors) is given to make. This can be changed  *"
echo "* with a new value of NPROCS.                                           *"
echo "*                                                                       *"
echo "*************************************************************************"

RECOMPILE=0
if [ "$#" -gt 0 ]; then
    RECOMPILE=$1
    if [ "$RECOMPILE" -lt 0 -o "$RECOMPILE" -gt 1 ]; then
        echo
        echo "The option for re-compilation must be either 0 or 1"
        echo "  $1 is invalid"
        echo
        exit 1
    fi
fi

FORCEDOWNLOAD=0
if [ "$#" -gt 1 ]; then
    FORCEDOWNLOAD=$2
    if [ "$FORCEDOWNLOAD" -lt 0 -o "$FORCEDOWNLOAD" -gt 1 ]; then
        echo
        echo "The option for downloading and compilation of additional packages must be either 0 or 1"
        echo "  $2 is invalid"
        echo
        exit 1
    fi
fi

MYDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

if [ ! -n "$PKG_HOME" ]; then
  PKG_HOME=$HOME  
fi

NPROCS=4
if [ "$#" -gt 2 ]; then
    NPROCS=$3
fi

test -d $PKG_HOME/pkg || mkdir $PKG_HOME/pkg

error_message() {
    echo
    echo
    echo "    [1;31m Error: $1 [0m"
    echo
    echo
}

METIS_VER=5.1.0

metis_patch_cmd() {
    patch CMakeLists.txt $MYDIR/metis-"$METIS_VER"_CMakeLists.txt.diff
    patch include/metis.h $MYDIR/metis-"$METIS_VER"_metis.h.diff
}

download_and_compile() {
    EXT=tar.gz
    IS_SVN=0
    CONF_PRMS=""
    DO_PATCH=0
    DO_CONF=0
    DO_CMAKECONF=0
    DO_MAKE=1
    DO_MAKE_INST=0
    case "$1" in
        metis)
            PKG=metis-$METIS_VER
            LOCATION=http://glaros.dtc.umn.edu/gkhome/fetch/sw/metis/metis-$METIS_VER.$EXT
            DO_CMAKECONF=1
            DO_PATCH=1
            PATCH_CMD=metis_patch_cmd
            DO_MAKE_INST=1
            ;;
        *)
            error_message "download_and_compile: __Internal_error__"
            exit 1
            ;;
    esac
    echo
    echo "********************************** ${1} ********************************"

    # change into the packages directory
    cd $PKG_HOME/pkg

    # package filename and directory
    PKG_FILENAME=$PKG.$EXT
    if [ -z "$PKG_DIR" ]; then PKG_DIR=$PKG; fi

    # (re)compile or return (erasing existing package) ?
    if [ "$IS_SVN" -eq 0 ]; then
        if [ -d "$PKG_HOME/pkg/$PKG_DIR" ]; then
            if [ "$RECOMPILE" -eq 1   -o   "$FORCEDOWNLOAD" -eq 1 ]; then
                echo "    Erasing existing $PKG_DIR"
                rm -rf $PKG_HOME/pkg/$PKG_DIR
            else
                echo "    Using existing $PKG_DIR"
                return
            fi
        fi
    else
        if [ -d "$PKG_HOME/pkg/$PKG_DIR" ]; then
            if [ "$RECOMPILE" -eq 1   -o   "$FORCEDOWNLOAD" -eq 1 ]; then
                echo "    Updating existing $PKG SVN repository"
                cd $PKG_DIR
                svn up
                cd $PKG_HOME/pkg
            else
                echo "    Using existing $PKG SVN repository in $PKG_DIR"
                return
            fi
        fi
    fi

    # download package
    if [ "$IS_SVN" -eq 0 ]; then
        if [ "$FORCEDOWNLOAD" -eq 1   -o   ! -e "$PKG_FILENAME" ]; then
            if [ -e "$PKG_FILENAME" ]; then
                echo "    Removing existing <$PKG_FILENAME>"
                rm $PKG_FILENAME
            fi
            echo "    Downloading <$PKG_FILENAME>"
            wget $LOCATION
        fi
    else
        if [ ! -d "$PKG_HOME/pkg/$PKG_DIR" ]; then
            echo "    Downloading new $PKG SVN repository"
            svn co $LOCATION $PKG
        fi
    fi

    # uncompress package
    if [ "$IS_SVN" -eq 0 ]; then
        echo "        . . . uncompressing . . ."
        if [ "$EXT" = "tar.bz2" ]; then
            tar xjf $PKG_FILENAME
        else
            tar xzf $PKG_FILENAME
        fi
    fi

    # change into the package directory
    cd $PKG_DIR

    # patch
    if [ "$DO_PATCH" -eq 1 ]; then
        echo "        . . . patching . . ."
        $PATCH_CMD
    fi

    # configure
    if [ "$DO_CONF" -eq 1 ]; then
        echo "        . . . configuring . . ."
        ./configure $CONF_PRMS 2> /dev/null
    fi

    # cmake configuration
    if [ "$DO_CMAKECONF" -eq 1 ]; then
        echo "        . . . configuring using cmake . . ."
        cmake . 2> /dev/null
    fi

    # compilation
    if [ "$DO_MAKE" -eq 1 ]; then
        echo "        . . . compiling . . ."
        make -j$NPROCS > /dev/null 2> /dev/null
    fi

    # make install
    if [ "$DO_MAKE_INST" -eq 1 ]; then
        echo "        . . . installing . . ."
        sudo make install > /dev/null 2> /dev/null
        echo "           . . ldconfig . ."
        sudo ldconfig
    fi

    # finished
    echo "        . . . finished . . . . . "
}

download_and_compile metis

echo
echo "Finished ###################################################################"
echo
