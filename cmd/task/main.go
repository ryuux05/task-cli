package main

import (
	"log"
	"os"

	"github.com/ryuux05/task-cli/storage"
	"github.com/ryuux05/task-cli/task"
)

func main() {
	//Init db
	db, err := storage.NewSqlite()
	if err != nil {
		log.Fatalf("Failed to connect to db: %v", err)
	}

	repo := task.NewTaskRepository(db)
	service := task.NewTaskService(repo)

	// 🔹 Check if user provided a command (Single Command Mode)
	if len(os.Args) > 1 {
		executeCommand(service, os.Args[1:])
		return
	}

	startInteractiveMode(service)
}