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

Running `cosr-front` on your local machine is very simple. You only need to have [Go 1.5+](https://golang.org/) installed.

Once Go is installed, just run:

```
make devserver
```

Then open http://localhost:9700/ in your browser.

## Tests

We have strong linting and some Go unit tests available, with many more to come.

To run everything:

```
make test
```

Check out the [Makefile](https://github.com/commonsearch/cosr-front) for additional test, lint & build commands!