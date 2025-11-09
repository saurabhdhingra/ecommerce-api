# E-Commerce API

A RESTful e-commerce API built with Go, featuring user authentication, product management, shopping cart functionality, and Stripe payment integration.

<img width="2048" height="1511" alt="image" src="https://github.com/user-attachments/assets/2b762db5-c1ff-4124-94c0-9730d5c1a599" />

## Features

- **User Authentication**: Signup and login with JWT-based authentication
- **Product Management**: Browse products with search functionality (admin-only product creation)
- **Shopping Cart**: Add items, view cart, and checkout
- **Payment Integration**: Stripe payment intent creation for checkout
- **Role-Based Access**: Admin and regular user roles with different permissions

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Stripe account (for payment processing)

## Environment Variables

Create a `.env` file in the root directory or set the following environment variables:

### Database Configuration
```bash
DB_HOST=localhost          # PostgreSQL host (default: localhost)
DB_USER=postgres           # PostgreSQL user (default: postgres)
DB_PASSWORD=yourpassword   # PostgreSQL password (default: mysecretpassword)
DB_NAME=ecommerce_db       # Database name (default: ecommerce_db)
DB_PORT=5432               # PostgreSQL port (default: 5432)
```

### Application Configuration
```bash
PORT=8080                  # Server port (default: 8080)
JWT_SECRET=your_secret_key  # Secret key for JWT signing (default: auto-generated)
```

### Stripe Configuration
```bash
STRIPE_SECRET_KEY=sk_test_...  # Your Stripe secret key (default: mocked for development)
```

### Admin User Configuration
```bash
ADMIN_USER=admin          # Admin username (default: ecommerce_admin)
ADMIN_PASS=adminpass123    # Admin password (default: SuperSecureAdminPass123)
```

**Note**: If environment variables are not set, the application will use the default values shown above. However, it's recommended to set proper values for production use.

## Installation

1. **Clone the repository** (if applicable):
   ```bash
   git clone <repository-url>
   cd ecommerce-api
   ```

2. **Install dependencies**:
   ```bash
   go mod download
   ```

3. **Set up PostgreSQL database**:
   ```bash
   # Create database
   createdb ecommerce_db
   
   # Or using psql
   psql -U postgres
   CREATE DATABASE ecommerce_db;
   ```

4. **Set environment variables** (create `.env` file or export variables):
   ```bash
   export DB_HOST=localhost
   export DB_USER=postgres
   export DB_PASSWORD=yourpassword
   export DB_NAME=ecommerce_db
   export DB_PORT=5432
   export JWT_SECRET=your_jwt_secret_key_here
   export STRIPE_SECRET_KEY=sk_test_your_stripe_key
   export PORT=8080
   export ADMIN_USER=admin
   export ADMIN_PASS=adminpass123
   ```

## Running the Application

1. **Start the server**:
   ```bash
   go run .
   ```

   Or build and run:
   ```bash
   go build -o ecommerce-api
   ./ecommerce-api
   ```

2. The server will start on `http://localhost:8080` (or the port specified in `PORT` environment variable).

3. The application will automatically:
   - Connect to the PostgreSQL database
   - Run database migrations (create tables if they don't exist)
   - Create an admin user if it doesn't exist

## API Endpoints

### Base URL
```
http://localhost:8080
```

### Authentication

All authenticated endpoints require a JWT token in the Authorization header:
```
Authorization: Bearer <your_jwt_token>
```

---

## Public Endpoints

### 1. User Signup

Create a new user account.

**Endpoint**: `POST /api/signup`

**Request Body**:
```json
{
  "username": "john_doe",
  "password": "securepassword123"
}
```

**Example**:
```bash
curl -X POST http://localhost:8080/api/signup \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "password": "securepassword123"
  }'
```

**Response** (201 Created):
```json
{
  "message": "User created",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

---

### 2. User Login

Authenticate and receive a JWT token.

**Endpoint**: `POST /api/login`

**Request Body**:
```json
{
  "username": "john_doe",
  "password": "securepassword123"
}
```

**Example**:
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "password": "securepassword123"
  }'
```

**Response** (200 OK):
```json
{
  "message": "Login successful",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "is_admin": false
}
```

---

### 3. Get Products

Retrieve all products with optional search query.

**Endpoint**: `GET /api/products`

**Query Parameters**:
- `q` (optional): Search query to filter products by name or description

**Example** (Get all products):
```bash
curl -X GET http://localhost:8080/api/products
```

**Example** (Search products):
```bash
curl -X GET "http://localhost:8080/api/products?q=laptop"
```

**Response** (200 OK):
```json
[
  {
    "id": 1,
    "name": "Laptop",
    "description": "High-performance laptop",
    "price": "999.99",
    "inventory": 50
  },
  {
    "id": 2,
    "name": "Mouse",
    "description": "Wireless mouse",
    "price": "29.99",
    "inventory": 100
  }
]
```

---

## Authenticated Endpoints (User)

### 4. Add Item to Cart

Add a product to the shopping cart.

**Endpoint**: `POST /api/cart/add`

**Headers**:
```
Authorization: Bearer <your_jwt_token>
```

**Request Body**:
```json
{
  "product_id": 1,
  "quantity": 2
}
```

**Example**:
```bash
curl -X POST http://localhost:8080/api/cart/add \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "product_id": 1,
    "quantity": 2
  }'
```

**Response** (200 OK):
```json
{
  "ID": 1,
  "CreatedAt": "2024-01-15T10:30:00Z",
  "UpdatedAt": "2024-01-15T10:30:00Z",
  "DeletedAt": null,
  "UserID": 1,
  "Items": [
    {
      "ID": 1,
      "CreatedAt": "2024-01-15T10:30:00Z",
      "UpdatedAt": "2024-01-15T10:30:00Z",
      "DeletedAt": null,
      "CartID": 1,
      "ProductID": 1,
      "Quantity": 2,
      "Name": "Laptop",
      "PriceCents": 99999
    }
  ]
}
```

---

### 5. View Cart

Retrieve the current user's shopping cart.

**Endpoint**: `GET /api/cart`

**Headers**:
```
Authorization: Bearer <your_jwt_token>
```

**Example**:
```bash
curl -X GET http://localhost:8080/api/cart \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**Response** (200 OK):
```json
{
  "cart": {
    "ID": 1,
    "CreatedAt": "2024-01-15T10:30:00Z",
    "UpdatedAt": "2024-01-15T10:30:00Z",
    "DeletedAt": null,
    "UserID": 1,
    "Items": [
      {
        "ID": 1,
        "CreatedAt": "2024-01-15T10:30:00Z",
        "UpdatedAt": "2024-01-15T10:30:00Z",
        "DeletedAt": null,
        "CartID": 1,
        "ProductID": 1,
        "Quantity": 2,
        "Name": "Laptop",
        "PriceCents": 99999
      }
    ]
  },
  "total_usd": "1999.98",
  "total_cents": 199998
}
```

---

### 6. Checkout

Process checkout and create a Stripe payment intent.

**Endpoint**: `POST /api/checkout`

**Headers**:
```
Authorization: Bearer <your_jwt_token>
```

**Example**:
```bash
curl -X POST http://localhost:8080/api/checkout \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**Response** (200 OK):
```json
{
  "message": "Checkout successful. Payment initiated.",
  "total_paid_cents": 199998,
  "payment_intent_id": "pi_1234567890",
  "client_secret": "pi_1234567890_secret_abc123"
}
```

**Note**: This endpoint will:
- Validate cart is not empty
- Check product inventory
- Update inventory for all products
- Create a Stripe payment intent
- Clear the user's cart

---

## Admin Endpoints

### 7. Create Product

Create a new product (Admin only).

**Endpoint**: `POST /api/admin/products`

**Headers**:
```
Authorization: Bearer <admin_jwt_token>
```

**Request Body**:
```json
{
  "name": "New Product",
  "description": "Product description",
  "price_cents": 4999,
  "inventory": 100
}
```

**Example**:
```bash
curl -X POST http://localhost:8080/api/admin/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "name": "New Product",
    "description": "Product description",
    "price_cents": 4999,
    "inventory": 100
  }'
```

**Response** (201 Created):
```json
{
  "ID": 3,
  "CreatedAt": "2024-01-15T10:30:00Z",
  "UpdatedAt": "2024-01-15T10:30:00Z",
  "DeletedAt": null,
  "Name": "New Product",
  "Description": "Product description",
  "PriceCents": 4999,
  "Inventory": 100
}
```

**Note**: 
- `price_cents` is the price in cents (e.g., 4999 = $49.99)
- Only users with `is_admin: true` can access this endpoint

---

## Error Responses

All endpoints may return error responses in the following format:

```json
{
  "error": "Error message description"
}
```

### Common HTTP Status Codes

- `200 OK`: Request successful
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid request data
- `401 Unauthorized`: Missing or invalid authentication token
- `403 Forbidden`: Insufficient permissions (e.g., non-admin accessing admin endpoint)
- `404 Not Found`: Resource not found
- `409 Conflict`: Resource conflict (e.g., username already taken)
- `500 Internal Server Error`: Server error

### Example Error Response

```json
{
  "error": "Invalid username or password"
}
```

---

## Database Schema

The application automatically creates the following tables:

- **users**: User accounts with authentication
- **products**: Product catalog
- **carts**: Shopping carts (one per user)
- **cart_items**: Items in shopping carts

---

## Testing the API

### Quick Test Flow

1. **Sign up a new user**:
   ```bash
   curl -X POST http://localhost:8080/api/signup \
     -H "Content-Type: application/json" \
     -d '{"username": "testuser", "password": "testpass123"}'
   ```
   Save the `token` from the response.

2. **Login** (alternative):
   ```bash
   curl -X POST http://localhost:8080/api/login \
     -H "Content-Type: application/json" \
     -d '{"username": "testuser", "password": "testpass123"}'
   ```

3. **Get products**:
   ```bash
   curl -X GET http://localhost:8080/api/products
   ```

4. **Add item to cart** (use token from step 1):
   ```bash
   curl -X POST http://localhost:8080/api/cart/add \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer YOUR_TOKEN_HERE" \
     -d '{"product_id": 1, "quantity": 1}'
   ```

5. **View cart**:
   ```bash
   curl -X GET http://localhost:8080/api/cart \
     -H "Authorization: Bearer YOUR_TOKEN_HERE"
   ```

6. **Checkout**:
   ```bash
   curl -X POST http://localhost:8080/api/checkout \
     -H "Authorization: Bearer YOUR_TOKEN_HERE"
   ```

---

## Project Structure

```
ecommerce-api/
├── config.go              # Configuration loading
├── main.go                # Application entry point and routing
├── domain/                 # Domain models and interfaces
│   ├── user.go
│   ├── product.go
│   ├── cart.go
│   ├── jwt_claims.go
│   ├── errors.go
│   ├── user_repo.go
│   ├── product_repo.go
│   └── cart_repo.go
├── repository/            # Database implementations
│   ├── postgres_repo.go
│   ├── user_repo.go
│   ├── product_repo.go
│   └── cart_repo.go
├── service/               # Business logic
│   ├── user_service.go
│   ├── product_service.go
│   ├── cart_service.go
│   ├── auth_service.go
│   └── stripe_service.go
├── handler/               # HTTP handlers
│   ├── handler.go
│   ├── user_handler.go
│   └── middleware.go
└── go.mod                 # Go dependencies
```

---

## Security Notes

⚠️ **Important for Production**:

1. **JWT Secret**: Use a strong, randomly generated secret key for `JWT_SECRET`
2. **Password Hashing**: The current implementation uses SHA256. For production, use `bcrypt` or `argon2`
3. **HTTPS**: Always use HTTPS in production
4. **Database**: Use strong database passwords and restrict database access
5. **Environment Variables**: Never commit `.env` files or secrets to version control
6. **Stripe Keys**: Use test keys for development and live keys only in production

---

## Troubleshooting

### Database Connection Issues

- Ensure PostgreSQL is running: `pg_isready`
- Verify database credentials in environment variables
- Check if database exists: `psql -l | grep ecommerce_db`

### Port Already in Use

- Change the `PORT` environment variable
- Or stop the process using the port: `lsof -ti:8080 | xargs kill`

### Migration Errors

- Ensure PostgreSQL user has CREATE TABLE permissions
- Check database connection string format

---

## License

MIT License

---

## Aknowledgement

https://roadmap.sh/projects/ecommerce-api
