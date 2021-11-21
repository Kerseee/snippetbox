package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecureHeaders(t *testing.T) {
	// Initialize a new response recorder and a request
	rr := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/", nil) 
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock next handler to call after calling the middleware.
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		w.Write([]byte("OK"))
	})

	// Pass the mock handler to secureHeaders.
	secureHeaders(next).ServeHTTP(rr, r)

	// Get the response.
	rs := rr.Result()

	// Check if secureHeaders correctly sets the X-Frame-Options header
	frameOptions := rs.Header.Get("X-Frame-Options")
	want := "deny"
	if frameOptions != want {
		t.Errorf("want %q; got %q", want, frameOptions)
	}

	// Check if secureHeaders correctly sets the X-XSS-Protection header
	xssProtection := rs.Header.Get("X-XSS-Protection")
	want = "1; mode-block"
	if xssProtection != want {
		t.Errorf("want %q; got %q", want, frameOptions)
	}

	// Check if secureHeaders correctly calls the next handler
	if rs.StatusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, rs.StatusCode)
	}

	// Check the response body. 
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != "OK" {
		t.Errorf("want body to eqyal %q", "OK")
	}
}