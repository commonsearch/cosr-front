package main

import (
	"github.com/NYTimes/gziphandler"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
	"path"
)

// CreateRouter creates a new router with all HTTP handlers and middlewares.
func CreateRouter() http.Handler {

	// Instanciate a router cf. https://github.com/julienschmidt/httprouter
	router := httprouter.New()

	// Middlewares shared by all dynamic pages
	commonMiddleware := alice.New(gziphandler.GzipHandler)

	// Main HTML search route
	router.Handler("GET", "/", commonMiddleware.ThenFunc(SearchHandler))

	// Main JSON search route
	router.Handler("GET", "/api/search", commonMiddleware.ThenFunc(APISearchHandler))

	// Static asset directories
	ServeStaticDirectory(router, "js", true)
	ServeStaticDirectory(router, "css", true)
	ServeStaticDirectory(router, "img", false)

	// We have a whitelist of allowed static files in the root directory
	staticFiles := []string{"/favicon.ico", "/apple-touch-icon-precomposed.png"}

	for _, file := range staticFiles {
		router.Handler("GET", file, http.FileServer(http.Dir(path.Join(Config.PathFront, "static/"))))
	}
	return router
}

// ServeStaticDirectory allows any file inside a directory to be served over HTTP.
func ServeStaticDirectory(r *httprouter.Router, directory string, gzip bool) {

	fileServer := http.FileServer(http.Dir(path.Join(Config.PathFront, "static/"+directory)))
	gzFileServer := gziphandler.GzipHandler(fileServer)

	handler := func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {

		req.URL.Path = ps.ByName("filepath")

		if gzip {
			gzFileServer.ServeHTTP(w, req)
		} else {
			fileServer.ServeHTTP(w, req)
		}
	}

	r.GET("/"+directory+"/*filepath", handler)

}
