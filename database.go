package main

import (
	"strings"
)

type Database interface {
	CommentService
}

var db Database

type DatabaseServer int

const (
	DB_SQLITE DatabaseServer = iota
	DB_MYSQL
	DB_MSSQL
	DB_POSTGRESQL
)

var DatabaseNames map[DatabaseServer]string = map[DatabaseServer]string{
	DB_SQLITE:     "SQLite",
	DB_MYSQL:      "MySQL",
	DB_MSSQL:      "Microsoft SQL Server",
	DB_POSTGRESQL: "PostgreSQL",
}

type dbInit func(map[string]string) (Database, error)

var registeredDatabases map[DatabaseServer]dbInit = make(map[DatabaseServer]dbInit)

// parseConnectionStr parses the given connectionStr and extracts two pieces
// of information: which database to use and the connection parameters for
// that database. For example, in sqlite3, a filename is sufficient. This will
// be encoded as:
//
//     connectionStr := "sqlite3:file=commento.sqlite3"
//
// Naturally, key=value pairs depend on the database in question. For MongoDB,
// this could be a URL. Multiple key=value pairs can be separated by a
// semicolon. To summarize, the canonical form of this strings is:
//
//     connectionStr := "database:key1=value1;key2=value2;key3=value3"
func parseConnectionStr(connectionStr string) (string, map[string]string, error) {
	dbPos := strings.Index(connectionStr, ":")
	if dbPos == -1 {
		return "", nil, errorList["err.db.conf.separator.missing"]
	}
	dbName := strings.TrimSpace(connectionStr[:dbPos])

	params := make(map[string]string)

	for _, param := range strings.Split(connectionStr[dbPos+1:], ";") {
		equalPos := strings.Index(param, "=")
		if equalPos != -1 {
			key := strings.TrimSpace(param[:equalPos])
			if len(key) == 0 {
				return "", nil, errorList["err.db.conf.key.missing"]
			}

			value := strings.TrimSpace(param[equalPos+1:])
			if len(value) == 0 {
				return "", nil, errorList["err.db.conf.value.missing"]
			}

			params[key] = value
		}
	}

	return dbName, params, nil
}

func LoadDatabase(connectionStr string) error {
	dbName, params, err := parseConnectionStr(connectionStr)
	if err != nil {
		return err
	}

	db = nil
	err = errorList["err.db.unimplemented"]

	// FIXME: Throw an error here if the configured database is not registered.
	//        Currently, a panic will happen for a nil reference.
	switch dbName {
	case "sqlite":
		db, err = registeredDatabases[DB_SQLITE](params)
	case "mysql":
		db, err = registeredDatabases[DB_MYSQL](params)
	}

	return err
}
