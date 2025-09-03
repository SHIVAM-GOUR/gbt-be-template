package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"gbt-be-template/internal/config"
	"gbt-be-template/internal/models"
	"gbt-be-template/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Command line flags
	email := flag.String("email", "", "Admin email (required)")
	username := flag.String("username", "", "Admin username (required)")
	password := flag.String("password", "", "Admin password (required)")
	firstName := flag.String("first-name", "", "Admin first name (required)")
	lastName := flag.String("last-name", "", "Admin last name (required)")
	flag.Parse()

	// Validate required fields
	if *email == "" || *username == "" || *password == "" || *firstName == "" || *lastName == "" {
		fmt.Println("Usage: go run cmd/seed/main.go -email=admin@example.com -username=admin -password=securepassword -first-name=Admin -last-name=User")
		os.Exit(1)
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	// logger := logger.New(cfg.Log.Level, cfg.Log.Format)
	// removed logger from below "db, err := repository.NewDatabase(cfg)" because newDatabase function accept only 1 argument

	// Initialize database
	db, err := repository.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repository
	userRepo := repository.NewUserRepository(db)

	// Check if admin already exists
	ctx := context.Background()
	existingUser, err := userRepo.GetByEmail(ctx, *email)
	if err != nil {
		log.Fatalf("Failed to check existing user: %v", err)
	}
	if existingUser != nil {
		log.Fatalf("User with email %s already exists", *email)
	}

	// Check if username is taken
	existingUser, err = userRepo.GetByUsername(ctx, *username)
	if err != nil {
		log.Fatalf("Failed to check existing username: %v", err)
	}
	if existingUser != nil {
		log.Fatalf("Username %s is already taken", *username)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	// Create admin user
	adminUser := &models.User{
		Email:     *email,
		Username:  *username,
		Password:  string(hashedPassword),
		FirstName: *firstName,
		LastName:  *lastName,
		IsActive:  true,
		IsAdmin:   true, // This is the key difference - set as admin
	}

	// Save to database
	if err := userRepo.Create(ctx, adminUser); err != nil {
		log.Fatalf("Failed to create admin user: %v", err)
	}

	fmt.Printf("✅ Admin user created successfully!\n")
	fmt.Printf("   Email: %s\n", adminUser.Email)
	fmt.Printf("   Username: %s\n", adminUser.Username)
	fmt.Printf("   ID: %d\n", adminUser.ID)
	fmt.Printf("   Is Admin: %t\n", adminUser.IsAdmin)
}
