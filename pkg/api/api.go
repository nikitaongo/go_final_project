package api

import (
	"encoding/json"
	"fmt"
	"gofinalproject/pkg/db"
	"net/http"
)

// Init runs nextDayHandler for /api/nextdate endpoint
func Init() {
	http.HandleFunc("/api/nextdate", nextDayHandler)
	http.HandleFunc("/api/task", taskHandler)
	http.HandleFunc("/api/tasks", tasksHandler)
}

func taskHandler(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		createTaskHandler(res, req)
	case http.MethodGet:
		getTaskHandler(res, req)
	case http.MethodPut:
		editTaskHandler(res, req)
	}
}

func createTaskHandler(res http.ResponseWriter, req *http.Request) {
	isNew := true
	addTaskHandler(res, req, isNew)
}

func editTaskHandler(res http.ResponseWriter, req *http.Request) {
	isNew := false
	addTaskHandler(res, req, isNew)
}

func getTaskHandler(res http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	task, err := db.GetTask(id)
	if err != nil {
		fmt.Println("gettask function error!")
		errJson(res, http.StatusInternalServerError, "gettask function error") // формулировка
		return
	}
	writeJson(res, http.StatusOK, task)
}

func writeJson(res http.ResponseWriter, status int, data any) {
	resp, err := json.Marshal(data)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	res.WriteHeader(status)
	res.Write(resp)
}

func errJson(res http.ResponseWriter, status int, msg string) {
	writeJson(res, status, map[string]string{"error": msg})
}
