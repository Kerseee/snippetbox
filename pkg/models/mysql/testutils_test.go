package mysql

import (
	"database/sql"
	"os"
	"testing"
)

// newTestDB create a new test DB and return the connection pool and an anonymous function
// that can close the connection pool and teardown the test DB.
func newTestDB(t *testing.T) (*sql.DB, func()) {
	// Open a DB connection pool.
	db, err := sql.Open("mysql", "test_web:pass@/test_snippetbox?parseTime=True&multiStatements=true")
	if err != nil {
		t.Fatal(err)
	}

	// Read the setup script and execute.
	script, err := os.ReadFile("./testdata/setup.sql")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(string(script))
	if err != nil {
		t.Fatal(err)
	}

	return db, func() {
		script, err := os.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}
		_, err = db.Exec(string(script))
		if err != nil {
			t.Fatal(err)
		}

		db.Close()
	}
}
