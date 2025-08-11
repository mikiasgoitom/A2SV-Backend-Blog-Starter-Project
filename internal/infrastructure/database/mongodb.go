package mongodb

import (
	"fmt"
	"log"
	"os"
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
	if err := createIndexes(ctx, client.Database(os.Getenv("MONGODB_DB_NAME"))); err != nil {
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

	// Unique index for user email
	usersCollection := db.Collection("users")
	emailIndex := mongo.IndexModel{
		Keys:    bson.M{"email": 1},
		Options: options.Index().SetUnique(true),
	}
	_, err = usersCollection.Indexes().CreateOne(ctx, emailIndex)
	if err != nil {
		return fmt.Errorf("failed to create unique index for users email: %w", err)
	}

	// Compound index for blogs: author_id + created_at (for author timeline queries)
	blogsCollection := db.Collection("blogs")
	authorCreatedIndex := mongo.IndexModel{
		Keys: bson.D{{Key: "author_id", Value: 1}},
	}
	_, err = blogsCollection.Indexes().CreateOne(ctx, authorCreatedIndex)
	if err != nil {
		return fmt.Errorf("failed to create compound index for blogs: %w", err)
	}

	// Text index for blogs: title and content (for search)
	blogTextIndex := mongo.IndexModel{
		Keys: bson.D{{Key: "title", Value: "text"}, {Key: "content", Value: "text"}},
	}
	_, err = blogsCollection.Indexes().CreateOne(ctx, blogTextIndex)
	if err != nil {
		return fmt.Errorf("failed to create text index for blogs: %w", err)
	}

	// Index for blogs._id (for fast lookup by blog id)
	blogIDIndex := mongo.IndexModel{
		Keys: bson.M{"_id": 1},
	}
	_, err = blogsCollection.Indexes().CreateOne(ctx, blogIDIndex)
	if err != nil {
		return fmt.Errorf("failed to create index for blogs._id: %w", err)
	}

	// Index for blogs.slug (for fast lookup by slug)
	slugIndex := mongo.IndexModel{
		Keys:    bson.M{"slug": 1},
		Options: options.Index().SetUnique(true),
	}
	_, err = blogsCollection.Indexes().CreateOne(ctx, slugIndex)
	if err != nil {
		return fmt.Errorf("failed to create index for blogs.slug: %w", err)
	}

	// Index for blog_tags.blog_id and blog_tags.tag_id (for tag lookups)
	blogTagsCollection := db.Collection("blog_tags")
	blogTagIndex := mongo.IndexModel{
		Keys: bson.D{{Key: "blog_id", Value: 1}, {Key: "tag_id", Value: 1}},
	}
	_, err = blogTagsCollection.Indexes().CreateOne(ctx, blogTagIndex)
	if err != nil {
		return fmt.Errorf("failed to create index for blog_tags: %w", err)
	}

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
