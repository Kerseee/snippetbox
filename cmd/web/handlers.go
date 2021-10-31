package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
)

// neuteredFileSystem is a wrapper of http.FileSystem to prevent listing files
// in directories without index.html.
type neuteredFileSystem struct {
	fs http.FileSystem
}

// Open implements the method that can be called when http.FilServer receives a request.
func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	// Open the file path.
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	// Check if this path is a directory.
	s, err := f.Stat()
	if s.IsDir() {
		// If it is a directory, then response the index.html in this directory.
		index := filepath.Join(path, "index.html")
		// If there is no index.html, then return a os.ErrNotExist error
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}
			return nil, err
		}
	}

	// If it is a valid file or a directory with an index file in it, then return 
	// the file or the directory and let http.FileServer to handle it in the following process. 
	return f, nil
}



// home is a handler function which writes a byte slice containing
// "Hello from Snippetbox" as the response body.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Send 404 response to client through http.NotFound if the path
	// does not exactly match "/".
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	// Initial the path of files that we want toe parse. 
	// Note that home.page must be the FIRST file in the slice. Otherwise nothing will be shown.
	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}
	
	// Read the template html files. Log errors and return if there is any error.
	// In here we use template.ParseFiles invoke only "base" template in base.layout.tmpl.
	// And then templates "title" and "main" will be invoked by template "base".
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}
	
	// Use template to write response body. Log errors if there is any.
	err = ts.Execute(w, nil)
	if err != nil {
		app.serverError(w, err)
	}
}

// showSnippet is a handler function which shows a specific snippet.
func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	// Extract the id in URL and parse to int.
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	// If id is wrong or invalid then response 404.
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	// Write the id to the http.ResponseWriter
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

// showSnippet is a handler function which creates a specific snippet.
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
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
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	
	w.Write([]byte("Create a specific snippet..."))
}