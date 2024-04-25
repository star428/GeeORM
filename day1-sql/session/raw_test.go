package session

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

var TestDB *sql.DB

func TestMain(m *testing.M) {
	TestDB, _ = sql.Open("sqlite3", "gee.db")
	code := m.Run()
	TestDB.Close()
	os.Exit(code)
}

func NewSession() *Session {
	return New(TestDB)
}

func TestSession_Exec(t *testing.T) {
	s := NewSession()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	result, _ := s.Raw("INSERT INTO User(`Name`) values(?), (?)", "Tom", "jack").Exec()
	if count, err := result.RowsAffected(); err != nil || count != 2 {
		t.Fatal("expect 2, but got", count)
	}
}

func TestSession_QueryRow(t *testing.T) {
	s := NewSession()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	_, _ = s.Raw("INSERT INTO User(`Name`) values(?), (?)", "Tom", "jack").Exec()

	row := s.Raw("SELECT Name FROM User LIMIT 1;").QueryRow()
	var name string
	if err := row.Scan(&name); err != nil || name != "Tom" {
		t.Fatal("failed to query db")
	}
}
