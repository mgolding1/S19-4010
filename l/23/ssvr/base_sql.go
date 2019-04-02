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
	stmt, data, _ = BindFixer(stmt, data)
	aRow = DB.QueryRow(stmt, data...)
	elapsed := time.Since(start)
	logQueries(stmt, nil, data, elapsed)
	return
}

// SqliteExec will run command that returns a resultSet (think insert).
func SqliteExec(stmt string, data ...interface{}) (resultSet sql.Result, err error) {
	start := time.Now()
	stmt, data, _ = BindFixer(stmt, data)
	resultSet, err = DB.Exec(stmt, data...)
	elapsed := time.Since(start)
	logQueries(stmt, err, data, elapsed)
	return
}

// SqliteQuery runs stmt and returns rows.
func SqliteQuery(stmt string, data ...interface{}) (resultSet *sql.Rows, err error) {
	start := time.Now()
	stmt, data, _ = BindFixer(stmt, data)
	resultSet, err = DB.Query(stmt, data...)
	elapsed := time.Since(start)
	logQueries(stmt, err, data, elapsed)
	return
}

// SqliteUpdate can run update statements that do not return data.
func SqliteUpdate(stmt string, data ...interface{}) (err error) {
	start := time.Now()
	// var resultSet *sql.Rows
	// resultSet, err = DB.Query(stmt, data...)
	fmt.Printf("SqliteUpdate: AT: %s\n", godebug.LF())
	sstmt, ddata, _ := BindFixer(stmt, data)
	fmt.Printf("SqliteUpdate: AT: %s stmtNew[%s] data %s\n", godebug.LF(), sstmt, godebug.SVar(ddata))
	statement, e0 := DB.Prepare(sstmt)
	err = e0
	if err == nil {
		fmt.Printf("SqliteUpdate: AT: %s\n", godebug.LF())
		statement.Exec(ddata...)
		statement.Close()
	}
	fmt.Printf("SqliteUpdate: AT: %s\n", godebug.LF())
	elapsed := time.Since(start)
	fmt.Printf("SqliteUpdate: AT: %s\n", godebug.LF())
	logQueries(sstmt, err, ddata, elapsed)
	// if err == nil && resultSet != nil {
	// 	resultSet.Close()
	// }
	return
}

// SqliteUInsert can run insert statements that do not return data.
func SqliteInsert(stmt string, data ...interface{}) (err error) {
	start := time.Now()
	// var resultSet *sql.Rows
	// resultSet, err = DB.Query(stmt, data...)
	sstmt, ddata, _ := BindFixer(stmt, data)
	fmt.Printf("%sSqliteInsert: AT: %s stmtNew[%s] data %s%s\n", MiscLib.ColorCyan, godebug.LF(), sstmt, godebug.SVar(ddata), MiscLib.ColorReset)
	statement, e0 := DB.Prepare(sstmt)
	err = e0
	if err == nil {
		statement.Exec(ddata...)
		statement.Close()
	}
	elapsed := time.Since(start)
	logQueries(sstmt, err, ddata, elapsed)
	//if err == nil && resultSet != nil {
	//	resultSet.Close()
	//}
	return
}

func ConnectToSqlite() {
	DbType = "SQLite"
	var err error
	DB, err = sql.Open("sqlite3", gCfg.DBSqlite)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%sFailed to connect to database: %s%s\n", MiscLib.ColorRed, err, MiscLib.ColorReset)
		os.Exit(1)
	}
}

func CreateTables() {
	var err error

	stmt := `CREATE TABLE IF NOT EXISTS documents (
		id 					text PRIMARY KEY,
		document_hash 		TEXT,
		email 				TEXT,
		real_name 			TEXT,
		phone_number 		TEXT,
		address_usps 		TEXT,
		document_file_name	TEXT,
		file_name			TEXT,
		orig_file_name		TEXT,
		url_file_name		TEXT,
		txid				TEXT,
		note 				TEXT,
 		updated 			DATETIME,
 		created 			DATETIME DEFAULT CURRENT_TIMESTAMP
	)`
	prep, err := DB.Prepare(stmt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%sStatment incorrect [%s] : %s%s\n", MiscLib.ColorRed, stmt, err, MiscLib.ColorReset)
		os.Exit(1)
	}
	prep.Exec()

	// Note: password should be encrypted or - better yet - use SRP6a and never have the server see the
	// password at all.
	stmt = `CREATE TABLE IF NOT EXISTS users (
		id 					text PRIMARY KEY,
		username 			TEXT,
		password 			TEXT,
 		updated 			DATETIME,
 		created 			DATETIME DEFAULT CURRENT_TIMESTAMP
	)`
	prep, err = DB.Prepare(stmt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%sStatment incorrect [%s] : %s%s\n", MiscLib.ColorRed, stmt, err, MiscLib.ColorReset)
		os.Exit(1)
	}
	prep.Exec()

}
