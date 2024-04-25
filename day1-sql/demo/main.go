package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, _ := sql.Open("sqlite3", "../database/gee.db")

	defer func() {
		_ = db.Close()
	}()

	db.Exec("DROP TABLE IF EXISTS User;")
	db.Exec("CREATE TABLE User(Name text);")
	r, err := db.Exec("INSERT INTO User(`Name`) values(?), (?)", "Tom", "Sam")
	if err == nil {
		affected, _ := r.RowsAffected()
		println(affected)
	}

	row := db.QueryRow("SELECT Name FROM User LIMIT 1")
	var name string

	if err := row.Scan(&name); err == nil {
		println(name)
	}
}
