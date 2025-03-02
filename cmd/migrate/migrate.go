package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "modernc.org/sqlite"
	"github.com/xwb1989/sqlparser" // SQL parser library
)

func main() {
	//Get the file path 
	filePath := flag.String("f", "", "Path to the migration SQL file")
	flag.Parse()

	if *filePath == "" {
		log.Fatal("You must provide the migration file path using the -f flag")
	}

	//Get the file name
	fileName := filepath.Base(*filePath)
	migrationName := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	
	//Read SQL file
	sqlBytes, err := os.ReadFile(*filePath)
	if err != nil {
		log.Fatalf("Failed to read migration file %v", *filePath, err)
	}
	sqlContent := string(sqlBytes)

	//Compute SHA256 checksum of the file
	hash := sha256.New()
	hash.Write(sqlBytes)
	checkSum := hex.EncodeToString(hash.Sum(nil))

	//Parse SQL into individual statement
	statements, err := sqlparser.SplitStatementToPieces(sqlContent)
	if err != nil {
		log.Fatalf("Error splitting SQL statements %v", err)
		statements = []string{sqlContent} //Fallback
	}

	steps := len(statements)
	if steps == 0 {
		steps = 1
	}

	// Retrieve connection details from environment variables
	// user := os.Getenv("DB_USERNAME")
	// password := os.Getenv("DB_PASSWORD")
	// host := os.Getenv("DB_HOST")
	// port := os.Getenv("DB_PORT")
	// dbname := os.Getenv("DB_NAME")

	// Construct the DSN (Data Source Name)
	// dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require TimeZone=UTC",
	// 	host, user, password, dbname, port)
	const dsn = "db/task.db"

	//Open DB connection
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		log.Fatalf("Failed to open database connection %v", *filePath, err)
	}
	//Close DB connection when the migration done
	defer db.Close()

	//Ensure that schema_migrations table exists.
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			id SERIAL PRIMARY KEY,
			migration_name VARCHAR(255) NOT NULL,
			checksum VARCHAR(64) NOT NULL,
			finished_at TIMESTAMPTZ,
			rolled_back_at TIMESTAMPTZ,
			applied_steps_count INT DEFAULT 0
		)
	`)
	if err != nil {
		log.Fatalf("Failed to create schema_migration table: %v", err)
	}

	//Check if the migration is already applied
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE migration_name = $1)", migrationName).Scan(&exists)
	if err != nil {
		log.Fatalf("Failed to check migrations status: %v", err)
	}
	if exists {
		log.Printf("Migrations %s has already been applied", migrationName)
		return
	}

	//Begin DB transaction
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Failed to begin transaction: %v", err)
	}

	//Execute SQL for each statement individually
	for _, stmt := range statements {
		trimmedStmt := strings.TrimSpace(stmt)
		if trimmedStmt == "" {
			continue
		}

		_, err := tx.Exec(trimmedStmt)
		if err != nil {
			tx.Rollback()
			log.Fatalf("Failed to execute statement: %v\nStatement: %s", err, trimmedStmt)
		}
	}

	//Record the migration in the schmea_migrations table
	_, err = tx.Exec(`
		INSERT INTO schema_migrations(
			migration_name, checksum, finished_at, applied_steps_count
		)
		VALUES($1, $2, $3, $4)`,
		migrationName, checkSum, time.Now(), steps,
	)

	if err != nil {
		tx.Rollback()
		log.Fatalf("Failed to record migration: %v", err)
	}

	//Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}

	log.Printf("Migration %s applied successfully with checksum %s and %d applied steps.\n", migrationName, checkSum, steps)
}