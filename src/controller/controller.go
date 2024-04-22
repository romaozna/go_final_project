package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"main/src/model"
	"main/src/service"
	"net/http"
	"time"
)

func GetNextDate(w http.ResponseWriter, r *http.Request) {
	now, err := time.Parse(service.DateFormat, r.FormValue("now"))
	if err != nil {
		http.Error(w, fmt.Sprintf(""), http.StatusBadRequest)
		return
	}

	date := r.FormValue("date")
	repeat := r.FormValue("repeat")
	nextDate, err := service.NextDate(now, date, repeat)

	if err != nil {
		http.Error(w, fmt.Sprintf(""), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(nextDate))

	if err != nil {
		http.Error(w, fmt.Errorf("ошибка при запросе следующей даты: %w", err).Error(), http.StatusBadRequest)
	}
}

func AddTask(w http.ResponseWriter, r *http.Request) {
	var task model.Task
	var buf bytes.Buffer

	if _, err := buf.ReadFrom(r.Body); err != nil {
		errorResponse(w, "ошибка при получении тела запроса", err)
		return
	}

	if err := json.Unmarshal(buf.Bytes(), &task); err != nil {
		errorResponse(w, "ошибка при десериализации JSON", err)
		return
	}

	validateTask, err := service.ValidateTask(&task)
	if err != nil {
		errorResponse(w, "неверный формат", err)
		return
	}
	taskId, err := service.InsertTask(validateTask)
	if err != nil {
		errorResponse(w, "ошибка при добавлении задачи в БД", err)
		return
	}

	resp, err := json.Marshal(model.TaskIdResponse{Id: taskId})
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(resp)
	log.Println(fmt.Sprintf("Добавлена задача с id=%d", taskId))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func GetTasks(w http.ResponseWriter, r *http.Request) {
	var tasks []model.Task

	err := errors.New("")
	tasks, err = service.GetTasks()
	if err != nil {
		errorResponse(w, "GetTasks: ошибка при получении списка задач", err)
		return
	}

	resp, err := json.Marshal(model.Tasks{Tasks: tasks})
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	log.Println(fmt.Sprintf("Получено %d задач", len(tasks)))

	if err != nil {
		errorResponse(w, "GetTasks: ошибка при записи ответа", err)
	}
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	task, err := service.GetTask(id)
	if err != nil {
		errorResponse(w, "GetTask: ошибка при получении задачи", err)
		return
	}

	resp, err := json.Marshal(task)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	log.Println(fmt.Sprintf("ошибка получения задачи с id=%s", id))

	if err != nil {
		errorResponse(w, "GetTask: ошибка при записи ответа", err)
	}
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	var task model.Task
	var buf bytes.Buffer

	if _, err := buf.ReadFrom(r.Body); err != nil {
		errorResponse(w, "ошибка при получении тела запроса", err)
		return
	}

	if err := json.Unmarshal(buf.Bytes(), &task); err != nil {
		errorResponse(w, "ошибка при десериализации JSON", err)
		return
	}

	validateTask, err := service.ValidateTask(&task)
	if err != nil {
		errorResponse(w, "неверный формат", err)
		return
	}
	task, err = service.PutTask(*validateTask)
	if err != nil {
		errorResponse(w, "ошибка при добавлении задачи в БД", err)
		return
	}

	resp, err := json.Marshal(task)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(resp)
	log.Println(fmt.Sprintf("Обновлена задача с id=%s", task.Id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	err := service.DeleteTask(id)
	if err != nil {
		errorResponse(w, "DeleteTask: ошибка при удалении задачи", err)
		return
	}

	resp, err := json.Marshal(struct{}{})
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	log.Println(fmt.Sprintf("ошибка удаления задачи с id=%s", id))

	if err != nil {
		errorResponse(w, "DeleteTask: ошибка при записи ответа", err)
	}
}

func MakeTaskAsDone(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	task, err := service.CheckAsDone(id)
	if err != nil {
		errorResponse(w, "ошибка при завершении задачи", err)
		return
	}

	resp, err := json.Marshal(struct{}{})
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	log.Println(fmt.Sprintf("задача с id=%s отмечена как выполненная", task.Id))

	if err != nil {
		errorResponse(w, "MakeTaskAsDone: ошибка при записи ответа", err)
	}
}

func errorResponse(w http.ResponseWriter, errorText string, err error) {
	errorResponse := model.ErrorResponse{
		Error: fmt.Errorf("%s: %w", errorText, err).Error()}
	errorData, _ := json.Marshal(errorResponse)
	w.WriteHeader(http.StatusBadRequest)
	_, err = w.Write(errorData)

	if err != nil {
		http.Error(w, fmt.Errorf("error: %w", err).Error(), http.StatusBadRequest)
	}
}
