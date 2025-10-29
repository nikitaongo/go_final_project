package api

import (
	"gofinalproject/pkg/db"
	"net/http"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(res http.ResponseWriter, req *http.Request) {
	tasks, err := db.Tasks(db.MaxTasks) // в параметре максимальное количество записей
	if err != nil {
		errJson(res, http.StatusInternalServerError, "tasks function error") // проверить формулировку
		return
	}
	writeJson(res, http.StatusOK, TasksResp{
		Tasks: tasks,
	})
}
