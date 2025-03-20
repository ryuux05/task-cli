package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ryuux05/task-cli/task"
)

// 🔹 Executes a single CLI command
func executeCommand(service task.TaskService, args []string) {
	command := args[0]

	switch command {
	case "add":
		if len(args) < 2 {
			fmt.Println("Usage: task add <task_description> [-c <collaborator>]")
			return
		}

		// Create a dedicated FlagSet for the add command
		addCmd := flag.NewFlagSet("add", flag.ContinueOnError)
		collaborator := addCmd.String("c", "", "Collaborator for this task")

		// Find where the task name ends and flags begin
		taskNameParts := []string{}
		flagStartIdx := -1

		for i, arg := range args[1:] {
			if arg == "-c" {
				flagStartIdx = i + 1
				break
			}
			taskNameParts = append(taskNameParts, arg)
		}

		// Process flags if any
		if flagStartIdx != -1 {
			err := addCmd.Parse(args[flagStartIdx+1:])
			if err != nil {
				fmt.Println("Error parsing flags:", err)
				return
			}
		}

		// Combine task name parts
		taskName := strings.Join(taskNameParts, " ")
		fmt.Printf("Adding task: '%s' with collaborator: '%s'\n", taskName, *collaborator)

		err := service.HandleAddTask(taskName, *collaborator)
		if err != nil {
			fmt.Printf("Error adding task: %v\n", err)
			return
		}

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
		collaborator := updateCmd.String("collaborator", "", "Collaborator for this task")
		updateCmd.Parse(args[2:])

		// Set up the completion status
		var isCompleted *bool = completed

		updateData := task.UpdateTaskSchema{
			ID:           id,
			Name:         updateCmd.Arg(0),
			Completed:    isCompleted,
			Collaborator: *collaborator,
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

	case "members":
		service.HandleListMembers()

	case "add-member":
		if len(args) < 2 {
			fmt.Println("Usage: task add-member <member_name>")
			return
		}

		memberName := args[1]
		err := service.HandleAddMember(memberName)
		if err != nil {
			fmt.Println("Error adding member:", err)
			return
		}
		fmt.Println("Member added successfully")

	case "switch-user":
		if len(args) < 2 {
			fmt.Println("Usage: task switch-user <member_name>")
			return
		}

		memberName := args[1]
		err := service.HandleSetCurrentMember(memberName)
		if err != nil {
			fmt.Println("Error switching user:", err)
			return
		}
		fmt.Printf("Switched to user: %s\n", memberName)

	case "connect":
		connectCmd := flag.NewFlagSet("connect", flag.ExitOnError)
		host := connectCmd.String("host", "", "Database host address")
		port := connectCmd.String("port", "", "Database port")
		dbName := connectCmd.String("db", "", "Database name")
		username := connectCmd.String("user", "", "Database username")
		password := connectCmd.String("pass", "", "Database password")
		team := connectCmd.String("team", "", "Team name (if connecting to a team database)")
		url := connectCmd.String("url", "", "Database connection URL (overrides individual connection parameters)")

		// Show usage if no arguments provided
		if len(os.Args) < 3 {
			fmt.Println("Usage: task connect [options]")
			fmt.Println("Options:")
			connectCmd.PrintDefaults()
			return
		}

		connectCmd.Parse(os.Args[2:])

		// Create a ConnectionDetails object
		details := task.ConnectionDetails{}

		// Check connection methods in order of precedence: URL > Team > Individual Parameters
		if *url != "" {
			details.URL = *url
		} else if *team != "" {
			details.Team = *team
		} else if *host != "" || *port != "" || *dbName != "" || *username != "" {
			// At least one individual parameter was specified, check if we have all required ones
			if *host == "" || *port == "" || *dbName == "" || *username == "" {
				fmt.Println("Error: Missing required connection parameters.")
				fmt.Println("For an external database connection, you must provide:")
				fmt.Println("  -host, -port, -db, and -user")
				return
			}

			details.Host = *host
			details.Port = *port
			details.Database = *dbName
			details.Username = *username
			details.Password = *password
		} else {
			fmt.Println("Error: No connection method specified.")
			fmt.Println("You must provide one of:")
			fmt.Println("  -url: A database connection URL")
			fmt.Println("  -team: A team name for a team-specific database")
			fmt.Println("  -host, -port, -db, -user: Individual connection parameters")
			return
		}

		// Connect to the database
		err := service.HandleConnect(details)
		if err != nil {
			fmt.Println(err)
		}

	default:
		fmt.Println("Unknown command:", command)
	}
}
