package mongodb

import (
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
)

type MongoDBClient struct {
	Client *mongo.Client
}

func NewMongoDBClient(uri string) (*MongoDBClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Println("Failed to connect to MongoDB:", err)
		return nil, err
	}

	// Ping the database to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		log.Println("Failed to ping MongoDB:", err)
		return nil, err
	}

	// Create indexes
	if err := createIndexes(ctx, client.Database("blogdb")); err != nil {
		log.Println("Failed to create indexes:", err)
		// We can choose to return the error or just log it
	}

	return &MongoDBClient{Client: client}, nil
}

func createIndexes(ctx context.Context, db *mongo.Database) error {
	// TTL index for blog_views
	blogViewsCollection := db.Collection("blog_views")
	ttlIndex := mongo.IndexModel{
		Keys:    bson.M{"viewed_at": 1},
		Options: options.Index().SetExpireAfterSeconds(24 * 60 * 60), // 24 hours
	}
	_, err := blogViewsCollection.Indexes().CreateOne(ctx, ttlIndex)
	if err != nil {
		return fmt.Errorf("failed to create TTL index for blog_views: %w", err)
	}

	// We can add other index creations here, for example, for users, blogs, etc.
	// Example: Unique index for user email
	// usersCollection := db.Collection("users")
	// emailIndex := mongo.IndexModel{
	// 	Keys:    bson.M{"email": 1},
	// 	Options: options.Index().SetUnique(true),
	// }
	// _, err = usersCollection.Indexes().CreateOne(ctx, emailIndex)
	// if err != nil {
	// 	return fmt.Errorf("failed to create unique index for users email: %w", err)
	// }

	log.Println("Successfully created database indexes.")
	return nil
}

// Disconnect disconnects the MongoDB client
func (m *MongoDBClient) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()
	return m.Client.Disconnect(ctx)
}

// GetCollection returns a collection from the database
func (m *MongoDBClient) GetCollection(dbName, collectionName string) *mongo.Collection {
	return m.Client.Database(dbName).Collection(collectionName)
}
