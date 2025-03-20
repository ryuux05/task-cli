# Task CLI

A powerful command-line task management application built with Go. Easily manage your tasks with both a command-line interface and beautiful HTML views.

## Features

- Create, view, update, and delete tasks
- Mark tasks as completed
- List all tasks with filtering options
- View individual tasks or all tasks in HTML format in your default browser
- Interactive CLI mode with command history
- Modern responsive UI for HTML views with filtering and sorting capabilities

## Installation

### Prerequisites

- Go 1.16 or higher
- SQLite3

### Building from Source

1. Clone the repository:
   ```
   git clone https://github.com/ryuux05/task-cli.git
   cd task-cli
   ```

2. Build the application:
   ```
   make build
   ```
   
   This will create a binary in the `bin` directory.

3. Run the database migrations:
   ```
   ./bin/migrate up
   ```

## Usage

The Task CLI supports both direct command execution and an interactive mode.

### Direct Command Execution

Run commands directly:

```
./bin/task <command> [arguments]
```

### Interactive Mode

Start the interactive CLI mode:

```
./bin/task
```

In interactive mode, you'll see a prompt `>` where you can enter commands.

## Commands

### Adding Tasks

Add a new task:

```
task add "Complete project documentation"
```

In interactive mode:
```
> add Complete project documentation
```

### Listing Tasks

List all tasks:

```
task list
```

List only completed tasks:
```
task list -c
```

List all tasks (including completed):
```
task list -a
```

### Completing Tasks

Mark a task as completed:

```
task done <task_id>
```

Example:
```
task done 1
```

### Updating Tasks

Update a task's description:

```
task update <task_id> "New task description"
```

Update a task and mark it as completed:
```
task update <task_id> -c "New task description"
```

### Viewing Tasks

View a single task in HTML format (opens in browser):

```
task view <task_id>
```

View a task in text format:
```
task view <task_id> --format text
```

View all tasks in HTML format (opens in browser):
```
task view-all
```

View all tasks in text format:
```
task view-all --format text
```

### Connecting to External Databases

Connect to a database using a URL (easiest method):

```
task connect -url <database_url>
```

Connect to an external database with individual parameters:

```
task connect -host <hostname> -port <port> -db <database_name> -user <username> -pass <password>
```

Connect to a team-specific database:
```
task connect -team <team_name>
```

This will create or connect to a team-specific SQLite database stored in the storage directory.

Example for URL connection:
```
task connect -url "sqlite:///path/to/database.db"
```

Example for individual parameters:
```
task connect -host localhost -port 5432 -db tasks -user taskuser -pass mypassword
```

Example for team database:
```
task connect -team engineering
```

### Collaborator and Member Management

The Task CLI now supports collaborators for tasks. When you connect to a database for the first time, you'll be prompted to enter your name, which will be stored as the current user.

#### Member Management

You can manage members with the following commands:

- **List Members**: View all members in the database
  ```
  task members
  ```

- **Add Member**: Add a new member to the database
  ```
  task add-member <name>
  ```

- **Switch User**: Change the current user to another member
  ```
  task switch-user <name>
  ```

#### Task Collaboration

When adding or updating tasks, you can specify a collaborator:

- **Add Task with Collaborator**: Add a new task and assign a collaborator
  ```
  task add <task_name> -c <collaborator>
  ```

- **Update Task Collaborator**: Update a task's collaborator
  ```
  task update <id> -c <collaborator>
  ```

Tasks are owned by the current user by default, but you can collaborate with other members in the same database.

### Deleting Tasks

Delete a task:

```
task delete <task_id>
```

## HTML View Features

When viewing tasks in HTML format, you can:

- Filter tasks by name using the search box
- Filter tasks by status (All, Pending, Completed)
- Toggle task status between Pending and Completed
- Edit task details
- Delete tasks
- See a count of total/filtered tasks

## Project Structure

```
task-cli/
├── bin/                 # Binary executables
├── cmd/                 # Command-line applications
│   ├── migrate/         # Database migration tool
│   └── task/            # Main task CLI application
├── db/                  # Database-related code and migrations
├── public/              # Public assets
│   ├── assets/          # Static assets (CSS, JS)
│   └── templates/       # HTML templates
├── storage/             # Persistent storage
└── task/                # Core task management logic
```

## Development

### Running Tests

```
make test
```

### Building for Development

```
make dev
```

### Running Linters

```
make lint
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Acknowledgments

- Built with [Go](https://golang.org/)
- UI built with [Bootstrap](https://getbootstrap.com/)

## Documentation

- [Installation Guide](docs/installation.md)
- [HTML View Documentation](docs/html_view.md)
- [Database Schema](docs/diagrams/database_schema.md)
- [README (Indonesian)](docs/README_ID.md)
- [README (Japanese)](docs/README_JP.md)