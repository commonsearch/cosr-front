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


#
# Tests & linting
#

# Lint everything
# Closure compiler performs a lint pass. Do we still need jshint or eslint?
lint: minify_js golint

# Lint and test everything
test: lint gotest

# Perform all available linting checks on the Go code
golint:
	go fmt ./server
	golint ./server
	go tool vet -all ./server
	aligncheck ./server
	structcheck ./server
	varcheck ./server
	errcheck -ignoretests ./server

# Dependencies for linting Go code
golint_deps:
	go get -u github.com/kisielk/errcheck
	go get -u github.com/golang/lint/golint
	go get github.com/opennota/check/cmd/aligncheck
	go get github.com/opennota/check/cmd/structcheck
	go get github.com/opennota/check/cmd/varcheck

# Run Go tests
gotest:
	COSR_TESTDATA=1 COSR_PATHFRONT="${PWD}" go test ./server

# Run Go benchmarks
gobench:
	COSR_PATHFRONT="${PWD}" go test ./server -bench=. -benchtime=5s



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

# Minify the JavaScript code
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
