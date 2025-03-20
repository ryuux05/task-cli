package task

type Task struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

type NewTaskSchema struct {
	Name string `json:"name"`
}

type UpdateTaskSchema struct {
	ID        int
	Name      string
	Completed *bool
}

type TaskRepository interface {
	AddTask(task Task) error
	GetTask() ([]Task, error)
	DoneTask(id int) error
	UpdateTask(id int, name string, status string) error
	GetTaskById(id int) (*Task, error)
}

type TaskService interface {
	HandleAdd(newTask NewTaskSchema)
	HandleList(completed bool, all bool)
	HandleDelete(id int)
	HandleDone(id int)
	HandleConnect()
	HandleUpdate(data UpdateTaskSchema)
	HandleSettings()
	HandlePriority()
	HandleViewTask(id int, format string)
	HandleViewAllTasks(format string)
}

func (t *NewTaskSchema) Validate() bool {
	return t.Name != ""
}
