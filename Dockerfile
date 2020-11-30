FROM ubuntu:20.04

# options
ARG GO_VERSION="1.15.5"

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
RUN curl https://dl.google.com/go/$GOFN -o /usr/local/$GOFN \
  && tar -xzf /usr/local/$GOFN -C /usr/local/ \
  && rm /usr/local/$GOFN

# set path for go executable
ENV PATH $PATH:/usr/local/go/bin
RUN go version

############################################################################################################
# install other tools as in:
# https://github.com/microsoft/vscode-dev-containers/blob/master/containers/go/.devcontainer/base.Dockerfile

# intall zsh and upgrade pkgs
ARG INSTALL_ZSH="true"
ARG UPGRADE_PACKAGES="true"

# install needed packages and setup non-root user. Use a separate RUN statement to add your own dependencies.
ARG USERNAME=vscode
ARG USER_UID=1000
ARG USER_GID=$USER_UID
COPY zscripts/microsoft/common-debian.sh /tmp/library-scripts/
RUN bash /tmp/library-scripts/common-debian.sh "${INSTALL_ZSH}" "${USERNAME}" "${USER_UID}" "${USER_GID}" "${UPGRADE_PACKAGES}" \
  && apt-get clean -y && rm -rf /var/lib/apt/lists/* /tmp/library-scripts

# install Go tools
ENV GO111MODULE=auto
COPY zscripts/microsoft/go-debian.sh /tmp/library-scripts/
RUN bash /tmp/library-scripts/go-debian.sh "none" "/usr/local/go" "${GOPATH}" "${USERNAME}" "false" \
  && apt-get clean -y && rm -rf /tmp/library-scripts
