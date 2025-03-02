package main

import (
	"fmt"
	"bufio"
	"flag"
	"os"
	"strings"
	"strconv"
	"github.com/ryuux05/cli-task/task"
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
			completed := flag.Bool("c", false, "Show only completed tasks")
			all := flag.Bool("a", false, "Show all tasks")
			flag.CommandLine.Parse(os.Args[2:])
			service.HandleList(*completed, *all)

		case "done":
			if len(args) < 2 {
				fmt.Println("Usage: done <task_id>")
				continue
			}
			id, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Println("Invalid task ID.")
				continue
			}
			service.HandleDone(id)

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