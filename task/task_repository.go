package task

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type TaskRepositoryImpl struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) TaskRepository {
	return &TaskRepositoryImpl{
		db: db,
	}
}

func (r *TaskRepositoryImpl) AddTask(task Task) error {
	fmt.Println("Adding task:", task.Name)

	// Check if the tasks table exists
	var tableExists int
	err := r.db.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='tasks'").Scan(&tableExists)
	if err != nil {
		return fmt.Errorf("failed to check if tasks table exists: %v", err)
	}

	if tableExists == 0 {
		fmt.Println("Tasks table does not exist, creating it...")
		_, err := r.db.Exec(`
			CREATE TABLE IF NOT EXISTS tasks (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				name TEXT NOT NULL,
				status INTEGER NOT NULL,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				owner TEXT NOT NULL,
				collaborator TEXT,
				FOREIGN KEY (status) REFERENCES status(id),
				FOREIGN KEY (owner) REFERENCES members(name),
				FOREIGN KEY (collaborator) REFERENCES members(name)
			)
		`)
		if err != nil {
			return fmt.Errorf("failed to create tasks table: %v", err)
		}

		// Also check if status table exists
		err = r.db.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='status'").Scan(&tableExists)
		if err != nil {
			return fmt.Errorf("failed to check if status table exists: %v", err)
		}

		if tableExists == 0 {
			fmt.Println("Status table does not exist, creating it...")
			_, err := r.db.Exec(`
				CREATE TABLE IF NOT EXISTS status (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					name TEXT NOT NULL UNIQUE
				)
			`)
			if err != nil {
				return fmt.Errorf("failed to create status table: %v", err)
			}

			// Add default statuses
			_, err = r.db.Exec("INSERT INTO status (id, name) VALUES (1, 'pending'), (2, 'done')")
			if err != nil {
				return fmt.Errorf("failed to add default statuses: %v", err)
			}
		}
	}

	// If owner is not set, get the current member
	if task.Owner == "" {
		owner, err := r.GetCurrentMember()
		if err != nil {
			return fmt.Errorf("failed to get current member: %v", err)
		}
		task.Owner = owner
	}

	// If collaborator is set, ensure the member exists
	if task.Collaborator != "" {
		if err := r.AddMember(task.Collaborator); err != nil {
			return err
		}
	}

	fmt.Printf("Adding task with owner: %s, collaborator: %s\n", task.Owner, task.Collaborator)

	query := `
		INSERT INTO tasks (name, status, owner, collaborator)
        SELECT ?, (SELECT id FROM status WHERE name = 'pending'), ?, ?
        WHERE NOT EXISTS (SELECT 1 FROM tasks WHERE name = ? AND owner = ?);
	`
	res, err := r.db.Exec(query, task.Name, task.Owner, task.Collaborator, task.Name, task.Owner)
	if err != nil {
		return fmt.Errorf("Failed to execute query: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("Failed to execute query: %v", err)
	}

	if rowsAffected == 0 {
		fmt.Println("Task already exists for this owner, skipping insert.")
	} else {
		fmt.Println("Task added successfully.")
	}
	return nil
}

func (r *TaskRepositoryImpl) GetTask() ([]Task, error) {
	query := `
		SELECT t.id, t.name, s.name, t.created_at, t.owner, IFNULL(t.collaborator, '')
        FROM tasks t
        JOIN status s ON t.status = s.id;
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return make([]Task, 0), fmt.Errorf("Failed to execute query: %v", err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.Id, &task.Name, &task.Status, &task.CreatedAt, &task.Owner, &task.Collaborator); err != nil {
			return make([]Task, 0), fmt.Errorf("Failed to scan result: %v", err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *TaskRepositoryImpl) DoneTask(id int) error {
	query := `
		UPDATE tasks 
		SET status = (SELECT id FROM status WHERE name = 'done')
		WHERE id = ?
	`

	res, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("Failed to mark task as done: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("Failed to get affected rows: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("No task found with ID %d", id)
	}

	return nil
}

func (r *TaskRepositoryImpl) UpdateTask(id int, name string, status string, collaborator string) error {
	// First, check if the collaborator exists and add them if needed
	if collaborator != "" {
		if err := r.AddMember(collaborator); err != nil {
			return err
		}
	}

	// Ensure status exists in the status table
	query := `
		UPDATE tasks 
		SET name = ?,
			status = (SELECT id FROM status WHERE name = ?),
			collaborator = ?
		WHERE id = ?
	`

	res, err := r.db.Exec(query, name, status, collaborator, id)
	if err != nil {
		return fmt.Errorf("Failed to update task: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("Failed to get affected rows: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("No task found with ID %d", id)
	}

	return nil
}

func (r *TaskRepositoryImpl) GetTaskById(id int) (*Task, error) {
	query := `
		SELECT t.id, t.name, s.name, t.created_at, t.owner, IFNULL(t.collaborator, '')
		FROM tasks t
		JOIN status s ON t.status = s.id
		WHERE t.id = ?
	`

	var task Task
	err := r.db.QueryRow(query, id).Scan(&task.Id, &task.Name, &task.Status, &task.CreatedAt, &task.Owner, &task.Collaborator)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("No task found with ID %d", id)
		}
		return nil, fmt.Errorf("Failed to query task: %v", err)
	}

	return &task, nil
}

// DeleteTask deletes a task with the given ID
func (r *TaskRepositoryImpl) DeleteTask(id int) error {
	query := `DELETE FROM tasks WHERE id = ?`
	res, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("Failed to delete task: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("Failed to get affected rows: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("No task found with ID %d", id)
	}

	return nil
}

// ConnectToExternalDB connects to an external database
func (r *TaskRepositoryImpl) ConnectToExternalDB(details ConnectionDetails) error {
	fmt.Println("Connecting to database...")

	// Close existing database connection if any
	if r.db != nil {
		fmt.Println("Closing existing database connection...")
		if err := r.db.Close(); err != nil {
			return fmt.Errorf("failed to close existing database connection: %v", err)
		}
	}

	// For now, just use SQLite in file mode for demonstration
	sqlitePath := "storage/tasks.db"
	if details.Team != "" {
		sqlitePath = fmt.Sprintf("storage/team_%s.db", details.Team)
	}

	// Ensure storage directory exists
	fmt.Println("Ensuring storage directory exists...")
	if err := os.MkdirAll("storage", 0755); err != nil {
		return fmt.Errorf("failed to create storage directory: %v", err)
	}

	fmt.Printf("Opening database connection to %s...\n", sqlitePath)
	// Open the database connection
	db, err := sql.Open("sqlite3", sqlitePath)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	// Test the connection
	fmt.Println("Testing database connection...")
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	// Set the connection in the repository
	r.db = db

	// Ensure required tables exist
	fmt.Println("Ensuring tables exist...")
	if err := r.ensureTablesExist(); err != nil {
		return fmt.Errorf("failed to ensure tables exist: %v", err)
	}

	return nil
}

// ensureTablesExist checks if the necessary tables exist in the database
// and creates them if they don't
func (r *TaskRepositoryImpl) ensureTablesExist() error {
	fmt.Println("Creating database tables if needed...")

	// Run each statement separately for better error handling
	statements := []string{
		// Create status table
		`CREATE TABLE IF NOT EXISTS status (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE
		)`,

		// Create members table
		`CREATE TABLE IF NOT EXISTS members (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		// Create current_member table
		`CREATE TABLE IF NOT EXISTS current_member (
			id INTEGER PRIMARY KEY,
			member_name TEXT NOT NULL,
			FOREIGN KEY (member_name) REFERENCES members(name)
		)`,

		// Create tasks table
		`CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			status INTEGER NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			owner TEXT NOT NULL,
			collaborator TEXT,
			FOREIGN KEY (status) REFERENCES status(id),
			FOREIGN KEY (owner) REFERENCES members(name),
			FOREIGN KEY (collaborator) REFERENCES members(name)
		)`,

		// Add default statuses if not exist
		`INSERT OR IGNORE INTO status (id, name) VALUES (1, 'pending'), (2, 'done')`,
	}

	// Execute each SQL statement
	for _, stmt := range statements {
		_, err := r.db.Exec(stmt)
		if err != nil {
			return fmt.Errorf("error executing SQL statement: %v\nStatement: %s", err, stmt)
		}
	}

	fmt.Println("Database tables created successfully.")
	return nil
}

// SetupMemberTable ensures the members table is set up and prompts for username if needed
func (r *TaskRepositoryImpl) SetupMemberTable() error {
	fmt.Println("Setting up member table...")

	// Check if the tables exist first
	var tableExists int
	err := r.db.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='members'").Scan(&tableExists)
	if err != nil {
		return fmt.Errorf("failed to check if members table exists: %v", err)
	}

	if tableExists == 0 {
		fmt.Println("Members table does not exist, creating it...")
		_, err := r.db.Exec(`
			CREATE TABLE IF NOT EXISTS members (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				name TEXT NOT NULL UNIQUE,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			)
		`)
		if err != nil {
			return fmt.Errorf("failed to create members table: %v", err)
		}
	}

	err = r.db.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='current_member'").Scan(&tableExists)
	if err != nil {
		return fmt.Errorf("failed to check if current_member table exists: %v", err)
	}

	if tableExists == 0 {
		fmt.Println("Current member table does not exist, creating it...")
		_, err := r.db.Exec(`
			CREATE TABLE IF NOT EXISTS current_member (
				id INTEGER PRIMARY KEY,
				member_name TEXT NOT NULL,
				FOREIGN KEY (member_name) REFERENCES members(name)
			)
		`)
		if err != nil {
			return fmt.Errorf("failed to create current_member table: %v", err)
		}
	}

	// Check if we have any members
	var count int
	err = r.db.QueryRow("SELECT COUNT(*) FROM members").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check members table: %v", err)
	}

	// Check if we have a current member set
	var currentMemberCount int
	err = r.db.QueryRow("SELECT COUNT(*) FROM current_member").Scan(&currentMemberCount)
	if err != nil {
		return fmt.Errorf("failed to check current member: %v", err)
	}

	fmt.Printf("Member setup complete. Members: %d, Current member set: %t\n", count, currentMemberCount > 0)
	return nil
}

// GetCurrentMember returns the name of the current member
func (r *TaskRepositoryImpl) GetCurrentMember() (string, error) {
	fmt.Println("Getting current member...")

	var name string

	// Check if current_member table has any rows
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM current_member").Scan(&count)
	if err != nil {
		return "", fmt.Errorf("failed to check current member: %v", err)
	}

	if count == 0 {
		return "", fmt.Errorf("no current member set")
	}

	// Get the current member name
	err = r.db.QueryRow("SELECT member_name FROM current_member LIMIT 1").Scan(&name)
	if err != nil {
		return "", fmt.Errorf("failed to get current member: %v", err)
	}

	fmt.Printf("Current member is: %s\n", name)
	return name, nil
}

// SetCurrentMember sets the current member
func (r *TaskRepositoryImpl) SetCurrentMember(name string) error {
	fmt.Printf("Setting current member to: %s\n", name)

	// First, ensure the member exists
	err := r.AddMember(name)
	if err != nil {
		return err
	}

	// Check if the current_member table exists
	var tableExists int
	err = r.db.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='current_member'").Scan(&tableExists)
	if err != nil {
		return fmt.Errorf("failed to check if current_member table exists: %v", err)
	}

	if tableExists == 0 {
		fmt.Println("Current member table does not exist, creating it...")
		_, err := r.db.Exec(`
			CREATE TABLE IF NOT EXISTS current_member (
				id INTEGER PRIMARY KEY,
				member_name TEXT NOT NULL,
				FOREIGN KEY (member_name) REFERENCES members(name)
			)
		`)
		if err != nil {
			return fmt.Errorf("failed to create current_member table: %v", err)
		}
	}

	// Clear existing current member
	fmt.Println("Clearing current member table...")
	_, err = r.db.Exec("DELETE FROM current_member")
	if err != nil {
		return fmt.Errorf("failed to clear current member: %v", err)
	}

	// Set new current member
	fmt.Println("Setting new current member...")
	_, err = r.db.Exec("INSERT INTO current_member (id, member_name) VALUES (1, ?)", name)
	if err != nil {
		return fmt.Errorf("failed to set current member: %v", err)
	}

	fmt.Printf("Current member set to: %s\n", name)
	return nil
}

// GetAllMembers returns all members
func (r *TaskRepositoryImpl) GetAllMembers() ([]Member, error) {
	query := "SELECT id, name, created_at FROM members ORDER BY name"

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query members: %v", err)
	}
	defer rows.Close()

	var members []Member
	for rows.Next() {
		var member Member
		if err := rows.Scan(&member.Id, &member.Name, &member.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan member: %v", err)
		}
		members = append(members, member)
	}

	return members, nil
}

// AddMember adds a new member if they don't already exist
func (r *TaskRepositoryImpl) AddMember(name string) error {
	fmt.Printf("Adding member: %s\n", name)

	// Check if the table exists first
	var tableExists int
	err := r.db.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='members'").Scan(&tableExists)
	if err != nil {
		return fmt.Errorf("failed to check if members table exists: %v", err)
	}

	if tableExists == 0 {
		fmt.Println("Members table does not exist, creating it...")
		_, err := r.db.Exec(`
			CREATE TABLE IF NOT EXISTS members (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				name TEXT NOT NULL UNIQUE,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			)
		`)
		if err != nil {
			return fmt.Errorf("failed to create members table: %v", err)
		}
	}

	if name == "" {
		return fmt.Errorf("member name cannot be empty")
	}

	// Check if member already exists
	var count int
	err = r.db.QueryRow("SELECT COUNT(*) FROM members WHERE name = ?", name).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check if member exists: %v", err)
	}

	// If member doesn't exist, add them
	if count == 0 {
		fmt.Printf("Member %s does not exist, adding...\n", name)
		_, err = r.db.Exec("INSERT INTO members (name) VALUES (?)", name)
		if err != nil {
			return fmt.Errorf("failed to add member: %v", err)
		}
		fmt.Printf("Member %s added successfully\n", name)
	} else {
		fmt.Printf("Member %s already exists\n", name)
	}

	return nil
}
