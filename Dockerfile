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

# copy scripts and patches
COPY zscripts /tmp/zscripts

# download the source code of MUMPS and compile it
RUN bash /tmp/zscripts/mumps/install.bash

# configure basic system
ARG INSTALL_ZSH="true"
ARG USERNAME="vscode"
ARG USER_UID="1000"
ARG USER_GID=$USER_UID
ARG UPGRADE_PACKAGES="true"
RUN bash /tmp/zscripts/microsoft/common-debian.sh "${INSTALL_ZSH}" "${USERNAME}" "${USER_UID}" "${USER_GID}" "${UPGRADE_PACKAGES}" \
  && apt-get clean -y && rm -rf /var/lib/apt/lists/*

# install Go tools
ARG GO_VERSION="latest"
ARG GOROOT="/usr/local/go"
ARG GOPATH="/go"
ARG UPDATE_RC="true"
ARG INSTALL_GO_TOOLS="true"
ENV GO111MODULE=auto
RUN bash /tmp/zscripts/microsoft/go-debian.sh "${GO_VERSION}" "${GOROOT}" "${GOPATH}" "${USERNAME}" "${UPDATE_RC}" "${INSTALL_GO_TOOLS}" \
  && apt-get clean -y

# clean up
RUN rm -rf /tmp/zscripts
