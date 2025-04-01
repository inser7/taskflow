package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"taskflow/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

type Server struct {
	db *gorm.DB
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func main() {
	// Database connection string
	dbURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"))

	// Connect to the database
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Perform migration
	err = models.Migrate(db)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// Initialize the server
	server := &Server{db: db}

	// Create routes
	mux := http.NewServeMux()
	mux.HandleFunc("/api/tasks", server.handleTasks)
	mux.HandleFunc("/api/login", server.handleLogin)

	// Configure CORS
	handler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"}, // Указываем адрес вашего фронтенда
		AllowedHeaders: []string{"Authorization", "Content-Type"},
		AllowedMethods: []string{"GET", "POST", "DELETE", "OPTIONS"},
	}).Handler(mux)

	// Start the server
	log.Println("Server started at http://localhost:8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

// handleTasks processes requests for retrieving, creating, and deleting tasks
func (s *Server) handleTasks(w http.ResponseWriter, r *http.Request) {
	// Check authorization via JWT
	tokenString := r.Header.Get("Authorization")
	if !checkAuthorization(tokenString) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		s.getTasks(w, r)
	case http.MethodPost:
		s.createTask(w, r)
	case http.MethodDelete:
		s.deleteTask(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// getTasks processes requests to retrieve the list of tasks
func (s *Server) getTasks(w http.ResponseWriter, r *http.Request) {
	var tasks []models.Task
	result := s.db.Find(&tasks)

	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(tasks)
}

// createTask processes requests to create a new task
func (s *Server) createTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&task); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	result := s.db.Create(&task)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

// deleteTask processes requests to delete a task
func (s *Server) deleteTask(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query().Get("id")
	if taskID == "" {
		http.Error(w, "Task ID is required", http.StatusBadRequest)
		return
	}

	result := s.db.Delete(&models.Task{}, taskID)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// handleLogin processes requests for authentication and JWT generation
func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var user models.User
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validate user credentials (this is a simple check, can be made more complex)
	if user.Username != "admin" || user.Password != "password" {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(), // Используйте Unix timestamp
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token using the secret key
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	// Send the token to the client
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}

// checkAuthorization verifies if the authorization header is present and valid
func checkAuthorization(tokenString string) bool {
	if tokenString == "" {
		return false
	}

	// Remove the "Bearer " prefix from the token
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return false
	}

	return true
}
