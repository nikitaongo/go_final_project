package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

const (
	Layout = "20060102"
)

var schema string = `CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(256) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(128) NOT NULL DEFAULT ""
);
CREATE INDEX scheduler_date ON scheduler (date);`
var Db *sql.DB

var Limit = 50

func Init(dbpath string) error {
	_, err := os.Stat(dbpath)
	var install bool
	if err != nil {
		install = true
	}
	Db, err = sql.Open("sqlite", dbpath)
	if err != nil {
		return err
	}

	if install {
		_, err = Db.Exec(schema)
		if err != nil {
			return err
		}
	}
	return nil
}
