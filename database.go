package main

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrDBConnStrParseNoSep = errors.New("Missing ':' in connection string")
	ErrDBConnStrParseNotImplemented = func(dbName string) error {
		return errors.New(fmt.Sprintf("The database '%s' is not implemented", dbName))
	}
)

type Database interface {
	CommentService
}

func parseConnectionStr(connectionStr string) (dbName string, dbParams map[string]interface{}, err error) {

	// Derive the database name dbName (sqlite, postgres, mongo etc.)
	dbNameSep := strings.Index(connectionStr, ":")
	if dbNameSep == -1 {
		err = ErrDBConnStrParseNoSep
		return
	}

	dbName = strings.TrimSpace(connectionStr[:dbNameSep])

	// Get the parameters to initialize the database
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
		return nil, ErrDBConnStrParseNotImplemented("postgres")
	case "mongo":
		return nil, ErrDBConnStrParseNotImplemented("mongo")
	default:
		return nil, ErrDBConnStrParseNotImplemented(dbName)
	}
}

var db Database
func LoadDatabase(dbFile string) error {

	var err error
	db, err = NewDatabase("sqlite:file=sqlite3.db")
	return err
}
