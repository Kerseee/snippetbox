package main

import (
	"fmt"
	"net/http"
)

// secureHeaders is a middleware that adds security measures for all requests.
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("X-XSS-Protection", "1; mode-block")
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, r)
	})
}

// logRequest is a middleware that write requests information to infoLog.
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

// recoverPanic is a middleware that recover the handler from panic
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
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
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){
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