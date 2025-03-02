package task

import (
	"fmt"
	"log"
)

type TaskServiceImpl struct {
	repository TaskRepository
}

func NewTaskService(repo TaskRepository) (TaskService) {
	return&TaskServiceImpl{
		repository: repo,
	}
}

func (s *TaskServiceImpl) HandleAdd(newTask NewTaskSchema) {
	if !newTask.Validate() {
		log.Println("Provide task to add")
	}
	task := Task{
		Name: newTask.Name,
	}

	if err := s.repository.AddTask(task); err != nil {
		log.Println("Failed to add task: %v", err)
	}

	log.Println("Task added")
}

func (s *TaskServiceImpl) HandleList(completed bool, all bool) {
	tasks, err := s.repository.GetTask()
	if err != nil {
		log.Println("Could't get task")
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
		fmt.Printf("%s [%d] %s (Priority: %s)\n", status, task.Id, task.Name)
	}
	
}

func (s *TaskServiceImpl) HandleDelete(id int) {

}	

func (s *TaskServiceImpl) HandleDone(id int) {

}