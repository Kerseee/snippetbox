package mysql

import (
	"database/sql"
	"errors"

	"kerseeeHuang.com/snippetbox/pkg/models"
)

// SnippetModel is a wrapper of sql.DB connection pool.
type SnippetModel struct {
	DB *sql.DB
}

// Insert inserts a new snippet into the database.
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	// stmt is a statement of inserting data into the database. 
	// '?'s are placeholder parameters.
	stmt := `INSERT INTO snippets (title, content, created, expires)
		VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
	
	// Use DB.Exec() to execute the statement with placeholder parameters and get the result.
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}
	
	// Get the id of snippet that we just insert.
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Get return a specific snippet based on given id.
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
		WHERE expires > UTC_TIMESTAMP() AND id = ?`
	
	// Use DB.QueryRow to retreive the data.
	row := m.DB.QueryRow(stmt, id)

	// Initial a pointer to a new zeroed snippet struct
	s := &models.Snippet{}

	// Use row.Scan to copy the value in the row into s. 
	// The number of arguments must be exactly the same as the number of columns
	// returned by DB.QueryRow.
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		// Check if the error is the sql.ErrNoRows error.
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		// If not, then return the error itself.
		return nil, err
	}

	return s, nil
}

// Latest return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
		WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`
	
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	// Close is needed for rows before return, otherwise it might cause that all the connections
	// in the pool being used up if something goes wrong here.
	defer rows.Close()

	snippets := []*models.Snippet{}
	// Iterate the rows and store all records into our data structure: snippets.
	// It will automatically close itself and frees-up the underlying database connection after
	// iterating all rows in it.
	for rows.Next() {
		s := &models.Snippet{}
		// This Scan scan the current row in this iteration.
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	// rows.Err() is needed to be call after iterate all the rows by rows.Next().
	// It will return any errors that happened during the iterating.
	if err = rows.Err(); err != nil {
		return nil, err
	}
	
	return snippets, nil
}