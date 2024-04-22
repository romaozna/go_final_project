package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"main/src/model"
	"main/src/service"
	"net/http"
	"strconv"
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

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validateTask, err := service.ValidateTask(&task)
	if err != nil {
		http.Error(w, fmt.Errorf("error: %w", err).Error(), http.StatusBadRequest)
		return
	}
	taskId, err := service.InsertTask(validateTask)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	resp, _ := json.Marshal(model.Task{
		Id:      strconv.Itoa(taskId),
		Date:    task.Date,
		Title:   task.Title,
		Comment: task.Comment,
		Repeat:  task.Repeat,
	})
	_, err = w.Write(resp)
	log.Println(fmt.Sprintf("Добавлена задача с id=%d", taskId))
	if err != nil {
		http.Error(w, fmt.Errorf("error: %w", err).Error(), http.StatusBadRequest)
	}
}
