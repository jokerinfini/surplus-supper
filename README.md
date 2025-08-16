# Surplus Supper 🍽️

A marketplace application designed to reduce food waste by connecting restaurants with surplus food to customers seeking discounted meals.

## 🌟 Features

### For Customers
- **Browse Nearby Restaurants**: Find restaurants with surplus food near your location
- **Surprise Bags**: Get mystery bags of surplus food at great discounts
- **Chef's Surprise**: AI-generated recipes using surplus ingredients
- **Real-time Notifications**: Get notified when new offers become available
- **Easy Ordering**: Simple checkout process with HTMX-powered interactions

### For Restaurants
- **Inventory Management**: Add and manage surplus food items
- **Profit Advisor**: AI-powered pricing recommendations
- **Real-time Dashboard**: Monitor incoming orders and sales
- **Order Management**: Track order status and customer communications

### AI Features
- **Profit Advisor**: Analyzes market data to suggest optimal discount prices
- **Creative Kitchen**: Generates recipes from surplus ingredients based on user preferences

## 🏗️ Architecture

### Backend (Go)
- **Microservices Architecture**: Separate services for users, restaurants, orders, AI, and notifications
- **GraphQL API**: Single entry point using gqlgen for type-safe queries and mutations
- **PostgreSQL Database**: Robust relational database with proper indexing
- **WebSocket Support**: Real-time notifications using Gorilla WebSocket
- **HTMX Integration**: Server-side rendering with dynamic interactions

### Frontend
- **HTMX**: Dynamic interactions without complex JavaScript
- **Alpine.js**: Lightweight client-side enhancements
- **Tailwind CSS**: Modern, responsive styling
- **Server-side Rendering**: Templates served directly by Go backend

## 📁 Project Structure

```
surplus-supper/
├── backend/
│   ├── aiService/           # AI features (Profit Advisor, Creative Kitchen)
│   ├── api/
│   │   ├── graph/          # GraphQL schema and resolvers
│   │   └── rest/           # HTMX handlers
│   ├── db/                 # Database connection and migrations
│   ├── notificationService/ # Real-time notifications
│   ├── orderService/       # Order management
│   ├── restaurantService/  # Restaurant and inventory management
│   ├── userService/        # User authentication and management
│   ├── go.mod             # Go dependencies
│   └── main.go            # Application entry point
├── frontend/
│   ├── assets/
│   │   ├── css/           # Tailwind CSS styles
│   │   └── js/            # Alpine.js components
│   ├── templates/         # HTML templates (placeholder)
│   └── tailwind.config.js # Tailwind configuration
└── README.md
```

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Node.js (for Tailwind CSS compilation, optional)

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd surplus-supper
   ```

2. **Set up the database**
   ```bash
   # Create PostgreSQL database
   createdb surplus_supper
   
   # Run migrations
   psql -d surplus_supper -f backend/db/migrations/001_initial_schema.sql
   ```

3. **Configure environment variables**
   ```bash
   export DATABASE_URL="postgres://username:password@localhost/surplus_supper?sslmode=disable"
   export PORT="8080"
   ```

4. **Install Go dependencies**
   ```bash
   cd backend
   go mod tidy
   ```

5. **Run the application**
   ```bash
   go run main.go
   ```

6. **Access the application**
   - Open your browser and navigate to `http://localhost:8080`
   - The application will serve both the API and frontend

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | `postgres://postgres:password@localhost/surplus_supper?sslmode=disable` |
| `PORT` | Server port | `8080` |
| `JWT_SECRET` | JWT signing secret | `your-secret-key` |

### Database Setup

The application includes sample data for testing:

- **Restaurants**: Tasty Bites, Spice Garden, Pizza Palace
- **Inventory Items**: Various surplus food items with discounts
- **Offers**: Surprise bags and special offers

## 📚 API Documentation

### GraphQL Endpoint
- **URL**: `http://localhost:8080/graphql`
- **Playground**: Available at the same URL when running in development

### Key Queries
```graphql
# Get nearby restaurants
query {
  nearbyRestaurants(latitude: 40.7128, longitude: -74.0060, radius: 10) {
    id
    name
    description
    address
    rating
  }
}

# Get restaurant inventory
query {
  availableInventoryItems(restaurantId: 1) {
    id
    name
    originalPrice
    surplusPrice
    quantity
    expiryTime
  }
}
```

### Key Mutations
```graphql
# Create an order
mutation {
  createOrder(input: {
    restaurantId: 1
    orderItems: [{
      inventoryItemId: 1
      quantity: 2
    }]
  }) {
    id
    totalAmount
    status
  }
}
```

## 🎨 Frontend Features

### HTMX Integration
- **Dynamic Search**: Real-time restaurant search without page reloads
- **Shopping Cart**: Add items to cart with instant feedback
- **Order Confirmation**: Seamless checkout process
- **Inventory Updates**: Real-time inventory status updates

### Alpine.js Components
- **Toast Notifications**: Success/error messages
- **Modal Dialogs**: Confirmation dialogs and forms
- **Dropdown Menus**: User navigation and settings
- **Form Validation**: Client-side validation with server feedback

### Responsive Design
- **Mobile-First**: Optimized for mobile devices
- **Progressive Enhancement**: Works without JavaScript
- **Accessibility**: WCAG compliant design

## 🤖 AI Features

### Profit Advisor
- Analyzes market trends and competitor pricing
- Suggests optimal discount percentages
- Considers item category and expiry time
- Provides confidence scores for recommendations

### Creative Kitchen
- Generates recipes from surplus ingredients
- Considers user dietary preferences
- Provides step-by-step instructions
- Includes nutritional information and cooking times

## 🔒 Security

- **JWT Authentication**: Secure user sessions
- **Password Hashing**: bcrypt for password security
- **SQL Injection Prevention**: Parameterized queries
- **CORS Configuration**: Proper cross-origin settings
- **Input Validation**: Server-side validation for all inputs

## 🧪 Testing

### Backend Testing
```bash
cd backend
go test ./...
```

### Database Testing
```bash
# Test database connection
go run main.go --test-db
```

## 📊 Monitoring

### Health Checks
- **Database**: Connection status and query performance
- **Services**: Individual service health endpoints
- **WebSocket**: Connection status and message throughput

### Logging
- Structured logging with request tracing
- Error tracking and alerting
- Performance metrics collection

## 🚀 Deployment

### Docker (Recommended)
```bash
# Build the application
docker build -t surplus-supper .

# Run with PostgreSQL
docker-compose up -d
```

### Manual Deployment
1. Set up PostgreSQL database
2. Configure environment variables
3. Build and run the Go application
4. Set up reverse proxy (nginx/Apache)
5. Configure SSL certificates

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- **HTMX**: For making dynamic web applications simple
- **Alpine.js**: For lightweight client-side interactivity
- **Tailwind CSS**: For rapid UI development
- **Go**: For fast, reliable backend development
- **PostgreSQL**: For robust data storage

## 📞 Support

For support and questions:
- Create an issue in the GitHub repository
- Check the documentation in the `/docs` folder
- Join our community discussions

---

**Made with ❤️ to reduce food waste and help the planet** 