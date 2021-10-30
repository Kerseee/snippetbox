package main

import (
	"fmt"
	"net/http"
	"strconv"
)

// home is a handler function which writes a byte slice containing
// "Hello from Snippetbox" as the response body.
func home(w http.ResponseWriter, r *http.Request) {
	// Send 404 response to client through http.NotFound if the path
	// does not exactly match "/".
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte("Hello from Snippetbox"))
}

// showSnippet is a handler function which shows a specific snippet.
func showSnippet(w http.ResponseWriter, r *http.Request) {
	// Extract the id in URL and parse to int.
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	// If id is wrong or invalid then response 404.
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	// Write the id to the http.ResponseWriter
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

// showSnippet is a handler function which creates a specific snippet.
func createSnippet(w http.ResponseWriter, r *http.Request) {
	// Check if the request is using "POST". If not then response 405.
	if r.Method != http.MethodPost {
		// Customize the headers. Tell client that "POST" is allowed.
		// This set the header map.
		w.Header().Set("Allow", http.MethodPost)

		// // It's only possible to call w.WriteHeader() once per response.
		// // After calling w.WriteHeader(), the response status can't be changed.
		// w.WriteHeader(405)
		// // Warning message should be send after write header. Otherwise it will
		// // automaticly send status 200.
		// w.Write([]byte("Method Not Allowed!"))
		
		// http.Error call w.WriteHeader() and w.Write() behind-of-scene.
		// This function is much used in practice than 
		// calling w.WriteHeader() and w.Write() directly.
		http.Error(w, "Method Not Allowed!", 405)
		return
	}
	
	w.Write([]byte("Create a specific snippet..."))
}