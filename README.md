# 🍽️ Surplus Supper

> A sustainable marketplace connecting restaurants with customers to reduce food waste by selling surplus food at discounted prices.

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org/)
[![Node.js Version](https://img.shields.io/badge/Node.js-18+-green.svg)](https://nodejs.org/)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Contributions Welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](CONTRIBUTING.md)

## 🌟 What is Surplus Supper?

Surplus Supper is a marketplace platform that helps restaurants reduce food waste by connecting them with customers who want to purchase surplus food at discounted prices. Our mission is to create a sustainable food ecosystem while helping restaurants reduce waste and customers save money.

### 🎯 Key Features

#### ✅ **Completed (Phase 1)**
- **🔐 User Authentication**: JWT-based authentication system with secure password hashing
- **🏪 Restaurant Discovery**: Location-based restaurant search with real-time geolocation
- **🎨 Modern UI/UX**: Glassmorphism design with responsive layout and smooth animations
- **📱 Mobile Responsive**: Optimized for all devices with touch-friendly interface
- **🌍 Location Services**: GPS integration with reverse geocoding and distance calculation
- **⚡ Real-time Search**: Instant restaurant search with location-based filtering
- **🎭 Interactive Elements**: 3D card effects, animated backgrounds, and smooth transitions

#### 🚧 **In Development (Phase 2)**
- **🏪 Restaurant Dashboard**: Complete restaurant management system
- **📦 Inventory Management**: Add, edit, and manage surplus food items
- **📊 Analytics Dashboard**: Business insights and performance metrics
- **🔔 Real-time Notifications**: Live order updates and alerts

#### 🎯 **Planned Features**
- **💳 Payment Processing**: Secure payment integration with Stripe
- **📱 Mobile Apps**: Native iOS and Android applications
- **🤖 AI Integration**: Smart pricing and waste prediction
- **⭐ Reviews & Ratings**: Customer feedback system
- **🎁 Loyalty Program**: Rewards for sustainable choices

## 🏗️ Tech Stack

### Backend
- **Language**: Go 1.21+
- **Framework**: Standard library with Gorilla Mux
- **Database**: PostgreSQL with migrations
- **Authentication**: JWT with bcrypt password hashing
- **API**: RESTful API with CORS support

### Frontend
- **Framework**: Next.js 14 with App Router
- **Language**: TypeScript
- **Styling**: Tailwind CSS with custom glassmorphism design
- **Animations**: Framer Motion and Three.js
- **Icons**: Lucide React
- **Components**: Shadcn/ui

### Infrastructure
- **Containerization**: Docker & Docker Compose
- **Database**: PostgreSQL with automatic migrations
- **Development**: Hot-reload for both frontend and backend

## 🚀 Quick Start

### Prerequisites
- **Go** 1.21+ (for backend)
- **Node.js** 18+ (for frontend)
- **Docker** (for database)
- **Git**

### 1. Clone the Repository
```bash
git clone https://github.com/YOUR_USERNAME/surplus-supper.git
cd surplus-supper
```

### 2. Start the Database
```bash
docker-compose up -d postgres
```

### 3. Setup Backend
```bash
cd backend
go mod download
go run main.go
```

### 4. Setup Frontend
```bash
cd frontend-next
npm install
npm run dev
```

### 5. Verify Setup
- **Backend API**: http://localhost:8080/health
- **Frontend App**: http://localhost:3000
- **Database**: PostgreSQL on port 5433

## 📁 Project Structure

```
surplus-supper/
├── backend/                 # Go backend API
│   ├── api/                # HTTP handlers
│   │   └── auth/           # Authentication endpoints
│   ├── middleware/         # Authentication & CORS
│   ├── userService/        # User business logic
│   ├── db/                 # Database migrations
│   └── main.go            # Application entry point
├── frontend-next/          # Next.js frontend
│   ├── src/
│   │   ├── app/           # App router pages
│   │   ├── components/    # React components
│   │   │   ├── auth/      # Authentication components
│   │   │   ├── features/  # Feature components
│   │   │   └── ui/        # UI components
│   │   ├── lib/           # Utility functions
│   │   └── types/         # TypeScript types
│   └── public/            # Static assets
├── .github/               # GitHub templates and workflows
├── CONTRIBUTING.md        # Contribution guidelines
├── ROADMAP.md            # Development roadmap
└── docker-compose.yml    # Development environment
```

## 🎨 Features in Detail

### 🔐 Authentication System
- **User Registration**: Secure account creation with validation
- **User Login**: JWT-based authentication with refresh tokens
- **Password Security**: bcrypt hashing with salt rounds
- **Session Management**: Automatic token refresh and logout

### 🏪 Restaurant Discovery
- **Location Search**: GPS-based restaurant finding
- **Real-time Geolocation**: Browser location API integration
- **Distance Calculation**: Haversine formula for accurate distances
- **Search Filtering**: Dynamic restaurant filtering and sorting

### 🎨 Modern UI/UX
- **Glassmorphism Design**: Modern glass-like interface elements
- **3D Card Effects**: Interactive hover animations
- **Animated Backgrounds**: Falling food particles with HTML5 Canvas
- **Smooth Transitions**: Framer Motion animations throughout
- **Responsive Layout**: Mobile-first design approach

### 📱 Mobile Experience
- **Touch Optimized**: Large touch targets and swipe gestures
- **Progressive Web App**: Installable web application
- **Offline Support**: Service worker for offline functionality
- **Fast Loading**: Optimized images and lazy loading

## 🛠️ Development

### Running in Development Mode
```bash
# Terminal 1: Backend with hot-reload
cd backend
go run main.go

# Terminal 2: Frontend with hot-reload
cd frontend-next
npm run dev

# Terminal 3: Database (if not using Docker)
docker-compose up -d postgres
```

### Testing
```bash
# Backend tests
cd backend
go test ./...

# Frontend tests
cd frontend-next
npm run test
```

### Building for Production
```bash
# Backend
cd backend
go build -o backend.exe

# Frontend
cd frontend-next
npm run build
```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### 🎯 Good First Issues
- [Restaurant Dashboard Feature](.github/ISSUE_TEMPLATE/restaurant_dashboard.md)
- UI/UX improvements
- Documentation updates
- Bug fixes

### 📋 How to Contribute
1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📊 Project Status

### ✅ Phase 1: Foundation (Completed)
- [x] Project architecture and setup
- [x] Go backend with REST API
- [x] Next.js frontend with modern UI
- [x] PostgreSQL database with migrations
- [x] Docker development environment
- [x] User authentication system
- [x] Restaurant search and discovery
- [x] Location-based services
- [x] Responsive design implementation

### 🚧 Phase 2: Restaurant Dashboard (In Progress)
- [ ] Restaurant registration and authentication
- [ ] Restaurant profile management
- [ ] Inventory management system
- [ ] Order management dashboard
- [ ] Real-time notifications
- [ ] Analytics and reporting

### 🎯 Phase 3: Customer Experience (Planned)
- [ ] Customer order placement
- [ ] Payment processing (Stripe)
- [ ] Order tracking and notifications
- [ ] Customer reviews and ratings
- [ ] Favorites and wishlist

### 🤖 Phase 4: AI & Intelligence (Planned)
- [ ] AI-powered recipe suggestions
- [ ] Smart pricing recommendations
- [ ] Food waste prediction
- [ ] Personalized recommendations

## 🏷️ Issue Labels

- `bug` - Something isn't working
- `enhancement` - New feature or request
- `good first issue` - Good for newcomers
- `help wanted` - Extra attention is needed
- `restaurant-dashboard` - Restaurant management features
- `frontend` - Frontend-related changes
- `backend` - Backend-related changes
- `database` - Database schema changes

## 📈 Roadmap

See our detailed [Development Roadmap](ROADMAP.md) for the complete feature timeline and technical milestones.

## 🐛 Reporting Bugs

Found a bug? Please use our [Bug Report Template](.github/ISSUE_TEMPLATE/bug_report.md) to report it.

## 💡 Suggesting Features

Have an idea for a new feature? Use our [Feature Request Template](.github/ISSUE_TEMPLATE/feature_request.md) to suggest it.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- **Design Inspiration**: Modern glassmorphism and sustainable design principles
- **Icons**: [Lucide React](https://lucide.dev/) for beautiful icons
- **Animations**: [Framer Motion](https://www.framer.com/motion/) for smooth interactions
- **3D Graphics**: [Three.js](https://threejs.org/) for 3D elements

## 📞 Support

- **Issues**: [GitHub Issues](https://github.com/YOUR_USERNAME/surplus-supper/issues)
- **Discussions**: [GitHub Discussions](https://github.com/YOUR_USERNAME/surplus-supper/discussions)
- **Documentation**: Check the [Contributing Guide](CONTRIBUTING.md)

---

**Made with ❤️ for a sustainable future** 🌱

*Help us reduce food waste, one meal at a time!* 