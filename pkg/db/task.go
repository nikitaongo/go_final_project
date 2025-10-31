package db

import (
	"fmt"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	var id int64
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES ($1, $2, $3, $4)`
	res, err := Db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err == nil {
		id, err = res.LastInsertId()
	} else {
		return 0, fmt.Errorf("addtask query error: %w", err)
	}

	return id, err
}

func Tasks(limit int) ([]*Task, error) {
	tasks := make([]*Task, 0, limit)
	query := `SELECT * FROM scheduler ORDER BY date DESC LIMIT $1`
	rows, err := Db.Query(query, limit)
	if err != nil {
		return []*Task{}, fmt.Errorf("tasks query error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		task := &Task{}
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return []*Task{}, fmt.Errorf("tasks scan error: %w", err)
		}

		tasks = append(tasks, task)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("tasks rows error: %w", err)
	}

	return tasks, nil
}

func Searcher(searchLn string, limit int) ([]*Task, error) {
	var err error
	tasks := make([]*Task, 0, limit)
	query := `SELECT * FROM scheduler WHERE title LIKE $1 OR comment LIKE $1 ORDER BY date LIMIT $2`
	if isDate(searchLn) {
		query = `SELECT * FROM scheduler WHERE date LIKE $1 LIMIT $2`
		searchLn, err = formatDate(searchLn)
		if err != nil {
			return []*Task{}, fmt.Errorf("tasks SearchLine query error: %w", err)
		}
	}
	rows, err := Db.Query(query, "%"+searchLn+"%", limit)
	if err != nil {
		return []*Task{}, fmt.Errorf("tasks SearchLine query error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		task := &Task{}
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return []*Task{}, fmt.Errorf("tasks SearchLine scan error: %w", err)
		}

		tasks = append(tasks, task)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("tasks SearchLine rows error: %w", err)
	}

	return tasks, nil
}

func isDate(searchLn string) bool {
	if len(searchLn) != 10 {
		return false
	}
	for i, r := range searchLn {
		switch i {
		case 2, 5:
			if r != '.' && r != '-' && r != '/' && r != ' ' && r != '_' {
				return false
			}
		default:
			if r < '0' || r > '9' {
				return false
			}
		}
	}
	return true
}

func formatDate(searchLn string) (string, error) {
	replacer := strings.NewReplacer("-", ".", "/", ".", " ", ".", "_", ".")
	searchLn = replacer.Replace(searchLn)
	parsed, err := time.Parse("02.01.2006", searchLn)
	if err != nil {
		return "", fmt.Errorf("%e:\ncan't parse date at search parameter - wrong input", err)
	}
	return parsed.Format("20060102"), nil //исп layout
}
