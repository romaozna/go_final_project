package store

import (
	"database/sql"
	"errors"
	"log"
	"main/src/model"
	"os"
)

var db *sql.DB

func GetPath() string {
	dbFile := "scheduler.db"
	envFile := os.Getenv("TODO_DBFILE")
	if len(envFile) > 0 {
		dbFile = envFile
	}
	return dbFile
}

func CreateTable(path string) {
	db, _ = openDatabase(path)
	_, err := db.Exec(
		"CREATE TABLE IF NOT EXISTS scheduler (id INTEGER PRIMARY KEY AUTOINCREMENT, date VARCHAR(8) NULL, title VARCHAR(64) NOT NULL, comment VARCHAR(255) NULL, repeat VARCHAR(128) NULL)")
	if err != nil {
		log.Fatal(err)
	}
	createIndex()
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
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task model.Task
		if err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func GetTaskById(id string) (model.Task, error) {
	var task model.Task

	row := db.QueryRow("SELECT * FROM scheduler WHERE id = :id",
		sql.Named("id", id))
	if err := row.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
		return model.Task{}, err
	}
	return task, nil
}

func UpdateTask(task model.Task) (model.Task, error) {
	result, err := db.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.Id))
	if err != nil {
		return model.Task{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return model.Task{}, err
	}

	if rowsAffected == 0 {
		return model.Task{}, errors.New("0 строчек было обновлено")
	}
	updatedTask, _ := GetTaskById(task.Id)

	return updatedTask, nil
}

func DeleteTaskById(id string) error {
	result, err := db.Exec("DELETE FROM scheduler WHERE id = :id",
		sql.Named("id", id))
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("0 строчек было удалено")
	}

	return err
}

func openDatabase(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createIndex() {
	_, err := db.Exec(
		"CREATE INDEX IF NOT EXISTS date_idx ON scheduler (date)")
	if err != nil {
		log.Fatal(err)
	}
}
