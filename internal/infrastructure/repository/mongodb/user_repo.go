package mongodb

import (
	// "github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/domain/contract"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/domain/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)
type MongoUserRepository struct {
	collection *mongo.Collection
}

func NewMongoUserRepository(collection *mongo.Collection) *MongoUserRepository {
	return &MongoUserRepository{collection: collection}
}

func (r *MongoUserRepository) CreateUser(ctx context.Context, user *entity.User) (error) {
	_,err := r.collection.InsertOne(ctx,user)
	return err
}

func (r *MongoUserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	var user entity.User
	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&user)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (r *MongoUserRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *MongoUserRepository) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *MongoUserRepository) GetByUserName(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (r *MongoUserRepository) UpdateUser(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()
	
	// Debug logging
	log.Printf("UpdateUser called with ID: %s", id.String())
	log.Printf("Updates map: %+v", updates)
	
	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"id": id},
		bson.M{"$set": updates},
	)
	if err != nil {
		log.Printf("UpdateOne error: %v", err)
		return err
	}
	
	log.Printf("UpdateOne result: MatchedCount=%d, ModifiedCount=%d", result.MatchedCount, result.ModifiedCount)
	
	if result.MatchedCount == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (r *MongoUserRepository) UpdateUserPassword(ctx context.Context, id uuid.UUID, hashedPassword string) error {
	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{"password_hash": hashedPassword}}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *MongoUserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	filter := bson.M{"id": id}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}