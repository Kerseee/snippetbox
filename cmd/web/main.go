package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

// application holds all the application-wide dependencies.
type application struct {
	errorLog	*log.Logger
	infoLog		*log.Logger
}

func main(){
	// 1. Parse the runtime configuration settings for the application.
	// addr is a flag to setting HTTP network address
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	// 2. Establishing the dependencies for the handlers
	// infoLog is a logger for writing information messages.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	// errorLog is a logger for writing error messages.
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	// Initial an application to hold all the dependencies and routes (mux).
	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
	}
	
	// 3. Running the HTTP server.
	// Initial the http server with addr, errorLog and handler defined above.
	// Otherwise the http default server will use stderr to output error.
	srv := &http.Server{
		Addr: *addr,
		ErrorLog: errorLog,
		Handler: app.routes(),	// Create a mux from app.routes()
	}
	// Use the http.ListenAndServe() function to start a new web server.
	// Call Fatal if there is any error.
	infoLog.Printf("Starting server on %s\n", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}