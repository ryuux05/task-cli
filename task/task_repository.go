package task

import (
	"database/sql"
	"fmt"
)

type TaskRepositoryImpl struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) (TaskRepository) {
	return &TaskRepositoryImpl{
		db: db,
	}
}

func (r *TaskRepositoryImpl) AddTask(task Task) (error) {
	query := `
		INSERT OR IGNORE INTO tasks (name) VALUES ($1)
	`
	res, err := r.db.Exec(query, task.Name)
	if err != nil {
		return fmt.Errorf("Failed to execute query: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("Failed to execute query: %v", err)
	}

	if rowsAffected == 0 {
		fmt.Println("Task already exists, skipping insert.")
	} else {
		fmt.Println("Task added successfully.")
	}
	return nil
}

func (r *TaskRepositoryImpl) GetTask() ([]Task, error) {
	query := `
		SELECT t.name, s.name 
        FROM tasks t
        JOIN status s ON t.status_id = s.id;
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return make([]Task, 0) ,fmt.Errorf("Failed to execute query: %v", err)		
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.Name, &task.Status); err != nil {
			return make([]Task, 0) ,fmt.Errorf("Failed to scan result: %v", err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

