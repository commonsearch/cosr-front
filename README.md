# cosr-front

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

## Alternate local install without Docker

If for some reason you don't want to use Docker, you might be able to use a local Go install to run `cosr-front`. Please note that this is an unsupported method and might break at any time.

After [installing Go](https://golang.org/doc/install), you should be able to do:

```
make devserver
```

Then open the service running on [http://127.0.0.1:9700](http://127.0.0.1:9700).
