package main

import (
	"GeeORM/day1-sql/geeorm"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	engine, _ := geeorm.NewEngine("sqlite3", "../database/gee.db")
	defer engine.Close()
	s := engine.NewSession()
	s.Raw("DROP TABLE IF EXISTS User;").Exec()
	s.Raw("CREATE TABLE User(Name text);").Exec()
	s.Raw("CREATE TABLE User(Name text);").Exec() // Run it twice to test if the table exists

	// Insert
	result, _ := s.Raw("INSERT INTO User(`Name`) VALUES (?), (?)", "Tom", "Sam").Exec()
	count, _ := result.RowsAffected()
	println(count)

	// Query rows
	rows, _ := s.Raw("SELECT Name FROM User;").QueryRows()

	for rows.Next() { // one scan, one next, even the first row
		var name string
		_ = rows.Scan(&name)
		println(name)
	}
}
