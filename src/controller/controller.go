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
		responseWithError(w, "ошибка при получении тела запроса", err)
		return
	}

	if err := json.Unmarshal(buf.Bytes(), &task); err != nil {
		responseWithError(w, "ошибка при десериализации JSON", err)
		return
	}

	validateTask, err := service.ValidateTask(&task)
	if err != nil {
		responseWithError(w, "неверный формат", err)
		return
	}
	taskId, err := service.InsertTask(validateTask)
	if err != nil {
		responseWithError(w, "ошибка при добавлении задачи в БД", err)
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

func responseWithError(w http.ResponseWriter, errorText string, err error) {
	errorResponse := model.ErrorResponse{
		Error: fmt.Errorf("%s: %w", errorText, err).Error()}
	errorData, _ := json.Marshal(errorResponse)
	w.WriteHeader(http.StatusBadRequest)
	_, err = w.Write(errorData)

	if err != nil {
		http.Error(w, fmt.Errorf("error: %w", err).Error(), http.StatusBadRequest)
	}
}

func GetTasks(w http.ResponseWriter, r *http.Request) {
	var tasks []model.Task

	err := errors.New("")
	tasks, err = service.GetTasks()
	if err != nil {
		responseWithError(w, "ошибка при получении списка задач", err)
		return
	}

	resp, err := json.Marshal(model.Tasks{Tasks: tasks})
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	log.Println(fmt.Sprintf("Read %d tasks", len(tasks)))

	if err != nil {
		responseWithError(w, "writing tasks error", err)
	}
}
