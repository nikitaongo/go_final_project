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

func addTaskHandler(res http.ResponseWriter, req *http.Request, isNew bool) {
	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		errJson(res, http.StatusBadRequest, "failed to read body")
		return
	}

	var task db.Task
	err = json.Unmarshal(body, &task)
	if err != nil {
		errJson(res, http.StatusBadRequest, "failed to unmarshal request")
		return
	}

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
	if isAfter(time.Now(), parsedDate) {
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
	if isNew {
		id, err := db.AddTask(&task)
		if err != nil {
			errJson(res, http.StatusInternalServerError, "failed to add task")
			return
		}
		res.Header().Set("Location", fmt.Sprintf("/api/tasks/%v", id))
		writeJson(res, http.StatusCreated, map[string]any{"id": id})
	} else {
		err = db.UpdateTask(&task)
		if err != nil {
			errJson(res, http.StatusInternalServerError, "failed to update task")
			return
		}
		writeJson(res, http.StatusOK, map[string]any{})
	}
}
