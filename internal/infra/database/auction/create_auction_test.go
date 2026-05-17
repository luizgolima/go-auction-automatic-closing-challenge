package auction

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/luigolima/go-auction-automatic-closing-challenge/internal/entity/auction_entity"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestAuctionClosing(t *testing.T) {
	os.Setenv("AUCTION_DURATION", "2s")
	ctx := context.Background()

	mongoURL := os.Getenv("MONGODB_URL")
	if mongoURL == "" {
		mongoURL = "mongodb://admin:admin@localhost:27017/auctions_test?authSource=admin"
	}
	
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		t.Skip("MongoDB not available")
	}
	
	err = client.Ping(ctx, nil)
	if err != nil {
		t.Skip("MongoDB not reachable")
	}
	defer client.Disconnect(ctx)

	db := client.Database("auctions_test")
	repo := NewAuctionRepository(db)

	auction := &auction_entity.Auction{
		Id:          "test-auction-id-" + time.Now().Format("150405"),
		ProductName: "Test Product",
		Category:    "Test Category",
		Description: "Test Description",
		Condition:   auction_entity.New,
		Status:      auction_entity.Active,
		Timestamp:   time.Now(),
	}

	errRepo := repo.CreateAuction(ctx, auction)
	assert.Nil(t, errRepo)

	// Verify it's active initially
	found, errRepo := repo.FindAuctionById(ctx, auction.Id)
	assert.Nil(t, errRepo)
	assert.Equal(t, auction_entity.Active, found.Status)

	// Wait for closure (duration + buffer)
	time.Sleep(3 * time.Second)

	// Verify it's closed
	found, errRepo = repo.FindAuctionById(ctx, auction.Id)
	assert.Nil(t, errRepo)
	assert.Equal(t, auction_entity.Completed, found.Status)
}
