process.env.USE_SAUCE_CONNECT = "1";

var config = require("./wdio.conf.local.js").config;

// https://wiki.saucelabs.com/display/DOCS/Platform+Configurator#/
// https://saucelabs.com/platforms/
// https://docs.travis-ci.com/user/sauce-connect/

var mainDesktopPlatforms = [
	"Linux",
	"OS X 10.11",
	"Windows 10",
	"Windows XP"
];
var mainDesktopBrowsers = [
	["chrome", "48.0"],
	["firefox", "44.0"],
];

config.capabilities = [];

mainDesktopPlatforms.forEach(function(platform) {
	mainDesktopBrowsers.forEach(function(browser) {
		config.capabilities.push({
			browserName: browser[0],
			version: browser[1],
			platform: platform
		})
	});
});

config.capabilities = config.capabilities.concat([

// IE<10 not yet supported
/*
{
    browserName: 'internet explorer',
	version: "8.0",
	platform: 'Windows XP',
},
{
    browserName: 'internet explorer',
	version: "9.0",
	platform: 'Windows 7',
},
*/

{
    browserName: 'internet explorer',
	version: "10.0",
	platform: 'Windows 7',
}, {
    browserName: 'internet explorer',
	version: "11.0",
	platform: 'Windows 7',
}

// Can't test on Safari until https://github.com/SeleniumHQ/selenium-google-code-issue-archive/issues/4136
/*
{
    browserName: 'safari',
	version: "9.0",
	platform: "OS X 10.11",
}, {
    browserName: 'safari',
	version: "8.0",
	platform: "OS X 10.10",
}, {
    browserName: 'safari',
	version: "7.0",
	platform: "OS X 10.9",
}
*/

]);

/*
Uncomment this to debug a specific browser
config.capabilities = [{
    browserName: 'safari',
	version: "7.0",
	platform: "OS X 10.9",
}];
*/

config.capabilities.forEach(function(capability) {
	// capability['tunnel-identifier'] = process.env.TRAVIS_JOB_NUMBER;
	capability.name = 'frontend-uitest';
	capability.build = "build-" + process.env.TRAVIS_BUILD_NUMBER;
	capability.public = true;
});

config.mochaOpts.timeout = config.mochaOpts.timeout * 10;

config.services = ["sauce"];

config.host = 'ondemand.saucelabs.com';
config.port = 80;
config.user = process.env.SAUCE_USERNAME;
config.key = process.env.SAUCE_ACCESS_KEY;
config.logLevel = 'verbose';


exports.config = config;