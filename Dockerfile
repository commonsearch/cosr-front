FROM debian:jessie

ENV LANG C.UTF-8

RUN apt-get update && apt-get upgrade -y && apt-get install -y --no-install-recommends \
    ca-certificates \
    curl \
    make \
    gcc \
    vim \
    libc6-dev \
    git \
    rubygems \
    openjdk-7-jre-headless


#
# Install Golang
#
RUN curl 'https://storage.googleapis.com/golang/go1.6.linux-amd64.tar.gz' | tar -C /usr/local -xzf -

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

ADD Makefile /Makefile

RUN go get github.com/tools/godep
RUN make golint_deps