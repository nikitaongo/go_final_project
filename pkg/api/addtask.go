package api

import (
	"encoding/json"
	"fmt"
	"gofinalproject/pkg/db"
	"io"
	"net/http"
	"strings"
	"time"
)

func parseJsonBody(req *http.Request, task any) error {
	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return fmt.Errorf("failed to read body: %v", err)
	}
	if err := json.Unmarshal(body, task); err != nil {
		return fmt.Errorf("failed to unmarshal request: %v", err)
	}
	return nil
}

func addTaskHandler(res http.ResponseWriter, req *http.Request) {
	var task db.Task

	if err := parseJsonBody(req, &task); err != nil {
		errJson(res, http.StatusBadRequest, err.Error())
		return
	}

	if strings.TrimSpace(task.Title) == "" {
		errJson(res, http.StatusInternalServerError, "task title can't be empty")
		return
	}

	if strings.TrimSpace(task.Date) == "" {
		task.Date = time.Now().Format(db.Layout)
	}

	parsedDate, err := time.Parse(db.Layout, task.Date)
	if err != nil {
		errJson(res, http.StatusInternalServerError, "task date parsing failed")
		return
	}
	if isAfter(time.Now(), parsedDate) {
		if strings.TrimSpace(task.Repeat) == "" {
			task.Date = time.Now().Format(db.Layout)
		} else {
			task.Date, err = NextDate(time.Now(), parsedDate.Format(db.Layout), task.Repeat)
			if err != nil {
				errJson(res, http.StatusInternalServerError, "wrong task repeat pattern")
				return
			}
		}
	}
	id, err := db.AddTask(&task)
	if err != nil {
		errJson(res, http.StatusInternalServerError, "failed to add task")
		return
	}
	res.Header().Set("Location", fmt.Sprintf("/api/tasks/%v", id))
	writeJson(res, http.StatusCreated, map[string]any{"id": id})
}
