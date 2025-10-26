package api

import (
	"encoding/json"
	"net/http"
)

// Init runs nextDayHandler for /api/nextdate endpoint
func Init() {
	http.HandleFunc("/api/nextdate", nextDayHandler)
	http.HandleFunc("/api/task", taskHandler)
}

func taskHandler(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		addTaskHandler(res, req)
	case http.MethodGet:
		getTaskHandler(res, req)
	case http.MethodDelete:
		deleteTaskHandler(res, req)
	}
}

func writeJson(res http.ResponseWriter, data any) {
	resp, err := json.Marshal(data)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	res.WriteHeader(http.StatusOK)
	res.Write(resp)
}
