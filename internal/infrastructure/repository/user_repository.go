package repository

// import (
// 	"context"
// 	"errors"

// 	"github.com/google/uuid"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	// "go.mongodb.org/mongo-driver/mongo/options"

// 	"a2sv-blog/internal/domain/contract"
// 	"a2sv-blog/internal/entity"
// 	"a2sv-blog/internal/infrastructure/repository/mongodb"
// )

// type userRepository struct {
// 	db         *mongodb.MongoDB
// 	collection *mongo.Collection
// }

// func NewUserRepository(db *mongodb.MongoDB) contract.IUserRepository {
// 	return &userRepository{
// 		db:         db,
// 		collection: db.Client.Database("your_database_name").Collection("users"),
// 	}
// }

// func (r *userRepository) CreateUser(ctx context.Context, user *entity.User) error {
// 	_, err := r.collection.InsertOne(ctx, user)
// 	return err
// }

// func (r *userRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
// 	filter := bson.M{"id": id}
// 	var user entity.User
// 	err := r.collection.FindOne(ctx, filter).Decode(&user)
// 	if err != nil {
// 		if errors.Is(err, mongo.ErrNoDocuments) {
// 			return nil, nil
// 		}
// 		return nil, err
// 	}
// 	return &user, nil
// }

// func (r *userRepository) UpdateUser(ctx context.Context, user *entity.User) error {
// 	filter := bson.M{"id": user.ID}
// 	update := bson.M{"$set": user}
// 	_, err := r.collection.UpdateOne(ctx, filter, update)
// 	return err
// }

// func (r *userRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
// 	filter := bson.M{"id": id}
// 	_, err := r.collection.DeleteOne(ctx, filter)
// 	return err
// }

// func (r *userRepository) FindUserByEmail(ctx context.Context, email string) (*entity.User, error) {
// 	filter := bson.M{"email": email}
// 	var user entity.User
// 	err := r.collection.FindOne(ctx, filter).Decode(&user)
// 	if err != nil {
// 		if errors.Is(err, mongo.ErrNoDocuments) {
// 			return nil, nil
// 		}
// 		return nil, err
// 	}
// 	return &user, nil
// }
