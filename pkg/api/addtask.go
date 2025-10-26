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
	var response any
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		errJson(res, http.StatusInternalServerError, "can't read bytes from body")
		return
	}

	err = json.Unmarshal(buf.Bytes(), &task)
	if err != nil {
		errJson(res, http.StatusInternalServerError, "can't unmarshal request")
		return
	}
	//check logic - does spaces at title field is valid?
	if strings.TrimSpace(task.Title) == "" {
		errJson(res, http.StatusInternalServerError, "task title can't be empty")
		return
	}

	if strings.TrimSpace(task.Date) == "" {
		task.Date = time.Now().Format(layout)
	}

	parsedDate, err := time.Parse(layout, task.Date)
	if err != nil {
		errJson(res, http.StatusInternalServerError, "task date parsing failed")
		return
	}

	if parsedDate.Before(time.Now()) {
		if strings.TrimSpace(task.Repeat) == "" {
			task.Date = time.Now().Format(layout)
		} else {
			task.Date, err = NextDate(time.Now(), parsedDate.Format(layout), task.Repeat)
			if err != nil {
				errJson(res, http.StatusInternalServerError, "wrong task repeat pattern")
				return
			}
		}
	}
	id, err := db.AddTask(&task)
	response = map[string]any{"id": id}

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJson(res, http.StatusCreated, response)
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
