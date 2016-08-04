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
    openjdk-7-jre-headless \
    xz-utils \
    bzip2 \
    libfreetype6 \
    libfontconfig


#
# Install Golang
#
RUN curl 'https://storage.googleapis.com/golang/go1.7rc5.linux-amd64.tar.gz' | tar -C /usr/local -xzf -

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH


#
# Install SASS
#
RUN gem install sass


#
# Install Node
#
ENV NODE_VERSION 6.3.1

RUN curl -SLO "https://nodejs.org/dist/v$NODE_VERSION/node-v$NODE_VERSION-linux-x64.tar.xz" \
  && tar -xJf "node-v$NODE_VERSION-linux-x64.tar.xz" -C /usr/local --strip-components=1 \
  && rm "node-v$NODE_VERSION-linux-x64.tar.xz"

RUN npm install -g eslint@3.2.2 jshint@2.9.2

#
# Install PhantomJS
#
ENV PHANTOM_JS_VERSION 1.9.8-linux-x86_64

RUN curl -sSL https://bitbucket.org/ariya/phantomjs/downloads/phantomjs-$PHANTOM_JS_VERSION.tar.bz2 | tar xjC /usr/local/ --strip-components=1



#
# Sauce Connect
# https://wiki.saucelabs.com/display/DOCS/Setting+Up+Sauce+Connect
#
RUN curl "https://saucelabs.com/downloads/sc-4.3.16-linux.tar.gz" | tar zxC /usr/local/ --strip-components=1



#
# Install dependencies
#

# Go
ADD Makefile /Makefile
RUN go get github.com/tools/godep
RUN make golint_deps

# UI tests
ADD tests/package.json /usr/local/lib/package.json
RUN cd /usr/local/lib/ && SAUCE_CONNECT_DOWNLOAD_ON_INSTALL=true npm install --no-optional