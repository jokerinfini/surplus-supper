-- Railway Migration for Surplus Supper
-- This file will be automatically run by Railway

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

-- Insert sample restaurants
INSERT INTO restaurants (name, description, address, latitude, longitude, phone, email, cuisine_type, rating) VALUES
('Pizza Palace', 'Authentic Italian pizza and pasta', '123 Main St, New York, NY', 40.7128, -74.0060, '+1-555-0123', 'info@pizzapalace.com', 'Italian', 4.5),
('Sushi Express', 'Fresh sushi and Japanese cuisine', '456 Broadway, New York, NY', 40.7589, -73.9851, '+1-555-0124', 'info@sushiexpress.com', 'Japanese', 4.3),
('Burger Joint', 'Classic American burgers and fries', '789 5th Ave, New York, NY', 40.7505, -73.9934, '+1-555-0125', 'info@burgerjoint.com', 'American', 4.1)
ON CONFLICT DO NOTHING;
