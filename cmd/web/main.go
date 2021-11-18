package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"kerseeeHuang.com/snippetbox/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql"	// We don't explicit need this, but database/sql need this.
	"github.com/golangcollege/sessions"
)

type contextKey string

const contextKeyIsAuthenticated = contextKey("isAuthenticated")

// application holds all the application-wide dependencies.
type application struct {
	errorLog		*log.Logger
	infoLog			*log.Logger
	session 		*sessions.Session
	snippets		*mysql.SnippetModel	// Our snippets model connected to the database.
	templateCache 	map[string]*template.Template	// template caches
	users			*mysql.UserModel
}

func main(){
	// Parse the runtime configuration settings for the application.
	// addr is a flag to set HTTP network address
	addr := flag.String("addr", ":4000", "HTTP network address")
	// dsn is a flag to set data source name.
	dsn := flag.String("dsn", "web:satoshi7442@/snippetbox?parseTime=true", "MySQL data source name")
	// secret is a flag to set the encrption key and will be used to authenticate session cookies
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGwhTzbpa@ge", "Secret key")
	flag.Parse()

	// Establishing the dependencies for the handlers
	// infoLog is a logger for writing information messages.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	// errorLog is a logger for writing error messages.
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	
	// Open the DB.
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// Initialize a new template cache.
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	// Initialize a session manager and set its lifetime.
	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	
	// Initialize an application to hold all the dependencies and routes (mux).
	app := &application{
		errorLog: 		errorLog,
		infoLog: 		infoLog,
		session: 		session,
		snippets: 		&mysql.SnippetModel{DB: db},
		templateCache: 	templateCache,
		users: 			&mysql.UserModel{DB: db},
	}

	// Config the curve preferences in TLS.
	tlsConfig := &tls.Config{
		PreferServerCipherSuites: 	true,
		CurvePreferences: 			[]tls.CurveID{tls.X25519, tls.CurveP256},
	}
	
	// Running the HTTP server.
	// Initialize the http server with addr, errorLog and handler defined above.
	// Otherwise the http default server will use stderr to output error.
	srv := &http.Server{
		Addr: *addr,
		ErrorLog: 		errorLog,
		Handler: 		app.routes(),	// Create a mux from app.routes()
		TLSConfig: 		tlsConfig,
		IdleTimeout: 	time.Minute,
		ReadTimeout: 	5 * time.Second,
		WriteTimeout: 	10 * time.Second,
	}
	// Use the http.ListenAndServe() function to start a new web server.
	// Call Fatal if there is any error.
	infoLog.Printf("Starting server on %s\n", *addr)
	// Open a HTTPS server.
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}

// openDB wraps sql.Open() and returns a sql.DB connection pool for a given DSN.
func openDB(dsn string) (*sql.DB, error) {
	// sql.Open does not actually connect to DB but only initialize the pool for future use.
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	// Thus we need db.Ping to test the connection to the db.
	if err = db.Ping(); err != nil {
		return nil , err
	}
	return db, nil
}