package main

import (
	"strings"
)

type Database interface {
	CommentService
}

// parseConnectionStr will parse a string which designates the type of database to connect to
// which are currently SQLite, Postgres and Mongo but, also the database parameters needed for
// the database to operate. Naturally, these depend on the database but if not parameters are
// specified then default values used for that database.
// A connection string is of the form
//    database:key1=value1;key2=value2
//
// For sqlite this would look like
//    sqlite:file=sqlite3.db
// but since sqlite3.db is the default database name we can use instead
//    sqlite:
func parseConnectionStr(connectionStr string) (dbName string, dbParams map[string]interface{}, err error) {

	// Derive the database name dbName
	//    Ex. sqlite, postgres, mongo etc.
	dbNameSep := strings.Index(connectionStr, ":")
	if dbNameSep == -1 {
		err = Error("err.conn.parse.no.separator")
		return
	}

	dbName = strings.TrimSpace(connectionStr[:dbNameSep])

	// Get the parameters to initialize the database
	//    Ex. filename=myDatabase;host=localhost:1234;
	dbParams = make(map[string]interface{})
	for _, param := range strings.Split(connectionStr[dbNameSep+1:], ";") {
		paramSep := strings.Index(param, "=")
		if paramSep != -1 {
			key := strings.TrimSpace(param[:paramSep])
			if len(key) > 0 {
				values := strings.Split(param[paramSep+1:], ",")
				if len(values) > 1 {
					for i, _ := range values {
						values[i] = strings.TrimSpace(values[i])
					}
					dbParams[key] = values
				} else {
					dbParams[key] = strings.TrimSpace(values[0])
				}
			} else {
				err = Error("err.conn.parse.key.missing")
				return
			}
		}
	}

	return
}

func NewDatabase(connectionStr string) (Database, error) {

	dbName, dbParams, err := parseConnectionStr(connectionStr)
	if err != nil {
		return nil, err
	}

	switch dbName {
	case "sqlite":
		return sqliteInit(dbParams)
	case "postgres":
		return nil, Error("err.conn.parse.db.not.implemented", "postgres")
	case "mongo":
		return nil, Error("err.conn.parse.db.not.implemented", "mongo")
	default:
		return nil, Error("err.conn.parse.db.not.implemented", dbName)
	}
}

var db Database
func LoadDatabase(dbConnectionStr string) error {

	var err error
	db, err = NewDatabase(dbConnectionStr)
	return err
}
