package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"surplus-supper/backend/api/auth"
	"surplus-supper/backend/middleware"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	// Initialize database connection
	db, err := initDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	if db != nil {
		defer db.Close()
	}

	// Initialize router with CORS
	r := mux.NewRouter()

	// Add CORS middleware
	corsMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get CORS origin from environment variable, default to all origins for development
			corsOrigin := os.Getenv("CORS_ORIGIN")
			if corsOrigin == "" {
				corsOrigin = "*"
			}

			w.Header().Set("Access-Control-Allow-Origin", corsOrigin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
	r.Use(corsMiddleware)

	// Initialize auth handler and middleware
	var authHandler *auth.AuthHandler
	var authMiddleware *middleware.AuthMiddleware

	if db != nil {
		authHandler = auth.NewAuthHandler(db)
		authMiddleware = middleware.NewAuthMiddleware()
	} else {
		log.Printf("Warning: Authentication disabled - no database connection")
	}

	// API endpoints
	api := r.PathPrefix("/api").Subrouter()

	// Public endpoints (no authentication required)
	api.HandleFunc("/restaurants", func(w http.ResponseWriter, r *http.Request) {
		restaurants, err := getRestaurants(db, r.URL.Query().Get("location"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(restaurants)
	}).Methods("GET", "OPTIONS")

	api.HandleFunc("/restaurant/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, _ := strconv.Atoi(vars["id"])
		restaurant, err := getRestaurantById(db, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(restaurant)
	}).Methods("GET", "OPTIONS")

	// Authentication endpoints
	if authHandler != nil {
		api.HandleFunc("/auth/register", authHandler.Register).Methods("POST", "OPTIONS")
		api.HandleFunc("/auth/login", authHandler.Login).Methods("POST", "OPTIONS")
		api.HandleFunc("/auth/refresh", authHandler.RefreshToken).Methods("POST", "OPTIONS")

		// Protected endpoints (authentication required)
		protected := api.PathPrefix("/auth").Subrouter()
		protected.Use(authMiddleware.Authenticate)
		protected.HandleFunc("/profile", authHandler.Profile).Methods("GET", "OPTIONS")
		protected.HandleFunc("/profile", authHandler.UpdateProfile).Methods("PUT", "OPTIONS")
	} else {
		// Mock auth endpoints for development
		api.HandleFunc("/auth/register", mockAuthHandler).Methods("POST", "OPTIONS")
		api.HandleFunc("/auth/login", mockAuthHandler).Methods("POST", "OPTIONS")
		api.HandleFunc("/auth/refresh", mockAuthHandler).Methods("POST", "OPTIONS")
		api.HandleFunc("/auth/profile", mockAuthHandler).Methods("GET", "OPTIONS")
	}

	// Health check endpoint
	r.HandleFunc("/health", healthCheckHandler(db)).Methods("GET")

	// WebSocket endpoint for real-time notifications
	r.HandleFunc("/ws", websocketHandler(db))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Add CORS middleware
	handler := corsMiddleware(r)

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

type Restaurant struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Address     string  `json:"address"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Phone       string  `json:"phone"`
	Email       string  `json:"email"`
	CuisineType string  `json:"cuisine_type"`
	Rating      float64 `json:"rating"`
	Distance    float64 `json:"distance"`
}

func getRestaurants(db *sql.DB, location string) ([]Restaurant, error) {
	query := `SELECT id, name, description, address, latitude, longitude, phone, email, cuisine_type, rating 
			 FROM restaurants`

	if location != "" {
		query += ` WHERE LOWER(address) LIKE LOWER($1)`
		rows, err := db.Query(query, "%"+location+"%")
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		return scanRestaurants(rows)
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanRestaurants(rows)
}

func getRestaurantById(db *sql.DB, id int) (*Restaurant, error) {
	var r Restaurant
	err := db.QueryRow(`
		SELECT id, name, description, address, latitude, longitude, phone, email, cuisine_type, rating 
		FROM restaurants WHERE id = $1`, id).Scan(
		&r.ID, &r.Name, &r.Description, &r.Address, &r.Latitude, &r.Longitude,
		&r.Phone, &r.Email, &r.CuisineType, &r.Rating)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func scanRestaurants(rows *sql.Rows) ([]Restaurant, error) {
	var restaurants []Restaurant
	for rows.Next() {
		var r Restaurant
		err := rows.Scan(
			&r.ID, &r.Name, &r.Description, &r.Address, &r.Latitude, &r.Longitude,
			&r.Phone, &r.Email, &r.CuisineType, &r.Rating)
		if err != nil {
			return nil, err
		}
		restaurants = append(restaurants, r)
	}
	return restaurants, nil
}

func initDB() (*sql.DB, error) {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres:password@localhost:5433/surplus_supper?sslmode=disable"
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		// For development, return a mock database if connection fails
		log.Printf("Warning: Database connection failed: %v", err)
		log.Printf("Running in development mode without database")
		return nil, nil
	}

	if err = db.Ping(); err != nil {
		log.Printf("Warning: Database ping failed: %v", err)
		log.Printf("Running in development mode without database")
		return nil, nil
	}

	return db, nil
}

func graphqlHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// GraphQL handler implementation
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"data": {"message": "GraphQL endpoint"}}`))
	}
}

func orderConfirmHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Handle order confirmation
		w.Header().Set("Content-Type", "text/html")
		// TODO: Process order and render confirmation
		w.Write([]byte("<h2>Order Confirmed!</h2>"))
	}
}

func restaurantLoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			// Handle login form submission
			w.Header().Set("Content-Type", "text/html")
			// TODO: Process login and redirect to dashboard
			w.Write([]byte("<h2>Login Successful</h2>"))
		} else {
			// Serve login page
			w.Header().Set("Content-Type", "text/html")
			// TODO: Render login template
			w.Write([]byte("<h2>Restaurant Login</h2><form method='POST'><input type='email' name='email' placeholder='Email'><input type='password' name='password' placeholder='Password'><button type='submit'>Login</button></form>"))
		}
	}
}

func restaurantDashboardHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Serve restaurant dashboard
		w.Header().Set("Content-Type", "text/html")
		// TODO: Render dashboard template
		w.Write([]byte("<h2>Restaurant Dashboard</h2><p>Manage your inventory and orders here.</p>"))
	}
}

func inventoryFormHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			// Handle inventory form submission
			w.Header().Set("Content-Type", "text/html")
			// TODO: Process inventory update
			w.Write([]byte("<h2>Inventory Updated</h2>"))
		} else {
			// Serve inventory form
			w.Header().Set("Content-Type", "text/html")
			// TODO: Render inventory form template
			w.Write([]byte("<h2>Add Inventory Item</h2><form method='POST'><input type='text' name='name' placeholder='Item Name'><input type='number' name='price' placeholder='Price'><button type='submit'>Add Item</button></form>"))
		}
	}
}

func websocketHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// WebSocket handler implementation
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("<h2>WebSocket Endpoint</h2><p>Real-time notifications will be handled here.</p>"))
	}
}

func mockAuthHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Mock auth handler called: %s %s", r.Method, r.URL.Path)

	w.Header().Set("Content-Type", "application/json")

	// Mock successful response for development
	response := map[string]interface{}{
		"token": "mock-jwt-token-for-development",
		"user": map[string]interface{}{
			"id":         1,
			"email":      "test@example.com",
			"first_name": "Test",
			"last_name":  "User",
		},
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func healthCheckHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if db != nil {
			// Test database connection
			err := db.Ping()
			if err != nil {
				http.Error(w, "Database connection failed: "+err.Error(), http.StatusInternalServerError)
				return
			}

			// Test a simple query
			var count int
			err = db.QueryRow("SELECT COUNT(*) FROM restaurants").Scan(&count)
			if err != nil {
				http.Error(w, "Database query failed: "+err.Error(), http.StatusInternalServerError)
				return
			}

			w.Write([]byte(`{"status": "healthy", "restaurant_count": ` + strconv.Itoa(count) + `}`))
		} else {
			// Return mock response when no database
			w.Write([]byte(`{"status": "healthy", "restaurant_count": 0, "mode": "development"}`))
		}
	}
}

// Add CORS middleware
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
