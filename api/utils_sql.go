package main

import ()

// scanner is a database/sql abstraction interface that can be used with both
// *sql.Row and *sql.Rows.
type sqlScanner interface {
	// Scan copies columns from the underlying query row(s) to the values
	// pointed to by dest.
	Scan(dest ...interface{}) error
}
