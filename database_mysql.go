// +build !no_mysql

package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type MySqlDatabase struct {
	*sql.DB
}

func init() {
	registeredDatabases[DB_MYSQL] = mysqlInit
}

func mysqlInit(params map[string]string) (Database, error) {
	server, ok := params["Server"]
	if !ok {
		return nil, errorList["err.db.conf.mysql.server.missing"]
	}

	port, ok := params["Port"]
	if !ok {
		port = "3306" // MySQL default port
	}

	user, ok := params["Uid"]
	if !ok {
		return nil, errorList["err.db.conf.mysql.uid.mising"]
	}

	password, ok := params["Password"]
	if !ok {
		return nil, errorList["err.db.conf.mysql.password.missing"]
	}

	database, ok := params["Database"]
	if !ok {
		return nil, errorList["err.db.conf.mysql.database.missing"]
	}

	// Connection string follows this format:
	// [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
    db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, server, port, database))
	if err != nil {
		return nil, err
	}

	statement := `
		CREATE TABLE IF NOT EXISTS comments (
            rowid int not null auto_increment,
			url varchar(2083) not null,
			name varchar(200) not null,
			comment varchar(3000) not null,
			depth int not null,
			time timestamp not null,
			parent int,
            primary key (rowid)
		)
	`
	if _, err = db.Exec(statement); err != nil {
		return nil, err
	} else {
		return &MySqlDatabase{db}, nil
	}
}

func (db *MySqlDatabase) CreateComment(c *Comment) error {
	statement := `
		SELECT depth, parent FROM comments WHERE rowid=?;
	`
	rows, err := db.Query(statement, c.Parent)
	if err != nil {
		return err
	}
	defer rows.Close()

	depth := 0
	for rows.Next() {
		var pParent int
		if err := rows.Scan(&depth, &pParent); err == nil {
			if depth+1 > 5 {
				c.Parent = pParent
			}
		}
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return err
	}

	statement = `
		INSERT INTO comments(url, name, comment, time, depth, parent) VALUES(?, ?, ?, ?, ?, ?);
	`
	_, err = db.Exec(statement, c.URL, c.Name, c.Comment, time.Now(), depth+1, c.Parent)
	return err
}

func (db *MySqlDatabase) GetComments(url string) ([]Comment, error) {
	statement := `
		SELECT rowid, url, comment, name, time, parent FROM comments WHERE url=?;
	`
	rows, err := db.Query(statement, url)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []Comment{}
	for rows.Next() {
		c := Comment{}
		var timeraw string
		if err = rows.Scan(&c.ID, &c.URL, &c.Comment, &c.Name, &timeraw, &c.Parent); err != nil {
			return nil, err
		}

		// Time needs to be parsed separately here.  MySQL returns a []uint8
		// for its date formats, and that cannot be directly converted to
		// *time.Time.
		t, err := time.Parse("2006-01-02 15:04:05.999", timeraw)
		if err != nil {
			return nil, err
		}
		c.Timestamp = t

		comments = append(comments, c)
	}

	return comments, rows.Err()
}
