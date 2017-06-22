#!/bin/bash

GP=${GOPATH//\\//}
D="$GP/src/github.com/cpmech/gosl/scripts/windows"

cp $D/fix-suitesparse/SuiteSparse_config.mk                 SuiteSparse/SuiteSparse_config/SuiteSparse_config.mk
cp $D/fix-suitesparse/SuiteSparse_config_Makefile.txt       SuiteSparse/SuiteSparse_config/Makefile
cp $D/fix-suitesparse/SuiteSparse_AMD_Lib_Makefile.txt      SuiteSparse/AMD/Lib/Makefile
cp $D/fix-suitesparse/SuiteSparse_BTF_Lib_Makefile.txt      SuiteSparse/BTF/Lib/Makefile
cp $D/fix-suitesparse/SuiteSparse_CAMD_Lib_Makefile.txt     SuiteSparse/CAMD/Lib/Makefile
cp $D/fix-suitesparse/SuiteSparse_CCOLAMD_Lib_Makefile.txt  SuiteSparse/CCOLAMD/Lib/Makefile
cp $D/fix-suitesparse/SuiteSparse_COLAMD_Lib_Makefile.txt   SuiteSparse/COLAMD/Lib/Makefile
cp $D/fix-suitesparse/SuiteSparse_CHOLMOD_Lib_Makefile.txt  SuiteSparse/CHOLMOD/Lib/Makefile
cp $D/fix-suitesparse/SuiteSparse_CXSparse_Lib_Makefile.txt SuiteSparse/CXSparse/Lib/Makefile
cp $D/fix-suitesparse/SuiteSparse_KLU_Lib_Makefile.txt      SuiteSparse/KLU/Lib/Makefile
cp $D/fix-suitesparse/SuiteSparse_LDL_Lib_Makefile.txt      SuiteSparse/LDL/Lib/Makefile
cp $D/fix-suitesparse/SuiteSparse_UMFPACK_Makefile.txt      SuiteSparse/UMFPACK/Makefile
cp $D/fix-suitesparse/SuiteSparse_UMFPACK_Lib_Makefile.txt  SuiteSparse/UMFPACK/Lib/Makefile
cp $D/fix-suitesparse/SuiteSparse_Makefile.txt              SuiteSparse/Makefile