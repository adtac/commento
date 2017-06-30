package main

import (
	"errors"
	. "fmt"
)

// This is the type of the error function you will create when
// you want to return a new error. Example shown below
type ErrFunc func(vs ...string) error

var errValues = map[string]interface{} {

	"err.conn.parse.no.separator" : func() error {
		return errors.New("Missing ':' in connection string")
	},

	"err.conn.parse.db.not.implemented" : func(dbName string) error {
		return errors.New(Sprintf("The database '%s' is not implemented", dbName))
	},

	"err.conn.parse.key.missing" : func() error {
		return errors.New("Configuration string has an empty key")
	},

	"err.conn.parse.sqlite.file.missing" : func() error {
		return errors.New("Add filename with 'file=[filename].db when using sqlite")
	},

	// Add new errors like follows:
	// "err.generic": func(someValue string) error {
	// 	return errors.New(Sprintf("This value '%s' has caused an error", someValue))
	// },
}

// Gets the actual 'error' value
func Error(errName string, vs ...string) error {
	return errValues[errName].(ErrFunc)(vs...)
}
