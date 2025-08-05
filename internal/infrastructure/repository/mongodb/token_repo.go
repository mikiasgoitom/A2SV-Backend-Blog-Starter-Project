package mongodb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
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

func (t *tokenDTO) ToEntity() *entity.Token {
	userID, _ := uuid.Parse(t.UserID) // handle error as needed
	id, _ := uuid.Parse(t.ID)         // handle error as needed
	return &entity.Token{
		ID:        id,
		UserID:    userID,
		TokenType: entity.TokenType(t.TokenType),
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

func (r *TokenRepository) CreateToken(ctx context.Context, token *entity.Token) error {
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

// GetByUserID fetches token by user's ID string
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

// GetTokenByUserID retrieves a token by user ID.
func (r *TokenRepository) GetTokenByUserID(ctx context.Context, userID uuid.UUID) (*entity.Token, error) {
	var dto tokenDTO
	// Query using the string representation of the UUID since that's how it's stored in MongoDB
	err := r.Collection.FindOne(ctx, bson.M{"user_id": userID.String()}).Decode(&dto)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("token not found")
		}
		return nil, err
	}
	return dto.ToEntity(), nil
}

func (r *TokenRepository) DeleteToken(ctx context.Context, id uuid.UUID) error {
	filter := bson.M{"_id": id.String()}
	_, err := r.Collection.DeleteOne(ctx, filter)
	return err
}

// UpdateToken updates the token hash and expiry
func (r *TokenRepository) UpdateToken(ctx context.Context, tokenID uuid.UUID, tokenHash string, expiry time.Time) error {
	filter := bson.M{"_id": tokenID.String()}
	update := bson.M{"$set": bson.M{"token_hash": tokenHash, "expires_at": expiry}}
	_, err := r.Collection.UpdateOne(ctx, filter, update)
	return err
}

// Revoke marks a token as revoked by its ID
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
