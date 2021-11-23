package main

import (
	"bytes"
	"net/http"
	"net/url"
	"testing"
)

func TestPing(t *testing.T) {
	// Create a test app
	app := newTestApplication(t)

	// Create a test server
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	// Mock a client to make a GET requst to "/ping"
	statusCode, _, body := ts.get(t, "/ping")

	// Check status code.
	if statusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, statusCode)
	}

	// Check response body.
	if string(body) != "OK" {
		t.Errorf("want body to equal %q", "OK")
	}
}

func TestShowSnippet(t *testing.T) {
	// Create test app and server
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	// Create test cases.
	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody []byte
	}{
		{"Valid ID", "/snippet/1", http.StatusOK, []byte("An old silent pond...")},
		{"Non-existent ID", "/snippet/2", http.StatusNotFound, nil},
		{"Negative ID", "/snippet/-1", http.StatusNotFound, nil},
		{"Decimal ID", "/snippet/1.58", http.StatusNotFound, nil},
		{"String ID", "/snippet/something", http.StatusNotFound, nil},
		{"Empty ID", "/snippet/", http.StatusNotFound, nil},
		{"Trailing slash", "/snippet/1/", http.StatusNotFound, nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			code, _, body := ts.get(t, test.urlPath)

			if code != test.wantCode {
				t.Errorf("want %d; got %d", test.wantCode, code)
			}

			if !bytes.Contains(body, test.wantBody) {
				t.Errorf("want body to contain %q", test.wantBody)
			}
		})
	}
}

func TestSignupUser(t *testing.T) {
	// Initialize a test app and server.
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	// Mock a client making a GET request to "/user/signup" and extract the csrfToken.
	_, _, body := ts.get(t, "/user/signup")
	csrfToken := extractCSRFToken(t, body)

	tests := []struct {
		name         string
		userName     string
		userEmail    string
		userPassword string
		csrfToken    string
		wantCode     int
		wantBody     []byte
	}{
		{"Valid submission", "Bob", "bob@example.com", "validPa$$word", csrfToken, http.StatusSeeOther, nil},
		{"Empty name", "", "bob@example.com", "validPa$$word", csrfToken, http.StatusOK, []byte("This field cannot be blank")},
		{"Empty email", "Bob", "", "validPa$$word", csrfToken, http.StatusOK, []byte("This field cannot be blank")},
		{"Empty password", "Bob", "bob@example.com", "", csrfToken, http.StatusOK, []byte("This field cannot be blank")},
		{"Invalid email (incomplete domain)", "Bob", "bob@example.", "validPa$$word", csrfToken, http.StatusOK, []byte("This field is invalid")},
		{"Invalid email (missing @)", "Bob", "bobexample.com", "validPa$$word", csrfToken, http.StatusOK, []byte("This field is invalid")},
		{"Invalid email (missing local part)", "Bob", "@example.com", "validPa$$word", csrfToken, http.StatusOK, []byte("This field is invalid")},
		{"Short password", "Bob", "bob@example.com", "pa$$word", csrfToken, http.StatusOK, []byte("This field is too short (minimum is 10 characters)")},
		{"Duplicate email", "Bob", "dup@example.com", "validPa$$word", csrfToken, http.StatusOK, []byte("Email address is already in use")},
		{"Invalid CSRF Token", "", "", "", "wrongToken", http.StatusBadRequest, nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("name", test.userName)
			form.Add("email", test.userEmail)
			form.Add("password", test.userPassword)
			form.Add("csrf_token", test.csrfToken)

			code, _, body := ts.postForm(t, "/user/signup", form)

			if code != test.wantCode {
				t.Errorf("want %d; got %d", test.wantCode, code)
			}

			if !bytes.Contains(body, test.wantBody) {
				t.Errorf("want body %s to contain %q", body, test.wantBody)
			}
		})
	}
}
