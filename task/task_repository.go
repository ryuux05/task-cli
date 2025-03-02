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
		INSERT INTO tasks (name, status)
        SELECT ?, (SELECT id FROM status WHERE name = 'pending')
        WHERE NOT EXISTS (SELECT 1 FROM tasks WHERE name = ?);
	`
	res, err := r.db.Exec(query, task.Name, task.Name)
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
		SELECT t.id, t.name, s.name 
        FROM tasks t
        JOIN status s ON t.status = s.id;
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return make([]Task, 0) ,fmt.Errorf("Failed to execute query: %v", err)		
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.Id, &task.Name, &task.Status); err != nil {
			return make([]Task, 0) ,fmt.Errorf("Failed to scan result: %v", err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

