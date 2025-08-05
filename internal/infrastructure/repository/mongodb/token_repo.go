package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/domain/contract"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/domain/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// ---------- DTO layer ------------------
type tokenDTO struct {
	ID        string    `bson:"_id,omitempty"`
	UserID    string    `bson:"user_id"`
	TokenType string    `bson:"token_type"`
	TokenHash string    `bson:"token_hash"`
	CreatedAt time.Time `bson:"created_at"`
	ExpiresAt time.Time `bson:"expires_at"`
	Revoke    bool      `bson:"revoke"`
}

// ...existing code...
func (t *tokenDTO) ToEntity() *entity.Token {
	return &entity.Token{
		ID:        t.ID,
		UserID:    t.UserID,
		TokenType: entity.TokenType(t.TokenType),
		TokenHash: t.TokenHash,
		CreatedAt: t.CreatedAt,
		ExpiresAt: t.ExpiresAt,
		Revoke:    t.Revoke,
	}
}

func FromTokenEntityToDTO(t *entity.Token) *tokenDTO {
	return &tokenDTO{
		ID:        t.ID,
		UserID:    t.UserID,
		TokenType: string(t.TokenType),
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
