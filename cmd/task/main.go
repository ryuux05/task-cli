package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ryuux05/cli-task/storage"
	"github.com/ryuux05/cli-task/task"
)

func main() {
	//Init db
	if len(os.Args) <2 {
		fmt.Println("Usage: mycli <command>")
		return
	}

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