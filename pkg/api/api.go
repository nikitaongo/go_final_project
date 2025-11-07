package api

import (
	"crypto/sha256"
	"encoding/json"
	"gofinalproject/pkg/db"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type user struct {
	Password string `json:"password"`
}

var secretKey []byte
var jwtToken *jwt.Token

// Init runs handlers for api endpoints
func Init(password string) {
	http.HandleFunc("/api/signin", authHandler(password))
	http.HandleFunc("/api/nextdate", nextDayHandler)
	http.HandleFunc("/api/task", auth(taskHandler, password))
	http.HandleFunc("/api/tasks", auth(tasksHandler, password))
	http.HandleFunc("/api/task/done", auth(taskDoneHandler, password))
}

func authHandler(password string) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
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
		if password == usr.Password || password == "" {
			jwtToken = jwt.New(jwt.SigningMethodHS256)
			secretKey = secretGen(password)
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
}

func secretGen(pass string) []byte {
	sumPass := sha256.Sum256([]byte(pass))
	return sumPass[:]
}

func auth(next http.HandlerFunc, password string) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if len(password) > 0 {
			var token string
			cookie, err := req.Cookie("token")
			if err == nil {
				token = cookie.Value
			}
			key := secretGen(password)
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
		addTaskHandler(res, req)
	case http.MethodGet:
		getTaskHandler(res, req)
	case http.MethodPut:
		editTaskHandler(res, req)
	case http.MethodDelete:
		deleteTaskHandler(res, req)
	default:
		res.Header().Set("Allow", "POST, GET, PUT, DELETE")
		http.Error(res, "HTTP method not allowed", http.StatusMethodNotAllowed)
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
		taskDate, err := time.Parse(db.Layout, task.Date)
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

func getTaskHandler(res http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	task, err := db.GetTask(id)
	if err != nil {
		errJson(res, http.StatusInternalServerError, "gettask function error")
		return
	}
	writeJson(res, http.StatusOK, task)
}

func editTaskHandler(res http.ResponseWriter, req *http.Request) {
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
	err = db.UpdateTask(&task)
	if err != nil {
		errJson(res, http.StatusInternalServerError, "failed to update task")
		return
	}
	writeJson(res, http.StatusOK, map[string]any{})

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
