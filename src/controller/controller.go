package controller

import (
	"bytes"
	"encoding/json"
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
		log.Println(fmt.Sprintf("Ошибка при записи ответа в функции GetNextDate"))
	}
}

func AddTask(w http.ResponseWriter, r *http.Request) {
	var task model.Task
	var buf bytes.Buffer

	if _, err := buf.ReadFrom(r.Body); err != nil {
		errorResponse(w, "ошибка при получении тела запроса", "500", err)
		return
	}

	if err := json.Unmarshal(buf.Bytes(), &task); err != nil {
		errorResponse(w, "ошибка при десериализации JSON", "500", err)
		return
	}

	validateTask, err := service.ValidateTask(&task)
	if err != nil {
		errorResponse(w, "неверный формат", "400", err)
		return
	}
	taskId, err := service.InsertTask(validateTask)
	if err != nil {
		errorResponse(w, "ошибка при добавлении задачи в БД", "500", err)
		return
	}

	resp, err := json.Marshal(model.TaskIdResponse{Id: taskId})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(resp)
	if err != nil {
		log.Println(fmt.Sprintf("Ошибка при записи ответа в функции AddTask"))
	}
	log.Println(fmt.Sprintf("Добавлена задача с id=%d", taskId))
}

func GetTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := service.GetTasks()
	if err != nil {
		errorResponse(w, "GetTasks: ошибка при получении списка задач", "500", err)
		return
	}

	resp, err := json.Marshal(model.Tasks{Tasks: tasks})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		log.Println(fmt.Sprintf("Ошибка при записи ответа в функции GetTasks"))
	}
	log.Println(fmt.Sprintf("Получено %d задач", len(tasks)))
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	task, err := service.GetTask(id)
	if err != nil {
		errorResponse(w, "GetTask: ошибка при получении задачи", "500", err)
		return
	}
	resp, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		log.Println(fmt.Sprintf("Ошибка при записи ответа в функции GetTask"))
	}
	log.Println(fmt.Sprintf("получена задача с id=%s", id))
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	var task model.Task
	var buf bytes.Buffer

	if _, err := buf.ReadFrom(r.Body); err != nil {
		errorResponse(w, "ошибка при получении тела запроса", "500", err)
		return
	}

	if err := json.Unmarshal(buf.Bytes(), &task); err != nil {
		errorResponse(w, "ошибка при десериализации JSON", "500", err)
		return
	}

	validateTask, err := service.ValidateTask(&task)
	if err != nil {
		errorResponse(w, "неверный формат", "500", err)
		return
	}
	task, err = service.PutTask(*validateTask)
	if err != nil {
		errorResponse(w, "ошибка при добавлении задачи в БД", "500", err)
		return
	}

	resp, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(resp)
	if err != nil {
		log.Println(fmt.Sprintf("Ошибка при записи ответа в функции UpdateTask"))
	}
	log.Println(fmt.Sprintf("Обновлена задача с id=%s", task.Id))
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	err := service.DeleteTask(id)
	if err != nil {
		errorResponse(w, "DeleteTask: ошибка при удалении задачи", "500", err)
		return
	}

	resp, err := json.Marshal(struct{}{})
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		log.Println(fmt.Sprintf("Ошибка при записи ответа в функции DeleteTask"))
	}
	log.Println(fmt.Sprintf("задача с id=%s удалена", id))
}

func MakeTaskAsDone(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	task, err := service.CheckAsDone(id)
	if err != nil {
		errorResponse(w, "ошибка при завершении задачи", "500", err)
		return
	}

	resp, err := json.Marshal(struct{}{})
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		log.Println(fmt.Sprintf("Ошибка при записи ответа в функции MakeTaskAsDone"))
	}
	log.Println(fmt.Sprintf("задача с id=%s отмечена как выполненная", task.Id))
}

func errorResponse(w http.ResponseWriter, errorText string, errorType string, err error) {
	errResponse := model.ErrorResponse{
		Error: fmt.Errorf("%s: %w", errorText, err).Error()}
	errorData, _ := json.Marshal(errResponse)
	switch errorType {
	case "400":
		w.WriteHeader(http.StatusBadRequest)
		break
	case "500":
		w.WriteHeader(http.StatusInternalServerError)
		break
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	_, err = w.Write(errorData)
	if err != nil {
		http.Error(w, fmt.Errorf("error: %w", err).Error(), http.StatusInternalServerError)
	}
}
