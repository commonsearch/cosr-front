# cosr-front

This repository contains the frontend for [Common Search](https://about.commonsearch.org/). A demo is currently hosted on [https://uidemo.commonsearch.org/](https://uidemo.commonsearch.org/)

The frontend has 2 main components:
 - A Go server that receives user requests (as HTTP GETs for page loads or AJAX calls), sends them to an ElasticSearch index, and then returns formatted results in HTML or JSON format.
 - An optional JavaScript/CSS layer that provides a fast, single-page progressive enhancement to the otherwise static result pages.

Help is welcome! You can use the [Issues page](https://github.com/commonsearch/cosr-front) to suggest improvements, report bugs, or send us Pull Requests!

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

Check out the [Makefile](https://github.com/commonsearch/cosr-front)