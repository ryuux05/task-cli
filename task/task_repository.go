package task

import "database/sql"

type TaskRepositoryImpl struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) (TaskRepository) {
	return &TaskRepositoryImpl{
		db: db,
	}
}

func (r *TaskRepositoryImpl) AddTask(newTask NewTaskSchema) (error) {

}

func (r *TaskRepositoryImpl) GetTask() ([]Task, error) {

}

