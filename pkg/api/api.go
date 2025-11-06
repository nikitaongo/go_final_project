package api

import (
	"crypto/sha256"
	"encoding/json"
	"gofinalproject/pkg/db"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type user struct {
	Password string `json:"password"`
}

var secretKey []byte
var jwtToken *jwt.Token

// Init runs handlers for api endpoints
func Init() {
	http.HandleFunc("/api/signin", authHandler)
	http.HandleFunc("/api/nextdate", nextDayHandler)
	http.HandleFunc("/api/task", auth(taskHandler))
	http.HandleFunc("/api/tasks", auth(tasksHandler))
	http.HandleFunc("/api/task/done", auth(taskDoneHandler))
}

func authHandler(res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		errJson(res, http.StatusBadRequest, "failed to read body")
		return
	}
	var usr user
	err = json.Unmarshal(body, &usr)
	if err != nil {
		errJson(res, http.StatusBadRequest, "failed to unmarshal request")
		return
	}
	secretPass := os.Getenv("TODO_PASSWORD")
	if secretPass == usr.Password || secretPass == "" {
		jwtToken = jwt.New(jwt.SigningMethodHS256)
		secretKey = secretGen(secretPass)
		token, err := jwtToken.SignedString(secretKey)
		if err != nil {
			errJson(res, http.StatusInternalServerError, "failed to sign token")
			return
		}
		writeJson(res, http.StatusOK, map[string]string{"token": token})
	} else {
		errJson(res, http.StatusUnauthorized, "wrong password")
		return
	}
}

func secretGen(pass string) []byte {
	sumPass := sha256.Sum256([]byte(pass))
	return sumPass[:]
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			var token string
			cookie, err := req.Cookie("token")
			if err == nil {
				token = cookie.Value
			}
			key := secretGen(pass)
			parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
				return key, nil
			})
			if err != nil || !parsedToken.Valid {
				errJson(res, http.StatusUnauthorized, "invalid token")
				return
			}
		}
		next(res, req)
	})
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
		errJson(res, http.StatusInternalServerError, "gettask function error")
		return
	}
	if task.Repeat == "" {
		db.DeleteTask(id)
		writeJson(res, http.StatusOK, struct{}{})
	} else {
		taskDate, err := time.Parse(layout, task.Date)
		if err != nil {
			errJson(res, http.StatusInternalServerError, "can't parse task_Date")
		}

		task.Date, err = NextDate(taskDate, task.Date, task.Repeat)
		if err != nil {
			errJson(res, http.StatusInternalServerError, "nextdate function error")
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
		errJson(res, http.StatusInternalServerError, "gettask function error")
		return
	}
	writeJson(res, http.StatusOK, task)
}

func deleteTaskHandler(res http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	err := db.DeleteTask(id)
	if err != nil {
		errJson(res, http.StatusInternalServerError, "deletetask function error")
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
	writeJson(res, status, map[string]string{"error": msg})
}
