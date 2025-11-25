package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"context"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/rs/cors"
	"github.com/joho/godotenv"
	"github.com/aashiq-04/session-management-system/backend/gateway/clients"
	"github.com/aashiq-04/session-management-system/backend/gateway/graph"
	"github.com/aashiq-04/session-management-system/backend/gateway/graph/generated"
	"github.com/aashiq-04/session-management-system/backend/gateway/middleware"
)

func main() {
	log.Println("Starting GraphQL Gateway...")

	godotenv.Load()

	// Load configuration
	config := loadConfig()

	// Initialize gRPC clients
	grpcClients, err := clients.NewGRPCClients(
		config.AuthServiceURL,
		config.SessionServiceURL,
		config.AuditServiceURL,
	)
	if err != nil {
		log.Fatalf("Failed to initialize gRPC clients: %v", err)
	}

	// Create resolver
	resolver := graph.NewResolver(grpcClients, config.JWTSecret)

	// Create GraphQL server
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: resolver,
	}))

	// Setup CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:3001"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	// Setup routes
	mux := http.NewServeMux()

	// GraphQL endpoint with auth middleware
	// mux.Handle("/graphql", middleware.AuthMiddleware(config.JWTSecret)(srv))
	mux.Handle("/graphql", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "httpRequest", r)
		middleware.AuthMiddleware(config.JWTSecret)(srv).ServeHTTP(w, r.WithContext(ctx))
	}))
	

	// GraphQL playground (development only)
	if config.Environment == "development" {
		mux.Handle("/playground", playground.Handler("GraphQL Playground", "/graphql"))
		log.Println("GraphQL Playground available at http://localhost:" + config.Port + "/playground")
	}

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Wrap with CORS
	handler := corsHandler.Handler(mux)

	// Start server
	log.Printf("GraphQL Gateway listening on port %s", config.Port)
	log.Printf("GraphQL endpoint: http://localhost:%s/graphql", config.Port)

	if err := http.ListenAndServe(":"+config.Port, handler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// Config holds application configuration
type Config struct {
	Port              string
	Environment       string
	AuthServiceURL    string
	SessionServiceURL string
	AuditServiceURL   string
	JWTSecret         string
}

// loadConfig loads configuration from environment variables
func loadConfig() Config {
	config := Config{
		Port:              getEnv("PORT", "8080"),
		Environment:       getEnv("ENVIRONMENT", "development"),
		AuthServiceURL:    getEnv("AUTH_SERVICE_URL", "localhost:50051"),
		SessionServiceURL: getEnv("SESSION_SERVICE_URL", "localhost:50052"),
		AuditServiceURL:   getEnv("AUDIT_SERVICE_URL", "localhost:50053"),
		JWTSecret:         getEnv("JWT_SECRET", ""),
	}

	// Validate required config
	if config.JWTSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
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
// Extract real IP from headers or remote address
func getRealIP(r *http.Request) string {
    // Check X-Forwarded-For first (common proxy header)
    forwarded := r.Header.Get("X-Forwarded-For")
    if forwarded != "" {
        return forwarded
    }

    // Check X-Real-IP (nginx header)
    realIP := r.Header.Get("X-Real-IP")
    if realIP != "" {
        return realIP
    }

    // Fall back to remote address
    ip, _, err := net.SplitHostPort(r.RemoteAddr)
    if err == nil {
        return ip
    }

    return "0.0.0.0"
}
