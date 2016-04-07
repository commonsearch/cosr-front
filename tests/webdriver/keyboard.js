describe('keyboard navigation', function() {

	// Key codes: https://seleniumhq.github.io/selenium/docs/api/py/webdriver/selenium.webdriver.common.keys.html

    it('tab & arrow keys from STATIC page', function* () {

        yield openSearchUrl({"q": "apple", "g": "en"});
        var focused = yield inspectElement(yield browser.elementActive());
        assert.equal(focused.tag, "body");

        // Arrow down should select the first result.
        // yield browser.keys(["ARROW_DOWN"]);
        // assert.equal(focused.tag, "a");
        // assert.ok(focused.text.indexOf("Apple") >= 0);

    });

});