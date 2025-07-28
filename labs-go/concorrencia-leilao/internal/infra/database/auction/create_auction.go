package auction

import (
	"context"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
}
type AuctionRepository struct {
	Collection *mongo.Collection
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	repo := &AuctionRepository{
		Collection: database.Collection("auctions"),
	}
	go repo.closeAuctionsRoutine(context.Background())
	return repo
}

func (ar *AuctionRepository) CreateAuction(
	ctx context.Context,
	auctionEntity *auction_entity.Auction) *internal_error.InternalError {
	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   auctionEntity.Timestamp.Unix(),
	}
	_, err := ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error trying to insert auction", err)
		return internal_error.NewInternalServerError("Error trying to insert auction")
	}

	return nil
}

func (ar *AuctionRepository) closeAuctionsRoutine(ctx context.Context) {
	auctionInterval := getAuctionInterval()
	// Usando um ticker para verificar a cada 5 segundos
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			logger.Info("Checking for expired auctions")
			auctions, err := ar.FindAuctions(ctx, auction_entity.Active, "", "")
			if err != nil {
				logger.Error("Error finding active auctions in routine", err, zap.String("journey", "closeAuctionRoutine"))
				continue
			}

			for _, auction := range auctions {
				auctionEndTime := auction.Timestamp.Add(auctionInterval)
				if time.Now().After(auctionEndTime) {
					logger.Info("Auction expired, closing it", zap.String("auctionId", auction.Id))
					if err := ar.UpdateAuctionStatus(ctx, auction.Id, auction_entity.Completed); err != nil {
						logger.Error("Error closing auction in routine", err, zap.String("auctionId", auction.Id))
					}
				}
			}
		case <-ctx.Done():
			logger.Info("Closing auction routine")
			return
		}
	}
}

func getAuctionInterval() time.Duration {
	durationStr := os.Getenv("AUCTION_INTERVAL")
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		// Valor padrão de 1 minuto caso a variável de ambiente não seja definida ou seja inválida
		return time.Minute
	}
	return duration
}

func (ar *AuctionRepository) UpdateAuctionStatus(
	ctx context.Context,
	id string,
	status auction_entity.AuctionStatus) *internal_error.InternalError {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"status": status}}

	_, err := ar.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Error("Error trying to update auction status", err)
		return internal_error.NewInternalServerError("Error trying to update auction status")
	}

	return nil
}
