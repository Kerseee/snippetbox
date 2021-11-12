package main

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"


	"kerseeeHuang.com/snippetbox/pkg/forms"
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
	// Render a blank form.
	app.render(w, r, "create.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

// showSnippet is a handler function which creates a specific snippet and store it into DB.
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	// Parse the from in the request and store it in r.PostForm.
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	// Retrieve data in the r.PostForm and validate the data
	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")

	// Redisplay the template and filled-in data if the form is not valid.
	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{Form: form})
		return
	}

	// Create a new snippet in db and get back the id of the new record.
	id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Add session data to show flash information.
	app.session.Put(r, "flash", "Snippet successfully created!")
	
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

// signupUserForm shows the sign up form to client.
func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

// signupUser create a new user and store user information into db.
func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	// Parse the form data.
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Validate the data.
	form := forms.New(r.PostForm)
	form.Required("name", "email", "password")
	form.MaxLength("name", 255)
	form.MaxLength("email", 255)
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 10)

	// Redisplay the template and filled-in data if the form is not valid.
	if !form.Valid() {
		app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		return
	}
	
	// Create an user if it is valid. Otherwise redisplay the signup form.
	err = app.users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.Errors.Add("email", "Email address is already in use")
			app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		} else {
			app.serverError(w, err)
		}
		return
	}

	// Add a confirmation flash message and redirect to the login page.
	app.session.Put(r, "flash", "Your signup was successful. Please log in.")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)

}

// loginUserForm show a login form to client.
func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Display the user login form...")
}

// loginUser let user login if the user is valid.
func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Authenticate and login the user...")
}

// logoutUser let user logout.
func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Logout the user...")
}