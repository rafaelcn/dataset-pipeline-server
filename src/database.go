package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

const (
	folder   string = "database"
	database string = "data.db"
	table    string = "CREATE TABLE IF NOT EXISTS TB_DATA" +
		"(DATA_FILENAME CHAR(1024), DATA_PK VARCHAR(1024) PRIMARY KEY, " +
		"DATA_SCORE VARCHAR(1024))"
)

/// The struct responsible to encapsulate the sql.DB struct in which
// callers might use to call functions that operate on the database.
type DatabaseHandler struct {
	database *sql.DB
}

// A struct that denotes a row in the database
type DataRow struct {
	filename string
	pk       string
	score    string
}

func InitDatabase() {
	// If the database folder doesn't exist, create one
	if !Exists(folder) {
		log.Println("[+] Creating the database folder")
		os.Mkdir(folder, os.FileMode(0775))
	}

	db := GetHandler()

	log.Printf("[+] Creating the storage table")

	stmt, err := db.database.Prepare(table)
	_, err = stmt.Exec()

	if err != nil {
		log.Fatalf("[!] Some error occurred trying to create the file table %v",
			err.Error())
	}

	stmt.Close()
}

// Gets a database handler, if any argument is provided it expects
// to be the test database
func GetHandler(params ...string) *DatabaseHandler {

	var db *sql.DB = nil
	var err error = nil

	if len(params) != 0 && params[0] == "test" {
		db, err = sql.Open("sqlite3", folder+"/test.db")
	} else {
		db, err = sql.Open("sqlite3", folder+"/"+database)
	}

	if err != nil {
		log.Fatalf("[!] Something went wrong trying to open the database %s",
			err.Error())
	}

	return &DatabaseHandler{db}
}

// Select information from the default database. Additional settings may
// be included using the params argument such as a where clause.
func (db DatabaseHandler) Select(params ...string) (*sql.Rows, error) {

	var query string = "SELECT * FROM TB_DATA "

	if len(params) != 0 {
		query = query + params[0]
	}

	rows, err := db.database.Query(query)

	if err != nil {
		return nil, err
	}

	return rows, nil
}

// Adds an entry to the database
func (db DatabaseHandler) Insert(filename, pk, score string) (sql.Result, error) {
	// A NULL value must be included to make usage of the FILE_PK variable that
	// is an auto increment field automatically in sqlite
	stmt, err := db.database.Prepare("INSERT INTO TB_DATA (DATA_FILENAME, " +
		"DATA_PK, DATA_SCORE) VALUES (?, ?, ?)")

	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	return stmt.Exec(filename, pk, score)
}
