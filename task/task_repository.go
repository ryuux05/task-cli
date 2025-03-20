package task

import (
	"database/sql"
	"fmt"
)

type TaskRepositoryImpl struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) TaskRepository {
	return &TaskRepositoryImpl{
		db: db,
	}
}

func (r *TaskRepositoryImpl) AddTask(task Task) error {
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
		SELECT t.id, t.name, s.name, t.created_at
        FROM tasks t
        JOIN status s ON t.status = s.id;
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return make([]Task, 0), fmt.Errorf("Failed to execute query: %v", err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.Id, &task.Name, &task.Status, &task.CreatedAt); err != nil {
			return make([]Task, 0), fmt.Errorf("Failed to scan result: %v", err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *TaskRepositoryImpl) DoneTask(id int) error {
	query := `
		UPDATE tasks 
		SET status = (SELECT id FROM status WHERE name = 'done')
		WHERE id = ?
	`

	res, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("Failed to mark task as done: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("Failed to get affected rows: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("No task found with ID %d", id)
	}

	return nil
}

func (r *TaskRepositoryImpl) UpdateTask(id int, name string, status string) error {
	// Ensure status exists in the status table
	query := `
		UPDATE tasks 
		SET name = ?,
			status = (SELECT id FROM status WHERE name = ?)
		WHERE id = ?
	`

	res, err := r.db.Exec(query, name, status, id)
	if err != nil {
		return fmt.Errorf("Failed to update task: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("Failed to get affected rows: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("No task found with ID %d", id)
	}

	return nil
}

func (r *TaskRepositoryImpl) GetTaskById(id int) (*Task, error) {
	query := `
		SELECT t.id, t.name, s.name, t.created_at
		FROM tasks t
		JOIN status s ON t.status = s.id
		WHERE t.id = ?
	`

	var task Task
	err := r.db.QueryRow(query, id).Scan(&task.Id, &task.Name, &task.Status, &task.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("No task found with ID %d", id)
		}
		return nil, fmt.Errorf("Failed to query task: %v", err)
	}

	return &task, nil
}
