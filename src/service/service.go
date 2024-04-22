package service

import (
	"errors"
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
	if len(task.Date) == 0 {
		task.Date = time.Now().Format(DateFormat)
	} else {
		date, err := time.Parse(DateFormat, task.Date)
		if err != nil {
			return task, errors.New("неверный формат даты")
		}

		if date.Before(time.Now()) {
			task.Date = time.Now().Format(DateFormat)
		}
	}

	if len(task.Title) == 0 {
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
		return 0, errors.New("ошибка при добавлении задачи в БД")
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
		return model.Task{}, errors.New("ошибка при обновлении задачи в БД")
	}
	return updatedTask, err
}

func DeleteTask(id string) error {
	err := store.DeleteTaskById(id)
	if err != nil {
		return errors.New("ошибка при обновлении задачи в БД")

	}
	return err
}

func CheckAsDone(id string) (model.Task, error) {
	task, err := store.GetTaskById(id)
	if err != nil {
		return model.Task{}, errors.New("ошибка при получении задачи из БД")
	}

	if len(task.Repeat) == 0 {
		err = store.DeleteTaskById(task.Id)
		if err != nil {
			return model.Task{}, errors.New("ошибка удаления задачи из БД")
		}
	} else {
		task.Date, err = NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			return model.Task{}, errors.New("ошибка при получении следующей даты")
		}

		_, err = store.UpdateTask(task)
		if err != nil {
			return model.Task{}, errors.New("ошибка при обновлении задачи в БД")
		}
	}
	return task, err
}
