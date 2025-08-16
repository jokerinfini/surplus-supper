package restaurantService

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"time"
)

// Restaurant represents a restaurant in the system
type Restaurant struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Address     string    `json:"address"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	Phone       string    `json:"phone"`
	Email       string    `json:"email"`
	CuisineType string    `json:"cuisine_type"`
	Rating      float64   `json:"rating"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// InventoryItem represents an inventory item in a restaurant
type InventoryItem struct {
	ID           int       `json:"id"`
	RestaurantID int       `json:"restaurant_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	OriginalPrice float64  `json:"original_price"`
	SurplusPrice float64   `json:"surplus_price"`
	Quantity     int       `json:"quantity"`
	Category     string    `json:"category"`
	ExpiryTime   time.Time `json:"expiry_time"`
	IsAvailable  bool      `json:"is_available"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Offer represents a special offer (Surprise Bag or Chef's Surprise)
type Offer struct {
	ID           int       `json:"id"`
	RestaurantID int       `json:"restaurant_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	OriginalPrice float64  `json:"original_price"`
	SurplusPrice float64   `json:"surplus_price"`
	OfferType    string    `json:"offer_type"`
	Ingredients  string    `json:"ingredients"`
	IsAvailable  bool      `json:"is_available"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CreateRestaurantInput represents the input for creating a new restaurant
type CreateRestaurantInput struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Address     string  `json:"address"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Phone       string  `json:"phone"`
	Email       string  `json:"email"`
	CuisineType string  `json:"cuisine_type"`
}

// UpdateRestaurantInput represents the input for updating a restaurant
type UpdateRestaurantInput struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Address     string  `json:"address"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Phone       string  `json:"phone"`
	Email       string  `json:"email"`
	CuisineType string  `json:"cuisine_type"`
	IsActive    bool    `json:"is_active"`
}

// CreateInventoryItemInput represents the input for creating a new inventory item
type CreateInventoryItemInput struct {
	RestaurantID  int       `json:"restaurant_id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	OriginalPrice float64   `json:"original_price"`
	SurplusPrice  float64   `json:"surplus_price"`
	Quantity      int       `json:"quantity"`
	Category      string    `json:"category"`
	ExpiryTime    time.Time `json:"expiry_time"`
}

// RestaurantService handles restaurant-related operations
type RestaurantService struct {
	db *sql.DB
}

// NewRestaurantService creates a new restaurant service
func NewRestaurantService(db *sql.DB) *RestaurantService {
	return &RestaurantService{db: db}
}

// CreateRestaurant creates a new restaurant
func (s *RestaurantService) CreateRestaurant(input CreateRestaurantInput) (*Restaurant, error) {
	var restaurant Restaurant
	err := s.db.QueryRow(`
		INSERT INTO restaurants (name, description, address, latitude, longitude, phone, email, cuisine_type)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, name, description, address, latitude, longitude, phone, email, cuisine_type, rating, is_active, created_at, updated_at
	`, input.Name, input.Description, input.Address, input.Latitude, input.Longitude, input.Phone, input.Email, input.CuisineType).Scan(
		&restaurant.ID, &restaurant.Name, &restaurant.Description, &restaurant.Address, &restaurant.Latitude, &restaurant.Longitude, &restaurant.Phone, &restaurant.Email, &restaurant.CuisineType, &restaurant.Rating, &restaurant.IsActive, &restaurant.CreatedAt, &restaurant.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create restaurant: %w", err)
	}

	return &restaurant, nil
}

// GetRestaurantByID retrieves a restaurant by ID
func (s *RestaurantService) GetRestaurantByID(id int) (*Restaurant, error) {
	var restaurant Restaurant
	err := s.db.QueryRow(`
		SELECT id, name, description, address, latitude, longitude, phone, email, cuisine_type, rating, is_active, created_at, updated_at
		FROM restaurants WHERE id = $1
	`, id).Scan(
		&restaurant.ID, &restaurant.Name, &restaurant.Description, &restaurant.Address, &restaurant.Latitude, &restaurant.Longitude, &restaurant.Phone, &restaurant.Email, &restaurant.CuisineType, &restaurant.Rating, &restaurant.IsActive, &restaurant.CreatedAt, &restaurant.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("restaurant not found")
		}
		return nil, fmt.Errorf("failed to get restaurant: %w", err)
	}

	return &restaurant, nil
}

// GetNearbyRestaurants retrieves restaurants within a specified radius
func (s *RestaurantService) GetNearbyRestaurants(latitude, longitude, radius float64) ([]*Restaurant, error) {
	// Using Haversine formula to calculate distance
	query := `
		SELECT id, name, description, address, latitude, longitude, phone, email, cuisine_type, rating, is_active, created_at, updated_at
		FROM restaurants 
		WHERE is_active = true 
		AND (
			6371 * acos(
				cos(radians($1)) * cos(radians(latitude)) * cos(radians(longitude) - radians($2)) + 
				sin(radians($1)) * sin(radians(latitude))
			) <= $3
		ORDER BY (
			6371 * acos(
				cos(radians($1)) * cos(radians(latitude)) * cos(radians(longitude) - radians($2)) + 
				sin(radians($1)) * sin(radians(latitude))
			)
		)
	`

	rows, err := s.db.Query(query, latitude, longitude, radius)
	if err != nil {
		return nil, fmt.Errorf("failed to get nearby restaurants: %w", err)
	}
	defer rows.Close()

	var restaurants []*Restaurant
	for rows.Next() {
		var restaurant Restaurant
		err := rows.Scan(
			&restaurant.ID, &restaurant.Name, &restaurant.Description, &restaurant.Address, &restaurant.Latitude, &restaurant.Longitude, &restaurant.Phone, &restaurant.Email, &restaurant.CuisineType, &restaurant.Rating, &restaurant.IsActive, &restaurant.CreatedAt, &restaurant.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan restaurant: %w", err)
		}
		restaurants = append(restaurants, &restaurant)
	}

	return restaurants, nil
}

// UpdateRestaurant updates a restaurant's information
func (s *RestaurantService) UpdateRestaurant(id int, input UpdateRestaurantInput) (*Restaurant, error) {
	query := `
		UPDATE restaurants SET 
		name = COALESCE($2, name),
		description = COALESCE($3, description),
		address = COALESCE($4, address),
		latitude = COALESCE($5, latitude),
		longitude = COALESCE($6, longitude),
		phone = COALESCE($7, phone),
		email = COALESCE($8, email),
		cuisine_type = COALESCE($9, cuisine_type),
		is_active = COALESCE($10, is_active),
		updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING id, name, description, address, latitude, longitude, phone, email, cuisine_type, rating, is_active, created_at, updated_at
	`

	var restaurant Restaurant
	err := s.db.QueryRow(query, id, input.Name, input.Description, input.Address, input.Latitude, input.Longitude, input.Phone, input.Email, input.CuisineType, input.IsActive).Scan(
		&restaurant.ID, &restaurant.Name, &restaurant.Description, &restaurant.Address, &restaurant.Latitude, &restaurant.Longitude, &restaurant.Phone, &restaurant.Email, &restaurant.CuisineType, &restaurant.Rating, &restaurant.IsActive, &restaurant.CreatedAt, &restaurant.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("restaurant not found")
		}
		return nil, fmt.Errorf("failed to update restaurant: %w", err)
	}

	return &restaurant, nil
}

// CreateInventoryItem creates a new inventory item
func (s *RestaurantService) CreateInventoryItem(input CreateInventoryItemInput) (*InventoryItem, error) {
	var item InventoryItem
	err := s.db.QueryRow(`
		INSERT INTO inventory_items (restaurant_id, name, description, original_price, surplus_price, quantity, category, expiry_time)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, restaurant_id, name, description, original_price, surplus_price, quantity, category, expiry_time, is_available, created_at, updated_at
	`, input.RestaurantID, input.Name, input.Description, input.OriginalPrice, input.SurplusPrice, input.Quantity, input.Category, input.ExpiryTime).Scan(
		&item.ID, &item.RestaurantID, &item.Name, &item.Description, &item.OriginalPrice, &item.SurplusPrice, &item.Quantity, &item.Category, &item.ExpiryTime, &item.IsAvailable, &item.CreatedAt, &item.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create inventory item: %w", err)
	}

	return &item, nil
}

// GetInventoryItems retrieves inventory items for a restaurant
func (s *RestaurantService) GetInventoryItems(restaurantID int, availableOnly bool) ([]*InventoryItem, error) {
	var query string
	var args []interface{}

	if availableOnly {
		query = `
			SELECT id, restaurant_id, name, description, original_price, surplus_price, quantity, category, expiry_time, is_available, created_at, updated_at
			FROM inventory_items 
			WHERE restaurant_id = $1 AND is_available = true AND expiry_time > NOW()
			ORDER BY created_at DESC
		`
		args = []interface{}{restaurantID}
	} else {
		query = `
			SELECT id, restaurant_id, name, description, original_price, surplus_price, quantity, category, expiry_time, is_available, created_at, updated_at
			FROM inventory_items 
			WHERE restaurant_id = $1
			ORDER BY created_at DESC
		`
		args = []interface{}{restaurantID}
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory items: %w", err)
	}
	defer rows.Close()

	var items []*InventoryItem
	for rows.Next() {
		var item InventoryItem
		err := rows.Scan(
			&item.ID, &item.RestaurantID, &item.Name, &item.Description, &item.OriginalPrice, &item.SurplusPrice, &item.Quantity, &item.Category, &item.ExpiryTime, &item.IsAvailable, &item.CreatedAt, &item.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan inventory item: %w", err)
		}
		items = append(items, &item)
	}

	return items, nil
}

// UpdateInventoryItem updates an inventory item
func (s *RestaurantService) UpdateInventoryItem(id int, updates map[string]interface{}) (*InventoryItem, error) {
	// Build dynamic query based on provided updates
	query := `
		UPDATE inventory_items SET 
		name = COALESCE($2, name),
		description = COALESCE($3, description),
		original_price = COALESCE($4, original_price),
		surplus_price = COALESCE($5, surplus_price),
		quantity = COALESCE($6, quantity),
		category = COALESCE($7, category),
		expiry_time = COALESCE($8, expiry_time),
		is_available = COALESCE($9, is_available),
		updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING id, restaurant_id, name, description, original_price, surplus_price, quantity, category, expiry_time, is_available, created_at, updated_at
	`

	var item InventoryItem
	err := s.db.QueryRow(query, id, updates["name"], updates["description"], updates["original_price"], updates["surplus_price"], updates["quantity"], updates["category"], updates["expiry_time"], updates["is_available"]).Scan(
		&item.ID, &item.RestaurantID, &item.Name, &item.Description, &item.OriginalPrice, &item.SurplusPrice, &item.Quantity, &item.Category, &item.ExpiryTime, &item.IsAvailable, &item.CreatedAt, &item.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("inventory item not found")
		}
		return nil, fmt.Errorf("failed to update inventory item: %w", err)
	}

	return &item, nil
}

// GetOffers retrieves offers for a restaurant
func (s *RestaurantService) GetOffers(restaurantID int, availableOnly bool) ([]*Offer, error) {
	var query string
	var args []interface{}

	if availableOnly {
		query = `
			SELECT id, restaurant_id, name, description, original_price, surplus_price, offer_type, ingredients, is_available, created_at, updated_at
			FROM offers 
			WHERE restaurant_id = $1 AND is_available = true
			ORDER BY created_at DESC
		`
		args = []interface{}{restaurantID}
	} else {
		query = `
			SELECT id, restaurant_id, name, description, original_price, surplus_price, offer_type, ingredients, is_available, created_at, updated_at
			FROM offers 
			WHERE restaurant_id = $1
			ORDER BY created_at DESC
		`
		args = []interface{}{restaurantID}
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get offers: %w", err)
	}
	defer rows.Close()

	var offers []*Offer
	for rows.Next() {
		var offer Offer
		err := rows.Scan(
			&offer.ID, &offer.RestaurantID, &offer.Name, &offer.Description, &offer.OriginalPrice, &offer.SurplusPrice, &offer.OfferType, &offer.Ingredients, &offer.IsAvailable, &offer.CreatedAt, &offer.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan offer: %w", err)
		}
		offers = append(offers, &offer)
	}

	return offers, nil
}

// CalculateDistance calculates the distance between two points using Haversine formula
func CalculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Earth's radius in kilometers

	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180

	dlat := lat2Rad - lat1Rad
	dlon := lon2Rad - lon1Rad

	a := math.Sin(dlat/2)*math.Sin(dlat/2) + math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dlon/2)*math.Sin(dlon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
} 