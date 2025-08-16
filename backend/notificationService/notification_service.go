package notificationService

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"
)

// Notification represents a notification in the system
type Notification struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	RestaurantID int    `json:"restaurant_id"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Type      string    `json:"type"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

// NotificationService handles notification-related operations
type NotificationService struct {
	db *sql.DB
	clients map[int]*Client
	mutex   sync.RWMutex
}

// Client represents a WebSocket client
type Client struct {
	ID       int
	UserID   int
	Send     chan []byte
	Hub      *NotificationService
}

// NewNotificationService creates a new notification service
func NewNotificationService(db *sql.DB) *NotificationService {
	return &NotificationService{
		db:      db,
		clients: make(map[int]*Client),
	}
}

// CreateNotification creates a new notification
func (s *NotificationService) CreateNotification(userID, restaurantID int, title, message, notificationType string) (*Notification, error) {
	var notification Notification
	err := s.db.QueryRow(`
		INSERT INTO notifications (user_id, restaurant_id, title, message, type)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, restaurant_id, title, message, type, is_read, created_at
	`, userID, restaurantID, title, message, notificationType).Scan(
		&notification.ID, &notification.UserID, &notification.RestaurantID, &notification.Title, &notification.Message, &notification.Type, &notification.IsRead, &notification.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	// Send real-time notification to connected clients
	s.sendToUser(userID, notification)

	return &notification, nil
}

// GetUserNotifications retrieves notifications for a user
func (s *NotificationService) GetUserNotifications(userID int, unreadOnly bool) ([]*Notification, error) {
	var query string
	var args []interface{}

	if unreadOnly {
		query = `
			SELECT id, user_id, restaurant_id, title, message, type, is_read, created_at
			FROM notifications WHERE user_id = $1 AND is_read = false
			ORDER BY created_at DESC
		`
		args = []interface{}{userID}
	} else {
		query = `
			SELECT id, user_id, restaurant_id, title, message, type, is_read, created_at
			FROM notifications WHERE user_id = $1
			ORDER BY created_at DESC
		`
		args = []interface{}{userID}
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get user notifications: %w", err)
	}
	defer rows.Close()

	var notifications []*Notification
	for rows.Next() {
		var notification Notification
		err := rows.Scan(
			&notification.ID, &notification.UserID, &notification.RestaurantID, &notification.Title, &notification.Message, &notification.Type, &notification.IsRead, &notification.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notification: %w", err)
		}
		notifications = append(notifications, &notification)
	}

	return notifications, nil
}

// MarkNotificationAsRead marks a notification as read
func (s *NotificationService) MarkNotificationAsRead(id int) (*Notification, error) {
	var notification Notification
	err := s.db.QueryRow(`
		UPDATE notifications SET is_read = true WHERE id = $1
		RETURNING id, user_id, restaurant_id, title, message, type, is_read, created_at
	`, id).Scan(
		&notification.ID, &notification.UserID, &notification.RestaurantID, &notification.Title, &notification.Message, &notification.Type, &notification.IsRead, &notification.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("notification not found")
		}
		return nil, fmt.Errorf("failed to mark notification as read: %w", err)
	}

	return &notification, nil
}

// DeleteNotification deletes a notification
func (s *NotificationService) DeleteNotification(id int) error {
	result, err := s.db.Exec("DELETE FROM notifications WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("notification not found")
	}

	return nil
}

// RegisterClient registers a new WebSocket client
func (s *NotificationService) RegisterClient(client *Client) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.clients[client.ID] = client
	log.Printf("Client %d registered", client.ID)
}

// UnregisterClient unregisters a WebSocket client
func (s *NotificationService) UnregisterClient(client *Client) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if _, ok := s.clients[client.ID]; ok {
		delete(s.clients, client.ID)
		close(client.Send)
		log.Printf("Client %d unregistered", client.ID)
	}
}

// sendToUser sends a notification to a specific user
func (s *NotificationService) sendToUser(userID int, notification Notification) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	notificationJSON, err := json.Marshal(notification)
	if err != nil {
		log.Printf("Failed to marshal notification: %v", err)
		return
	}

	// Send to all clients of this user
	for _, client := range s.clients {
		if client.UserID == userID {
			select {
			case client.Send <- notificationJSON:
			default:
				// Client buffer is full, skip this notification
				log.Printf("Client %d buffer full, skipping notification", client.ID)
			}
		}
	}
}

// BroadcastToRestaurant sends a notification to all users of a restaurant
func (s *NotificationService) BroadcastToRestaurant(restaurantID int, title, message, notificationType string) error {
	// Get all users who have ordered from this restaurant
	rows, err := s.db.Query(`
		SELECT DISTINCT user_id FROM orders WHERE restaurant_id = $1 AND user_id IS NOT NULL
	`, restaurantID)
	if err != nil {
		return fmt.Errorf("failed to get restaurant users: %w", err)
	}
	defer rows.Close()

	var userIDs []int
	for rows.Next() {
		var userID int
		err := rows.Scan(&userID)
		if err != nil {
			return fmt.Errorf("failed to scan user ID: %w", err)
		}
		userIDs = append(userIDs, userID)
	}

	// Create notifications for all users
	for _, userID := range userIDs {
		_, err := s.CreateNotification(userID, restaurantID, title, message, notificationType)
		if err != nil {
			log.Printf("Failed to create notification for user %d: %v", userID, err)
		}
	}

	return nil
}

// SendOrderNotification sends a notification about an order
func (s *NotificationService) SendOrderNotification(orderID int, message string) error {
	// Get order details
	var userID, restaurantID int
	err := s.db.QueryRow("SELECT user_id, restaurant_id FROM orders WHERE id = $1", orderID).Scan(&userID, &restaurantID)
	if err != nil {
		return fmt.Errorf("failed to get order details: %w", err)
	}

	// Send notification to user
	if userID > 0 {
		_, err = s.CreateNotification(userID, restaurantID, "Order Update", message, "order_update")
		if err != nil {
			log.Printf("Failed to send order notification to user: %v", err)
		}
	}

	// Send notification to restaurant
	_, err = s.CreateNotification(0, restaurantID, "New Order", fmt.Sprintf("New order #%d received", orderID), "order_update")
	if err != nil {
		log.Printf("Failed to send order notification to restaurant: %v", err)
	}

	return nil
}

// SendOfferNotification sends a notification about a new offer
func (s *NotificationService) SendOfferNotification(restaurantID int, offerName string) error {
	title := "New Surplus Offer Available"
	message := fmt.Sprintf("A new offer '%s' is now available at a restaurant near you!", offerName)

	return s.BroadcastToRestaurant(restaurantID, title, message, "new_offer")
}

// GetUnreadCount gets the count of unread notifications for a user
func (s *NotificationService) GetUnreadCount(userID int) (int, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND is_read = false", userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get unread count: %w", err)
	}

	return count, nil
} 