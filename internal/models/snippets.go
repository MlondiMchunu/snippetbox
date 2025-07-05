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

type SnippetModel struct {
	DB *sql.DB
}
