const spawn = require('child_process').spawn;
const execSync = require('child_process').execSync;


PHANTOMJS_PROCESS = null;
COSR_FRONT_PROCESS = null;
SAUCE_CONNECT_PROCESS = null;

exports.config = {

    //
    // ==================
    // Specify Test Files
    // ==================
    // Define which test specs should run. The pattern is relative to the directory
    // from which `wdio` was called. Notice that, if you are calling `wdio` from an
    // NPM script (see https://docs.npmjs.com/cli/run-script) then the current working
    // directory is where your package.json resides, so `wdio` will be called from there.
    //
    specs: [
        'tests/webdriver/*.js'
    ],
    // Patterns to exclude.
    exclude: [
        // 'path/to/excluded/files'
    ],
    //
    // ============
    // Capabilities
    // ============
    // Define your capabilities here. WebdriverIO can run multiple capabilities at the same
    // time. Depending on the number of capabilities, WebdriverIO launches several test
    // sessions. Within your capabilities you can overwrite the spec and exclude options in
    // order to group specific specs to a specific capability.
    //
    // First, you can define how many instances should be started at the same time. Let's
    // say you have 3 different capabilities (Chrome, Firefox, and Safari) and you have
    // set maxInstances to 1; wdio will spawn 3 processes. Therefore, if you have 10 spec
    // files and you set maxInstances to 10, all spec files will get tested at the same time
    // and 30 processes will get spawned. The property handles how many capabilities
    // from the same test should run tests.
    //
    capabilities: [{
        browserName: 'phantomjs'
    }],
    //
    // ===================
    // Test Configurations
    // ===================
    // Define all options that are relevant for the WebdriverIO instance here
    //
    // Level of logging verbosity: silent | verbose | command | data | result | error
    logLevel: "error",

    // Enables colors for log output.
    coloredLogs: true,
    //
    // Saves a screenshot to a given path if a command fails.
    screenshotPath: './errorShots/',
    //
    // Set a base URL in order to shorten url command calls. If your url parameter starts
    // with "/", then the base url gets prepended.
    baseUrl: 'http://localhost:9700',
    //
    // Default timeout for all waitFor* commands.
    waitforTimeout: 10000,
    //
    // Default timeout in milliseconds for request
    // if Selenium Grid doesn't send response
    connectionRetryTimeout: 90000,
    //
    // Default request retries count
    connectionRetryCount: 3,
    //
    // Initialize the browser instance with a WebdriverIO plugin. The object should have the
    // plugin name as key and the desired plugin options as properties. Make sure you have
    // the plugin installed before running any tests. The following plugins are currently
    // available:
    // WebdriverCSS: https://github.com/webdriverio/webdrivercss
    // WebdriverRTC: https://github.com/webdriverio/webdriverrtc
    // Browserevent: https://github.com/webdriverio/browserevent
    // plugins: {
    //     webdrivercss: {
    //         screenshotRoot: 'my-shots',
    //         failedComparisonsRoot: 'diffs',
    //         misMatchTolerance: 0.05,
    //         screenWidth: [320,480,640,1024]
    //     },
    //     webdriverrtc: {},
    //     browserevent: {}
    // },
    //
    // Test runner services
    // Services take over a specific job you don't want to take care of. They enhance
    // your test setup with almost no effort. Unlike plugins, they don't add new
    // commands. Instead, they hook themselves up into the test process.
    services: [],
    //
    // Framework you want to run your specs with.
    // The following are supported: Mocha, Jasmine, and Cucumber
    // see also: http://webdriver.io/guide/testrunner/frameworks.html
    //
    // Make sure you have the wdio adapter package for the specific framework installed
    // before running any tests.
    framework: 'mocha',
    //
    // Test reporter for stdout.
    // The following are supported: dot (default), spec, and xunit
    // see also: http://webdriver.io/guide/testrunner/reporters.html
    reporters: ['dot'],

    //
    // Options to be passed to Mocha.
    // See the full list at http://mochajs.org/
    mochaOpts: {
        ui: 'bdd',
        timeout: 20000
    },
    //
    // =====
    // Hooks
    // =====
    // WedriverIO provides several hooks you can use to interfere with the test process in order to enhance
    // it and to build services around it. You can either apply a single function or an array of
    // methods to it. If one of them returns with a promise, WebdriverIO will wait until that promise got
    // resolved to continue.
    //
    // Gets executed once before all workers get launched.
    onPrepare: function (config, capabilities) {

        console.log("Spawning local cosr-front process...");
        execSync("make gobuild")
        COSR_FRONT_PROCESS = spawn("build/cosr-front.bin");

        if (process.env.USE_SAUCE_CONNECT) {
            console.log("Spawning local Sauce Connect proxy...");
            SAUCE_CONNECT_PROCESS =  spawn("sc");
            execSync('while ! timeout 1 bash -c "echo > /dev/tcp/localhost/4445" 2>/dev/null; do sleep 0.1; done');
        } else {
            console.log("Spawning local PhantomJS process...");
            PHANTOMJS_PROCESS = spawn("phantomjs", ["--webdriver=4444"]);
            execSync('while ! timeout 1 bash -c "echo > /dev/tcp/localhost/4444" 2>/dev/null; do sleep 0.1; done');
        }

        // Wait synchronously for the ports to be open
        execSync('while ! timeout 1 bash -c "echo > /dev/tcp/localhost/9700" 2>/dev/null; do sleep 0.1; done');

        console.log("Local ports are open!");

    },


    //
    // Gets executed before test execution begins. At this point you can access all global
    // variables, such as `browser`. It is the perfect place to define custom commands.
    before: function (capabilities, specs) {

        DOCKER_HOST = "127.0.0.1";
        COSR_PORT = 9700;

        openSearchUrl = function(opts) {
            var url = "http://" + DOCKER_HOST + ":" + COSR_PORT + "/";
            var qs = [];
            for (key in opts) {
                qs.push(key + "=" + encodeURIComponent(opts[key]));
            }
            qs.sort();
            if (qs.length > 1) {
                url += "?" + qs.join("&");
            }
            console.log("Opening " + url);
            return browser.url(url);
        };
        assert = require('assert');
        sleep = function(ms) {
            return function(done) {
                setTimeout(done, ms);
            }
        };
        getPath = function() {
            var urlmodule = require("url");
            return browser.url().then(function(res) {
                return urlmodule.parse(res.value).path;
            });
        };
        inspectElement = function(elt) {
            var inspect = {};

            // TODO: is there an easier way to do that?
            return browser.elementIdName(elt.value.ELEMENT).then(function(res) {
                inspect.tag = res.value;
                return browser.elementIdText(elt.value.ELEMENT);
            }).then(function(res) {
                inspect.text = res.value;
                return browser.elementIdLocation(elt.value.ELEMENT);
            }).then(function(res) {
                inspect.location = res.value;
                console.log("Inspect: <" + inspect.tag + "> at ("+inspect.location.x+","+inspect.location.y+")");
                return inspect;
            })
        };
    },
    //
    // Hook that gets executed before the suite starts
    // beforeSuite: function (suite) {
    // },
    //
    // Hook that gets executed _before_ a hook within the suite starts (e.g. runs before calling
    // beforeEach in Mocha)
    // beforeHook: function () {
    // },
    //
    // Hook that gets executed _after_ a hook within the suite starts (e.g. runs after calling
    // afterEach in Mocha)
    // afterHook: function () {
    // },
    //
    // Function to be executed before a test (in Mocha/Jasmine) or a step (in Cucumber) starts.
    // beforeTest: function (test) {
    // },
    //
    // Runs before a WebdriverIO command gets executed.
    // beforeCommand: function (commandName, args) {
    // },
    //
    // Runs after a WebdriverIO command gets executed
    // afterCommand: function (commandName, args, result, error) {
    // },
    //
    // Function to be executed after a test (in Mocha/Jasmine) or a step (in Cucumber) starts.
    // afterTest: function (test) {
    // },
    //
    // Hook that gets executed after the suite has ended
    // afterSuite: function (suite) {
    // },
    //
    // Gets executed after all tests are done. You still have access to all global variables from
    // the test.
    // after: function (capabilities, specs) {

    // },

    // Gets executed after all workers got shut down and the process is about to exit. It is not
    // possible to defer the end of the process using a promise.
    onComplete: function(exitCode) {
        if (PHANTOMJS_PROCESS) {
            PHANTOMJS_PROCESS.kill();
        }
        if (COSR_FRONT_PROCESS) {
            COSR_FRONT_PROCESS.kill();
        }
        if (SAUCE_CONNECT_PROCESS) {
            SAUCE_CONNECT_PROCESS.kill();
        }
    }
}
