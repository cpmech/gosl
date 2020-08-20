FROM ubuntu:20.04

# set go version
ARG GOVER=1.15

# disable tzdata questions
ENV DEBIAN_FRONTEND=noninteractive

# use bash
SHELL ["/bin/bash", "-c"]

# install apt-utils
RUN apt-get update -y && \
    apt-get install -y apt-utils 2> >( grep -v 'debconf: delaying package configuration, since apt-utils is not installed' >&2 ) \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

# essential net tools
RUN apt-get update -y && apt-get install -y --no-install-recommends \
    ca-certificates \
    netbase \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

# download tools and compilers
RUN apt-get update -y && apt-get install -y --no-install-recommends \
    curl \
    git \
    gcc \
    gfortran \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

# required libraries for gosl
RUN apt-get update -y && apt-get install -y --no-install-recommends \
    libopenmpi-dev \
    libhwloc-dev \
    liblapacke-dev \
    libopenblas-dev \
    libmetis-dev \
    libsuitesparse-dev \
    libmumps-dev \
    libfftw3-dev \
    libfftw3-mpi-dev \
    && rm -rf /var/lib/apt/lists/*

# download go
ARG GOFN=go$GOVER.linux-amd64.tar.gz
RUN curl https://dl.google.com/go/$GOFN -o /usr/local/$GOFN
RUN tar -xzf /usr/local/$GOFN -C /usr/local/
RUN rm /usr/local/$GOFN

# set path for go apps
ENV GOPATH /mygo
ENV PATH $PATH:/mygo/bin:/usr/local/go/bin
RUN go version

# build gosl
WORKDIR /mygo
#RUN bash all.bash
