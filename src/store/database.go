package store

import (
	"database/sql"
	"log"
	"main/src/model"
	"os"
	"path/filepath"
)

var db *sql.DB

func GetPath() string {
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

func openDatabase(path string) *sql.DB {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil
	}

	return db
}

func CreateTable(path string) {
	db = openDatabase(path)
	_, err := db.Exec(
		"CREATE TABLE IF NOT EXISTS scheduler (id INTEGER PRIMARY KEY AUTOINCREMENT, date VARCHAR(8) NULL, title VARCHAR(64) NOT NULL, comment VARCHAR(255) NULL, repeat VARCHAR(128) NULL)")
	if err != nil {
		log.Fatal(err)
	}
	createIndex()
}

func createIndex() {
	_, err := db.Exec(
		"CREATE INDEX IF NOT EXISTS date_idx ON scheduler (date)")
	if err != nil {
		log.Fatal(err)
	}
}

func InsertTask(task *model.Task) (int, error) {
	result, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func GetAllTasks() ([]model.Task, error) {
	var tasks []model.Task

	rows, err := db.Query("SELECT * FROM scheduler ORDER BY date LIMIT '10'")
	if err != nil {
		return []model.Task{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var task model.Task
		if err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return []model.Task{}, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return []model.Task{}, err
	}

	if tasks == nil {
		tasks = []model.Task{}
	}

	return tasks, nil
}
