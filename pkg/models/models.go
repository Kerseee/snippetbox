package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecord = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail = errors.New("models: duplicate email")
)

// Snippet define the structure of a snippet retrieved from the database.
type Snippet struct {
	ID 		int
	Title 	string
	Content string
	Created time.Time
	Expires time.Time
}

// User define the structure of a user retrieved from the database.
type User struct {
	ID				int
	Name			string
	Email			string
	HashedPassword	[]byte
	Created			time.Time
	Active 			bool
}