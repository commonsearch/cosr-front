# cosr-front

[![Build Status](https://travis-ci.org/commonsearch/cosr-front.svg?branch=master)](https://travis-ci.org/commonsearch/cosr-front) [![Apache License 2.0](https://img.shields.io/github/license/commonsearch/cosr-front.svg)](LICENSE)

This repository contains the frontend for [Common Search](https://about.commonsearch.org/). A demo is currently hosted on [uidemo.commonsearch.org](https://uidemo.commonsearch.org/)

Help is welcome! We have a complete guide on [how to contribute](CONTRIBUTING.md).



## Understand the project

The frontend has 2 main components:

 - A [Go server](https://github.com/commonsearch/cosr-front/tree/master/server) that receives user queries (as HTTP GETs for page loads or AJAX calls), sends them to an Elasticsearch index, and then returns results as HTML or JSON.
 - An optional [JavaScript/CSS layer](https://github.com/commonsearch/cosr-front/tree/master/static) that provides a fast, single-page search experience to the otherwise static result pages.

Here is how they fit in our [general architecture](https://about.commonsearch.org/developer/architecture):

![General technical architecture of Common Search](https://about.commonsearch.org/images/developer/architecture-2016-02.svg)



## Local install

A complete guide available in [INSTALL.md](INSTALL.md).



## Launching the tests

You can run our full test suite easily:

```
make docker_test
```

Check out the [Makefile](https://github.com/commonsearch/cosr-front/blob/master/Makefile) for additional test, lint & build commands!



## How to contribute

Everything you need to know is in [CONTRIBUTING.md](CONTRIBUTING.md). We also have a tutorial on [how to send your first Frontend patch](https://about.commonsearch.org/developer/tutorials/first-frontend-patch).

Thanks for joining the adventure!
