package mongodb

import (
	// "blog/internal/domain/contract"
	"blog/internal/domain/entity"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)
type MongoUserRepository struct {
	collection *mongo.Collection
}

func NewMongoUserRepository(collection *mongo.Collection) *MongoUserRepository {
	return &MongoUserRepository{collection: collection}
}

func (r *MongoUserRepository) CreateUser(ctx context.Context, user entity.User) (error) {
	_,err := r.collection.InsertOne(ctx,user)
	return err
}

func (r *MongoUserRepository) GetUserByID(ctx context.Context, id primitive.ObjectID) (*entity.User, error) {
	var user entity.User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, errors.New("user not found")
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

func (r *MongoUserRepository) UpdateUser(ctx context.Context, user entity.User) (*entity.User, error) {
	filter := bson.M{"_id": user.ID}
	setFields := bson.M{}

	if user.Username != "" {
		setFields["username"] = user.Username
	}
	if user.Email != "" {
		setFields["email"] = user.Email
	}
	if user.PasswordHash != "" {
		setFields["password_hash"] = user.PasswordHash
	}
	if user.Role != "" {
		setFields["role"] = user.Role
	}
	if user.PackageID != nil {
		setFields["package_id"] = user.PackageID
	}
	if user.PackageExpiry != nil {
		setFields["package_expiry"] = user.PackageExpiry
	}
	setFields["is_active"] = user.IsActive // bool zero value is false, so always update
	if user.FirstName != nil {
		setFields["first_name"] = user.FirstName
	}
	if user.LastName != nil {
		setFields["last_name"] = user.LastName
	}
	if user.AvatarURL != nil {
		setFields["avatar_url"] = user.AvatarURL
	}
	setFields["updated_at"] = user.UpdatedAt

	update := bson.M{
		"$set": setFields,
	}

	result := r.collection.FindOneAndUpdate(ctx, filter, update)
	if result.Err() != nil {
		return nil, result.Err()
	}

	var updatedUser entity.User
	err := result.Decode(&updatedUser)
	if err != nil {
		return nil, err
	}

	return &updatedUser, nil
}

func (r *MongoUserRepository) DeleteUser(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}