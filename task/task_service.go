package task

import (
	"fmt"
	"log"
)

type TaskServiceImpl struct {
	repository TaskRepository
}

func NewTaskService(repo TaskRepository) TaskService {
	return &TaskServiceImpl{
		repository: repo,
	}
}

func (s *TaskServiceImpl) HandleAdd(newTask NewTaskSchema) {
	if !newTask.Validate() {
		log.Println("Provide task to add")
		return
	}
	task := Task{
		Name: newTask.Name,
	}

	if err := s.repository.AddTask(task); err != nil {
		log.Println("Failed to add task: %v", err)
		return
	}
}

func (s *TaskServiceImpl) HandleList(completed bool, all bool) {
	tasks, err := s.repository.GetTask()
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
	err := s.repository.DoneTask(id)
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

func (s *TaskServiceImpl) HandleConnect() {
	// Implementation for HandleConnect
	fmt.Println("Connect functionality not implemented yet")
}

func (s *TaskServiceImpl) HandleSettings() {
	// Implementation for HandleSettings
	fmt.Println("Settings functionality not implemented yet")
}

func (s *TaskServiceImpl) HandlePriority() {
	// Implementation for HandlePriority
	fmt.Println("Priority functionality not implemented yet")
}

// Actual implementation that will be used
func (s *TaskServiceImpl) HandleUpdate(data UpdateTaskSchema) {
	// Get all tasks
	tasks, err := s.repository.GetTask()
	if err != nil {
		fmt.Println("Error retrieving tasks:", err)
		return
	}

	// Look for the task with the matching ID
	found := false
	for _, task := range tasks {
		if task.Id == data.ID {
			found = true

			// Determine the status based on completed flag
			status := "pending"
			if data.Completed != nil && *data.Completed {
				status = "done"
			}

			// Update task using repository method
			err := s.repository.UpdateTask(data.ID, data.Name, status)
			if err != nil {
				fmt.Printf("Error updating task: %v\n", err)
				return
			}

			fmt.Printf("Task %d updated successfully.\n", data.ID)
			break
		}
	}

	if !found {
		fmt.Printf("Task with ID %d not found.\n", data.ID)
	}
}

func (s *TaskServiceImpl) HandleViewTask(id int, format string) {
	// Get all tasks
	tasks, err := s.repository.GetTask()
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
	tasks, err := s.repository.GetTask()
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
