package main

import (
	"api-stockflow/internal/domain"
	"api-stockflow/internal/security"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
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

	log.Println("Starting database seeding...")

	passwordHash := security.NewBcryptPasswordHash()

	// Seed admin user
	seedUser(db, passwordHash, &domain.User{
		ID:       "user-" + uuid.New().String(),
		Username: "admin",
		Password: "admin123",
		Fullname: "Administrator",
		Role:     domain.RoleAdmin,
	})

	// Seed manager user
	seedUser(db, passwordHash, &domain.User{
		ID:       "user-" + uuid.New().String(),
		Username: "manager",
		Password: "manager123",
		Fullname: "Manager User",
		Role:     domain.RoleManager,
	})

	// Seed staff user
	seedUser(db, passwordHash, &domain.User{
		ID:       "user-" + uuid.New().String(),
		Username: "staff",
		Password: "staff123",
		Fullname: "Staff User",
		Role:     domain.RoleStaff,
	})

	log.Println("Database seeding completed successfully!")
	log.Println("\nDefault Users:")
	log.Println("1. Admin    - username: admin,   password: admin123")
	log.Println("2. Manager  - username: manager, password: manager123")
	log.Println("3. Staff    - username: staff,   password: staff123")
}

func seedUser(db *sql.DB, passwordHash security.PasswordHash, user *domain.User) {
	// Check if user already exists
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", user.Username).Scan(&exists)
	if err != nil {
		log.Printf("Error checking user existence: %v\n", err)
		return
	}

	if exists {
		log.Printf("User '%s' already exists, skipping...\n", user.Username)
		return
	}

	// Hash password
	hashedPassword, err := passwordHash.Hash(user.Password)
	if err != nil {
		log.Printf("Failed to hash password for user '%s': %v\n", user.Username, err)
		return
	}

	// Insert user
	query := `
		INSERT INTO users (id, username, password, fullname, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
	`

	_, err = db.Exec(query, user.ID, user.Username, hashedPassword, user.Fullname, user.Role)
	if err != nil {
		log.Printf("Failed to insert user '%s': %v\n", user.Username, err)
		return
	}

	log.Printf("Successfully seeded user: %s (%s)\n", user.Username, user.Role)
}
