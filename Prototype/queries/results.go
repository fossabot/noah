package queries

import (
	"database/sql"
)

type QueryResult struct {
	NodeID int
	Rows   *sql.Rows
	Error  error
}
