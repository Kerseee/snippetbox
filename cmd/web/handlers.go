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

// ping writes "OK" to http.Responsewriter.
func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

// home is a handler function which renders the home page.
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

// about is a handler function which renders the about page.
func (app *application) about(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "about.page.tmpl", &templateData{})
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
		return
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
	app.render(w, r, "login.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

// loginUser let user login if the user is valid.
func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Check if credentials are valid. Redisplay the login page if there is any error.
	form := forms.New(r.PostForm)
	id, err := app.users.Authenticate(form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.Errors.Add("generic", "Email or Password is incorrect")
			app.render(w, r, "login.page.tmpl", &templateData{Form: form})
		} else {
			app.serverError(w, err)
		}
		return
	}

	// Add the user id to the session, so that this user is logged in.
	app.session.Put(r, "authenticatedUserID", id)

	// Redirect to the origin path that this client want to before login, if exist.
	redirectLoc := app.session.PopString(r, "redirectLocation")
	if redirectLoc != "" {
		http.Redirect(w, r, redirectLoc, http.StatusSeeOther)
		return
	}

	// Redirect the user to the create snippet page as default.
	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

// logoutUser let user logout.
func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	// Remove the authenticatedUserID of user
	app.session.Remove(r, "authenticatedUserID")
	// Inform user that they are succesfully logged out
	app.session.Put(r, "flash", "You've been logged out succesfully!")
	// Redirect to home page.
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// userProfile show the profile of given user.
func (app *application) userProfile(w http.ResponseWriter, r *http.Request) {
	// Get the id of the authenticated user.
	id := app.session.GetInt(r, "authenticatedUserID")

	// Retreive the data from db.
	user, err := app.users.Get(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Show the user profile.
	app.render(w, r, "profile.page.tmpl", &templateData{
		User: user,
	})
}

// changePasswordForm show the form for users to change their passwords.
func (app *application) changePasswordForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "password.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

// changePassword change the password of the user and redirect to the profile if success.
func (app *application) changePassword(w http.ResponseWriter, r *http.Request) {
	// Parse the form data.
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Validate the data.
	form := forms.New(r.PostForm)
	form.Required("currentPassword", "newPassword", "confirmPassword")
	form.MinLength("newPassword", 10)

	// Confirm password.
	if form.Get("newPassword") != form.Get("confirmPassword") {
		form.Errors.Add("confirmPassword", "Password do not match")
	}

	// Redisplay the template and filled-in data if the form is not valid.
	if !form.Valid() {
		app.render(w, r, "password.page.tmpl", &templateData{Form: form})
		return
	}

	// Update the password of this user.
	id := app.session.GetInt(r, "authenticatedUserID")
	err = app.users.ChangePassword(id, form.Get("currentPassword"), form.Get("newPassword"))
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.Errors.Add("currentPassword", "Wrong password")
			app.render(w, r, "password.page.tmpl", &templateData{Form: form})
		} else {
			app.serverError(w, err)
		}
		return
	}

	// Add a confirmation flash message and redirect to the login page.
	app.session.Put(r, "flash", "Change password successfully!")
	http.Redirect(w, r, "/user/profile", http.StatusSeeOther)
}
