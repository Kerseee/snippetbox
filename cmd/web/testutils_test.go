package main

import (
	"html"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"
	"time"

	"kerseeeHuang.com/snippetbox/pkg/models/mock"

	"github.com/golangcollege/sessions"
)

// csrfTokenRX is a regular expression which captures the CSRF token.
var csrfTokenRX = regexp.MustCompile(`<input type='hidden' name='csrf_token' value='(.+)'>`)

// extractCSRFToken extract the CSRF token from the html body.
func extractCSRFToken(t *testing.T, body []byte) string {
	matches := csrfTokenRX.FindSubmatch(body)
	if len(matches) < 2 {
		t.Fatal("no csrf token found in body")
	}
	return html.UnescapeString(string(matches[1]))
}

// newTestApplication return an application struct for test.
func newTestApplication(t *testing.T) *application {
	// Create a test templateCache.
	templateCache, err := newTemplateCache()
	if err != nil {
		t.Fatal(err)
	}

	// Create a test session manager.
	session := sessions.New([]byte("3dSmsje8xh19sj38cnsl2i38Sja29Si2"))
	session.Lifetime = 12 * time.Hour
	session.Secure = true

	return &application{
		errorLog:      log.New(io.Discard, "", 0),
		infoLog:       log.New(io.Discard, "", 0),
		session:       session,
		snippets:      &mock.SnippetModel{},
		templateCache: templateCache,
		users:         &mock.UserModel{},
	}
}

// testServer is a wrapper of httptest.Server.
type testServer struct {
	*httptest.Server
}

// newTestServer return a pointer of testServer
func newTestServer(t *testing.T, h http.Handler) *testServer {
	// Initialize a test server.
	ts := httptest.NewTLSServer(h)

	// Initialize a cookie jar.
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	// Add the cookie jar to client.
	ts.Client().Jar = jar

	// Set the client to interrupt redirection.
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

// get makes a GET request to a given url on the test server,
// and return the response status code, headers and body
func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, []byte) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	return rs.StatusCode, rs.Header, body
}

// postForm sends POST request to a given url on the test server.
func (ts *testServer) postForm(t *testing.T, urlPath string, form url.Values) (int, http.Header, []byte) {
	rs, err := ts.Client().PostForm(ts.URL+urlPath, form)
	if err != nil {
		t.Fatal(err)
	}

	// Read the response body.
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	return rs.StatusCode, rs.Header, body
}


