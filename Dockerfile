FROM ubuntu:20.04

# set go and gosl version
ARG GO_VERSION=1.15
ARG GOSL_VERSION=1.2.0
ARG GOSL_BRANCH=trim-gosl

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
    libopenmpi-dev \
    libhwloc-dev \
    liblapacke-dev \
    libopenblas-dev \
    libmetis-dev \
    libsuitesparse-dev \
    libmumps-dev \
    libfftw3-dev \
    libfftw3-mpi-dev \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

# download go
ARG GOFN=go$GO_VERSION.linux-amd64.tar.gz
RUN curl https://dl.google.com/go/$GOFN -o /usr/local/$GOFN
RUN tar -xzf /usr/local/$GOFN -C /usr/local/
RUN rm /usr/local/$GOFN

# set path for go apps
ENV GOPATH /mygo
ENV PATH $PATH:/mygo/bin:/usr/local/go/bin
RUN go version

# build gosl
WORKDIR /mygo/src
RUN git clone -b $GOSL_BRANCH --single-branch --depth 1 https://github.com/cpmech/gosl.git
WORKDIR /mygo/src/gosl
RUN bash ./all.bash
