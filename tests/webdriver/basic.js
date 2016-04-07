
describe('basic usage', function() {

    it('type from homepage and find apple.com', function* () {

    	// Open homepage
        yield openSearchUrl({});
        sleep(1000);
        assert.equal(yield getPath(), "/");
        assert.equal(yield browser.getTitle(), "Common Search");

        // Make sure the input is focused
        var focused = yield browser.elementActive();
        var elt = yield browser.elementIdAttribute(focused.value.ELEMENT, "id");
        assert.equal(elt.value, "q");

        // Start typing a query
        yield browser.keys(["a"]);
        yield browser.waitForExist("#hits .r", 5000);

        // Some first element should appear
        var hits = yield browser.elements("#hits .r");
        assert.equal(hits.value.length, 25);

        // Type the rest of the query
        yield browser.keys(["p", "p", "l", "e"]);
        yield sleep(2000);

        var hits = yield browser.elements("#hits .r");
        assert.equal(hits.value.length, 25);

        var first_url = yield browser.elementIdElement(hits.value[0].ELEMENT, ".u");
        var first_url_text = yield browser.elementIdText(first_url.value.ELEMENT);
        assert.equal(first_url_text.value, "apple.com");

        // And the URL should have been updated after the history timeout of 2s
        yield sleep(3000);
        assert.equal(yield getPath(), "/?g=en&q=apple");
    });
});