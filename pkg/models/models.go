package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record found")

// Snippet define the structure of a snippet retreived from the database.
type Snippet struct {
	ID int
	Title string
	Content string
	Create time.Time
	Expires time.Time
}