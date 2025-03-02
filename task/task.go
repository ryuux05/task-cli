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
	AddTask(task Task) (error)
	GetTask() ([]Task, error)
}

type TaskService interface {
	HandleAdd(newTask NewTaskSchema) 
	HandleList(completed bool, all bool) 
	HandleDelete(id int) 
	HandleDone(id int)
}

func(t *NewTaskSchema) Validate() (bool) {
	return t.Name != ""
}