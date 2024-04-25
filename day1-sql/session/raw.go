package session

import (
	sqllog "GeeORM/day1-sql/log"
	"database/sql"
	"strings"
)

// Session is a struct that contains a database pointer, sql language, and sql variables
type Session struct {
	// db is a database pointer
	db *sql.DB
	// sql is sql language such as "select * from user where id = ?"
	sql strings.Builder
	// sqlVars is the variables in the sql language, such as 1, 2, 3
	sqlVars []any
}

// New creates a new Session struct
func New(db *sql.DB) *Session {
	return &Session{db: db}
}

// Clear clears the sql language and sql variables
func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
}

// DB returns the database pointer
func (s *Session) DB() *sql.DB {
	return s.db
}

// Raw sets the sql language and sql variables
func (s *Session) Raw(sql string, values ...any) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ") // add a space after one sql language to another
	s.sqlVars = append(s.sqlVars, values...)
	return s
}

// Exec executes the sql language and sql variables
func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	sqllog.Info(s.sql.String(), s.sqlVars)
	if result, err = s.DB().Exec(s.sql.String(), s.sqlVars...); err != nil {
		sqllog.Error(err)
	}
	return
}

// QueryRow executes the sql language and sql variables and returns a row
func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	sqllog.Info(s.sql.String(), s.sqlVars)
	return s.DB().QueryRow(s.sql.String(), s.sqlVars...) // 返回一条记录
}

// QueryRows executes the sql language and sql variables and returns rows
func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
	sqllog.Info(s.sql.String(), s.sqlVars)
	if rows, err = s.DB().Query(s.sql.String(), s.sqlVars...); err != nil {
		sqllog.Error(err)
	}
	return
}
