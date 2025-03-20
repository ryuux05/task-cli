//go:generate mockgen -source=task.go -destination=mock_task.go -package=task
package task

import (
	"errors"
)

// Task represents a task in the task list
type Task struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Status       string `json:"status"`
	CreatedAt    string `json:"created_at"`
	Owner        string `json:"owner"`
	Collaborator string `json:"collaborator"`
}

// Member represents a user in the system
type Member struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

// NewTaskSchema is the schema for adding a new task
type NewTaskSchema struct {
	Name         string `json:"name"`
	Collaborator string `json:"collaborator,omitempty"`
}

// UpdateTaskSchema is the schema for updating a task
type UpdateTaskSchema struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Status       string `json:"status"`
	Completed    *bool  `json:"completed,omitempty"`
	Collaborator string `json:"collaborator,omitempty"`
}

// ConnectionDetails is the details for connecting to an external database
type ConnectionDetails struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
	Team     string
	URL      string
}

// Validate validates the connection details
func (c *ConnectionDetails) Validate() error {
	// If we have a URL, it overrides all other fields
	if c.URL != "" {
		return nil
	}

	// If we have a team name, it's valid
	if c.Team != "" {
		return nil
	}

	// Need to have host, port, database name, username, and password
	if c.Host == "" {
		return errors.New("host is required")
	}
	if c.Port == "" {
		return errors.New("port is required")
	}
	if c.Database == "" {
		return errors.New("database is required")
	}
	if c.Username == "" {
		return errors.New("username is required")
	}
	if c.Password == "" {
		return errors.New("password is required")
	}

	return nil
}

// TaskRepository is the repository for task-related operations
type TaskRepository interface {
	// Task operations
	AddTask(task Task) error
	GetTask() ([]Task, error)
	UpdateTask(id int, name string, status string, collaborator string) error
	DoneTask(id int) error
	GetTaskById(id int) (*Task, error)
	DeleteTask(id int) error

	// Database operations
	ConnectToExternalDB(details ConnectionDetails) error

	// Member operations
	SetupMemberTable() error
	GetCurrentMember() (string, error)
	SetCurrentMember(name string) error
	GetAllMembers() ([]Member, error)
	AddMember(name string) error
}

// TaskService is the service for task-related operations
type TaskService interface {
	// Task management
	HandleAdd(name string)
	HandleList(completed bool, all bool)
	HandleDone(id int)
	HandleUpdate(data UpdateTaskSchema)
	HandleDelete(id int)
	HandleViewTask(id int, format string)
	HandleViewAllTasks(format string)

	// Database connection
	HandleConnect(details ConnectionDetails) error

	// New methods with collaborators
	HandleAddTask(name, collaborator string) error
	HandleListTasks() error
	HandleUpdateTask(id, name, status, collaborator string) error
	HandleGetTask(id string) error

	// Member management
	HandleSetupMember() error
	PromptForUsername() error
	HandleListMembers() error
	HandleAddMember(name string) error
	HandleSetCurrentMember(name string) error
}

// ViewService is the service for rendering tasks in different formats
type ViewService interface {
	RenderTaskHTML(task Task) (string, error)
	RenderAllTasksHTML(tasks []Task) (string, error)
}

func (t *NewTaskSchema) Validate() bool {
	return t.Name != ""
}
