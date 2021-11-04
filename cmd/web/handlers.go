package main

import (
	"errors"
	"fmt"
	// "html/template"
	"net/http"
	"path/filepath"
	"strconv"

	"kerseeeHuang.com/snippetbox/pkg/models"
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

	// Show the latest snippets in the database.
	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	for _, snippet := range s {
		fmt.Fprintf(w, "%v\n", snippet)
	}

	// // Initial the path of files that we want toe parse. 
	// // Note that home.page must be the FIRST file in the slice. Otherwise nothing will be shown.
	// files := []string{
	// 	"./ui/html/home.page.tmpl",
	// 	"./ui/html/base.layout.tmpl",
	// 	"./ui/html/footer.partial.tmpl",
	// }
	
	// // Read the template html files. Log errors and return if there is any error.
	// // In here we use template.ParseFiles invoke only "base" template in base.layout.tmpl.
	// // And then templates "title" and "main" will be invoked by template "base".
	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// 	app.serverError(w, err)
	// 	return
	// }
	
	// // Use template to write response body. Log errors if there is any.
	// err = ts.Execute(w, nil)
	// if err != nil {
	// 	app.serverError(w, err)
	// }
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

	// Get data via SnippetModel connected to the database based on given id.
	// If no matching record is found, return a 404 Not Found response.
	s, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
	}

	// Write the id to the http.ResponseWriter
	fmt.Fprintf(w, "%v", s)
}

// showSnippet is a handler function which creates a specific snippet.
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	// Check if the request is using "POST". If not then response 405.
	if r.Method != http.MethodPost {
		// Customize the headers. Tell client that "POST" is allowed.
		// This set the header map.
		w.Header().Set("Allow", http.MethodPost)

		// http.Error call w.WriteHeader() and w.Write() behind-of-scene.
		// This function is much used in practice than 
		// calling w.WriteHeader() and w.Write() directly.
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	// Some dummy data
	title := "0 snail"
	content := "0 snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n- Kobayashi Issa"
	expires := "7"

	// Pass the data to the SnippetModel.Insert() and get back the id of the new record.
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
	}
	
	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}