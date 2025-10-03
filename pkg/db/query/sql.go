package query

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/cresendoo/decidash-backend/pkg/errorx"
)

type stmtType int8

const (
	stmtInsert stmtType = iota + 1
	stmtInsertMulti
	stmtUpdate
	stmtUpdateByID
	stmtUpsert
	stmtSelect
	stmtSelectByID
)

type KeyValue map[string]any

type QueryString struct {
	PreStmt      string
	Stmt         string
	Vals         []any
	OnDuplicates []string
}

/*
	Suport SQL Values : ? --> real value

EX) SELECT * FROM WHERE id = ? --> SELECT * FROM WHERE id = 1
*/
func (q *QueryString) AddValue(vals ...any) *QueryString {
	for _, v := range vals {
		val := reflect.ValueOf(v)
		switch val.Kind() {
		case reflect.Slice:
			for i := 0; i < val.Len(); i++ {
				q.Vals = append(q.Vals, val.Index(i).Interface())
			}
		default:
			q.Vals = append(q.Vals, v)
		}
	}
	return q
}

// Support SQL : IN (?,?,? ...)
func (q *QueryString) AddInClause(col string, values any) *QueryString {
	v := reflect.ValueOf(values)
	if v.Kind() != reflect.Slice || v.Len() == 0 {
		return q
	}

	q.Stmt = q.Stmt + fmt.Sprintf(" %s IN %s", col, makeValuesString(v.Len()))
	return q.AddValue(values)
}

func (q *QueryString) Add(clause string, vals ...any) *QueryString {
	q.Stmt += clause
	return q.AddValue(vals...)
}

func (q *QueryString) AddOnDuplicateUpdate(col string, preVal string, val any) *QueryString {
	q.OnDuplicates = append(q.OnDuplicates, fmt.Sprintf("%s = %s ?", col, preVal))
	return q.AddValue(val)
}

// NewQuery makes new QueryString
func NewQuery(stmt string, v ...any) *QueryString {
	return &QueryString{
		Stmt: stmt,
		Vals: v,
	}
}

// NewQueryWithPreStmt makes new QueryString
func NewQueryWithPreStmt(preStmt, stmt string, v ...any) *QueryString {
	return &QueryString{
		PreStmt: preStmt,
		Stmt:    stmt,
		Vals:    v,
	}
}

/*
makeValuesString for values placeholder

EX) (?,?,?,?)
*/
func makeValuesString(n int) string {
	return fmt.Sprintf("(%s)", strings.Join(stringRepeat("?", n), ", "))
}

func whereInQuery(k string, paramCount int) (string, error) {
	if paramCount < 1 {
		return "", errorx.New("Parameter count must be grater than 0").With("param_count", paramCount)
	}
	return fmt.Sprintf("%s IN (?%s)", k, strings.Repeat(",?", paramCount-1)), nil
}

func makeQueryStringWithCondition(table string, cols []string, andConditions []string) string {
	return fmt.Sprintf("UPDATE %s SET %s=? WHERE %s", table, strings.Join(cols, "=?, "), strings.Join(andConditions, " AND "))
}

func makeQueryString(table string, cols []string, t stmtType) (string, []string) {
	pk := cols[0]
	var stmt string
	switch t {
	case stmtSelect:
		stmt = fmt.Sprintf("SELECT %s FROM %s", strings.Join(cols, ", "), table)
	case stmtSelectByID:
		stmt = fmt.Sprintf("SELECT %s FROM %s WHERE %s=?", strings.Join(cols, ", "), table, pk)
	case stmtInsert:
		stmt = fmt.Sprintf("INSERT INTO %s (%s) VALUES %s",
			table,
			strings.Join(cols, ", "),
			makeValuesString(len(cols)))
	case stmtInsertMulti:
		stmt = fmt.Sprintf("INSERT INTO %s (%s) VALUES ",
			table,
			strings.Join(cols, ", "))
	case stmtUpdateByID:
		stmt = fmt.Sprintf("UPDATE %s SET %s=? WHERE %s=?", table, strings.Join(cols[1:], "=?, "), pk)
		cols = append(cols[1:], pk)
	case stmtUpdate:
		stmt = fmt.Sprintf("UPDATE %s SET %s=?", table, strings.Join(cols, "=?, "))
	case stmtUpsert:
		var onDuplicateStmt strings.Builder
		onDuplicateStmt.WriteString(" ON DUPLICATE KEY UPDATE ")
		var updateCols []string
		for _, col := range cols[1:] {
			updateCols = append(updateCols, fmt.Sprintf("%s=VALUES(%s)", col, col))
		}
		onDuplicateStmt.WriteString(strings.Join(updateCols, ", "))
		stmt = onDuplicateStmt.String()
		cols = cols[1:]
	}
	return stmt, cols
}

func stringRepeat(s string, n int) []string {
	val := make([]string, n)
	for i := 0; i < n; i++ {
		val[i] = s
	}
	return val
}
