FROM debian:sid

# disable tzdata questions
ENV DEBIAN_FRONTEND=noninteractive

# use bash
SHELL ["/bin/bash", "-c"]

# enable non-free
RUN sed -i "s#deb http://deb.debian.org/debian sid main#deb http://deb.debian.org/debian sid main non-free#g" /etc/apt/sources.list

# install apt-utils
RUN apt-get update -y && \
  apt-get install -y apt-utils 2> >( grep -v 'debconf: delaying package configuration, since apt-utils is not installed' >&2 ) \
  && apt-get clean && rm -rf /var/lib/apt/lists/*

# essential tools
RUN apt-get update -y && apt-get install -y --no-install-recommends \
  ca-certificates \
  netbase \
  curl \
  git \
  && apt-get clean && rm -rf /var/lib/apt/lists/*

# required compilers and libraries for gosl
RUN apt-get update -y && apt-get install -y --no-install-recommends \
  build-essential \
  gcc \
  gfortran \
  libopenmpi-dev \
  libhwloc-dev \
  liblapacke-dev \
  libopenblas-dev \
  libmetis-dev \
  libparmetis-dev \
  libscotch-dev \
  libptscotch-dev \
  libatlas-base-dev \
  libscalapack-mpi-dev \
  libsuitesparse-dev \
  libfftw3-dev \
  libfftw3-mpi-dev \
  intel-mkl-full \
  && apt-get clean && rm -rf /var/lib/apt/lists/*

# download the source code of MUMPS and compile it
ARG MUMPS_VERSION="5.3.5"
ARG MUMPS_GZ=mumps_${MUMPS_VERSION}.orig.tar.gz
ARG MUMPS_DIR=/tmp/MUMPS_${MUMPS_VERSION}
RUN curl http://deb.debian.org/debian/pool/main/m/mumps/${MUMPS_GZ} -o /tmp/${MUMPS_GZ}
RUN cd /tmp/ && tar xzf ${MUMPS_GZ}
COPY zscripts/patches/mumps/PORD.lib.Makefile.diff ${MUMPS_DIR}/PORD/lib/Makefile.diff
COPY zscripts/patches/mumps/src.Makefile.diff ${MUMPS_DIR}/src/Makefile.diff
COPY zscripts/patches/mumps/Makefile.diff ${MUMPS_DIR}/Makefile.diff
COPY zscripts/patches/mumps/Makefile.inc ${MUMPS_DIR}/Makefile.inc
RUN cd ${MUMPS_DIR} \
   && patch -u PORD/lib/Makefile PORD/lib/Makefile.diff \
   && patch -u src/Makefile src/Makefile.diff \
   && patch -u Makefile Makefile.diff \
   && make all \
   && cp include/*.h /usr/include/ \
   && cp -av lib/*.so /usr/lib/ \
   && cd .. && rm ${MUMPS_GZ} && rm -rf ${MUMPS_DIR} \
   && ldconfig 

############################################################################################################
# install other tools as in:
# https://github.com/microsoft/vscode-dev-containers/blob/master/containers/go/.devcontainer/base.Dockerfile

# install needed packages and setup non-root user. Use a separate RUN statement to add your own dependencies.
ARG INSTALL_ZSH="true"
ARG USERNAME="vscode"
ARG USER_UID="1000"
ARG USER_GID=$USER_UID
ARG UPGRADE_PACKAGES="true"
COPY zscripts/microsoft/common-debian.sh /tmp/library-scripts/
RUN bash /tmp/library-scripts/common-debian.sh "${INSTALL_ZSH}" "${USERNAME}" "${USER_UID}" "${USER_GID}" "${UPGRADE_PACKAGES}" \
  && apt-get clean -y && rm -rf /var/lib/apt/lists/* /tmp/library-scripts

# install Go tools
ARG GO_VERSION="latest"
ARG GOROOT="/usr/local/go"
ARG GOPATH="/go"
ARG UPDATE_RC="true"
ARG INSTALL_GO_TOOLS="true"
ENV GO111MODULE=auto
COPY zscripts/microsoft/go-debian.sh /tmp/library-scripts/
RUN bash /tmp/library-scripts/go-debian.sh "${GO_VERSION}" "${GOROOT}" "${GOPATH}" "${USERNAME}" "${UPDATE_RC}" "${INSTALL_GO_TOOLS}" \
  && apt-get clean -y && rm -rf /tmp/library-scripts
