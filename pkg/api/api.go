package api

import (
	"encoding/json"
	"gofinalproject/pkg/db"
	"net/http"
	"time"
)

// Init runs nextDayHandler for /api/nextdate endpoint
func Init() {
	http.HandleFunc("/api/nextdate", nextDayHandler)
	http.HandleFunc("/api/task", taskHandler)
	http.HandleFunc("/api/tasks", tasksHandler)
	http.HandleFunc("/api/task/done", taskDoneHandler)
}

func taskHandler(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		createTaskHandler(res, req)
	case http.MethodGet:
		getTaskHandler(res, req)
	case http.MethodPut:
		editTaskHandler(res, req)
	case http.MethodDelete:
		deleteTaskHandler(res, req)
	}
}
func taskDoneHandler(res http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	task, err := db.GetTask(id)
	if err != nil {
		errJson(res, http.StatusInternalServerError, "gettask function error") // формулировка
		return
	}
	if task.Repeat == "" {
		db.DeleteTask(id)
		writeJson(res, http.StatusOK, struct{}{})
	} else {
		taskDate, err := time.Parse(layout, task.Date)
		if err != nil {
			errJson(res, http.StatusInternalServerError, "can't parse task_Date") // формулировка
		}

		task.Date, err = NextDate(taskDate.AddDate(0, 0, 1), task.Date, task.Repeat)
		if err != nil {
			errJson(res, http.StatusInternalServerError, "nextdate function error") // формулировка
			return
		}
		db.UpdateDate(task)
		writeJson(res, http.StatusOK, struct{}{})
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
		errJson(res, http.StatusInternalServerError, "gettask function error") // формулировка
		return
	}
	writeJson(res, http.StatusOK, task)
}

func deleteTaskHandler(res http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	err := db.DeleteTask(id)
	if err != nil {
		errJson(res, http.StatusInternalServerError, "deletetask function error") // формулировка
		return
	}
	writeJson(res, http.StatusOK, struct{}{})
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
	writeJson(res, status, map[string]string{"error": msg}) //мб type APIResponse struct {} ???
}
