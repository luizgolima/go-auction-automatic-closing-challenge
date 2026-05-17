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

## 📊 Testing the Auction System

Since the base project does not return the ID in the creation response, follow this flow to test the automatic closing:

### 1. Create a New Auction
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
*Response: `201 Created`*

### 2. List Active Auctions to Find the ID
The system requires the `status=0` (Active) parameter to list auctions.
```bash
curl -i "http://localhost:8080/auction?status=0"
```
*Look for the `"id"` of the auction you just created.*

### 3. Check Status and Wait for Closing
Initially, the status will be `0` (Active). After the `AUCTION_DURATION` expires (default 5s in .env), the auction will move to status `1` (Completed).

**Check specific auction:**
```bash
# Replace <ID> with the ID from step 2
curl -i http://localhost:8080/auction/<ID>
```

**Wait for expiration and check again:**
The status should now be `"status": 1`.

### 4. Place a Bid
```bash
curl -i -X POST http://localhost:8080/bid \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "any-user-id",
    "auction_id": "<ID>",
    "amount": 1500.0
  }'
```

### 5. Find the Winning Bid
Once the auction is completed, you can check who won:
```bash
curl -i http://localhost:8080/auction/winner/<ID>
```

---

## 🛠️ Implementation Details

The core logic for automatic closing resides in `internal/infra/database/auction/create_auction.go`. Upon inserting a new auction, a Goroutine is spawned:

```go
go ar.scheduleAuctionClose(auctionEntity.Id)
```

This Goroutine waits for the configured `AUCTION_DURATION` and then updates the auction status in MongoDB.
