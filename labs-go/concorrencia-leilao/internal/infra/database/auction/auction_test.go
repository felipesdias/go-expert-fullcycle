package auction

import (
	"context"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupTestDB(t *testing.T) *mongo.Database {
	os.Setenv("AUCTION_INTERVAL", "5s")
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://admin:admin@localhost:27017"))
	assert.NoError(t, err)
	return client.Database("auctions_test")
}

func TestAuctionRepository_CloseExpiredAuctions(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		err := db.Drop(context.Background())
		if err != nil {
			t.Errorf("Failed to drop test database: %v", err)
		}
	}()

	repo := NewAuctionRepository(db)

	auction, err := auction_entity.CreateAuction("test product", "test category", "xxxxxxxxxxxxxxxxxxxxxxxxxxxx.", auction_entity.New)

	assert.Nil(t, err)

	err = repo.CreateAuction(context.Background(), auction)
	assert.Nil(t, err)

	time.Sleep(3 * time.Second)
	updatedAuction, findErr := repo.FindAuctionById(context.Background(), auction.Id)
	assert.Nil(t, findErr)
	assert.NotNil(t, updatedAuction)
	assert.Equal(t, auction_entity.Active, updatedAuction.Status)

	// Aguarda o leilão expirar e ser fechado pela rotina
	time.Sleep(5 * time.Second)

	updatedAuction, findErr = repo.FindAuctionById(context.Background(), auction.Id)
	assert.Nil(t, findErr)
	assert.NotNil(t, updatedAuction)
	assert.Equal(t, auction_entity.Completed, updatedAuction.Status)
}
