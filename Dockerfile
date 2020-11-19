FROM ubuntu:20.04

# options
ARG DEV_IMG="false"
ARG GOSL_VERSION="2.0.0"
ARG GO_VERSION="1.15.2"
ARG MUMPS_VERSION="5.3.5"

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
  build-essential \
  libmetis-dev libparmetis-dev libscotch-dev libptscotch-dev libatlas-base-dev \
  openmpi-bin libopenmpi-dev libscalapack-openmpi-dev \
  libhwloc-dev \
  liblapacke-dev \
  libopenblas-dev \
  libsuitesparse-dev \
  libfftw3-dev \
  libfftw3-mpi-dev \
  && apt-get clean && rm -rf /var/lib/apt/lists/*

# download the source code of MUMPS and compile it
ARG MUMPS_FN=mumps_$MUMPS_VERSION.orig.tar.gz
RUN curl http://deb.debian.org/debian/pool/main/m/mumps/$MUMPS_FN -o /tmp/$MUMPS_FN
RUN cd /tmp/ \
    && tar xzvf $MUMPS_FN \
    && cd MUMPS_$MUMPS_VERSION \
    && cp Make.inc/Makefile.debian.PAR ./Makefile.inc \
    && sed -i 's/-lblacs-openmpi//g' Makefile.inc \
    && make all \
    && cp lib/*.a /usr/lib
 
# download go
ARG GOFN=go$GO_VERSION.linux-amd64.tar.gz
RUN curl https://dl.google.com/go/$GOFN -o /usr/local/$GOFN
RUN tar -xzf /usr/local/$GOFN -C /usr/local/
RUN rm /usr/local/$GOFN

# set path for go executable
ENV PATH $PATH:/usr/local/go/bin
RUN go version

# build gosl
COPY zscripts/gosl-clone-and-build.bash /tmp/library-scripts/
RUN /bin/bash /tmp/library-scripts/gosl-clone-and-build.bash "${DEV_IMG}" "${GOSL_VERSION}"

##################################################################################################
#                                                                                                #
#   The code below is copied from:                                                               #
#      https://github.com/microsoft/vscode-remote-try-go/blob/master/.devcontainer/Dockerfile    #
#                                                                                                #
#   NOTE: remember to fix zscripts/common-debian.sh                                              #
#                                                                                                #
##################################################################################################

# Options for setup script
ARG INSTALL_ZSH="true"
ARG UPGRADE_PACKAGES="false"
ARG USERNAME=vscode
ARG USER_UID=1000
ARG USER_GID=$USER_UID

# Install needed packages and setup non-root user. Use a separate RUN statement to add your own dependencies.
COPY zscripts/common-debian.sh /tmp/library-scripts/
RUN apt-get update \
  && /bin/bash /tmp/library-scripts/common-debian.sh "${INSTALL_ZSH}" "${USERNAME}" "${USER_UID}" "${USER_GID}" "${UPGRADE_PACKAGES}" \
  && apt-get autoremove -y && apt-get clean -y && rm -rf /var/lib/apt/lists/* /tmp/library-scripts

# Install Go tools
ARG GO_TOOLS_WITH_MODULES="\
  golang.org/x/tools/gopls \
  honnef.co/go/tools/... \
  golang.org/x/tools/cmd/gorename \
  golang.org/x/tools/cmd/goimports \
  golang.org/x/tools/cmd/guru \
  golang.org/x/lint/golint \
  github.com/mdempsky/gocode \
  github.com/cweill/gotests/... \
  github.com/haya14busa/goplay/cmd/goplay \
  github.com/sqs/goreturns \
  github.com/josharian/impl \
  github.com/davidrjenni/reftools/cmd/fillstruct \
  github.com/uudashr/gopkgs/v2/cmd/gopkgs \
  github.com/ramya-rao-a/go-outline \
  github.com/acroca/go-symbols \
  github.com/godoctor/godoctor \
  github.com/rogpeppe/godef \
  github.com/zmb3/gogetdoc \
  github.com/fatih/gomodifytags \
  github.com/mgechev/revive \
  github.com/go-delve/delve/cmd/dlv"
RUN mkdir -p /tmp/gotools \
  && cd /tmp/gotools \
  && export GOPATH=/tmp/gotools \
  # Go tools w/module support
  && export GO111MODULE=on \
  && (echo "${GO_TOOLS_WITH_MODULES}" | xargs -n 1 go get -x )2>&1 \
  # gocode-gomod
  && export GO111MODULE=auto \
  && go get -x -d github.com/stamblerre/gocode 2>&1 \
  && go build -o gocode-gomod github.com/stamblerre/gocode \
  # golangci-lint
  && curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b /usr/local/bin 2>&1 \
  # Move Go tools into path and clean up
  && mv /tmp/gotools/bin/* /usr/local/bin/ \
  && mv gocode-gomod /usr/local/bin/ \
  && rm -rf /tmp/gotools

ENV GO111MODULE=auto

##################################################################################################
##################################################################################################
  