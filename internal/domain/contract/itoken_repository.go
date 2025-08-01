package contract

import (
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/domain/entity"
	"context"
)

type ITokenRepository interface {
	Create(ctx context.Context, token *entity.Token) error
	GetByID(ctx context.Context, id string) (*entity.Token, error)
	GetByUserID(ctx context.Context, userID string) (*entity.Token, error)
	Revoke(ctx context.Context, id string) error
}
