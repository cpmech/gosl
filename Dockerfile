FROM ubuntu:18.04

ARG GOVER=1.12
ARG GOFN=go$GOVER.linux-amd64

WORKDIR /gosl

COPY . /gosl

RUN apt-get update

RUN DEBIAN_FRONTEND=noninteractive apt-get -y install tzdata
RUN ln -fs /usr/share/zoneinfo/Australia/Brisbane /etc/localtime
RUN dpkg-reconfigure --frontend noninteractive tzdata

RUN apt-get -y install wget git gcc \
    libopenmpi-dev libhwloc-dev libsuitesparse-dev libmumps-dev \
    gfortran python-scipy python-matplotlib dvipng \
    libfftw3-dev libfftw3-mpi-dev libmetis-dev \
    liblapacke-dev libopenblas-dev libhdf5-dev \
    libvtk7-dev

RUN wget https://dl.google.com/go/$GOFN.tar.gz -O /opt/$GOFN.tar.gz
RUN tar xf /opt/$GOFN.tar.gz -C /opt/
ENV PATH "$PATH:/opt/go/bin"
RUN go version

RUN USE_VTK=1 bash all.bash