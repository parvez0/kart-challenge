# Food Ordering System API

A RESTful API for a food ordering system built with Go, using SQLite as the database.

## Prerequisites

- Go 1.24.3 or higher
- SQLite3

## Getting Started

1. Clone the repository
2. Navigate to the project directory
3. Install dependencies:
```bash
go mod download
```
4. Run the application:
```bash
go run main.go
```

The server will start on port 8080.

## API Endpoints

### Health Check
- **GET** `/health`
- Returns a simple "OK" response to verify the server is running
- Response: `200 OK`

### Products

#### Get All Products
- **GET** `/products`
- Returns a list of all available products
- Response: `200 OK`
```json
[
  {
    "id": 1,
    "name": "Margherita Pizza",
    "price": 12.99,
    "category": "Pizza"
  }
]
```

#### Get Product by ID
- **GET** `/product/{productId}`
- Returns details of a specific product
- Response: `200 OK`
```json
{
  "id": 1,
  "name": "Margherita Pizza",
  "price": 12.99,
  "category": "Pizza"
}
```

### Orders

#### Get All Orders
- **GET** `/orders`
- Returns a list of all orders with their items and products
- Response: `200 OK`
```json
[
  {
    "id": 1,
    "items": [
      {
        "productId": "1",
        "quantity": 2
      }
    ],
    "products": [
      {
        "id": 1,
        "name": "Margherita Pizza",
        "price": 12.99,
        "category": "Pizza"
      }
    ]
  }
]
```

#### Create Order
- **POST** `/order`
- Creates a new order
- Request Body:
```json
{
  "couponCode": "optional_coupon_code",
  "items": [
    {
      "productId": "1",
      "quantity": 2
    }
  ]
}
```
- Response: `200 Created`
```json
{
  "id": 1,
  "items": [
    {
      "productId": "1",
      "quantity": 2
    }
  ],
  "products": [
    {
      "id": 1,
      "name": "Margherita Pizza",
      "price": 12.99,
      "category": "Pizza"
    }
  ]
}
```

## Error Responses

The API uses standard HTTP status codes:

- `400 Bad Request` - Invalid request parameters
- `404 Not Found` - Resource not found
- `422 Unprocessable Entity` - Invalid coupon code
- `500 Internal Server Error` - Server-side error

## Database

The application uses SQLite in-memory database for simplicity. The database is automatically seeded with:
- Sample products (pizzas, salads, sides, desserts)
- Coupon codes from files in the `data` directory

## Project Structure

```
.
├── data/           # Contains coupon code files
├── pkg/            # Core package with business logic
│   ├── db.go       # Database setup and configuration
│   ├── handler.go  # HTTP request handlers
│   ├── models.go   # Data models
│   └── seeder.go   # Database seeding logic
├── utils/          # Utility functions
│   ├── helper.go   # Helper functions
│   └── logger.go   # Logging configuration
├── go.mod          # Go module file
├── go.sum          # Go module checksum
└── main.go         # Application entry point
```

## Logging

The application uses structured logging with logrus. Logs are output to both console and file (if configured).

## Middleware

The API includes two middleware components:
1. URL Logging - Logs request URLs and response times
2. Authorization - Currently disabled but can be enabled for API key authentication
