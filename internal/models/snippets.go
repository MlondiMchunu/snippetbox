package models

import (
	"database/sql"
	"time"
)

/*Define a Snippet type to hold the data for an individual snippet*/
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// define a SnippetModel type which wraps a sql.DB connection pool
type SnippetModel struct {
	DB *sql.DB
}

// this will insert a new snippeet into the database
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	return 0, nil
}

// this will return a specific snippet based on its id
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	return nil, nil
}
