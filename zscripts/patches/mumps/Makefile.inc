#
#  This file is part of MUMPS 5.3.5, released
#  on Thu Oct 22 09:29:08 UTC 2020
#
# These settings for a PC under Debian/linux with standard packages :
# metis (parmetis), scotch (ptscotch), openmpi, gfortran

# packages installation: 
# apt-get install libmetis-dev libparmetis-dev libscotch-dev libptscotch-dev libatlas-base-dev openmpi-bin libopenmpi-dev liblapack-dev

# Begin orderings
LSCOTCHDIR = /usr/lib
ISCOTCH    = -I/usr/include/scotch
LSCOTCH    = -L$(LSCOTCHDIR) -lptesmumps -lptscotch -lptscotcherr

LPORDDIR = $(topdir)/PORD/lib/
IPORD    = -I$(topdir)/PORD/include/
LPORD    = -L$(LPORDDIR) -lpord

LMETISDIR = /usr/lib 
IMETIS    = -I/usr/include/parmetis
LMETIS    = -L$(LMETISDIR) -lparmetis -lmetis

# Corresponding variables reused later
ORDERINGSF = -Dmetis -Dpord -Dparmetis -Dscotch -Dptscotch
ORDERINGSC  = $(ORDERINGSF)

LORDERINGS  = $(LMETIS) $(LPORD) $(LSCOTCH)
IORDERINGSF = $(ISCOTCH)
IORDERINGSC = $(IMETIS) $(IPORD) $(ISCOTCH)
# End orderings
################################################################################

PLAT    =
LIBEXT  = .so
OUTC    = -o 
OUTF    = -o 
RM = /bin/rm -f
CC = mpicc
FC = mpif90 -fallow-argument-mismatch -w
FL = mpif90 -fallow-argument-mismatch -w
AR = ar vr 
RANLIB = ranlib
LAPACK = -llapack
SCALAP = -lscalapack-openmpi

INCPAR = # not needed with mpif90/mpicc:  -I/usr/include/openmpi

LIBPAR = $(SCALAP) $(LAPACK) # not needed with mpif90/mpicc: -lmpi_mpifh -lmpi

INCSEQ = -I$(topdir)/libseq
LIBSEQ  = $(LAPACK) -L$(topdir)/libseq -lmpiseq

LIBBLAS = -lblas
LIBOTHERS = -lpthread

#Preprocessor defs for calling Fortran from C (-DAdd_ or -DAdd__ or -DUPPER)
CDEFS   = -DAdd_

#Begin Optimized options
OPTF    = -O -fopenmp
OPTL    = -O -fopenmp
OPTC    = -O -fopenmp
#End Optimized options
 
INCS = $(INCPAR)
LIBS = $(LIBPAR)
LIBSEQNEEDED =
