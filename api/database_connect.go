package main

import (
	"time"
	"database/sql"
	_ "github.com/lib/pq"
	"os"
)

func connectDB(retriesLeft int) error {
	con := os.Getenv("POSTGRES")
	logger.Infof("opening connection to postgres: %s", con)

	var err error
	db, err = sql.Open("postgres", con)
	if err != nil {
		logger.Errorf("cannot open connection to postgres: %v", err)
		return err
	}

	err = db.Ping()
	if err != nil {
		if retriesLeft > 0 {
			logger.Errorf("cannot talk to postgres, retrying in 10 seconds (%d attempts left): %v", retriesLeft-1, err)
			time.Sleep(10 * time.Second)
			return connectDB(retriesLeft - 1)
		} else {
			logger.Errorf("cannot talk to postgres, last attempt failed: %v", err)
			return err
		}
	}

	statement := `
    CREATE TABLE IF NOT EXISTS migrations (
      filename TEXT NOT NULL UNIQUE
    );
  `
	_, err = db.Exec(statement)
	if err != nil {
		logger.Errorf("cannot create migrations table: %v", err)
		return err
	}

	// At most 1000 database connections will be left open in the idle state. This
	// was found to be important when benchmarking with `wrk`: if this was unset,
	// too many open idle connections were present, resulting in dropped requests
	// due to the limit on the number of file handles. On benchmarking, around
	// 100 was found to be pretty optimal.
	db.SetMaxIdleConns(100)

	return nil
}
