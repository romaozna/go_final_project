package store

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
)

var db *sql.DB

func getPath() string {
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")
	envFile := os.Getenv("TODO_DBFILE")
	if len(envFile) > 0 {
		dbFile = envFile
	}

	return dbFile
}

func openDatabase(pathToDatabase string) *sql.DB {
	db, err := sql.Open("sqlite", pathToDatabase)
	if err != nil {
		return nil
	}

	return db
}

func createTable(db *sql.DB) {
	_, err := db.Exec(
		"CREATE TABLE IF NOT EXISTS scheduler (id INTEGER PRIMARY KEY AUTOINCREMENT, date VARCHAR(8) NULL, title VARCHAR(64) NOT NULL, comment VARCHAR(255) NULL, repeat VARCHAR(128) NULL)")
	if err != nil {
		log.Fatal(err)
	}
}

func createIndex(db *sql.DB) {
	_, err := db.Exec(
		"CREATE INDEX IF NOT EXISTS date_idx ON scheduler (date)")
	if err != nil {
		log.Fatal(err)
	}
}

func CreateDatabase() {
	path := getPath()
	_, err := os.Stat(path)

	if err != nil {
		_, err := os.Create(path)
		if err != nil {
			log.Fatal(err)
		}
	}

	db = openDatabase(path)
	createTable(db)
	createIndex(db)
}
