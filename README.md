# Flyola Hotel Services Backend

A professional, modular Go backend service for hotel booking management system built with Gin framework and GORM.

## 🏗️ Architecture

This project follows a clean, modular architecture with clear separation of concerns:

```
Flyola-Services-Backend/
├── cmd/
│   └── server/           # Application entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── database/        # Database initialization
│   ├── handlers/        # HTTP request handlers
│   ├── middleware/      # HTTP middleware
│   ├── models/          # Data models (separated by domain)
│   ├── routes/          # Route definitions (separated by domain)
│   ├── router/          # Main router setup
│   └── services/        # Business logic layer
├── pkg/                 # Public packages (if any)
├── .air.toml           # Hot reload configuration
├── Dockerfile          # Docker configuration
├── Makefile           # Build and development commands
└── go.mod             # Go module definition
```

## 📁 Project Structure Details

### Models (Domain-Separated)
- `city.go` - City/destination models
- `hotel.go` - Hotel models
- `room.go` - Room, RoomCategory, RoomAvailability models
- `booking.go` - Booking and guest models
- `payment.go` - Payment models
- `review.go` - Review models
- `meal_plan.go` - Meal plan models

### Services (Business Logic)
Each service handles the business logic for its respective domain:
- `city_service.go`
- `hotel_service.go`
- `room_service.go`
- `room_category_service.go`
- `room_availability_service.go`
- `booking_service.go`
- `payment_service.go`
- `review_service.go`
- `meal_plan_service.go`

### Handlers (HTTP Layer)
HTTP request handlers that use services:
- `city_handler.go`
- `hotel_handler.go`
- `room_handler.go`
- `room_category_handler.go`
- `room_availability_handler.go`
- `booking_handler.go`
- `payment_handler.go`
- `review_handler.go`
- `meal_plan_handler.go`

### Routes (Domain-Separated)
Route definitions separated by domain:
- `city_routes.go`
- `hotel_routes.go`
- `room_routes.go`
- `booking_routes.go`
- `payment_routes.go`
- `review_routes.go`
- `meal_plan_routes.go`

## 🚀 Getting Started

### Prerequisites
- Go 1.21 or higher
- MySQL 8.0 or higher
- Git

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd Flyola-Services-Backend
   ```

2. **Install dependencies**
   ```bash
   make deps
   ```

3. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your database credentials
   ```

4. **Run the application**
   ```bash
   make run
   ```

## 🛠️ Development

### Available Make Commands

```bash
make build         # Build the application
make run           # Run the application
make dev           # Run with hot reload (requires air)
make clean         # Clean build artifacts
make test          # Run tests
make test-coverage # Run tests with coverage
make fmt           # Format code
make lint          # Lint code (requires golangci-lint)
make tidy          # Tidy dependencies
make deps          # Install dependencies
make docker-build  # Build Docker image
make docker-run    # Run Docker container
make help          # Show help message
```

### Hot Reload Development

For development with hot reload:

1. **Install air**
   ```bash
   go install github.com/cosmtrek/air@latest
   ```

2. **Run with hot reload**
   ```bash
   make dev
   ```

### Database Setup

1. **Create MySQL database**
   ```sql
   CREATE DATABASE flyola;
   ```

2. **Run migrations**
   The application will auto-migrate tables on startup (development mode).

## 📡 API Endpoints

### Cities
- `GET /api/v1/cities` - Get all cities
- `POST /api/v1/cities` - Create city
- `GET /api/v1/cities/:id` - Get city by ID
- `PUT /api/v1/cities/:id` - Update city
- `DELETE /api/v1/cities/:id` - Delete city

### Hotels
- `GET /api/v1/hotels` - Get all hotels
- `POST /api/v1/hotels` - Create hotel
- `GET /api/v1/hotels/:id` - Get hotel by ID
- `PUT /api/v1/hotels/:id` - Update hotel
- `DELETE /api/v1/hotels/:id` - Delete hotel
- `GET /api/v1/hotels/city/:cityId` - Get hotels by city

### Rooms
- `GET /api/v1/rooms` - Get all rooms
- `POST /api/v1/rooms` - Create room
- `GET /api/v1/rooms/:id` - Get room by ID
- `PUT /api/v1/rooms/:id` - Update room
- `DELETE /api/v1/rooms/:id` - Delete room
- `GET /api/v1/rooms/hotel/:hotelId` - Get rooms by hotel

### Room Categories
- `GET /api/v1/room-categories` - Get all room categories
- `POST /api/v1/room-categories` - Create room category
- `GET /api/v1/room-categories/:id` - Get room category by ID
- `PUT /api/v1/room-categories/:id` - Update room category
- `DELETE /api/v1/room-categories/:id` - Delete room category

### Room Availability
- `GET /api/v1/room-availability` - Get room availability
- `POST /api/v1/room-availability` - Create room availability
- `GET /api/v1/room-availability/:id` - Get room availability by ID
- `PUT /api/v1/room-availability/:id` - Update room availability
- `DELETE /api/v1/room-availability/:id` - Delete room availability

### Bookings
- `GET /api/v1/bookings` - Get all bookings
- `POST /api/v1/bookings` - Create booking
- `GET /api/v1/bookings/:id` - Get booking by ID
- `PUT /api/v1/bookings/:id` - Update booking
- `DELETE /api/v1/bookings/:id` - Delete booking
- `PUT /api/v1/bookings/:id/cancel` - Cancel booking

### Payments
- `POST /api/v1/payments/process` - Process payment
- `GET /api/v1/payments/:id` - Get payment by ID
- `GET /api/v1/payments/booking/:bookingId` - Get payment by booking

### Reviews
- `GET /api/v1/reviews` - Get all reviews
- `POST /api/v1/reviews` - Create review
- `GET /api/v1/reviews/hotel/:hotelId` - Get hotel reviews
- `PUT /api/v1/reviews/:id/status` - Update review status
- `DELETE /api/v1/reviews/:id` - Delete review

### Meal Plans
- `GET /api/v1/meal-plans` - Get all meal plans
- `POST /api/v1/meal-plans` - Create meal plan
- `GET /api/v1/meal-plans/:id` - Get meal plan by ID
- `PUT /api/v1/meal-plans/:id` - Update meal plan
- `DELETE /api/v1/meal-plans/:id` - Delete meal plan

## 🐳 Docker

### Build and run with Docker

```bash
# Build Docker image
make docker-build

# Run Docker container
make docker-run
```

### Using Docker Compose

```bash
docker-compose up -d
```

## 🧪 Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage
```

## 📝 Configuration

Environment variables:

- `PORT` - Server port (default: 8080)
- `DATABASE_URL` - MySQL connection string
- `ENVIRONMENT` - Environment (development/production)

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests and linting
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License.

## 🏨 Features

- ✅ Modular architecture with clean separation of concerns
- ✅ RESTful API design
- ✅ GORM for database operations
- ✅ Gin framework for HTTP routing
- ✅ Hot reload development setup
- ✅ Docker support
- ✅ Comprehensive hotel management
- ✅ Room availability and pricing
- ✅ Booking management
- ✅ Payment processing
- ✅ Review system
- ✅ Meal plan management
- ✅ Professional error handling
- ✅ CORS middleware
- ✅ Request logging
- ✅ Database migrations# Flyola-Mp-Services
