package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postman",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

func getAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	tasksList := make([]Task, 0, len(tasks))
	for _, task := range tasks {
		tasksList = append(tasksList, task)
	}

	resp, err := json.Marshal(tasksList)
	if err != nil {
		http.Error(w, "Ошибка кодировки", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func postTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var newTask Task

	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
		http.Error(w, "Тело запроса некорректно", http.StatusBadRequest)
		return
	}

	if _, exists := tasks[newTask.ID]; exists {
		http.Error(w, "Задача с таким ID уже существует", http.StatusBadRequest)
		return
	}

	tasks[newTask.ID] = newTask
	w.WriteHeader(http.StatusCreated)
}

func getTaskById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	id := chi.URLParam(r, "id")

	task, exists := tasks[id]
	if !exists {
		http.Error(w, "Нет задачи", http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(task)
	if err != nil {
		http.Error(w, "Ошибка кодирования в JSON", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	id := chi.URLParam(r, "id")

	if _, exists := tasks[id]; !exists {
		http.Error(w, "Нет задачи", http.StatusBadRequest)
		return
	}

	delete(tasks, id)
	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	r.Get("/tasks", getAllTasksHandler)
	r.Post("/tasks", postTaskHandler)
	r.Get("/tasks/{id}", getTaskById)
	r.Delete("/tasks/{id}", deleteTaskHandler)

	if err := http.ListenAndServe(":8086", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
