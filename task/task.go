package task

type Task struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Status string `json:"status"`
	CreatedAt string `json:"created_at"`
}

type NewTaskSchema struct {
	Name string `json:"name"`
}

type TaskRepository interface {
	AddTask(newTask NewTaskSchema) (error)
	GetTask() ([]Task)
}

type TaskService interface {
	HandleAdd(newTask NewTaskSchema) (error)
	HandleList(completed bool, all bool) ([]Task, error)
	HandleDelete(id int) (error)
	HandleDone(id int) (error)
}