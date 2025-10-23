package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

var schema string = `CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(256) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(128) NOT NULL DEFAULT ""
);
CREATE INDEX scheduler_date ON scheduler (date);`
var db *sql.DB

func Init(dbpath string) error {
	_, err := os.Stat(dbpath)
	var install bool
	if err != nil {
		install = true
	}
	db, err = sql.Open("sqlite", dbpath)
	if err != nil {
		return err
	}
	defer db.Close()

	if install {
		_, err = db.Exec(schema)
		if err != nil {
			return err
		}
	}
	return nil
}
