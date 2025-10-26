package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gofinalproject/pkg/db"
	"net/http"
	"strings"
	"time"
)

func addTaskHandler(res http.ResponseWriter, req *http.Request) {
	var task db.Task
	var buf bytes.Buffer
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(buf.Bytes(), &task)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	//check logic - does spaces at title field is valid?
	if strings.TrimSpace(task.Title) == "" {
		http.Error(res, "missing required field: task", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(task.Date) == "" {
		task.Date = time.Now().Format(layout)
	}

	parsedDate, err := time.Parse(layout, task.Date)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	if parsedDate.Before(time.Now()) {
		if strings.TrimSpace(task.Repeat) == "" {
			task.Date = time.Now().Format(layout)
		} else {
			task.Date, err = NextDate(time.Now(), parsedDate.Format(layout), task.Repeat)
			if err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
	fmt.Println(task)
	id, err := db.AddTask(&task)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJson(res, id)
}
func getTaskHandler(res http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")

	if id == "" {
		http.Error(res, fmt.Sprintf("wrong id: %q", id), http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("getTaskHandler() is under construction"))
}

func deleteTaskHandler(res http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")

	if id == "" {
		http.Error(res, fmt.Sprintf("wrong id: %q", id), http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("deleteTaskHandler() is under construction"))
}
