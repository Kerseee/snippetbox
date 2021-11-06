package main

import (
	"net/http"

	"github.com/justinas/alice"

)	
// routes return a http.Handler that routes all requests to corresponding handlers.
func (app *application) routes() http.Handler {
	// Create the standard chain of middleware.
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	
	// Use the http.NewServeMux() function to initialize a new servemux, then
	// register the home function as the handler for the "/" URL pattern.
	mux := http.NewServeMux()
	// mux.HandleFunc actually equals that we first turn our handle method into
	// http.HandlerFunc type by http.HandlerFunc(handle) then call mux.Handle.
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	// fileServer serve the static files in ./ut/static directory.
	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static/")})
	// Handle all the request with /static/ prefix. Because fileServer only serve
	// the files under /static/, so we need to strip the "/static" in the request
	// so that the fileServer can find the right path.
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}