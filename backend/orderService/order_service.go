package orderService

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// Order represents an order in the system
type Order struct {
	ID                int       `json:"id"`
	UserID            int       `json:"user_id"`
	RestaurantID      int       `json:"restaurant_id"`
	TotalAmount       float64   `json:"total_amount"`
	Status            string    `json:"status"`
	PickupTime        time.Time `json:"pickup_time"`
	SpecialInstructions string  `json:"special_instructions"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	ID              int     `json:"id"`
	OrderID         int     `json:"order_id"`
	InventoryItemID int     `json:"inventory_item_id"`
	OfferID         int     `json:"offer_id"`
	Quantity        int     `json:"quantity"`
	UnitPrice       float64 `json:"unit_price"`
	TotalPrice      float64 `json:"total_price"`
	CreatedAt       time.Time `json:"created_at"`
}

// CartItem represents an item in the shopping cart
type CartItem struct {
	ID              int     `json:"id"`
	InventoryItemID int     `json:"inventory_item_id"`
	OfferID         int     `json:"offer_id"`
	Quantity        int     `json:"quantity"`
	UnitPrice       float64 `json:"unit_price"`
	TotalPrice      float64 `json:"total_price"`
}

// CreateOrderInput represents the input for creating a new order
type CreateOrderInput struct {
	UserID              int           `json:"user_id"`
	RestaurantID        int           `json:"restaurant_id"`
	OrderItems          []OrderItemInput `json:"order_items"`
	SpecialInstructions string        `json:"special_instructions"`
}

// OrderItemInput represents the input for an order item
type OrderItemInput struct {
	InventoryItemID int `json:"inventory_item_id"`
	OfferID         int `json:"offer_id"`
	Quantity        int `json:"quantity"`
}

// PaymentInput represents payment information
type PaymentInput struct {
	OrderID     int     `json:"order_id"`
	Amount      float64 `json:"amount"`
	PaymentMethod string `json:"payment_method"`
	StripeToken string `json:"stripe_token"`
}

// OrderService handles order-related operations
type OrderService struct {
	db *sql.DB
}

// NewOrderService creates a new order service
func NewOrderService(db *sql.DB) *OrderService {
	return &OrderService{db: db}
}

// CreateOrder creates a new order
func (s *OrderService) CreateOrder(input CreateOrderInput) (*Order, error) {
	// Start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Calculate total amount
	var totalAmount float64
	for _, item := range input.OrderItems {
		if item.InventoryItemID > 0 {
			// Get inventory item price
			var price float64
			err := tx.QueryRow("SELECT surplus_price FROM inventory_items WHERE id = $1 AND is_available = true", item.InventoryItemID).Scan(&price)
			if err != nil {
				return nil, fmt.Errorf("failed to get inventory item price: %w", err)
			}
			totalAmount += price * float64(item.Quantity)
		} else if item.OfferID > 0 {
			// Get offer price
			var price float64
			err := tx.QueryRow("SELECT surplus_price FROM offers WHERE id = $1 AND is_available = true", item.OfferID).Scan(&price)
			if err != nil {
				return nil, fmt.Errorf("failed to get offer price: %w", err)
			}
			totalAmount += price * float64(item.Quantity)
		}
	}

	// Create order
	var order Order
	err = tx.QueryRow(`
		INSERT INTO orders (user_id, restaurant_id, total_amount, status, special_instructions)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, restaurant_id, total_amount, status, pickup_time, special_instructions, created_at, updated_at
	`, input.UserID, input.RestaurantID, totalAmount, "pending", input.SpecialInstructions).Scan(
		&order.ID, &order.UserID, &order.RestaurantID, &order.TotalAmount, &order.Status, &order.PickupTime, &order.SpecialInstructions, &order.CreatedAt, &order.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Create order items
	for _, item := range input.OrderItems {
		var unitPrice float64
		if item.InventoryItemID > 0 {
			err := tx.QueryRow("SELECT surplus_price FROM inventory_items WHERE id = $1", item.InventoryItemID).Scan(&unitPrice)
			if err != nil {
				return nil, fmt.Errorf("failed to get inventory item price: %w", err)
			}
		} else if item.OfferID > 0 {
			err := tx.QueryRow("SELECT surplus_price FROM offers WHERE id = $1", item.OfferID).Scan(&unitPrice)
			if err != nil {
				return nil, fmt.Errorf("failed to get offer price: %w", err)
			}
		}

		totalPrice := unitPrice * float64(item.Quantity)

		_, err = tx.Exec(`
			INSERT INTO order_items (order_id, inventory_item_id, offer_id, quantity, unit_price, total_price)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, order.ID, item.InventoryItemID, item.OfferID, item.Quantity, unitPrice, totalPrice)
		if err != nil {
			return nil, fmt.Errorf("failed to create order item: %w", err)
		}

		// Update inventory quantity
		if item.InventoryItemID > 0 {
			_, err = tx.Exec("UPDATE inventory_items SET quantity = quantity - $1 WHERE id = $2", item.Quantity, item.InventoryItemID)
			if err != nil {
				return nil, fmt.Errorf("failed to update inventory quantity: %w", err)
			}
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &order, nil
}

// GetOrderByID retrieves an order by ID
func (s *OrderService) GetOrderByID(id int) (*Order, error) {
	var order Order
	err := s.db.QueryRow(`
		SELECT id, user_id, restaurant_id, total_amount, status, pickup_time, special_instructions, created_at, updated_at
		FROM orders WHERE id = $1
	`, id).Scan(
		&order.ID, &order.UserID, &order.RestaurantID, &order.TotalAmount, &order.Status, &order.PickupTime, &order.SpecialInstructions, &order.CreatedAt, &order.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("order not found")
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	return &order, nil
}

// GetOrderItems retrieves items for an order
func (s *OrderService) GetOrderItems(orderID int) ([]*OrderItem, error) {
	rows, err := s.db.Query(`
		SELECT id, order_id, inventory_item_id, offer_id, quantity, unit_price, total_price, created_at
		FROM order_items WHERE order_id = $1
	`, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order items: %w", err)
	}
	defer rows.Close()

	var items []*OrderItem
	for rows.Next() {
		var item OrderItem
		err := rows.Scan(
			&item.ID, &item.OrderID, &item.InventoryItemID, &item.OfferID, &item.Quantity, &item.UnitPrice, &item.TotalPrice, &item.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order item: %w", err)
		}
		items = append(items, &item)
	}

	return items, nil
}

// UpdateOrderStatus updates the status of an order
func (s *OrderService) UpdateOrderStatus(id int, status string) (*Order, error) {
	var order Order
	err := s.db.QueryRow(`
		UPDATE orders SET status = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING id, user_id, restaurant_id, total_amount, status, pickup_time, special_instructions, created_at, updated_at
	`, id, status).Scan(
		&order.ID, &order.UserID, &order.RestaurantID, &order.TotalAmount, &order.Status, &order.PickupTime, &order.SpecialInstructions, &order.CreatedAt, &order.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("order not found")
		}
		return nil, fmt.Errorf("failed to update order status: %w", err)
	}

	return &order, nil
}

// GetUserOrders retrieves all orders for a user
func (s *OrderService) GetUserOrders(userID int) ([]*Order, error) {
	rows, err := s.db.Query(`
		SELECT id, user_id, restaurant_id, total_amount, status, pickup_time, special_instructions, created_at, updated_at
		FROM orders WHERE user_id = $1 ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user orders: %w", err)
	}
	defer rows.Close()

	var orders []*Order
	for rows.Next() {
		var order Order
		err := rows.Scan(
			&order.ID, &order.UserID, &order.RestaurantID, &order.TotalAmount, &order.Status, &order.PickupTime, &order.SpecialInstructions, &order.CreatedAt, &order.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, &order)
	}

	return orders, nil
}

// GetRestaurantOrders retrieves all orders for a restaurant
func (s *OrderService) GetRestaurantOrders(restaurantID int, status string) ([]*Order, error) {
	var query string
	var args []interface{}

	if status != "" {
		query = `
			SELECT id, user_id, restaurant_id, total_amount, status, pickup_time, special_instructions, created_at, updated_at
			FROM orders WHERE restaurant_id = $1 AND status = $2 ORDER BY created_at DESC
		`
		args = []interface{}{restaurantID, status}
	} else {
		query = `
			SELECT id, user_id, restaurant_id, total_amount, status, pickup_time, special_instructions, created_at, updated_at
			FROM orders WHERE restaurant_id = $1 ORDER BY created_at DESC
		`
		args = []interface{}{restaurantID}
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get restaurant orders: %w", err)
	}
	defer rows.Close()

	var orders []*Order
	for rows.Next() {
		var order Order
		err := rows.Scan(
			&order.ID, &order.UserID, &order.RestaurantID, &order.TotalAmount, &order.Status, &order.PickupTime, &order.SpecialInstructions, &order.CreatedAt, &order.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, &order)
	}

	return orders, nil
}

// ProcessPayment processes payment for an order (placeholder for Stripe integration)
func (s *OrderService) ProcessPayment(input PaymentInput) error {
	// This is a placeholder for Stripe payment processing
	// In a real implementation, this would integrate with Stripe API
	
	// For now, just update the order status to "paid"
	_, err := s.db.Exec("UPDATE orders SET status = 'paid', updated_at = CURRENT_TIMESTAMP WHERE id = $1", input.OrderID)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	return nil
}

// CancelOrder cancels an order
func (s *OrderService) CancelOrder(id int) (*Order, error) {
	// Start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Get order items to restore inventory
	rows, err := tx.Query("SELECT inventory_item_id, quantity FROM order_items WHERE order_id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("failed to get order items: %w", err)
	}
	defer rows.Close()

	// Restore inventory quantities
	for rows.Next() {
		var inventoryItemID, quantity int
		err := rows.Scan(&inventoryItemID, &quantity)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order item: %w", err)
		}
		if inventoryItemID > 0 {
			_, err = tx.Exec("UPDATE inventory_items SET quantity = quantity + $1 WHERE id = $2", quantity, inventoryItemID)
			if err != nil {
				return nil, fmt.Errorf("failed to restore inventory quantity: %w", err)
			}
		}
	}

	// Update order status
	var order Order
	err = tx.QueryRow(`
		UPDATE orders SET status = 'cancelled', updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING id, user_id, restaurant_id, total_amount, status, pickup_time, special_instructions, created_at, updated_at
	`, id).Scan(
		&order.ID, &order.UserID, &order.RestaurantID, &order.TotalAmount, &order.Status, &order.PickupTime, &order.SpecialInstructions, &order.CreatedAt, &order.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("order not found")
		}
		return nil, fmt.Errorf("failed to cancel order: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &order, nil
} 