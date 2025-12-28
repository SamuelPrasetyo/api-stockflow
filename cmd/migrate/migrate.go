package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run migrate.go [up|down]")
	}

	command := os.Args[1]

	// Database connection
	database := os.Getenv("BLUEPRINT_DB_DATABASE")
	password := os.Getenv("BLUEPRINT_DB_PASSWORD")
	username := os.Getenv("BLUEPRINT_DB_USERNAME")
	port := os.Getenv("BLUEPRINT_DB_PORT")
	host := os.Getenv("BLUEPRINT_DB_HOST")
	schema := os.Getenv("BLUEPRINT_DB_SCHEMA")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s",
		username, password, host, port, database, schema)

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	migrationsDir := "internal/database/migrations"

	switch command {
	case "up":
		runMigrationsUp(db, migrationsDir)
	case "down":
		runMigrationsDown(db, migrationsDir)
	default:
		log.Fatal("Unknown command. Use 'up' or 'down'")
	}
}

func runMigrationsUp(db *sql.DB, dir string) {
	files, err := filepath.Glob(filepath.Join(dir, "*.up.sql"))
	if err != nil {
		log.Fatal("Failed to read migration files:", err)
	}

	sort.Strings(files)

	for _, file := range files {
		log.Printf("Running migration: %s\n", filepath.Base(file))

		content, err := os.ReadFile(file)
		if err != nil {
			log.Fatal("Failed to read file:", err)
		}

		if _, err := db.Exec(string(content)); err != nil {
			log.Fatalf("Failed to execute migration %s: %v\n", filepath.Base(file), err)
		}

		log.Printf("Successfully executed: %s\n", filepath.Base(file))
	}

	log.Println("All migrations completed successfully!")
}

func runMigrationsDown(db *sql.DB, dir string) {
	files, err := filepath.Glob(filepath.Join(dir, "*.down.sql"))
	if err != nil {
		log.Fatal("Failed to read migration files:", err)
	}

	// Sort in reverse order for down migrations
	sort.Sort(sort.Reverse(sort.StringSlice(files)))

	for _, file := range files {
		log.Printf("Rolling back migration: %s\n", filepath.Base(file))

		content, err := os.ReadFile(file)
		if err != nil {
			log.Fatal("Failed to read file:", err)
		}

		if _, err := db.Exec(string(content)); err != nil {
			log.Fatalf("Failed to rollback migration %s: %v\n", filepath.Base(file), err)
		}

		log.Printf("Successfully rolled back: %s\n", filepath.Base(file))
	}

	log.Println("All migrations rolled back successfully!")
}
