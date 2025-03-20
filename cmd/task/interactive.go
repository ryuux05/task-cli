package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ryuux05/task-cli/task"
)

func startInteractiveMode(service task.TaskService) {
	fmt.Println("Task Manager CLI")
	fmt.Println("Type 'help' for commands or 'exit' to quit.")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ") // CLI Prompt
		scanner.Scan()
		input := scanner.Text()
		args := strings.Fields(input)

		if len(args) == 0 {
			continue
		}

		command := args[0]

		switch command {
		case "add":
			if len(args) < 2 {
				fmt.Println("Usage: add <task_description>")
				continue
			}
			taskText := strings.Join(args[1:], " ")
			service.HandleAdd(task.NewTaskSchema{Name: taskText})

		case "list":
			listCmd := flag.NewFlagSet("list", flag.ExitOnError)
			completed := listCmd.Bool("c", false, "Show only completed tasks")
			all := listCmd.Bool("a", false, "Show all tasks")
			flag.CommandLine.Parse(os.Args[1:])
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
			// Parse from the remaining interactive arguments
			if len(args) > 2 {
				viewCmd.Parse(args[2:])
			}

			service.HandleViewTask(id, *format)

		case "view-all":
			viewAllCmd := flag.NewFlagSet("view-all", flag.ContinueOnError)
			format := viewAllCmd.String("format", "html", "Output format (html or text)")
			viewAllCmd.Parse(args[1:])

			service.HandleViewAllTasks(*format)

		case "delete":
			if len(args) < 2 {
				fmt.Println("Usage: delete <task_id>")
				continue
			}
			id, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Println("Invalid task ID.")
				continue
			}
			service.HandleDelete(id)

		case "help":
			fmt.Println("Available commands:")
			fmt.Println("  add <task_description> - Add a new task")
			fmt.Println("  list - Show all tasks")
			fmt.Println("  done <task_id> - Mark a task as completed")
			fmt.Println("  update <task_id> <new_name> - Update a task")
			fmt.Println("  view <task_id> - View task details")
			fmt.Println("  view-all - View all tasks")
			fmt.Println("  delete <task_id> - Delete a task")
			fmt.Println("  exit - Exit the CLI")

		case "exit":
			fmt.Println("Goodbye!")
			return

		default:
			fmt.Println("Unknown command:", command)
		}
	}
}
