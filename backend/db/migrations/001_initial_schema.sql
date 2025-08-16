-- Initial database schema for Surplus Supper application

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone VARCHAR(20),
    address TEXT,
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Restaurants table
CREATE TABLE IF NOT EXISTS restaurants (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    address TEXT NOT NULL,
    latitude DECIMAL(10, 8) NOT NULL,
    longitude DECIMAL(11, 8) NOT NULL,
    phone VARCHAR(20),
    email VARCHAR(255),
    cuisine_type VARCHAR(100),
    rating DECIMAL(3, 2) DEFAULT 0.0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Restaurant staff table
CREATE TABLE IF NOT EXISTS restaurant_staff (
    id SERIAL PRIMARY KEY,
    restaurant_id INTEGER REFERENCES restaurants(id) ON DELETE CASCADE,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) DEFAULT 'staff',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Inventory items table
CREATE TABLE IF NOT EXISTS inventory_items (
    id SERIAL PRIMARY KEY,
    restaurant_id INTEGER REFERENCES restaurants(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    original_price DECIMAL(10, 2) NOT NULL,
    surplus_price DECIMAL(10, 2) NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 1,
    category VARCHAR(100),
    expiry_time TIMESTAMP,
    is_available BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Offers table (for "Surprise Bags" and "Chef's Surprise")
CREATE TABLE IF NOT EXISTS offers (
    id SERIAL PRIMARY KEY,
    restaurant_id INTEGER REFERENCES restaurants(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    original_price DECIMAL(10, 2) NOT NULL,
    surplus_price DECIMAL(10, 2) NOT NULL,
    offer_type VARCHAR(50) NOT NULL, -- 'surprise_bag' or 'chef_surprise'
    ingredients TEXT, -- JSON array of ingredients for Chef's Surprise
    is_available BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Orders table
CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    restaurant_id INTEGER REFERENCES restaurants(id) ON DELETE CASCADE,
    total_amount DECIMAL(10, 2) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending', -- pending, confirmed, preparing, ready, delivered, cancelled
    pickup_time TIMESTAMP,
    special_instructions TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Order items table
CREATE TABLE IF NOT EXISTS order_items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER REFERENCES orders(id) ON DELETE CASCADE,
    inventory_item_id INTEGER REFERENCES inventory_items(id) ON DELETE SET NULL,
    offer_id INTEGER REFERENCES offers(id) ON DELETE SET NULL,
    quantity INTEGER NOT NULL DEFAULT 1,
    unit_price DECIMAL(10, 2) NOT NULL,
    total_price DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Notifications table
CREATE TABLE IF NOT EXISTS notifications (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    type VARCHAR(50) DEFAULT 'info', -- info, success, warning, error
    is_read BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- AI recommendations table
CREATE TABLE IF NOT EXISTS ai_recommendations (
    id SERIAL PRIMARY KEY,
    restaurant_id INTEGER REFERENCES restaurants(id) ON DELETE CASCADE,
    inventory_item_id INTEGER REFERENCES inventory_items(id) ON DELETE CASCADE,
    recommended_price DECIMAL(10, 2) NOT NULL,
    confidence_score DECIMAL(3, 2) DEFAULT 0.0,
    reasoning TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_restaurants_location ON restaurants(latitude, longitude);
CREATE INDEX IF NOT EXISTS idx_inventory_items_restaurant ON inventory_items(restaurant_id);
CREATE INDEX IF NOT EXISTS idx_inventory_items_available ON inventory_items(is_available);
CREATE INDEX IF NOT EXISTS idx_offers_restaurant ON offers(restaurant_id);
CREATE INDEX IF NOT EXISTS idx_offers_available ON offers(is_available);
CREATE INDEX IF NOT EXISTS idx_orders_user ON orders(user_id);
CREATE INDEX IF NOT EXISTS idx_orders_restaurant ON orders(restaurant_id);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
CREATE INDEX IF NOT EXISTS idx_notifications_user ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_unread ON notifications(is_read);

-- Insert sample data
INSERT INTO restaurants (name, description, address, latitude, longitude, phone, email, cuisine_type, rating) VALUES
('Pizza Palace', 'Authentic Italian pizza and pasta', '123 Main St, New York, NY', 40.7128, -74.0060, '+1-555-0123', 'info@pizzapalace.com', 'Italian', 4.5),
('Sushi Express', 'Fresh sushi and Japanese cuisine', '456 Broadway, New York, NY', 40.7589, -73.9851, '+1-555-0124', 'info@sushiexpress.com', 'Japanese', 4.3),
('Burger Joint', 'Classic American burgers and fries', '789 5th Ave, New York, NY', 40.7505, -73.9934, '+1-555-0125', 'info@burgerjoint.com', 'American', 4.1),
('Taco Town', 'Authentic Mexican tacos and burritos', '321 Lexington Ave, New York, NY', 40.7484, -73.9857, '+1-555-0126', 'info@tacotown.com', 'Mexican', 4.2),
('Green Garden', 'Fresh salads and healthy options', '654 Park Ave, New York, NY', 40.7628, -73.9745, '+1-555-0127', 'info@greengarden.com', 'Healthy', 4.4);

-- Insert sample inventory items
INSERT INTO inventory_items (restaurant_id, name, description, original_price, surplus_price, quantity, category, expiry_time) VALUES
(1, 'Margherita Pizza', 'Classic tomato and mozzarella pizza', 18.00, 9.00, 5, 'Pizza', NOW() + INTERVAL '2 hours'),
(1, 'Pepperoni Pizza', 'Spicy pepperoni pizza', 20.00, 10.00, 3, 'Pizza', NOW() + INTERVAL '1 hour'),
(2, 'California Roll', 'Crab, avocado, and cucumber roll', 12.00, 6.00, 8, 'Sushi', NOW() + INTERVAL '3 hours'),
(2, 'Salmon Nigiri', 'Fresh salmon over rice', 8.00, 4.00, 12, 'Sushi', NOW() + INTERVAL '2 hours'),
(3, 'Classic Burger', 'Beef burger with lettuce and tomato', 15.00, 7.50, 6, 'Burger', NOW() + INTERVAL '1 hour'),
(4, 'Chicken Tacos', 'Grilled chicken tacos with salsa', 10.00, 5.00, 10, 'Tacos', NOW() + INTERVAL '2 hours'),
(5, 'Caesar Salad', 'Fresh romaine with Caesar dressing', 12.00, 6.00, 4, 'Salad', NOW() + INTERVAL '1 hour');

-- Insert sample offers
INSERT INTO offers (restaurant_id, name, description, original_price, surplus_price, offer_type) VALUES
(1, 'Pizza Surprise Bag', 'Mystery pizza selection - could be any of our delicious pizzas!', 25.00, 12.50, 'surprise_bag'),
(2, 'Sushi Chef''s Surprise', 'Chef''s selection of fresh sushi rolls', 30.00, 15.00, 'chef_surprise'),
(3, 'Burger Combo Surprise', 'Burger with fries and drink - surprise flavor!', 20.00, 10.00, 'surprise_bag'),
(4, 'Taco Fiesta Bag', 'Assorted tacos with sides', 18.00, 9.00, 'surprise_bag'),
(5, 'Healthy Bowl Surprise', 'Chef''s choice of healthy ingredients', 22.00, 11.00, 'chef_surprise'); 