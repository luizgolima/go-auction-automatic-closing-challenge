# 🚀 Full Cycle Challenge: Auction Automatic Closing in Go

[![Go Version](https://img.shields.io/badge/Go-1.26.3-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Docker](https://img.shields.io/badge/Docker-2496ED?style=flat&logo=docker)](https://www.docker.com/)
[![MongoDB](https://img.shields.io/badge/MongoDB-47A248?style=flat&logo=mongodb)](https://www.mongodb.com/)
[![Status](https://img.shields.io/badge/Status-Completed-success?style=flat)]()

This project is an **Auction System** implemented in Go, with a focus on automatic auction closing using Goroutines. It was developed as a graduation project for the Full Cycle Go Expert course.

## 🧠 Architecture

The system follows Clean Architecture principles, ensuring decoupling and testability.

1.  **Automatic Closing**: When a new auction is created, a Goroutine is started to monitor its duration.
2.  **Concurrency Control**: Uses Goroutines and `time.Sleep` to handle the timing without blocking the main execution flow.
3.  **Persistence**: MongoDB is used to store auctions, bids, and user information.
4.  **Status Management**: Auctions start with `Active` status and automatically transition to `Completed` once the configured duration expires.

---

## 📁 Project Structure

```text
.
├── cmd/
│   └── auction/         # Application entry point and .env
├── internal/
│   ├── entity/          # Business logic and interfaces
│   ├── usecase/         # Application services
│   ├── infra/           # Infrastructure implementations
│   │   ├── database/    # MongoDB repositories
│   │   └── api/         # Web controllers and routing
│   └── internal_error/  # Custom error handling
├── configuration/       # App configurations (logger, db connection)
├── docker-compose.yml   # Container orchestration
├── Dockerfile           # Application Docker image
└── README.md            # Documentation
```

---

## ⚙️ Configuration

All settings are managed via environment variables (in `cmd/auction/.env`):

| Variable | Description | Default |
|----------|-------------|---------|
| `AUCTION_DURATION` | Duration for each auction (e.g., 30s, 1m) | `5m` |
| `MONGODB_URL` | MongoDB connection string | `mongodb://admin:admin@mongodb:27017/auctions?authSource=admin` |
| `MONGODB_DB` | MongoDB database name | `auctions` |

---

## 🚀 How to Run

### 1. Start the System
```bash
docker compose up --build
```

This will start:
- The Go application on `http://localhost:8080`
- A MongoDB instance.

### 2. Running Tests
To run the automated tests, including the automatic closing verification:
```bash
go test ./internal/infra/database/auction/...
```
*Note: The tests require a running MongoDB instance on localhost:27017 or the MONGODB_URL env var.*

---

## 📡 API Reference

### Auctions
| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/auction` | Create a new auction |
| `GET` | `/auction` | List auctions (Filter by `status`, `category`, `productName`) |
| `GET` | `/auction/:auctionId` | Get auction details by ID |
| `GET` | `/auction/winner/:auctionId` | Get the winning bid for an auction |

### Bids
| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/bid` | Place a new bid |
| `GET` | `/bid/:auctionId` | List all bids for a specific auction |

### Users
| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/user/:userId` | Get user details by ID |

---

## 📊 Complete Testing Flow

Follow these steps to test the full lifecycle of an auction, including the **Automatic Closing** feature.

### 1. Create a New Auction
Initially, the auction is created with status `Active (0)`.
```bash
curl -i -X POST http://localhost:8080/auction \
  -H "Content-Type: application/json" \
  -d '{
    "product_name": "Macbook Pro M3",
    "category": "Eletrônicos",
    "description": "64GB RAM, 1TB SSD",
    "condition": 0
  }'
```

### 2. Retrieve the Auction ID
Since the creation does not return the ID, list the active auctions:
```bash
curl -i "http://localhost:8080/auction?status=0"
```
*Copy the `"id"` field from the response.*

### 3. Place Bids (While Active)
Place one or more bids before the timer expires.
```bash
curl -i -X POST http://localhost:8080/bid \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user-uuid-here",
    "auction_id": "your-auction-id-here",
    "amount": 5000.0
  }'
```

### 4. Wait for Automatic Closing
Wait for the `AUCTION_DURATION` (configured in `.env`) to pass. The background Goroutine will automatically update the status to `Completed (1)`.

Check the status:
```bash
curl -i http://localhost:8080/auction/your-auction-id-here
```
*Expected: `"status": 1`*

### 5. Verify the Winner
Once closed, you can retrieve the winning bid:
```bash
curl -i http://localhost:8080/auction/winner/your-auction-id-here
```

### 6. List All Bids
View all bids placed during the auction:
```bash
curl -i http://localhost:8080/bid/your-auction-id-here
```

---

## 🛠️ Setup & Seed Data

Since the system does not have a user creation endpoint, you need to create a user manually in the database to place bids:

### Create a Test User
Run this command while the containers are active:
```bash
docker compose exec mongodb mongosh auctions -u admin -p admin --authenticationDatabase admin --eval "db.users.insertOne({_id: '1', name: 'User Test'})"
```

---

## 🛠️ Implementation Details

The core logic for automatic closing resides in `internal/infra/database/auction/create_auction.go`. Upon inserting a new auction, a Goroutine is spawned:

```go
go ar.scheduleAuctionClose(auctionEntity.Id)
```

This Goroutine waits for the configured `AUCTION_DURATION` and then updates the auction status in MongoDB.
