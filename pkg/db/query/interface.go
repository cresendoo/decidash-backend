package query

import "database/sql"

// Support database model
type Model interface {
	TableName() string
}

type AutoIncrementModel interface {
	Model
	SetID(id int64) // Support insert query autoincrement result ID
}

type Conn interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
}

type Beginner interface {
	Begin() (*sql.Tx, error)
}
