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
			addCmd := flag.NewFlagSet("add", flag.ContinueOnError)
			addCmd.Usage = func() {
				fmt.Println("Usage: add <task_name> [-c <collaborator>]")
				fmt.Println("Add a new task with the given name")
			}
			addCollaborator := addCmd.String("c", "", "Collaborator for this task")

			err := addCmd.Parse(args[1:])
			if err != nil {
				if err == flag.ErrHelp {
					continue
				}
				fmt.Println("Error parsing add command:", err)
				continue
			}

			if addCmd.NArg() < 1 {
				fmt.Println("Error: Task name is required")
				addCmd.Usage()
				continue
			}

			taskName := addCmd.Arg(0)
			err = service.HandleAddTask(taskName, *addCollaborator)
			if err != nil {
				fmt.Println("Error adding task:", err)
			}
			continue

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
			updateCmd := flag.NewFlagSet("update", flag.ContinueOnError)
			updateCmd.Usage = func() {
				fmt.Println("Usage: update <id> -name <new_name> -status <new_status> [-c <collaborator>]")
				fmt.Println("Update a task")
			}
			updateName := updateCmd.String("name", "", "New task name")
			updateStatus := updateCmd.String("status", "", "New task status (pending/done)")
			updateCollaborator := updateCmd.String("c", "", "Collaborator for this task")

			err := updateCmd.Parse(args[1:])
			if err != nil {
				if err == flag.ErrHelp {
					continue
				}
				fmt.Println("Error parsing update command:", err)
				continue
			}

			if updateCmd.NArg() < 1 {
				fmt.Println("Error: Task ID is required")
				updateCmd.Usage()
				continue
			}

			idStr := updateCmd.Arg(0)

			if *updateName == "" && *updateStatus == "" && *updateCollaborator == "" {
				fmt.Println("Error: At least one field to update must be provided")
				updateCmd.Usage()
				continue
			}

			err = service.HandleUpdateTask(idStr, *updateName, *updateStatus, *updateCollaborator)
			if err != nil {
				fmt.Println("Error updating task:", err)
			} else {
				fmt.Println("Task updated successfully")
			}
			continue

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

		case "connect":
			connectCmd := flag.NewFlagSet("connect", flag.ContinueOnError)
			connectCmd.Usage = func() {
				fmt.Println("Usage: connect [options]")
				fmt.Println("Connect to an external database using one of these methods:")
				fmt.Println("  1. URL: connect -url <database_url>")
				fmt.Println("  2. Team: connect -team <team_name>")
				fmt.Println("  3. Individual parameters: connect -host <host> -port <port> -db <database> -user <username> -pass <password>")
			}

			// Connection options
			url := connectCmd.String("url", "", "Database connection URL (overrides individual connection parameters)")
			team := connectCmd.String("team", "", "Team name (for predefined connections)")
			host := connectCmd.String("host", "", "Database host")
			port := connectCmd.String("port", "", "Database port")
			dbName := connectCmd.String("db", "", "Database name")
			user := connectCmd.String("user", "", "Database username")
			pass := connectCmd.String("pass", "", "Database password")

			err := connectCmd.Parse(args[1:])
			if err != nil {
				if err == flag.ErrHelp {
					continue
				}
				fmt.Println("Error parsing connect command:", err)
				continue
			}

			// Create connection details
			details := task.ConnectionDetails{
				URL:      *url,
				Team:     *team,
				Host:     *host,
				Port:     *port,
				Database: *dbName,
				Username: *user,
				Password: *pass,
			}

			// Connect to the database
			err = service.HandleConnect(details)
			if err != nil {
				fmt.Println(err)
			}

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

		case "members":
			membersCmd := flag.NewFlagSet("members", flag.ContinueOnError)
			membersCmd.Usage = func() {
				fmt.Println("Usage: members")
				fmt.Println("List all members")
			}

			err := membersCmd.Parse(args[1:])
			if err != nil {
				if err == flag.ErrHelp {
					continue
				}
				fmt.Println("Error parsing members command:", err)
				continue
			}

			err = service.HandleListMembers()
			if err != nil {
				fmt.Println("Error listing members:", err)
			}
			continue

		case "add-member":
			addMemberCmd := flag.NewFlagSet("add-member", flag.ContinueOnError)
			addMemberCmd.Usage = func() {
				fmt.Println("Usage: add-member <member_name>")
				fmt.Println("Add a new member with the given name")
			}

			err := addMemberCmd.Parse(args[1:])
			if err != nil {
				if err == flag.ErrHelp {
					continue
				}
				fmt.Println("Error parsing add-member command:", err)
				continue
			}

			if addMemberCmd.NArg() < 1 {
				fmt.Println("Error: Member name is required")
				addMemberCmd.Usage()
				continue
			}

			memberName := addMemberCmd.Arg(0)
			err = service.HandleAddMember(memberName)
			if err != nil {
				fmt.Println("Error adding member:", err)
			} else {
				fmt.Println("Member added successfully")
			}
			continue

		case "switch-user":
			switchUserCmd := flag.NewFlagSet("switch-user", flag.ContinueOnError)
			switchUserCmd.Usage = func() {
				fmt.Println("Usage: switch-user <member_name>")
				fmt.Println("Switch to another user")
			}

			err := switchUserCmd.Parse(args[1:])
			if err != nil {
				if err == flag.ErrHelp {
					continue
				}
				fmt.Println("Error parsing switch-user command:", err)
				continue
			}

			if switchUserCmd.NArg() < 1 {
				fmt.Println("Error: Member name is required")
				switchUserCmd.Usage()
				continue
			}

			memberName := switchUserCmd.Arg(0)
			err = service.HandleSetCurrentMember(memberName)
			if err != nil {
				fmt.Println("Error switching user:", err)
			} else {
				fmt.Printf("Switched to user: %s\n", memberName)
			}
			continue

		case "help":
			fmt.Println("Available commands:")
			fmt.Println("  add <task_name> [-c <collaborator>] - Add a new task")
			fmt.Println("  list [--all] [--completed] - List tasks")
			fmt.Println("  done <id> - Mark a task as done")
			fmt.Println("  update <id> -name <new_name> -status <new_status> [-c <collaborator>] - Update a task")
			fmt.Println("  view <id> [-format html|text] - View details of a task")
			fmt.Println("  view-all [-format html|text] - View all tasks")
			fmt.Println("  delete <id> - Delete a task")
			fmt.Println("  members - List all members")
			fmt.Println("  switch-space <name> - Switch to another task space")
			fmt.Println("  connect -host <host> -port <port> -db <dbname> -user <username> -pass <password> - Connect to an external database")
			fmt.Println("  connect -team <team_name> - Connect to a team database")
			fmt.Println("  connect -url <connection_url> - Connect using a database URL")
			fmt.Println("  help - Show this help text")
			fmt.Println("  exit - Exit the program")

		case "exit":
			fmt.Println("Goodbye!")
			return

		default:
			fmt.Println("Unknown command:", command)
		}
	}
}
