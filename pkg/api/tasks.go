package api

import (
	"gofinalproject/pkg/db"
	"net/http"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(res http.ResponseWriter, req *http.Request) {
	searchLn := req.FormValue("search")
	switch {
	case searchLn != "":
		tasks, err := db.Searcher(searchLn, db.MaxTasks)
		if err != nil {
			errJson(res, http.StatusInternalServerError, "searchline function error")
			return
		}
		writeJson(res, http.StatusOK, TasksResp{
			Tasks: tasks,
		})
	default:
		tasks, err := db.Tasks(db.MaxTasks)
		if err != nil {
			errJson(res, http.StatusInternalServerError, "tasks function error")
			return
		}
		writeJson(res, http.StatusOK, TasksResp{
			Tasks: tasks,
		})
	}

}
