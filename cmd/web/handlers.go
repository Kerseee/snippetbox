package main

import (
	"errors"
	"fmt"
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
	// Show the latest snippets in the database.
	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Render the html with template and data.
	app.render(w, r, "home.page.tmpl", &templateData{
		Snippets: s,
	})
}

// showSnippet is a handler function which shows a specific snippet.
func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	// Extract the id in URL and parse to int.
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
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

	// Render the html with template and data.
	app.render(w, r, "show.page.tmpl", &templateData{
		Snippet: s,
	})
}

// createSnippetForm create the form for client to create a snippet.
func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Create a new snippet..."))
}

// showSnippet is a handler function which creates a specific snippet and store it into DB.
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	// Some dummy data
	title := "0 snail"
	content := "0 snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n- Kobayashi Issa"
	expires := "7"

	// Pass the data to the SnippetModel.Insert() and get back the id of the new record.
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}