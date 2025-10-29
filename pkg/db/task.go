package db

import (
	"fmt"

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
		fmt.Println("tasks rows error")
		return nil, fmt.Errorf("tasks rows error: %w", err)
	}

	return tasks, nil
}
