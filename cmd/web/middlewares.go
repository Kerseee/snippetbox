package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"kerseeeHuang.com/snippetbox/pkg/models"

	"github.com/justinas/nosurf"
)

// secureHeaders is a middleware that adds security measures for all requests.
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode-block")
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, r)
	})
}

// logRequest is a middleware that write requests information to infoLog.
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

// recoverPanic is a middleware that recover the handler from panic
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Recover from the panic if there is any
		defer func() {
			if err := recover(); err != nil {
				// Set connection-close header on the response
				w.Header().Set("Connection", "close")
				// Call serverError
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// requireAuthentication is a middleware that redirects unauthenticated user to the login page.
func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Redirect unauthenticated user to the login page and return from the middleware chain.
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Prevent user's browser from caching pages that require authentication.
		w.Header().Add("Cache-Control", "no-store")

		// Move on the next handler
		next.ServeHTTP(w, r)
	})
}

// noSurf is a middleware that wraps the next handler with a customized CSRF cookie.
func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})

	return csrfHandler
}

// authenticate is a middleware that create a copy of the request context with
// authenticated key and pass the copy to the next handler if the current user
// is authenticated and active. Otherwise directly move on to the next handler.
func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if current user is authenticated.
		exist := app.session.Exists(r, "authenticatedUserID")
		if !exist {
			next.ServeHTTP(w, r)
			return
		}

		// Check if this user exists in DB.
		user, err := app.users.Get(app.session.GetInt(r, "authenticatedUserID"))
		if errors.Is(err, models.ErrNoRecord) {
			app.session.Remove(r, "authenticatedUserID")
			next.ServeHTTP(w, r)
			return
		}
		if err != nil {
			app.serverError(w, err)
			return
		}

		// Check if this user is deactive.
		if !user.Active {
			app.session.Remove(r, "authenticatedUserID")
			next.ServeHTTP(w, r)
			return
		}

		// Mark the request from this user so that the request indicates it is from an
		// authenticated and active user.
		ctx := context.WithValue(r.Context(), contextKeyIsAuthenticated, true)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
