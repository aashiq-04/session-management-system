package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/aashiq-04/session-management-system/backend/services/session-service/internal/handlers"
	pb "github.com/aashiq-04/session-management-system/backend/services/session-service/proto"
)

func main() {
	log.Println("Starting Session Service...")
	godotenv.Load()

	// Load configuration from environment variables
	config := loadConfig()

	// Connect to database
	db, err := connectDatabase(config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Successfully connected to database")

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Register session service
	sessionHandler := handlers.NewSessionHandler(db)
	pb.RegisterSessionServiceServer(grpcServer, sessionHandler)

	// Enable reflection for grpcurl/grpc-ui
	reflection.Register(grpcServer)

	// Start listening
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", config.GRPCPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Session Service listening on port %s", config.GRPCPort)

	// Handle graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan

		log.Println("Shutting down Session Service...")
		grpcServer.GracefulStop()
		log.Println("Session Service stopped")
	}()

	// Start serving
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// Config holds the application configuration
type Config struct {
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	GRPCPort       string
	AuthServiceURL string
}

// loadConfig loads configuration from environment variables
func loadConfig() Config {
	config := Config{
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBUser:         getEnv("DB_USER", "admin"),
		DBPassword:     getEnv("DB_PASSWORD", "admin123"),
		DBName:         getEnv("DB_NAME", "session_management"),
		GRPCPort:       getEnv("GRPC_PORT", "50052"),
		AuthServiceURL: getEnv("AUTH_SERVICE_URL", "localhost:50051"),
	}

	return config
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// connectDatabase establishes a connection to PostgreSQL
func connectDatabase(config Config) (*sql.DB, error) {
	// Build connection string
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost,
		config.DBPort,
		config.DBUser,
		config.DBPassword,
		config.DBName,
	)

	// Open database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test connection with retries
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		err = db.Ping()
		if err == nil {
			return db, nil
		}

		log.Printf("Failed to ping database (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("failed to connect to database after %d attempts", maxRetries)
}