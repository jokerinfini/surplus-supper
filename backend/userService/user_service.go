package userService

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the system
type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateUserInput represents the input for creating a new user
type CreateUserInput struct {
	Email     string  `json:"email"`
	Password  string  `json:"password"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Phone     string  `json:"phone"`
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// UpdateUserInput represents the input for updating a user
type UpdateUserInput struct {
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Phone     string  `json:"phone"`
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// LoginInput represents the input for user login
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserService handles user-related operations
type UserService struct {
	db *sql.DB
}

// NewUserService creates a new user service
func NewUserService(db *sql.DB) *UserService {
	return &UserService{db: db}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(input CreateUserInput) (*User, error) {
	// Check if user already exists
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", input.Email).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}
	if exists {
		return nil, errors.New("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Insert user
	var user User
	err = s.db.QueryRow(`
		INSERT INTO users (email, password_hash, first_name, last_name, phone, address, latitude, longitude)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, email, first_name, last_name, phone, address, latitude, longitude, created_at, updated_at
	`, input.Email, string(hashedPassword), input.FirstName, input.LastName, input.Phone, input.Address, input.Latitude, input.Longitude).Scan(
		&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Phone, &user.Address, &user.Latitude, &user.Longitude, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(id int) (*User, error) {
	var user User
	err := s.db.QueryRow(`
		SELECT id, email, first_name, last_name, phone, address, latitude, longitude, created_at, updated_at
		FROM users WHERE id = $1
	`, id).Scan(
		&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Phone, &user.Address, &user.Latitude, &user.Longitude, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(email string) (*User, error) {
	var user User
	err := s.db.QueryRow(`
		SELECT id, email, first_name, last_name, phone, address, latitude, longitude, created_at, updated_at
		FROM users WHERE email = $1
	`, email).Scan(
		&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Phone, &user.Address, &user.Latitude, &user.Longitude, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// UpdateUser updates a user's information
func (s *UserService) UpdateUser(id int, input UpdateUserInput) (*User, error) {
	// Build dynamic query based on provided fields
	query := `
		UPDATE users SET 
		first_name = COALESCE($2, first_name),
		last_name = COALESCE($3, last_name),
		phone = COALESCE($4, phone),
		address = COALESCE($5, address),
		latitude = COALESCE($6, latitude),
		longitude = COALESCE($7, longitude),
		updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING id, email, first_name, last_name, phone, address, latitude, longitude, created_at, updated_at
	`

	var user User
	err := s.db.QueryRow(query, id, input.FirstName, input.LastName, input.Phone, input.Address, input.Latitude, input.Longitude).Scan(
		&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Phone, &user.Address, &user.Latitude, &user.Longitude, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &user, nil
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(id int) error {
	result, err := s.db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

// AuthenticateUser authenticates a user with email and password
func (s *UserService) AuthenticateUser(input LoginInput) (*User, error) {
	var user User
	var passwordHash string

	err := s.db.QueryRow(`
		SELECT id, email, password_hash, first_name, last_name, phone, address, latitude, longitude, created_at, updated_at
		FROM users WHERE email = $1
	`, input.Email).Scan(
		&user.ID, &user.Email, &passwordHash, &user.FirstName, &user.LastName, &user.Phone, &user.Address, &user.Latitude, &user.Longitude, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invalid email or password")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(input.Password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	return &user, nil
}

// GetAllUsers retrieves all users (for admin purposes)
func (s *UserService) GetAllUsers() ([]*User, error) {
	rows, err := s.db.Query(`
		SELECT id, email, first_name, last_name, phone, address, latitude, longitude, created_at, updated_at
		FROM users ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Phone, &user.Address, &user.Latitude, &user.Longitude, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, &user)
	}

	return users, nil
} 