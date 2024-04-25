package geeorm

import (
	sqllog "GeeORM/day1-sql/log"
	"GeeORM/day1-sql/session"
	"database/sql"
)

type Engine struct {
	db *sql.DB
}

// NewEngine creates a new Engine struct
func NewEngine(driver string, source string) (e *Engine, err error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		sqllog.Error(err)
		return
	}

	// Send a ping to the database to check if the connection is alive
	if err = db.Ping(); err != nil {
		sqllog.Error(err)
		return
	}

	e = &Engine{db: db}
	sqllog.Info("Connect database success")
	return
}

// Close closes the database
func (e *Engine) Close() {
	if err := e.db.Close(); err != nil {
		sqllog.Error("Failed to close database")
	}
	sqllog.Info("Close database success")
}

// NewSession creates a new Session struct
func (e *Engine) NewSession() *session.Session {
	return session.New(e.db)
}
