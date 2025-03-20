package task

import (
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// TaskViewModel enhances Task data for template rendering
type TaskViewModel struct {
	Id          int
	Name        string
	StatusText  string
	StatusClass string
	CreatedAt   string
}

// TasksListViewModel represents a list of tasks for the template
type TasksListViewModel struct {
	Tasks      []TaskViewModel
	TotalTasks int
}

// FromTask converts a regular Task to a view model
func FromTask(task Task) TaskViewModel {
	statusText := "Pending"
	statusClass := "badge-warning"
	if task.Status == "done" {
		statusText = "Completed"
		statusClass = "badge-success"
	}

	return TaskViewModel{
		Id:          task.Id,
		Name:        task.Name,
		StatusText:  statusText,
		StatusClass: statusClass,
		CreatedAt:   task.CreatedAt,
	}
}

// GenerateAndDisplayHTML creates an HTML view for a task and opens it in a browser
func GenerateAndDisplayHTML(task Task) error {
	// Get the template path
	templatePath, err := getTemplatePath("task_view.html")
	if err != nil {
		return err
	}

	// Create a temporary HTML file
	tempFile, err := os.CreateTemp("", "task-*.html")
	if err != nil {
		return fmt.Errorf("error creating temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name()) // Clean up temp file when done

	// Parse and execute the template
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("error parsing template: %v", err)
	}

	// Convert the task to a view model
	viewModel := FromTask(task)

	// Render the template to the temp file
	if err := tmpl.Execute(tempFile, viewModel); err != nil {
		return fmt.Errorf("error rendering template: %v", err)
	}

	tempFile.Close()

	// Open the HTML file in the default browser
	return openInBrowser(tempFile.Name())
}

// GenerateAndDisplayTaskList creates an HTML view for all tasks and opens it in a browser
func GenerateAndDisplayTaskList(tasks []Task) error {
	// Get the template path
	templatePath, err := getTemplatePath("tasks_list.html")
	if err != nil {
		return err
	}

	// Create file path in user's home directory for better visibility
	homeDir, err := os.UserHomeDir()
	var filePath string
	if err != nil {
		// Fallback to temp directory
		tempFile, err := os.CreateTemp("", "tasks-list-*.html")
		if err != nil {
			return fmt.Errorf("error creating temporary file: %v", err)
		}
		filePath = tempFile.Name()
		tempFile.Close()
		fmt.Println("DEBUG: Using temporary file:", filePath)
	} else {
		filePath = filepath.Join(homeDir, "task_cli_tasks_list.html")
		fmt.Println("DEBUG: Using home directory file:", filePath)
	}

	// Create/open the file for writing
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating HTML file: %v", err)
	}
	defer file.Close()

	// Parse the template
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("error parsing template: %v", err)
	}

	// Convert tasks to view models
	var taskViewModels []TaskViewModel
	for _, task := range tasks {
		taskViewModels = append(taskViewModels, FromTask(task))
	}

	// Create the view model for the template
	viewModel := TasksListViewModel{
		Tasks:      taskViewModels,
		TotalTasks: len(taskViewModels),
	}

	// Render the template to the file
	if err := tmpl.Execute(file, viewModel); err != nil {
		return fmt.Errorf("error rendering template: %v", err)
	}

	// Also save a copy to the current working directory for easier access
	debugFilePath := "debug_tasks_list.html"
	debugContent, _ := os.ReadFile(filePath)
	_ = os.WriteFile(debugFilePath, debugContent, 0644)
	fmt.Println("DEBUG: Also saved a debug copy at:", debugFilePath)

	// Open the HTML file in the default browser
	return openInBrowser(filePath)
}

// getTemplatePath finds the template file by name
func getTemplatePath(templateName string) (string, error) {
	fmt.Println("DEBUG: Looking for template:", templateName)

	// Get the executable directory to find the template
	execPath, err := os.Executable()
	if err != nil {
		fmt.Println("DEBUG: Error getting executable path:", err)
		return "", fmt.Errorf("error getting executable path: %v", err)
	}
	fmt.Println("DEBUG: Executable path:", execPath)

	// Find template relative to the executable
	basePath := filepath.Dir(filepath.Dir(execPath))
	templatePath := filepath.Join(basePath, "public", "templates", templateName)
	fmt.Println("DEBUG: First template path attempt:", templatePath)

	// Check if template exists
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		fmt.Println("DEBUG: Template not found in first location, trying fallback")
		// Fallback to current working directory if template not found
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("DEBUG: Error getting working directory:", err)
			return "", fmt.Errorf("error getting working directory: %v", err)
		}
		fmt.Println("DEBUG: Working directory:", cwd)
		templatePath = filepath.Join(cwd, "public", "templates", templateName)
		fmt.Println("DEBUG: Fallback template path:", templatePath)

		// Check if template exists in the fallback location
		if _, err := os.Stat(templatePath); os.IsNotExist(err) {
			fmt.Println("DEBUG: Template not found in fallback location either")
			return "", fmt.Errorf("template %s not found", templateName)
		}
	}

	fmt.Println("DEBUG: Template found at:", templatePath)
	return templatePath, nil
}

// openInBrowser opens the specified file in the default browser
func openInBrowser(filePath string) error {
	fmt.Println("DEBUG: Opening file in browser:", filePath)

	// Ensure the file has a .html extension
	if !strings.HasSuffix(filePath, ".html") {
		newPath := filePath + ".html"
		if err := os.Rename(filePath, newPath); err != nil {
			fmt.Println("DEBUG: Error renaming file:", err)
		} else {
			filePath = newPath
			fmt.Println("DEBUG: Renamed file to:", filePath)
		}
	}

	// Make the file readable by everyone
	if err := os.Chmod(filePath, 0644); err != nil {
		fmt.Println("DEBUG: Error setting file permissions:", err)
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		fmt.Println("DEBUG: Detected macOS, using 'open' command")
		cmd = exec.Command("open", filePath)
	case "windows":
		fmt.Println("DEBUG: Detected Windows, using 'start' command")
		cmd = exec.Command("cmd", "/c", "start", filePath)
	default: // Linux and others
		fmt.Println("DEBUG: Detected Linux/other, using 'xdg-open' command")
		cmd = exec.Command("xdg-open", filePath)
	}

	fmt.Println("DEBUG: Executing command:", cmd.String())
	if err := cmd.Run(); err != nil {
		fmt.Printf("DEBUG: Error opening browser: %v\n", err)
		return fmt.Errorf("error opening browser: %v", err)
	}

	fmt.Println("DEBUG: Browser opened successfully")

	// Print instructions for manual opening
	fmt.Printf("\nHTML file generated at: %s\n", filePath)
	fmt.Println("If the browser didn't open automatically, please open this file manually.")

	return nil
}
