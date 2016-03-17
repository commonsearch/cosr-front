/*
  Main JavaScript entry point for CommonSearch result pages.

  This allows rendering results as soon as the user types the first letters of the query.

  Constraints for this JS code are:
   - Progressive enhancement: this whole file is optional
   - Small size: this will be sent to each client. Each byte counts in the final uglified, gzipped file.
   - Compatibility with Closure Compiler's advanced level: https://developers.google.com/closure/compiler/docs/api-tutorial3
   - UI speed: optimize for the feeling of speed from the end user, using tricks to mask network latency
   - Portability: must run the same on a wide range of modern browsers
   - Readability: let our toolchain minify this file, it must be understandable by a regular web developer.
*/

(function(window, document, navigator) {

  // If the browser doesn't support JSON decoding (IE<8), we don't use JS at all.
  // Static version will be fine for them and will save us further headaches.
  if (!window.JSON) return;

  // Sends a HTTP request to retrieve a JSON file
  var requestJSON = function(method, url, options, callback) {
    if (!options) {
      options = {};
    }
    var http = (
      window.XMLHttpRequest ?
      function() {
        return new window.XMLHttpRequest();
      } :
      function() {
        return new window.ActiveXObject('Microsoft.XMLHTTP');
      }
    )();
    http.open(method, url, true);

    // TODO http.setRequestHeader?
    http.onreadystatechange = function() {
      if (http && http["readyState"] == 4) {
        if (http["status"] == 200) {
          callback(null, window.JSON.parse(http.responseText));
        } else {
          callback(http);
        }
      }
    };
    http.send(options.data || "");
    return http;
  };

  // Returns an element from the DOM by its ID.
  var $id = function(id) {
    return document.getElementById(id);
  };

  // Returns an html-safe version of a string
  var htmlSafe = function(str) {
    return String(str).replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');
  };

  // Returns the preferred user language from the browser settings
  var getUserLanguage = function() {
    if (navigator.languages) {
      return navigator.languages[0];
    } else {
      return (navigator.language || (navigator||{})["userLanguage"]);
    }
  };

  // Equivalent of underscorejs's _.debounce. Returns a function that will be called
  // only after the wrapping function has stopped being called for 'wait' milliseconds.
  var debounce = function(func, wait) {
    var timeout;
    return function() {
      if (timeout) window.clearTimeout(timeout);
      timeout = window.setTimeout(func, wait);
    };
  };

  // Guess if a key event is associated to a printable character
  // https://developer.mozilla.org/en-US/docs/Web/API/KeyboardEvent/keyCode
  var isPrintableKeyEvent = function(event) {
	 
	 //When "Esc" is press send focus back to HTML document
	 if(event.keyCode == 27){
		 // Remove focus from any focused element
		 // This action will transfer focus back HTML documents
		 if (document.activeElement) {
		 	document.activeElement.blur();
		 }
	 }
	 
	 // Traverse Search link indirection based on Up and down arrow key press event
	traverseElementOnKeyDown(event);
	    
    return (
      (!event.metaKey && !event.altKey && !event.ctrlKey) &&
      (event.keyCode > 46) &&
      (event.keyCode < 112 || event.keyCode > 145) &&
      (event.keyCode < 91 || event.keyCode > 95)
    );
  };
  
    // Function which is used to traverse search links based on Up and down arrow key
  var traverseElementOnKeyDown = function(event) {
	
	// Find active Element tabIndex
	var tabIndexElement = document.activeElement.tabIndex;
	
	if(tabIndexElement < 8){
		return;
	}
	// Based on Up and Down Key press move focus from next or previous search links
	if(event.keyCode == 38)	{ // up key press
		findElementByAttributeValue('tabIndex',--tabIndexElement).focus();
	}
	if(event.keyCode == 40)	{ // down key press
		findElementByAttributeValue('tabIndex',++tabIndexElement).focus();
	}	
	
  };

  // Find element from document based on Attribute and value
  var findElementByAttributeValue =	function(attribute, value) {
	  var All = document.getElementsByTagName('*');
	  for (var i = 0; i < All.length; i++) {
		if (All[i].getAttribute(attribute) == value) { 
			return All[i]; 
		}
	  }
  };

  // Get the associated URL to a Search object
  // Same function is used on the server side
  // If pageDiff is false, don't include page numbers in URLs
  // Parameters must be alphabetically sorted for caching
  var getSearchHref = function(search, pageDiff) {

    var components = [];

    if (search["g"]) {
      components.push("g=" + search["g"]);
    }

    if (search["q"] && pageDiff !== false && (search["p"] + pageDiff) > 1) {
      components.push("p=" + (search["p"] + pageDiff));
    }

    if (search["q"]) {
      components.push("q=" + encodeURIComponent(search["q"]).replace(/%20/g, "+"));
    }

    if (!components.length) {
      return "/";
    }

    return "/?" + components.join("&");

  };

  // Gets the page title associated to a Search object
  // Same function is done in the Go template
  var makePageTitle = function(search) {
    if (!search["q"]) return "Common Search";
    return search["q"] + " | Common Search";
  };

  // Get references to all the elements we'll use from the DOM
  var eltForm = $id("f"),
      eltSubmit = $id("s"),
      eltSearchInput = $id("q"),
      eltHits = $id("hits"),
      eltPagination = $id("pager"),
      eltDebug = $id("dbg"),
      eltLang = $id("g").childNodes[0],
      eltTitle = document.getElementsByTagName('title')[0];

  // Page layout (are we on the homepage or search results?) is controlled by a single CSS class
  var currentPageLayout = document.body.className;

  // The last search object we sent.
  // TODO: Should this be replaced by a local cache of the N previous results? Or can we rely on
  // the browser's network cache directly in all cases?
  var lastSentSearch = window.JSON.parse(eltForm.getAttribute("data-init") || "{}");

  // The last search object we considered acting upon
  var lastConsideredSearch = lastSentSearch;

  // Current in-flight XMLHttpRequest
  var currentHttpRequest = null;

  // Returns the Search object currently input by the user (not necessarily displayed yet)
  var getCurrentSearch = function() {
    return {
      "q": eltSearchInput.value.trim(),
      "p": parseInt(eltPagination.getAttribute("data-page"), 10) || 1,
      "g": eltLang.value
    };
  };

  // Set a new page layout
  var setPageLayout = function(layoutType) {
    if (layoutType != currentPageLayout) {
      document.body.className = layoutType;
      currentPageLayout = layoutType;
    }
  };

  // Sets a new page title
  var setPageTitle = function(title) {
    eltTitle.innerHTML = htmlSafe(title);
    document.title = title;
  };

  // If the HTML5 history API is available, we can keep track of the current search in the URL
  var historyPushState = function() {};

  if (window.history && window.history.pushState) {

    historyPushState = function() {

      var search = lastConsideredSearch;
      var title = makePageTitle(search);

      if (getSearchHref(search, 0) != lastConsideredSearch) {
        window.history.pushState(
          {"s": search},
          title,
          getSearchHref(search, false)
        );
      }

      setPageTitle(title);
    };

    // On a first page load that's not from history we replace the history with the init state
    if (!window.history.state) {
      window.history.replaceState(
        {"s": lastConsideredSearch},
        makePageTitle(lastConsideredSearch),
        getSearchHref(lastConsideredSearch, 0)
      );
    }

  }

  var historyPushStateDebounced = debounce(historyPushState, 2000);

  // Convert URL for display in search results
  // Same regexp is used on the server side
  var simplifyURL = function(url) {
    return url.replace(/(.*?:\/\/)(([^\/]+)(\/.+)?).*/, "$2").substring(0, 100);
  };

  // Render some search results in JSON form to the page
  // We should aim to have the same result whether we are rending from here
  // or from the Go template!
  var renderHits = function(search, result) {
	
	// TabIndex starts from 9 we already got 8 elements on html page
	var tabIndexCount = 9;
	
    var html = "";
    for (var i = 0; i < (result["h"] || []).length; i++) {
      var hit = result["h"][i];
      html += "<div class='r'>" +
                "<h3><a href='"+hit["u"]+"' tabindex='"+(tabIndexCount+=1)+"'>"+htmlSafe(hit["t"])+"</a></h3>" +
                "<div class='u'><a href='"+hit["u"]+"' tabindex='"+(tabIndexCount+=1)+"'>" + simplifyURL(hit["u"]) + "</a></div>" +
                "<div class='s'>"+htmlSafe(hit["s"])+"</div>" +
              "</div>";
    }
    if (!html && search["q"]) {
      html = "<div class='z'>We didn't find any results for this search, sorry!</div>";
    }
    eltHits.innerHTML = html;

    var paginationHTML = "";
    if (search.p && search.p > 1) {
      paginationHTML = '<a href="' + getSearchHref(search, -1) + '">&laquo; Previous</a>';
    }
    // Do we have more results?
    if (result.m) {
      paginationHTML += '<a href="' + getSearchHref(search, 1) + '">Next &raquo;</a>';
    }
    eltPagination.innerHTML = paginationHTML;

    eltPagination.setAttribute("data-page", search.p);

    if (!result["t"]) {
      eltDebug.innerHTML = "";
    } else {
      var t = result["t"];
      eltDebug.innerHTML = "Text: <span>"+t["tq"]+" / "+t["tr"] + "us</span><br/>" +
                           "Docs: <span>"+t["dq"]+" / "+t["dr"] + "us</span><br/>" +
                           "Total: <span>"+t["o"] + "us</span><br/>";
    }

  };

  // Sends a search request right away
  //  object search: the current Search object
  //  bool fromSubmit: true if it's a 'final' search, false if the user may still type more
  var sendSearch = function(search, fromSubmit) {

    lastSentSearch = search;

    // Empty search will have empty results, no need to send the request
    if (!search["q"]) {
      renderHits(search, {});
      return;
    }

    // We should be in result page in all cases by now
    setPageLayout("hits");

    if (currentHttpRequest) currentHttpRequest.abort();

    currentHttpRequest = requestJSON("GET", "/api/search" + getSearchHref(search, 0).substring(1), {}, function(err, result) {

      currentHttpRequest = null;
      if (err) return; // TODO feedback

      // If the result page called for a straight redirect (bangs may do that), do it right away
      if (result["r"] && fromSubmit) {
        window.location = result["r"];
        return;
      }

      renderHits(search, result);
    });
  };

  // Something happened from the user, we might send a new search request.
  //  bool fromSubmit: true if it's a 'final' search, false if the user may still type more
  var newSearch = function(fromSubmit, fromHistory) {

    lastConsideredSearch = getCurrentSearch();

    if (fromSubmit && !fromHistory) {
      historyPushState();
    }

    // We might avoid doing the request completely if it's similar to the last one we sent
    if (getSearchHref(lastConsideredSearch, 0) == getSearchHref(lastSentSearch, 0)) {
      return;
    }

    if (!fromSubmit && !fromHistory && currentPageLayout != "full") {
      historyPushStateDebounced();
    }

    // Currently, any new search starts from page 1 again. Will be changed with infinite scrolling
    lastConsideredSearch["p"] = 1;

    // TODO: intelligently cancel and space out requests
    return sendSearch(lastConsideredSearch, fromSubmit);

  };

  // We debounce keypress events: start a search only when the user does a pause when typing
  var newSearchDebounced = debounce(function() {
    newSearch(false, false);
  }, 150);


  // Before binding events, we autodetect the language if it was empty on load!
  if (!lastSentSearch["g"]) {
    eltLang.value = (getUserLanguage() || "en").substring(0, 2).toLowerCase();

    // If the option didn't exist, the value won't be actually saved.
    // TODO: is this safe in all browsers?
    if (!eltLang.value) {
      eltLang.value = "en";
    }
  }

  // The final stage is binding event handlers to DOM elements

  // Bind on the "keydown" handler. It's the first keyboard event we can receive
  eltSearchInput.onkeydown = function(event) {

    // Some non-printable keys shouldn't trigger an immediate shift to a result page
    if (isPrintableKeyEvent(event)) {
      setPageLayout("hits");
    }

    newSearchDebounced();

    // No need to propagate further and trigger onchange() events etc.
    event.stopPropagation();

  };

  // Submitting the form causes a search right away
  eltForm.onsubmit = function(evt) {
    newSearch(true, false);
    return false;
  };

  // This will be called without any prior onkeydown in some edge cases
  // (For instance: copypasting, third-party manipulation, ... )
  eltSearchInput.onchange = function() {
    if (lastConsideredSearch["q"] == getCurrentSearch()["q"]) return;
    newSearchDebounced();
  };

  eltSearchInput.onmouseup = eltSearchInput.onchange;

  // Any printable keypress outside of the text area brings us back in focus
  document.body.onkeydown = function(event) {
    if (isPrintableKeyEvent(event)) {
      eltSearchInput.value = eltSearchInput.value.replace(/\s$/, "") + " ";
      eltSearchInput.focus();
      window.scrollTo(0, 0);
    }
  };

  // Lang dropdown was used
  eltLang.onchange = function() {
    if (lastConsideredSearch["g"] == getCurrentSearch()["g"]) return;
    newSearchDebounced();
  };

  // User has hit the back or forward history button
  window.onpopstate = function(event) {
    if (!event.state || !event.state["s"]) return;
    eltSearchInput.value = htmlSafe(event.state["s"]["q"]);
    newSearch(true, true);
  };

  // If the user added some input before the JS was fully loaded, we want to capture it
  newSearch(false, false);

})(window, document, navigator);
