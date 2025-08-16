---
name: 🏪 Restaurant Dashboard Feature
about: Build restaurant dashboard for managing surplus food inventory
title: '[FEATURE] Restaurant Dashboard - Inventory Management'
labels: ['enhancement', 'restaurant-dashboard', 'good first issue', 'help wanted']
assignees: ''
---

## 🏪 Restaurant Dashboard Feature

### 📝 Description
Build a comprehensive restaurant dashboard for restaurant owners to manage their surplus food inventory, orders, and business analytics. This is a core feature that enables restaurants to participate in the Surplus Supper marketplace.

### 🎯 Goals
- [ ] Restaurant registration and authentication system
- [ ] Restaurant profile management (name, address, cuisine type, hours)
- [ ] Inventory management (add/edit/delete surplus food items)
- [ ] Order management and tracking system
- [ ] Analytics dashboard (sales, popular items, waste reduction)
- [ ] Real-time notifications for new orders
- [ ] Responsive dashboard UI/UX

### 🛠️ Technical Requirements

#### Backend (Go)
- **Restaurant Service**: Business logic for restaurant operations
- **API Endpoints**: 
  - `POST /api/restaurant/register` - Restaurant registration
  - `POST /api/restaurant/login` - Restaurant authentication
  - `GET /api/restaurant/profile` - Get restaurant profile
  - `PUT /api/restaurant/profile` - Update restaurant profile
  - `GET /api/restaurant/inventory` - Get inventory items
  - `POST /api/restaurant/inventory` - Add inventory item
  - `PUT /api/restaurant/inventory/:id` - Update inventory item
  - `DELETE /api/restaurant/inventory/:id` - Delete inventory item
  - `GET /api/restaurant/orders` - Get restaurant orders
  - `PUT /api/restaurant/orders/:id/status` - Update order status

#### Frontend (Next.js)
- **Dashboard Layout**: Responsive sidebar navigation
- **Authentication**: Restaurant login/register forms
- **Inventory Management**: CRUD operations for food items
- **Order Management**: View and manage incoming orders
- **Analytics**: Charts and reports for business insights
- **Profile Management**: Restaurant information editing

#### Database Schema
```sql
-- Restaurants table
CREATE TABLE restaurants (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    address TEXT,
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    cuisine_type VARCHAR(100),
    phone VARCHAR(20),
    opening_hours JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Inventory items table
CREATE TABLE inventory_items (
    id SERIAL PRIMARY KEY,
    restaurant_id INTEGER REFERENCES restaurants(id),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    original_price DECIMAL(10, 2),
    surplus_price DECIMAL(10, 2),
    quantity INTEGER DEFAULT 1,
    category VARCHAR(100),
    expiry_date TIMESTAMP,
    image_url VARCHAR(500),
    is_available BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### 📁 Files to Create/Modify

#### Backend Files
- `backend/restaurantService/restaurant_service.go` - Restaurant business logic
- `backend/restaurantService/auth.go` - Restaurant authentication
- `backend/api/restaurant/restaurant_handler.go` - Restaurant API handlers
- `backend/api/restaurant/inventory_handler.go` - Inventory API handlers
- `backend/api/restaurant/order_handler.go` - Order management handlers
- `backend/middleware/restaurant_auth.go` - Restaurant authentication middleware
- `backend/db/migrations/002_restaurant_schema.sql` - Database migrations

#### Frontend Files
- `frontend-next/src/components/restaurant/` - Restaurant dashboard components
- `frontend-next/src/app/restaurant/` - Restaurant dashboard pages
- `frontend-next/src/lib/restaurant.ts` - Restaurant API client
- `frontend-next/src/types/restaurant.ts` - Restaurant type definitions

### 🎨 UI/UX Considerations
- **Design System**: Follow existing glassmorphism design theme
- **Responsive**: Mobile-first design for restaurant staff on mobile devices
- **Accessibility**: WCAG 2.1 AA compliance
- **Performance**: Fast loading times for inventory management
- **User Experience**: Intuitive workflow for adding/managing inventory

### 📋 Acceptance Criteria
- [ ] Restaurant owners can register and login to the system
- [ ] Restaurant profile can be created and edited
- [ ] Inventory items can be added, edited, and deleted
- [ ] Orders are displayed in real-time with status updates
- [ ] Dashboard shows key metrics (sales, inventory, orders)
- [ ] Mobile-responsive design works on all devices
- [ ] Proper error handling and validation throughout
- [ ] Integration with existing authentication system

### 🏷️ Labels
- `enhancement`
- `restaurant-dashboard`
- `good first issue`
- `help wanted`
- `backend`
- `frontend`
- `database`

### 📚 Additional Notes
- Should integrate seamlessly with existing user authentication system
- Follow existing code patterns and styling conventions
- Include comprehensive error handling and input validation
- Consider performance implications for real-time features
- Implement proper security measures for restaurant data
