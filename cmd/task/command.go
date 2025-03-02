package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/ryuux05/cli-task/task"
)

// 🔹 Executes a single CLI command
func executeCommand(service task.TaskService, args []string) {
	command := args[0]

	switch command {
	case "add":
		if len(args) < 2 {
			fmt.Println("Usage: task add <task_description>")
			return
		}
		flag.CommandLine.Parse(os.Args[2:])
		args := flag.Args()
		newTask := task.NewTaskSchema{
			Name: args[0], 
		}
		service.HandleAdd(newTask)
	
	case "list":
		completed := flag.Bool("c", false, "Show only completed tasks")
		all := flag.Bool("a", false, "Show all tasks")
		flag.CommandLine.Parse(os.Args[2:])
		service.HandleList(*completed, *all)

	case "done":
		if len(args) < 2 {
			fmt.Println("Usage: task done <task_id>")
			return
		}
		id, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Invalid task ID.")
			return
		}
		service.HandleDone(id)

	case "delete":
		if len(args) < 2 {
			fmt.Println("Usage: task delete <task_id>")
			return
		}
		id, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Invalid task ID.")
			return
		}
		service.HandleDelete(id)

	default:
		fmt.Println("Unknown command:", command)
	}
}