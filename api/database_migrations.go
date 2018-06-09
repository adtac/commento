package main

import (
	"io/ioutil"
	"os"
	"strings"
)

func performMigrations() error {
	return performMigrationsFromDir(os.Getenv("STATIC") + "/db")
}

func performMigrationsFromDir(dir string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		logger.Errorf("cannot read directory for migrations: %v", err)
		return err
	}

	statement := `
    SELECT filename
    FROM migrations;
  `
	rows, err := db.Query(statement)
	if err != nil {
		logger.Errorf("cannot query migrations: %v", err)
		return err
	}

	defer rows.Close()

	filenames := make(map[string]bool)
	for rows.Next() {
		var filename string
		if err = rows.Scan(&filename); err != nil {
			logger.Errorf("cannot scan filename: %v", err)
			return err
		}

		filenames[filename] = true
	}

	completed := 0
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			if !filenames[file.Name()] {
				f := dir + string(os.PathSeparator) + file.Name()
				contents, err := ioutil.ReadFile(f)
				if err != nil {
					logger.Errorf("cannot read file %s: %v", file.Name(), err)
					return err
				}

				if _, err = db.Exec(string(contents)); err != nil {
					logger.Errorf("cannot execute the SQL in %s: %v", f, err)
					return err
				}

				statement = `
          INSERT INTO
          migrations (filename)
          VALUES     ($1      );
        `
				_, err = db.Exec(statement, file.Name())
				if err != nil {
					logger.Errorf("cannot insert filename into the migrations table: %v", err)
					return err
				}

				completed++
			}
		}
	}

	if completed > 0 {
		logger.Infof("%d migrations found, %d new migrations completed (%d total)", len(filenames), completed, len(filenames)+completed)
	}

	return nil
}
