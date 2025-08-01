package mongodb

import (
	"blog/internal/domain/contract"
	"blog/internal/domain/entity"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// ---------- DTO layer ------------------
type tokenDTO struct {
	ID        string    `bson:"_id,omitempty"`
	UserID    string    `bson:"user_id"`
	TokenHash string    `bson:"token_hash"`
	CreatedAt time.Time `bson:"created_at"`
	ExpiresAt time.Time `bson:"expires_at"`
	Revoke    bool      `bson:"revoke"`
}

// ...existing code...
func (t *tokenDTO) ToEntity() *entity.Token {
	userID, _ := uuid.Parse(t.UserID) // handle error as needed
	id, _ := uuid.Parse(t.ID)         // handle error as needed
	return &entity.Token{
		ID:        id,
		UserID:    userID,
		TokenHash: t.TokenHash,
		CreatedAt: t.CreatedAt,
		ExpiresAt: t.ExpiresAt,
		Revoke:    t.Revoke,
	}
}

func FromTokenEntityToDTO(t *entity.Token) *tokenDTO {
	return &tokenDTO{
		ID:        t.ID.String(),
		UserID:    t.UserID.String(),
		TokenHash: t.TokenHash,
		CreatedAt: t.CreatedAt,
		ExpiresAt: t.ExpiresAt,
		Revoke:    t.Revoke,
	}
}

// ---------------------------------------

type TokenRepository struct {
	Collection *mongo.Collection
}

// check in compile time if TokenRepository implements ITokenRepository
var _ contract.ITokenRepository = (*TokenRepository)(nil)

func NewTokenRepository(colln *mongo.Collection) *TokenRepository {
	return &TokenRepository{
		Collection: colln,
	}
}

func (r *TokenRepository) Create(ctx context.Context, token *entity.Token) error {
	dto := FromTokenEntityToDTO(token)
	_, err := r.Collection.InsertOne(ctx, dto)
	if err != nil {
		return err
	}

	return nil
}

func (r *TokenRepository) GetByID(ctx context.Context, id string) (*entity.Token, error) {
	filter := bson.M{"_id": id}
	var dto tokenDTO
	err := r.Collection.FindOne(ctx, filter).Decode(&dto)
	if err != nil {
		return nil, err
	}
	token := dto.ToEntity()

	return token, nil
}

func (r *TokenRepository) GetByUserID(ctx context.Context, userID string) (*entity.Token, error) {
	filter := bson.M{"user_id": userID}
	var dto tokenDTO
	err := r.Collection.FindOne(ctx, filter).Decode(&dto)
	if err != nil {
		return nil, err
	}
	token := dto.ToEntity()

	return token, nil
}

func (r *TokenRepository) Revoke(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"revoke": true}}
	result, err := r.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("failed to revoke token with: %v", id)
	}

	return nil
}
