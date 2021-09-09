FROM ubuntu:20.04

# disable tzdata questions
ENV DEBIAN_FRONTEND=noninteractive

# use bash
SHELL ["/bin/bash", "-c"]

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
  gcc \
  gfortran \
  libfftw3-dev \
  liblapacke-dev \
  libmetis-dev \
  libmumps-seq-dev \
  libopenblas-dev \
  libsuitesparse-dev \
  && apt-get clean && rm -rf /var/lib/apt/lists/*

# configure basic system
ARG INSTALL_ZSH="true"
ARG USERNAME="vscode"
ARG USER_UID="1000"
ARG USER_GID=$USER_UID
ARG UPGRADE_PACKAGES="true"
COPY zscripts/microsoft/common-debian.sh /tmp/
RUN bash /tmp/common-debian.sh "${INSTALL_ZSH}" "${USERNAME}" "${USER_UID}" "${USER_GID}" "${UPGRADE_PACKAGES}" \
  && apt-get clean -y && rm -rf /var/lib/apt/lists/* && rm /tmp/common-debian.sh

# install Go tools
ARG GO_VERSION="latest"
ARG GOROOT="/usr/local/go"
ARG GOPATH="/go"
ARG UPDATE_RC="true"
ARG INSTALL_GO_TOOLS="true"
ENV GO111MODULE=auto
COPY zscripts/microsoft/go-debian.sh /tmp/
RUN bash /tmp/go-debian.sh "${GO_VERSION}" "${GOROOT}" "${GOPATH}" "${USERNAME}" "${UPDATE_RC}" "${INSTALL_GO_TOOLS}" \
  && apt-get clean -y && rm /tmp/go-debian.sh
