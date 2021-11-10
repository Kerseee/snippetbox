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

	// dynamicMiddleware is a chan that contains all middleware specific to dynamic application routes
	dynamicMiddleware := alice.New(app.session.Enable)
	
	// Create a mux with third-party package.
	mux := pat.New()
	// Register handlers with the allowed method. The order of statement below MATTERS!
	// Pat will match patterns in the order that these handler are registered.
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/snippet/create", dynamicMiddleware.ThenFunc(app.createSnippetForm))
	mux.Post("/snippet/create", dynamicMiddleware.ThenFunc(app.createSnippet))
	mux.Get("/snippet/:id", dynamicMiddleware.ThenFunc(app.showSnippet))

	// fileServer serve the static files in ./ut/static directory.
	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static/")})
	// Handle all the request with /static/ prefix. Because fileServer only serve
	// the files under /static/, so we need to strip the "/static" in the request
	// so that the fileServer can find the right path.
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}