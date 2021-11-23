package main

import (
	"net/http"

	"kerseeeHuang.com/snippetbox/ui"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

// routes return a http.Handler that routes all requests to corresponding handlers.
func (app *application) routes() http.Handler {
	// Create the standard chain of middleware.
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// dynamicMiddleware is a chan that contains all middleware specific to dynamic application routes
	dynamicMiddleware := alice.New(app.session.Enable, noSurf, app.authenticate)

	// authenticatedMiddleware is a chan for pages needed user authentication
	authenticatedMiddleware := dynamicMiddleware.Append(app.requireAuthentication)

	// Create a mux with third-party package.
	mux := pat.New()
	// Register handlers with the allowed method. The order of statement below MATTERS!
	// Pat will match patterns in the order that these handler are registered.
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/snippet/create", authenticatedMiddleware.ThenFunc(app.createSnippetForm))
	mux.Post("/snippet/create", authenticatedMiddleware.ThenFunc(app.createSnippet))
	mux.Get("/snippet/:id", dynamicMiddleware.ThenFunc(app.showSnippet))

	// Add routes about user authentication.
	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Post("/user/logout", authenticatedMiddleware.ThenFunc(app.logoutUser))

	// Add ping just for test.
	mux.Get("/ping", http.HandlerFunc(ping))

	// fileServer serve the static files in ./ut/static directory from the ui.Files
	// embedded file system.
	fileServer := http.FileServer(neuteredFileSystem{http.FS(ui.Files)})
	// Handle all the request with /static/ prefix.
	mux.Get("/static/", fileServer)

	return standardMiddleware.Then(mux)
}
