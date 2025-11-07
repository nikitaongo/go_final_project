package db

import (
	"fmt"
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
		return 0, fmt.Errorf("addtask exec error: %w", err)
	}

	return id, err
}

func GetTask(id string) (*Task, error) {
	var task Task
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = $1`
	err := Db.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, fmt.Errorf("gettask queryrow error: %w", err)
	}
	return &task, nil
}

func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = $1`
	res, err := Db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("deletetask exec error: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("deletetask rowsaffected error: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("no rows deleted: %w \nid: %s", err, id)
	}
	return nil
}

func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date = $2, title = $3, comment = $4, repeat = $5 WHERE id = $1`
	res, err := Db.Exec(query, task.ID, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil
}

func UpdateDate(task *Task) error {
	query := `UPDATE scheduler SET date = $2 WHERE id = $1`
	res, err := Db.Exec(query, task.ID, task.Date)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task date`)
	}
	return nil
}

func Tasks(limit int) ([]*Task, error) {
	tasks := make([]*Task, 0, limit)
	query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date DESC LIMIT $1`
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

func Search(searchLn string, limit int) ([]*Task, error) {
	var err error
	tasks := make([]*Task, 0, limit)
	query := `SELECT id, date, title, comment, repeat
	FROM scheduler WHERE title LIKE $1
	OR comment LIKE $1 ORDER BY date LIMIT $2`
	parsed, err := time.Parse("02.01.2006", searchLn)
	if err == nil {
		query = `SELECT id, date, title, comment, repeat FROM scheduler WHERE date LIKE $1 LIMIT $2`
		searchLn = parsed.Format(Layout)
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
