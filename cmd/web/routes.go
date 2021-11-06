package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)	

// routes return a http.Handler that routes all requests to corresponding handlers.
func (app *application) routes() http.Handler {
	// Create the standard chain of middleware.
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	
	// Create a mux with third-party package.
	mux := pat.New()
	// Register handlers with the allowed method. The order of statement below MATTERS!
	// Pat will match patterns in the order that these handler are registered.
	mux.Get("/", http.HandlerFunc(app.home))
	mux.Get("/snippet/create", http.HandlerFunc(app.createSnippetForm))
	mux.Post("/snippet/create", http.HandlerFunc(app.createSnippet))
	mux.Get("/snippet/:id", http.HandlerFunc(app.showSnippet))

	// fileServer serve the static files in ./ut/static directory.
	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static/")})
	// Handle all the request with /static/ prefix. Because fileServer only serve
	// the files under /static/, so we need to strip the "/static" in the request
	// so that the fileServer can find the right path.
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}