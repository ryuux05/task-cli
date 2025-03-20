package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/ryuux05/task-cli/task"
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

	case "update":
		if len(args) < 3 {
			fmt.Println("Usage: task update <task_id> <new_name>")
			return
		}
		id, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Invalid task ID.")
			return
		}
		updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
		completed := updateCmd.Bool("c", false, "Mark as completed")
		updateCmd.Parse(args[2:])

		updateData := task.UpdateTaskSchema{
			ID:        id,
			Name:      updateCmd.Arg(0),
			Completed: completed,
		}

		service.HandleUpdate(updateData)

	case "view":
		if len(args) < 2 {
			fmt.Println("Usage: task view <task_id>")
			return
		}
		id, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Invalid task ID.")
			return
		}
		viewCmd := flag.NewFlagSet("view", flag.ContinueOnError)
		format := viewCmd.String("format", "html", "Output format (html or text)")
		if len(os.Args) > 3 {
			viewCmd.Parse(os.Args[3:])
		}

		service.HandleViewTask(id, *format)

	case "view-all":
		viewAllCmd := flag.NewFlagSet("view-all", flag.ContinueOnError)
		format := viewAllCmd.String("format", "html", "Output format (html or text)")
		if len(os.Args) > 2 {
			viewAllCmd.Parse(os.Args[2:])
		}

		service.HandleViewAllTasks(*format)

	default:
		fmt.Println("Unknown command:", command)
	}
}
