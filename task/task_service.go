package task

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type TaskServiceImpl struct {
	repo TaskRepository
}

func NewTaskService(repo TaskRepository) TaskService {
	return &TaskServiceImpl{
		repo: repo,
	}
}

func (s *TaskServiceImpl) ensureConnect() error {
	if s.repo == nil {
		return fmt.Errorf("not connected to a repository")
	}
	return nil
}

func (s *TaskServiceImpl) PromptForUsername() error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter your name: ")
	name, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read input: %v", err)
	}

	// Trim whitespace and newlines
	name = strings.TrimSpace(name)

	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	// Set as current member
	return s.repo.SetCurrentMember(name)
}

func (s *TaskServiceImpl) HandleSetupMember() error {
	if err := s.ensureConnect(); err != nil {
		return err
	}

	// Try to get current member
	_, err := s.repo.GetCurrentMember()
	if err != nil {
		// If no current member, prompt for username
		fmt.Println("No username set for this database. Let's set one up.")
		return s.PromptForUsername()
	}

	return nil
}

func (s *TaskServiceImpl) HandleAddTask(name, collaborator string) error {
	if err := s.ensureConnect(); err != nil {
		return err
	}

	return s.repo.AddTask(Task{
		Name:         name,
		Collaborator: collaborator,
	})
}

func (s *TaskServiceImpl) HandleListTasks() error {
	if err := s.ensureConnect(); err != nil {
		return err
	}

	tasks, err := s.repo.GetTask()
	if err != nil {
		return err
	}

	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return nil
	}

	fmt.Println("Tasks:")
	for _, task := range tasks {
		collaboratorInfo := ""
		if task.Collaborator != "" {
			collaboratorInfo = fmt.Sprintf(" (Collaborator: %s)", task.Collaborator)
		}
		fmt.Printf("- ID: %d, Name: %s, Status: %s, Owner: %s%s\n", task.Id, task.Name, task.Status, task.Owner, collaboratorInfo)
	}

	return nil
}

func (s *TaskServiceImpl) HandleUpdateTask(id, name, status, collaborator string) error {
	if err := s.ensureConnect(); err != nil {
		return err
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("Invalid ID: %s", id)
	}

	return s.repo.UpdateTask(idInt, name, status, collaborator)
}

func (s *TaskServiceImpl) HandleGetTask(id string) error {
	if err := s.ensureConnect(); err != nil {
		return err
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("Invalid ID: %s", id)
	}

	task, err := s.repo.GetTaskById(idInt)
	if err != nil {
		return err
	}

	collaboratorInfo := ""
	if task.Collaborator != "" {
		collaboratorInfo = fmt.Sprintf("Collaborator: %s\n", task.Collaborator)
	}

	fmt.Printf("Task Details:\n")
	fmt.Printf("ID: %d\n", task.Id)
	fmt.Printf("Name: %s\n", task.Name)
	fmt.Printf("Status: %s\n", task.Status)
	fmt.Printf("Created At: %s\n", task.CreatedAt)
	fmt.Printf("Owner: %s\n", task.Owner)
	fmt.Printf("%s", collaboratorInfo)

	return nil
}

func (s *TaskServiceImpl) HandleListMembers() error {
	if err := s.ensureConnect(); err != nil {
		return err
	}

	members, err := s.repo.GetAllMembers()
	if err != nil {
		return err
	}

	if len(members) == 0 {
		fmt.Println("No members found.")
		return nil
	}

	// Get current member for highlighting
	currentMember, err := s.repo.GetCurrentMember()
	if err == nil {
		fmt.Printf("Current user: %s\n\n", currentMember)
	}

	fmt.Println("Members:")
	for _, member := range members {
		isCurrent := ""
		if currentMember == member.Name {
			isCurrent = " (you)"
		}
		fmt.Printf("- %s%s (joined: %s)\n", member.Name, isCurrent, member.CreatedAt)
	}

	return nil
}

func (s *TaskServiceImpl) HandleAddMember(name string) error {
	if err := s.ensureConnect(); err != nil {
		return err
	}

	return s.repo.AddMember(name)
}

func (s *TaskServiceImpl) HandleSetCurrentMember(name string) error {
	if err := s.ensureConnect(); err != nil {
		return err
	}

	return s.repo.SetCurrentMember(name)
}

// HandleAdd handles the add command
func (s *TaskServiceImpl) HandleAdd(name string) {
	task := NewTaskSchema{
		Name: name,
	}

	if !task.Validate() {
		fmt.Println("Task name cannot be empty")
		return
	}

	if err := s.repo.AddTask(Task{
		Name: task.Name,
	}); err != nil {
		log.Println("Failed to add task: %v", err)
		return
	}

	fmt.Println("Task added successfully")
}

func (s *TaskServiceImpl) HandleList(completed bool, all bool) {
	tasks, err := s.repo.GetTask()
	if err != nil {
		log.Println("Could't get task: %v", err)
		return
	}

	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return
	}

	fmt.Println("Tasks:")
	for _, task := range tasks {
		status := "[ ]" // Default: Not completed
		if task.Status == "done" {
			status = "[✔]"
		}
		fmt.Printf("%s [%d] %s \n", status, task.Id, task.Name)
	}

}

func (s *TaskServiceImpl) HandleDone(id int) {
	err := s.repo.DoneTask(id)
	if err != nil {
		fmt.Printf("Error marking task %d as done: %v\n", id, err)
		return
	}
	fmt.Printf("Task %d marked as done successfully.\n", id)
}

func (s *TaskServiceImpl) HandleDelete(id int) {
	// Implementation of delete task
	fmt.Printf("Deleting task with ID %d\n", id)
}

// HandleConnect handles the connect command
func (s *TaskServiceImpl) HandleConnect(details ConnectionDetails) error {
	// Connect to the external database
	err := s.repo.ConnectToExternalDB(details)
	if err != nil {
		return fmt.Errorf("error connecting to database: %v", err)
	}

	fmt.Println("Successfully connected to external database!")

	// After connecting, check if we need to set up a user
	err = s.HandleSetupMember()
	if err != nil {
		return fmt.Errorf("error setting up member: %v", err)
	}

	// Get current member to confirm
	member, err := s.repo.GetCurrentMember()
	if err == nil {
		fmt.Printf("Connected as: %s\n", member)
	}

	return nil
}

func (s *TaskServiceImpl) HandleSettings() {
	// Implementation for HandleSettings
	fmt.Println("Settings functionality not implemented yet")
}

func (s *TaskServiceImpl) HandlePriority() {
	// Implementation for HandlePriority
	fmt.Println("Priority functionality not implemented yet")
}

// HandleUpdate handles the update command
func (s *TaskServiceImpl) HandleUpdate(data UpdateTaskSchema) {
	// Get all tasks
	tasks, err := s.repo.GetTask()
	if err != nil {
		fmt.Println("Error retrieving tasks:", err)
		return
	}

	// Check if the task ID exists
	var found bool
	for _, task := range tasks {
		if task.Id == data.ID {
			found = true
			break
		}
	}

	if !found {
		fmt.Printf("No task found with ID %d\n", data.ID)
		return
	}

	// Determine the status based on the Completed field
	status := "pending"
	if data.Completed != nil && *data.Completed {
		status = "done"
	} else if data.Status != "" {
		status = data.Status
	}

	// Update task using repository method
	err = s.repo.UpdateTask(data.ID, data.Name, status, data.Collaborator)
	if err != nil {
		fmt.Printf("Error updating task: %v\n", err)
		return
	}

	fmt.Println("Task updated successfully")
}

func (s *TaskServiceImpl) HandleViewTask(id int, format string) {
	// Get all tasks
	tasks, err := s.repo.GetTask()
	if err != nil {
		fmt.Println("Error retrieving tasks:", err)
		return
	}

	// Find the task with the matching ID
	var task *Task
	for i, t := range tasks {
		if t.Id == id {
			task = &tasks[i]
			break
		}
	}

	if task == nil {
		fmt.Printf("Task with ID %d not found.\n", id)
		return
	}

	if format == "html" {
		// Convert to view.Task
		viewTask := Task{
			Id:        task.Id,
			Name:      task.Name,
			Status:    task.Status,
			CreatedAt: task.CreatedAt,
		}

		if err := GenerateAndDisplayHTML(viewTask); err != nil {
			fmt.Printf("Error displaying HTML view: %v\n", err)
		}
	} else {
		// Display in text format
		status := "Pending"
		if task.Status == "done" {
			status = "Completed"
		}
		fmt.Printf("Task ID: %d\n", task.Id)
		fmt.Printf("Description: %s\n", task.Name)
		fmt.Printf("Status: %s\n", status)
		fmt.Printf("Created At: %s\n", task.CreatedAt)
	}
}

func (s *TaskServiceImpl) HandleViewAllTasks(format string) {
	// Get all tasks
	fmt.Println("DEBUG: HandleViewAllTasks called with format:", format)
	tasks, err := s.repo.GetTask()
	if err != nil {
		fmt.Println("Error retrieving tasks:", err)
		return
	}

	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return
	}

	fmt.Println("DEBUG: Found", len(tasks), "tasks")

	if format == "html" {
		fmt.Println("DEBUG: Generating HTML view")
		// Convert tasks to view.Task
		var viewTasks []Task
		for _, task := range tasks {
			viewTask := Task{
				Id:        task.Id,
				Name:      task.Name,
				Status:    task.Status,
				CreatedAt: task.CreatedAt,
			}
			viewTasks = append(viewTasks, viewTask)
		}

		fmt.Println("DEBUG: Calling GenerateAndDisplayTaskList with", len(viewTasks), "tasks")
		if err := GenerateAndDisplayTaskList(viewTasks); err != nil {
			fmt.Printf("Error displaying HTML view: %v\n", err)
		} else {
			fmt.Println("DEBUG: GenerateAndDisplayTaskList completed successfully")
		}
	} else {
		// Display in text format
		fmt.Println("Tasks:")
		for _, task := range tasks {
			status := "[ ]" // Default: Not completed
			if task.Status == "done" {
				status = "[✔]"
			}
			fmt.Printf("%s [%d] %s\n", status, task.Id, task.Name)
		}
	}
}
