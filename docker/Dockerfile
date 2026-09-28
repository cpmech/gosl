# syntax=docker/dockerfile:1.7
#
# Gosl's build-and-test environment: Ubuntu 24.04 + a pinned Go release.
# There is no Go toolchain on the host; everything happens in here.
#
#   docker build -f docker/Dockerfile -t gosl:ci .
#   ./gosl all
#
# Change the pinned versions only by editing them here.
ARG BASE=ubuntu:24.04
FROM ${BASE}

ENV DEBIAN_FRONTEND=noninteractive
SHELL ["/bin/bash", "-c"]

# The C/Fortran libraries Gosl links against -- see la/flags.go, fun/fftw/flags.go, gm/metis.go.
RUN --mount=type=cache,target=/var/cache/apt,sharing=locked \
    apt-get update \
 && apt-get install -y --no-install-recommends \
      ca-certificates curl git make bash \
      gcc gfortran \
      libopenblas-dev liblapacke-dev libfftw3-dev \
      libsuitesparse-dev libmetis-dev libmumps-seq-dev \
 && rm -rf /var/lib/apt/lists/*

# Record what this image ships: the first question asked about any numerical bug report.
RUN dpkg-query -W \
      libopenblas-dev liblapacke-dev libfftw3-dev libsuitesparse-dev \
      libmetis-dev libmumps-seq-dev libgfortran5 \
    | sort -o /etc/gosl-libs.txt

ARG GO_VERSION=1.27.1
RUN curl -fsSL "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz" -o /tmp/go.tgz \
 && tar -C /usr/local -xzf /tmp/go.tgz \
 && rm /tmp/go.tgz
ENV PATH=/usr/local/go/bin:/go/bin:${PATH} \
    GOTOOLCHAIN=local

WORKDIR /gosl
CMD ["bash", "-lc", "cat /etc/gosl-libs.txt && go version"]
