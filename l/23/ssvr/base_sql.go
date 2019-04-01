package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pschlump/MiscLib"
	"github.com/pschlump/godebug"
)

// LogQueries is called with all statments to log them to a file.
func logQueries(stmt string, err error, data []interface{}, elapsed time.Duration) {
	if logFilePtr != nil {
		if err != nil {
			fmt.Fprintf(logFilePtr, "Error: %s stmt: %s data: %v elapsed: %s called from: %s\n", err, stmt, data, elapsed, godebug.LF(3))
		} else {
			fmt.Fprintf(logFilePtr, "stmt: %s data: %v elapsed: %s\n", stmt, data, elapsed)
		}
	}
}

// SqliteQueryRow queries a single row and returns that data.
func SqliteQueryRow(stmt string, data ...interface{}) (aRow *sql.Row) {
	start := time.Now()
	aRow = DB.QueryRow(stmt, data...)
	elapsed := time.Since(start)
	logQueries(stmt, nil, data, elapsed)
	return
}

// SqliteExec will run command that returns a resultSet (think insert).
func SqliteExec(stmt string, data ...interface{}) (resultSet sql.Result, err error) {
	start := time.Now()
	resultSet, err = DB.Exec(stmt, data...)
	elapsed := time.Since(start)
	logQueries(stmt, err, data, elapsed)
	return
}

// SqliteQuery runs stmt and returns rows.
func SqliteQuery(stmt string, data ...interface{}) (resultSet *sql.Rows, err error) {
	start := time.Now()
	resultSet, err = DB.Query(stmt, data...)
	elapsed := time.Since(start)
	logQueries(stmt, err, data, elapsed)
	return
}

// SqliteUpdate can run update statements that do not return data.
func SqliteUpdate(stmt string, data ...interface{}) (err error) {
	start := time.Now()
	var resultSet *sql.Rows
	resultSet, err = DB.Query(stmt, data...)
	elapsed := time.Since(start)
	logQueries(stmt, err, data, elapsed)
	if err == nil && resultSet != nil {
		resultSet.Close()
	}
	return
}

// SqliteUInsert can run insert statements that do not return data.
func SqliteInsert(stmt string, data ...interface{}) (err error) {
	start := time.Now()
	var resultSet *sql.Rows
	resultSet, err = DB.Query(stmt, data...)
	elapsed := time.Since(start)
	logQueries(stmt, err, data, elapsed)
	if err == nil && resultSet != nil {
		resultSet.Close()
	}
	return
}

func ConnectToSqlite() {
	var err error
	DB, err = sql.Open("sqlite3", gCfg.DBSqlite)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%sFailed to connect to database: %s%s\n", MiscLib.ColorRed, err, MiscLib.ColorReset)
		os.Exit(1)
	}
}

func CreateTables() {
	var err error
	stmt := "CREATE TABLE IF NOT EXISTS documents (id INTEGER PRIMARY KEY, name TEXT, hash TEXT)"
	prep, err := DB.Prepare(stmt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%sStatment incorrect [%s] : %s%s\n", MiscLib.ColorRed, stmt, err, MiscLib.ColorReset)
		os.Exit(1)
	}
	prep.Exec()
}
