package service

import (
	"errors"
	"fmt"
	"log"
	"main/src/model"
	"main/src/store"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const DateFormat string = "20060102"

func CreateDatabase() {
	path := store.GetPath()
	_, err := os.Stat(path)

	if err != nil {
		_, err := os.Create(path)
		if err != nil {
			log.Fatal(err)
		}
	}

	store.CreateTable(path)
}

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if len(repeat) == 0 {
		return "", errors.New("в колонке repeat пусто")
	}

	dayMatch, _ := regexp.MatchString(`d \d{1,3}`, repeat)
	yearMatch, _ := regexp.MatchString(`y`, repeat)

	if dayMatch {
		days, err := strconv.Atoi(strings.TrimPrefix(repeat, "d "))
		if err != nil {
			return "", err
		}

		if days > 400 {
			return "", errors.New("максимально допустимое число равно 400")
		}

		parsedDate, err := time.Parse(DateFormat, date)
		if err != nil {
			return "", err
		}

		newDate := parsedDate.AddDate(0, 0, days)

		for newDate.Before(now) {
			newDate = newDate.AddDate(0, 0, days)
		}

		return newDate.Format(DateFormat), nil
	} else if yearMatch {
		parsedDate, err := time.Parse(DateFormat, date)
		if err != nil {
			return "", err
		}

		newDate := parsedDate.AddDate(1, 0, 0)

		for newDate.Before(now) {
			newDate = newDate.AddDate(1, 0, 0)
		}

		return newDate.Format(DateFormat), nil
	}

	return "", errors.New("время в переменной date не может быть преобразовано в корректную дату")
}

func ValidateTask(task *model.Task) (*model.Task, error) {
	if task.Id != "" {
		_, err := strconv.Atoi(task.Id)
		if err != nil {
			return task, fmt.Errorf("неверный формат id: %w", err)
		}
	}

	if task.Date == "" {
		task.Date = time.Now().Format(DateFormat)
	} else {
		date, err := time.Parse(DateFormat, task.Date)
		if err != nil {
			return task, fmt.Errorf("неверный формат даты: %w", err)
		}

		if date.Before(time.Now()) {
			task.Date = time.Now().Format(DateFormat)
		}
	}

	if task.Title == "" {
		return task, errors.New("заголовок задачи не может быть пустым")
	}

	if len(task.Repeat) > 0 {
		if _, err := NextDate(time.Now(), task.Date, task.Repeat); err != nil {
			return task, errors.New("неверный формат даты в repeat")
		}
	}
	return task, nil
}

func InsertTask(task *model.Task) (int, error) {
	taskId, err := store.InsertTask(task)
	if err != nil {
		return 0, fmt.Errorf("ошибка при добавлении задачи: %w", err)
	}
	return taskId, err
}

func GetTasks() ([]model.Task, error) {
	tasks, err := store.GetAllTasks()
	if err != nil {
		return tasks, err
	}

	if tasks == nil {
		tasks = []model.Task{}
	}

	return tasks, err
}

func GetTask(id string) (model.Task, error) {
	task, err := store.GetTaskById(id)
	if err != nil {
		return model.Task{}, err
	}
	return task, err
}

func PutTask(task model.Task) (model.Task, error) {
	updatedTask, err := store.UpdateTask(task)
	if err != nil {
		return model.Task{}, fmt.Errorf("ошибка при обновлении задачи в БД: %w", err)
	}
	return updatedTask, err
}

func DeleteTask(id string) error {
	err := store.DeleteTaskById(id)
	if err != nil {
		return fmt.Errorf("ошибка при удалении задачи из БД: %w", err)

	}
	return err
}

func CheckAsDone(id string) (model.Task, error) {
	task, err := store.GetTaskById(id)
	if err != nil {
		return model.Task{}, fmt.Errorf("ошибка при получении задачи из БД: %w", err)
	}

	if task.Repeat == "" {
		err = store.DeleteTaskById(task.Id)
		if err != nil {
			return model.Task{}, fmt.Errorf("ошибка удаления задачи из БД: %w", err)
		}
		return task, nil
	}

	task.Date, err = NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		return model.Task{}, fmt.Errorf("ошибка при получении следующей даты: %w", err)
	}

	task, err = store.UpdateTask(task)
	if err != nil {
		return model.Task{}, fmt.Errorf("ошибка при обновлении задачи в БД: %w", err)
	}
	return task, err
}
