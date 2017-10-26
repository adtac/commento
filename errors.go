package main

import (
	"errors"
)

var errorList = map[string]error{
	"err.internal": errors.New("Some internal error occurred"),

	"err.request.method.invalid": errors.New("Invalid request method"),

	"err.request.field.missing": errors.New("Missing one or more required fields"),

	"err.request.field.invalid": errors.New("One or more fields is invalid"),

	"err.db.unimplemented": errors.New("Database type not implemented"),

	"err.db.conf.separator.missing": errors.New("Missing separator in connection string"),

	"err.db.conf.key.missing": errors.New("Missing DB configuration key"),

	"err.db.conf.value.missing": errors.New("Missing DB configuration value"),

	"err.db.conf.sqlite.filename.missing": errors.New("sqlite: Filename missing"),
}
