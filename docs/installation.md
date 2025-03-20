# Task CLI Installation Guide

This guide will walk you through the process of installing and setting up the Task CLI application.

## System Requirements

- **Operating System**: Windows, macOS, or Linux
- **Go**: Version 1.16 or higher
- **SQLite3**: For database storage
- **Git**: For cloning the repository

## Installation Steps

### 1. Install Prerequisites

#### Installing Go

1. Download Go from [golang.org/dl](https://golang.org/dl/)
2. Follow the installation instructions for your operating system
3. Verify the installation by running:
   ```
   go version
   ```

#### Installing SQLite3

**macOS**:
```
brew install sqlite3
```

**Ubuntu/Debian**:
```
sudo apt-get install sqlite3
```

**Windows**:
Download the precompiled binaries from [sqlite.org/download.html](https://sqlite.org/download.html)

### 2. Clone the Repository

```
git clone https://github.com/ryuux05/task-cli.git
cd task-cli
```

### 3. Build the Application

Using the included Makefile:

```
make build
```

This will:
1. Download all required dependencies
2. Compile the application
3. Place the binary in the `bin` directory

If you don't have `make` installed, you can build manually:

```
go build -o bin/task cmd/task/main.go cmd/task/command.go cmd/task/interactive.go
go build -o bin/migrate cmd/migrate/main.go
```

### 4. Initialize the Database

Run the database migrations to set up the SQLite database:

```
./bin/migrate up
```

This will:
1. Create the SQLite database file in the `storage` directory
2. Set up the necessary tables for task management
3. Initialize the status table with "pending" and "done" values

### 5. Test the Installation

Verify that the CLI is working correctly:

```
./bin/task list
```

You should see output indicating that no tasks are found or listing your existing tasks.

## Configuration

### Database Location

By default, the Task CLI uses a SQLite database located in the `storage` directory. You can modify this by editing the database connection string in the main.go file.

### Template Customization

The HTML templates are located in the `public/templates` directory. You can customize these to change the appearance of the HTML views.

### CSS and JavaScript

The CSS and JavaScript files for the HTML views are located in `public/assets/css` and `public/assets/js` respectively. You can modify these to customize the styling and behavior of the HTML views.

## Adding to PATH (Optional)

To use the Task CLI from anywhere on your system, you can add the binary to your PATH:

### macOS/Linux

Add this to your `~/.bashrc`, `~/.zshrc`, or equivalent shell configuration file:

```bash
export PATH="/path/to/task-cli/bin:$PATH"
```

Replace `/path/to/task-cli` with the actual path to the cloned repository.

### Windows

1. Right-click on "This PC" or "My Computer" and select "Properties"
2. Click on "Advanced system settings"
3. Click on "Environment Variables"
4. Under "System variables", find and select "Path", then click "Edit"
5. Click "New" and add the path to the bin directory (e.g., `C:\path\to\task-cli\bin`)
6. Click "OK" to close all dialogs

After adding to PATH, you can use the commands from anywhere:

```
task list
task add "New task"
```

## Troubleshooting

### Command Not Found

If you get a "command not found" error when trying to run the CLI, make sure:

1. You've built the application (`make build`)
2. You're running the command from the correct directory or have added the bin directory to your PATH

### Database Errors

If you encounter database errors:

1. Make sure you've run the migrations (`./bin/migrate up`)
2. Check that the `storage` directory exists and is writable
3. Ensure SQLite3 is installed on your system

### HTML View Issues

If the HTML views don't open in your browser:

1. Check the console output for the path to the HTML file
2. Try opening the file manually in your browser
3. Make sure your default browser is set correctly
4. Look for any error messages in the terminal output

## Updating

To update to the latest version:

```
git pull
make build
./bin/migrate up
```

## Uninstallation

To uninstall, simply delete the cloned repository directory:

```
rm -rf /path/to/task-cli
```

Replace `/path/to/task-cli` with the actual path to the cloned repository. 