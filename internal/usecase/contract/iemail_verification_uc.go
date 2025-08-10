package usecasecontract

import (
	"context"

	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/domain/entity"
)

type IEmailVerificationUC interface {
	RequestVerificationEmail(ctx context.Context, user *entity.User) error
	VerifyEmailToken(ctx context.Context, verifier, plainToken string) (*entity.User, error)
}
