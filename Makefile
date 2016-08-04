PWD := $(shell pwd)

#
# Day-to-day usage commands
#

# Run server for local development
devserver: gobuild
	GODEBUG=gctrace=1 COSR_DEBUG=1 ./build/cosr-front.bin

# Run server for production
runserver: gobuild
	./build/cosr-front.bin

# Save all Go dependencies to the vendor/ directory
godep_save:
	GO15VENDOREXPERIMENT=1 godep save -v ./server

# Logins into the container
docker_shell:
	docker run -e DOCKER_HOST --rm -v "$(PWD):/go/src/github.com/commonsearch/cosr-front:rw" -w /go/src/github.com/commonsearch/cosr-front -p 9700:9700 -i -t commonsearch/local-front bash

# Run server for local development in a container
docker_devserver:
	docker run -e DOCKER_HOST --rm -v "$(PWD):/go/src/github.com/commonsearch/cosr-front:rw" -w /go/src/github.com/commonsearch/cosr-front -p 9700:9700 -i -t commonsearch/local-front make devserver

# Starts the local services needed by cosr-front
start_services:
	docker run -e DOCKER_HOST -d -p 39200:9200 -p 39300:9300 commonsearch/local-elasticsearch

# Starts the local services needed by cosr-front with the devindex
start_services_devindex:
	docker run -e DOCKER_HOST -d -p 39200:9200 -p 39300:9300 commonsearch/local-elasticsearch-devindex

# Stops local services
stop_services:
	bash -c 'docker ps | tail -n +2 | grep -E "((commonsearch/local-elasticsearch))" | cut -d " " -f 1 | xargs docker stop -t=0'


#
# Tests & linting
#

# Lint everything
lint: jslint golint

# Lint and test everything
test: lint gotest

# Lint and test everything inside Docker
docker_test:
	docker run -e DOCKER_HOST --rm -v "$(PWD):/go/src/github.com/commonsearch/cosr-front:rw" -w /go/src/github.com/commonsearch/cosr-front -i -t commonsearch/local-front make test

# Perform static checks on the Go code
# See new ideas at https://github.com/alecthomas/gometalinter
golint:
	go fmt -n ./server/ | sed -e "s/gofmt -l/gofmt -s -l/g" | sh
	golint ./server
	test -z "$$(golint ./server)"
	go tool vet -all -shadow ./server
	aligncheck ./server
	structcheck ./server
	varcheck ./server
	errcheck -ignoretests ./server
	ineffassign ./server/
	unconvert ./server/
	gosimple ./server/
	staticcheck ./server/

# Dependencies for linting Go code
golint_deps:
	go get -u github.com/kisielk/errcheck
	go get -u github.com/golang/lint/golint
	go get github.com/opennota/check/cmd/aligncheck
	go get github.com/opennota/check/cmd/structcheck
	go get github.com/opennota/check/cmd/varcheck
	go get github.com/gordonklaus/ineffassign
	go get github.com/mdempsky/unconvert
	go get honnef.co/go/simple/cmd/gosimple
	go get honnef.co/go/staticcheck/cmd/staticcheck

# Run Go tests
gotest:
	COSR_TESTDATA=1 COSR_PATHFRONT="${PWD}" go test ./server

# Run Go benchmarks
gobench:
	COSR_PATHFRONT="${PWD}" go test ./server -bench=. -benchtime=5s

# Static linting on the JS code
jslint:
	node_modules/jshint/bin/jshint static/js/*.js
	node_modules/eslint/bin/eslint.js static/js/index.js
	java -jar tools/closure-compiler/compiler.jar --warning_level VERBOSE --summary_detail_level 3 --compilation_level ADVANCED --use_types_for_optimization --language_in ECMASCRIPT5_STRICT --js static/js/index.js > /dev/null

# Run local UI tests with PhantomJS
uitest:
	wdio tests/wdio.conf.local.js

docker_uitest:
	docker run -e DOCKER_HOST --rm -v "$(PWD):/go/src/github.com/commonsearch/cosr-front:rw" -w /go/src/github.com/commonsearch/cosr-front -i -t commonsearch/local-front make uitest

# Run UI tests on Sauce Labs
uitest_sauce:
	wdio tests/wdio.conf.sauce.js

docker_uitest_sauce:
	docker run -e DOCKER_HOST -e SAUCE_USERNAME -e SAUCE_ACCESS_KEY -e TRAVIS_BUILD_NUMBER --rm -v "$(PWD):/go/src/github.com/commonsearch/cosr-front:rw" -w /go/src/github.com/commonsearch/cosr-front -i -t commonsearch/local-front make uitest_sauce


#
# Build commands
#

# Lint then build everything
build: lint build_static

# Build static content, including images
build_static: minify
	mkdir -p build/static/img
	cp -R static/img build/static/
	cp static/*.png static/*.ico build/static/

# Minify everything
minify: minify_js minify_css

# Minify the JavaScript code (will also do some more linting)
minify_js:

	mkdir -p build/static/js/

	# v20150315 is the last version that doesn't crash. TODO, report it!
	# https://github.com/google/closure-compiler/wiki/Binary-Downloads
	java -jar tools/closure-compiler/compiler.jar --warning_level VERBOSE --summary_detail_level 3 --compilation_level ADVANCED --use_types_for_optimization --language_in ECMASCRIPT5_STRICT --js_output_file build/static/js/index.js --js static/js/index.js
	ls -la build/static/js/index.js

# Minify the CSS code
minify_css:
	mkdir -p build/static/css/

	cat static/css/global.css > build/static/css/index.scss
	cat static/css/header.css >> build/static/css/index.scss
	cat static/css/footer.css >> build/static/css/index.scss
	cat static/css/hits.css >> build/static/css/index.scss
	cat static/css/responsive.css >> build/static/css/index.scss

	sass --scss build/static/css/index.scss build/static/css/index.css --style compressed --sourcemap=none --no-cache
	rm build/static/css/index.scss
	ls -la build/static/css/index.css

# Build the Go server code
gobuild:
	mkdir -p build
	go build -o build/cosr-front.bin ./server

# Build local Docker images
docker_build:
	docker build -t commonsearch/local-front .

# Pull Docker images from the registry
docker_pull:
	docker version
	docker info
	docker pull commonsearch/local-front
	docker pull commonsearch/local-elasticsearch
	docker pull commonsearch/local-elasticsearch-devindex
