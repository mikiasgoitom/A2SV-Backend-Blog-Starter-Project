package mongodb

// import (
// 	"context"
// 	"time"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// 	"log"
// )

// type MongoDB struct {
// 	Client *mongo.Client
// }

// func NewMongoDB(ctx context.Context, uri string) (*MongoDB, error) {
// 	clientOptions := options.Client().ApplyURI(uri)
// 	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
// 	defer cancel()

// 	client, err := mongo.Connect(ctx, clientOptions)
// 	if err != nil {
// 		log.Printf("Failed to connect to MongoDB: %v", err)
// 		return nil, err
// 	}

// 	err = client.Ping(ctx, nil)
// 	if err != nil {
// 		log.Printf("Failed to ping MongoDB: %v", err)
// 		return nil, err
// 	}

// 	return &MongoDB{
// 		Client: client,
// 	}, nil
// }

// func (db *MongoDB) Disconnect(ctx context.Context) error {
// 	return db.Client.Disconnect(ctx)
// }
